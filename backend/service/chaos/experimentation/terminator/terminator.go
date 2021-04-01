package terminator

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
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

var CriteriaFactories = map[string]CriteriaFactory{
	TypeUrl(&terminatorv1.MaxTimeTerminationCriteria{}): &maxTimeTerminationFactory{},
}

type CriteriaFactory interface {
	Create(cfg *any.Any) (TerminationCriteria, error)
}

type TerminationCriteria interface {
	// ShouldTerminate determines whether the provided experiment should be terminated.
	// To signal that termination should occur, return a non-empty string with a nil error.
	ShouldTerminate(experiment *experimentationv1.Experiment, experimentConfig proto.Message) (string, error)
}

func New(cfg *any.Any, logger *zap.Logger, scope tally.Scope) (service.Service, error) {
	m, err := NewMonitor(cfg, logger, scope)
	if err != nil {
		m.Run(context.Background())
	}

	return m, nil
}

func NewMonitor(cfg *any.Any, logger *zap.Logger, scope tally.Scope) (*Monitor, error) {
	typedConfig := &terminatorv1.Config{}
	err := ptypes.UnmarshalAny(cfg, typedConfig)
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

	terminationCriteria := map[string][]TerminationCriteria{}

	for configType, perConfigTypeConfig := range typedConfig.PerConfigTypeConfiguration {
		perConfigCriterias := []TerminationCriteria{}

		for _, c := range perConfigTypeConfig.TerminationCriteria {
			factory, ok := CriteriaFactories[c.TypeUrl]
			if !ok {
				return nil, fmt.Errorf("terminator module configured with unknown criteria '%s'", c.TypeUrl)
			}

			criteria, err := factory.Create(c)
			if err != nil {
				return nil, fmt.Errorf("failed to create termination criteria '%s': %s", c.TypeUrl, err)
			}

			perConfigCriterias = append(perConfigCriterias, criteria)
		}

		terminationCriteria[configType] = perConfigCriterias
	}

	outerLoopInterval, err := ptypes.Duration(typedConfig.OuterLoopInterval)
	if err != nil {
		return nil, err
	}

	perExperimentCheckInterval, err := ptypes.Duration(typedConfig.PerExperimentCheckInterval)
	if err != nil {
		return nil, err
	}

	return &Monitor{
		store:                        storer,
		terminationCriteriaByTypeUrl: terminationCriteria,
		outerLoopInterval:            outerLoopInterval,
		perExperimentCheckInterval:   perExperimentCheckInterval,
		log:                          logger.Sugar(),
		activeMonitoringRoutines:     trackingGauge{gauge: scope.Gauge("active_monitoring_routines")},
		criteriaEvaluationSuccess:    scope.Counter("criteria_success"),
		criteriaEvaluationFailure:    scope.Counter("criteria_failure"),
		terminationCount:             scope.Counter("terminations"),
		marshallingErrors:            scope.Counter("unpack_error"),
	}, nil
}

type Monitor struct {
	store                        experimentstore.Storer
	terminationCriteriaByTypeUrl map[string][]TerminationCriteria

	outerLoopInterval          time.Duration
	perExperimentCheckInterval time.Duration

	log *zap.SugaredLogger

	criteriaEvaluationSuccess tally.Counter
	criteriaEvaluationFailure tally.Counter
	activeMonitoringRoutines  trackingGauge
	terminationCount          tally.Counter
	marshallingErrors         tally.Counter
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
		go func(configType string, criteria []TerminationCriteria) {
			trackedExperiments := map[uint64]context.CancelFunc{}
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
func (m *Monitor) monitorNewExperiments(es []*experimentationv1.Experiment, trackedExperiments map[uint64]context.CancelFunc, criterias []TerminationCriteria) map[uint64]struct{} {
	// For each active experiment, create a monitoring goroutine if necessary.
	activeExperiments := map[uint64]struct{}{}
	for _, e := range es {
		activeExperiments[e.Id] = struct{}{}
		if _, ok := trackedExperiments[e.Id]; !ok {
			ctx, cancel := context.WithCancel(context.Background())
			trackedExperiments[e.Id] = cancel

			m.activeMonitoringRoutines.inc()
			go func() {
				defer m.activeMonitoringRoutines.dec()
				m.monitorSingleExperiment(ctx, e, criterias)
			}()
		}
	}
	return activeExperiments
}

func (m *Monitor) monitorSingleExperiment(ctx context.Context, e *experimentationv1.Experiment, criterias []TerminationCriteria) {
	ticker := time.NewTicker(m.perExperimentCheckInterval)
	terminated := false

	// TODO(snowp): I suspect that we might run into issues here if we try to unpack an "old" proto
	// which is no longer linked into the program due to how the type registry works, but this should
	// be rare enough that logging an error is probably fine.
	unpackedConfig, err := anypb.UnmarshalNew(e.Config, proto.UnmarshalOptions{})
	if err != nil {
		m.log.Errorw("failed to unmarshal experiment", "id", e.Id)
		m.marshallingErrors.Inc(1)
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
			for _, c := range criterias {
				terminationReason, err := c.ShouldTerminate(e, unpackedConfig)
				// TODO(snowp): The logs here might get spammy, rate limit or terminate montioring routine somehow?
				if err != nil {
					m.criteriaEvaluationFailure.Inc(1)
					m.log.Errorw("error while evaluating termination criteria", "id", e.Id, "error", err)
					continue
				}

				m.criteriaEvaluationSuccess.Inc(1)

				if terminationReason != "" {
					err = m.store.CancelExperimentRun(context.Background(), e.Id, terminationReason)
					if err != nil {
						m.log.Errorw("failed to terminate experiment", "err", err, "experimentId", e.Id)
						continue
					}
					m.log.Errorw("terminated experiment", "experimentId", e.Id)
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
