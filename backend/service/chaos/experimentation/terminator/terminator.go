package terminator

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/uber-go/tally"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	experimentationv1 "github.com/lyft/clutch/backend/api/chaos/experimentation/v1"
	"github.com/lyft/clutch/backend/service/chaos/experimentation/experimentstore"
)

type TerminationCriteria interface {
	// ShouldTerminate determines whether the provided experiment should be terminated.
	// To signal that termination should occur, return a non-empty string with a nil error.
	ShouldTerminate(experiment *experimentationv1.Experiment, experimentConfig proto.Message) (string, error)
}

type Monitor interface {
	Run(ctx context.Context)
}

// TODO(snowp): Remove this once we have a proper service object that we can create.
func NewTestMonitor(store experimentstore.Storer, enabledConfigTypes []string, criterias []TerminationCriteria, log *zap.SugaredLogger, stats tally.Scope) Monitor {
	return &monitor{
		store:                      store,
		enabledConfigTypes:         enabledConfigTypes,
		criterias:                  criterias,
		outerLoopInterval:          1,
		perExperimentCheckInterval: 1,
		log:                        log,
		activeMonitoringRoutines:   trackingGauge{gauge: stats.Gauge("active_monitoring_routines")},
		criteriaEvaluationSuccess:  stats.Counter("criteria_success"),
		criteriaEvaluationFailure:  stats.Counter("criteria_failure"),
		terminationCount:           stats.Counter("terminations"),
		marshallingErrors:          stats.Counter("unpack_error"),
	}
}

type monitor struct {
	store              experimentstore.Storer
	enabledConfigTypes []string
	criterias          []TerminationCriteria

	outerLoopInterval          time.Duration
	perExperimentCheckInterval time.Duration

	log *zap.SugaredLogger

	criteriaEvaluationSuccess tally.Counter
	criteriaEvaluationFailure tally.Counter
	activeMonitoringRoutines  trackingGauge
	terminationCount          tally.Counter
	marshallingErrors         tally.Counter
}

func (m *monitor) Run(ctx context.Context) {
	for _, configType := range m.enabledConfigTypes {
		// For each monitored config type, start a single goroutine that polls the active experiments at a fixed interval.
		// Whenever a new (new to this goroutine) experiment is found, open up a goroutine that periodically evaluates all
		// the termination criteria for the experiment.
		//
		// This approach ensures provides a steady DB pressure (mostly outer loop, some from triggering termination) and relatively high
		// fairness, as checking the termination conditions for one experiment should not be delaying the checks for another experiment
		// unless we're under very heavy load.
		go func(configType string) {
			trackedExperiments := map[string]context.CancelFunc{}
			ticker := time.NewTicker(m.outerLoopInterval)

			for {
				select {
				case <-ctx.Done():
					return
				case <-ticker.C:
					es, err := m.store.GetExperiments(context.Background(), configType, experimentationv1.GetExperimentsRequest_STATUS_RUNNING)
					if err != nil {
						m.log.Errorw("failed to retrieve experiments from experiment store", "err", err, "enableConfigTypes", m.enabledConfigTypes)
						continue
					}

					activeExperiments := m.monitorNewExperiments(es, trackedExperiments)

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
		}(configType)
	}
}

// Iterates over all the provided experiments, spawning a goroutine to montior each experiment that doesn't already have
// a monitoring routine. Returns a set containing all the active experiment ids for further processing.
func (m *monitor) monitorNewExperiments(es []*experimentationv1.Experiment, trackedExperiments map[string]context.CancelFunc) map[string]struct{} {
	// For each active experiment, create a monitoring goroutine if necessary.
	activeExperiments := map[string]struct{}{}
	for _, e := range es {
		activeExperiments[e.RunId] = struct{}{}
		if _, ok := trackedExperiments[e.RunId]; !ok {
			ctx, cancel := context.WithCancel(context.Background())
			trackedExperiments[e.RunId] = cancel

			m.activeMonitoringRoutines.inc()
			go func() {
				defer m.activeMonitoringRoutines.dec()
				m.monitorSingleExperiment(ctx, e)
			}()
		}
	}
	return activeExperiments
}

func (m *monitor) monitorSingleExperiment(ctx context.Context, e *experimentationv1.Experiment) {
	ticker := time.NewTicker(m.perExperimentCheckInterval)
	terminated := false

	// TODO(snowp): I suspect that we might run into issues here if we try to unpack an "old" proto
	// which is no longer linked into the program due to how the type registry works, but this should
	// be rare enough that logging an error is probably fine.
	unpackedConfig, err := anypb.UnmarshalNew(e.Config, proto.UnmarshalOptions{})
	if err != nil {
		m.log.Errorw("failed to unmarshal experiment", "runId", e.RunId)
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
			for _, c := range m.criterias {
				terminationReason, err := c.ShouldTerminate(e, unpackedConfig)
				// TODO(snowp): The logs here might get spammy, rate limit or terminate montioring routine somehow?
				if err != nil {
					m.criteriaEvaluationFailure.Inc(1)
					m.log.Errorw("error while evaluating termination criteria", "runId", e.RunId, "error", err)
					continue
				}

				m.criteriaEvaluationSuccess.Inc(1)

				if terminationReason != "" {
					err = m.store.CancelExperimentRun(context.Background(), e.RunId, terminationReason)
					if err != nil {
						m.log.Errorw("failed to terminate experiment", "err", err, "experimentRunId", e.RunId)
						continue
					}
					m.log.Errorw("terminated experiment", "experimentRunId", e.RunId)
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
