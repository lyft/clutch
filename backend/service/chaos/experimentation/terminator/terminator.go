package terminator

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/uber-go/tally"
	"go.uber.org/zap"

	experimentationv1 "github.com/lyft/clutch/backend/api/chaos/experimentation/v1"
	"github.com/lyft/clutch/backend/service/chaos/experimentation/experimentstore"
)

type terminationCriteria interface {
	ShouldTerminate(experimentStarted time.Time, config interface{}) error
}

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

type monitor struct {
	store              experimentstore.Storer
	enabledConfigTypes []string
	criterias          []terminationCriteria

	outerLoopInterval          time.Duration
	perExperimentCheckInterval time.Duration

	log *zap.SugaredLogger

	activeMonitoringRoutines trackingGauge
	terminationCount         tally.Counter
}

func (t *monitor) run(ctx context.Context) {
	// Start a single goroutine that monitors all the active experiments at a fixed interval. Whenever a new (new to this goroutine)
	// experiment is found, open up a goroutine that periodically evaluates all the termination criteria for the experiment.
	//
	// This approach ensures provides a steady DB pressure (mostly outer loop, some from triggering termination) and relatively high
	// fairness, as checking the termination conditions for one experiment should not be delaying the checks for another experiment
	// unless we're under very heavy load.
	go func(store experimentstore.Storer, enabledConfigTypes []string, criterias []terminationCriteria, activeMonitoringRoutines *trackingGauge) {
		trackedExperiments := map[uint64]context.CancelFunc{}
		ticker := time.NewTicker(t.outerLoopInterval)

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				es, err := store.GetExperiments(context.Background(), enabledConfigTypes, experimentationv1.GetExperimentsRequest_STATUS_RUNNING)
				if err != nil {
					t.log.Errorw("failed to retrieve experiments from experiment store", "err", err, "enableConfigTypes", enabledConfigTypes)
					continue
				}

				// For each active experiment, create a monitoring goroutine if necessary.
				activeExperiments := map[uint64]struct{}{}
				for _, e := range es {
					activeExperiments[e.Id] = struct{}{}
					if _, ok := trackedExperiments[e.Id]; !ok {
						ctx, cancel := context.WithCancel(context.Background())
						trackedExperiments[e.Id] = cancel

						activeMonitoringRoutines.inc()
						go func() {
							defer t.activeMonitoringRoutines.dec()

							ticker := time.NewTicker(t.perExperimentCheckInterval)
							terminated := false

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
										tt, _ := ptypes.Timestamp(e.StartTime)
										err := c.ShouldTerminate(tt, e)
										if err != nil {
											err := store.TerminateExperiment(context.Background(), e.Id, err.Error())
											if err != nil {
												t.log.Errorw("failed to terminate experiment", "err", err, "experimentId", e.Id)
											} else {
												t.terminationCount.Inc(1)
												terminated = true
											}
										}
									}
								case <-ctx.Done():
									return
								}
							}
						}()
					}
				}

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
	}(t.store, t.enabledConfigTypes, t.criterias, &t.activeMonitoringRoutines)
}
