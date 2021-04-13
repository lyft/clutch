package terminator

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/uber-go/tally"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	experimentationv1 "github.com/lyft/clutch/backend/api/chaos/experimentation/v1"
	terminatorv1 "github.com/lyft/clutch/backend/api/config/service/chaos/experimentation/terminator/v1"
	"github.com/lyft/clutch/backend/service"
	"github.com/lyft/clutch/backend/service/chaos/experimentation/experimentstore"
)

const Name = "clutch.module.chaos.experimentation.termination"

func TypeUrl(message proto.Message) string {
	return "type.googleapis.com/" + string(message.ProtoReflect().Descriptor().FullName())
}

var CriterionFactories = map[string]CriterionFactory{
	TypeUrl(&terminatorv1.MaxTimeTerminationCriterion{}): &maxTimeTerminationFactory{},
}

type CriterionFactory interface {
	Create(cfg *anypb.Any) (TerminationCriterion, error)
}

type TerminationCriterion interface {
	// ShouldTerminate determines whether the provided experiment should be terminated.
	// To signal that termination should occur, return a non-empty string with a nil error.
	ShouldTerminate(experiment *experimentstore.Experiment, experimentConfig proto.Message) (string, error)
}

func New(cfg *anypb.Any, logger *zap.Logger, scope tally.Scope) (service.Service, error) {
	m, err := NewMonitor(cfg, logger, scope)
	if err != nil {
		m.Run(context.Background())
	}

	return m, nil
}

func NewMonitor(cfg *anypb.Any, logger *zap.Logger, scope tally.Scope) (*Monitor, error) {
	typedConfig := &terminatorv1.Config{}
	err := cfg.UnmarshalTo(typedConfig)
	if err != nil {
		return nil, err
	}

	store, ok := service.Registry[experimentstore.Name]
	if !ok {
		return nil, errors.New("could not find experiment store service")
	}

	storer, ok := store.(experimentstore.Storer)
	if !ok {
		return nil, errors.New("service was not the correct type")
	}

	terminationCriteria := map[string][]TerminationCriterion{}

	for configType, perConfigTypeConfig := range typedConfig.PerConfigTypeConfiguration {
		perConfigCriteria := []TerminationCriterion{}

		for _, c := range perConfigTypeConfig.TerminationCriteria {
			factory, ok := CriterionFactories[c.TypeUrl]
			if !ok {
				return nil, fmt.Errorf("terminator module configured with unknown criterion '%s'", c.TypeUrl)
			}

			criterion, err := factory.Create(c)
			if err != nil {
				return nil, fmt.Errorf("failed to create termination criteria '%s': %s", c.TypeUrl, err)
			}

			perConfigCriteria = append(perConfigCriteria, criterion)
		}

		terminationCriteria[configType] = perConfigCriteria
	}

	return &Monitor{
		store:                           storer,
		terminationCriteriaByTypeUrl:    terminationCriteria,
		outerLoopInterval:               typedConfig.OuterLoopInterval.AsDuration(),
		perExperimentCheckInterval:      typedConfig.PerExperimentCheckInterval.AsDuration(),
		log:                             logger.Sugar(),
		activeMonitoringRoutines:        trackingGauge{gauge: scope.Gauge("active_monitoring_routines")},
		criterionEvaluationSuccessCount: scope.Counter("criterion_success"),
		criterionEvaluationFailureCount: scope.Counter("criterion_failure"),
		terminationCount:                scope.Counter("terminations"),
		marshallingErrorCount:           scope.Counter("unpack_error"),
	}, nil
}

type Monitor struct {
	store                        experimentstore.Storer
	terminationCriteriaByTypeUrl map[string][]TerminationCriterion

	outerLoopInterval          time.Duration
	perExperimentCheckInterval time.Duration

	log *zap.SugaredLogger

	activeMonitoringRoutines        trackingGauge
	criterionEvaluationSuccessCount tally.Counter
	criterionEvaluationFailureCount tally.Counter
	terminationCount                tally.Counter
	marshallingErrorCount           tally.Counter
}

func (m *Monitor) Run(ctx context.Context) {
	for configType, criteria := range m.terminationCriteriaByTypeUrl {
		// For each monitored config type, start a single goroutine that polls the active experiments at a fixed interval.
		// Whenever a new (new to this goroutine) experiment is found, open up a goroutine that periodically evaluates all
		// the termination criteria for the experiment.
		//
		// This approach ensures provides a steady DB pressure (mostly outer loop, some from triggering termination) and relatively high
		// fairness, as checking the termination conditions for one experiment should not be delaying the checks for another experiment
		// unless we're under very heavy load.
		go func(configType string, criteria []TerminationCriterion) {
			trackedExperiments := map[string]context.CancelFunc{}
			ticker := time.NewTicker(m.outerLoopInterval)

			for {
				select {
				case <-ctx.Done():
					return
				case <-ticker.C:
					es, err := m.store.GetExperiments(context.Background(), configType, experimentationv1.GetExperimentsRequest_STATUS_RUNNING)
					if err != nil {
						m.log.Errorw("failed to retrieve experiments from experiment store", "err", err, "configType", configType)
						continue
					}

					activeExperiments := m.monitorNewExperiments(es, trackedExperiments, criteria)

					// For all experiments that we're tracking that no longer appear to be active, cancel the goroutine
					// and clean it up.
					for e, cancel := range trackedExperiments {
						if _, ok := activeExperiments[e]; !ok {
							cancel()
							delete(trackedExperiments, e)
						}
					}
				}
			}
		}(configType, criteria)
	}
}

// Iterates over all the provided experiments, spawning a goroutine to montior each experiment that doesn't already have
// a monitoring routine. Returns a set containing all the active experiment ids for further processing.
func (m *Monitor) monitorNewExperiments(es []*experimentstore.Experiment, trackedExperiments map[string]context.CancelFunc, criteria []TerminationCriterion) map[string]struct{} {
	// For each active experiment, create a monitoring goroutine if necessary.
	activeExperiments := map[string]struct{}{}
	for _, e := range es {
		activeExperiments[e.Run.Id] = struct{}{}
		if _, ok := trackedExperiments[e.Run.Id]; !ok {
			ctx, cancel := context.WithCancel(context.Background())
			trackedExperiments[e.Run.Id] = cancel

			m.activeMonitoringRoutines.inc()
			go func() {
				defer m.activeMonitoringRoutines.dec()
				m.monitorSingleExperiment(ctx, e, criteria)
			}()
		}
	}
	return activeExperiments
}

func (m *Monitor) monitorSingleExperiment(ctx context.Context, e *experimentstore.Experiment, criteria []TerminationCriterion) {
	ticker := time.NewTicker(m.perExperimentCheckInterval)
	terminated := false

	// TODO(snowp): I suspect that we might run into issues here if we try to unpack an "old" proto
	// which is no longer linked into the program due to how the type registry works, but this should
	// be rare enough that logging an error is probably fine.
	unpackedConfig, err := anypb.UnmarshalNew(e.Config.Config, proto.UnmarshalOptions{})
	if err != nil {
		m.log.Errorw("failed to unmarshal experiment", "runId", e.Run.Id)
		m.marshallingErrorCount.Inc(1)
		return
	}

	for {
		select {
		case <-ticker.C:
			if terminated {
				// Note that we don't exit the goroutine here and instead wait for the outer
				// monitoring loop to detect that this experiment is no longer active, which will
				// cause this goroutine to be canceled. Exiting here is not helpful since the outer
				// loop can race and restart this goroutine.
				continue
			}
			for _, c := range criteria {
				terminationReason, err := c.ShouldTerminate(e, unpackedConfig)
				// TODO(snowp): The logs here might get spammy, rate limit or terminate montioring routine somehow?
				if err != nil {
					m.criterionEvaluationFailureCount.Inc(1)
					m.log.Errorw("error while evaluating termination criterion", "runId", e.Run.Id, "error", err)
					continue
				}

				m.criterionEvaluationSuccessCount.Inc(1)

				if terminationReason != "" {
					err = m.store.CancelExperimentRun(context.Background(), e.Run.Id, terminationReason)
					if err != nil {
						m.log.Errorw("failed to terminate experiment", "err", err, "experimentRunId", e.Run.Id)
						continue
					}
					m.log.Errorw("terminated experiment", "experimentRunId", e.Run.Id)
					m.terminationCount.Inc(1)
					terminated = true
				}
			}
		case <-ctx.Done():
			return
		}
	}
}

// Helper type for tracking an atomic value that updates a gauge whenever it changes.
type trackingGauge struct {
	gauge tally.Gauge
	value uint32
}

func (t *trackingGauge) inc() {
	t.gauge.Update(float64(atomic.AddUint32(&t.value, 1)))
}

func (t *trackingGauge) dec() {
	t.gauge.Update(float64(atomic.AddUint32(&t.value, ^uint32(0))))
}
