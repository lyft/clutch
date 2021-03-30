package terminator

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/stretchr/testify/assert"
	"github.com/uber-go/tally"
	"go.uber.org/zap"

	experimentationv1 "github.com/lyft/clutch/backend/api/chaos/experimentation/v1"
	serverexperimentationv1 "github.com/lyft/clutch/backend/api/chaos/serverexperimentation/v1"
	"github.com/lyft/clutch/backend/mock/service/chaos/experimentation/experimentstoremock"
)

func TestTerminator(t *testing.T) {
	l, err := zap.NewDevelopment()
	assert.NoError(t, err)

	testScope := tally.NewTestScope("", map[string]string{})

	store := &experimentstoremock.SimpleStorer{}
	criteria := &testCriteria{}
	monitor := monitor{
		store:                      store,
		enabledConfigTypes:         []string{"type.googleapis.com/clutch.chaos.serverexperimentation.v1.HTTPFaultConfig"},
		criterias:                  []TerminationCriteria{criteria},
		outerLoopInterval:          time.Millisecond,
		perExperimentCheckInterval: time.Millisecond,
		log:                        l.Sugar(),
		activeMonitoringRoutines: trackingGauge{
			gauge: testScope.Gauge("active_routines"),
			value: 0,
		},
		terminationCount: testScope.Counter("terminations"),
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	monitor.Run(ctx)

	createTestExperiment(t, 500, store)

	awaitGaugeValue(ctx, t, testScope, "active_routines", 1)

	es, err := store.GetExperiments(ctx, "type.googleapis.com/clutch.chaos.serverexperimentation.v1.HTTPFaultConfig", experimentationv1.GetExperimentsRequest_STATUS_RUNNING)
	assert.NoError(t, err)
	assert.Len(t, es, 1)

	criteria.update(true)

	// Ensure that we tear down the monitoring goroutines once we're done with an experiment.
	awaitGaugeValue(ctx, t, testScope, "active_routines", 0)
	assert.Equal(t, int64(1), testScope.Snapshot().Counters()["terminations+"].Value())

	es, err = store.GetExperiments(ctx, "type.googleapis.com/clutch.chaos.serverexperimentation.v1.HTTPFaultConfig", experimentationv1.GetExperimentsRequest_STATUS_RUNNING)
	assert.NoError(t, err)
	assert.Len(t, es, 0)
}

func createTestExperiment(t *testing.T, faultHttpStatus int, storer *experimentstoremock.SimpleStorer) *experimentationv1.Experiment {
	now := time.Now()
	config := serverexperimentationv1.HTTPFaultConfig{
		Fault: &serverexperimentationv1.HTTPFaultConfig_AbortFault{
			AbortFault: &serverexperimentationv1.AbortFault{
				Percentage: &serverexperimentationv1.FaultPercentage{
					Percentage: 100,
				},
				AbortStatus: &serverexperimentationv1.FaultAbortStatus{
					HttpStatusCode: uint32(faultHttpStatus),
				},
			},
		},
		FaultTargeting: &serverexperimentationv1.FaultTargeting{
			Enforcer: &serverexperimentationv1.FaultTargeting_UpstreamEnforcing{
				UpstreamEnforcing: &serverexperimentationv1.UpstreamEnforcing{
					UpstreamType: &serverexperimentationv1.UpstreamEnforcing_UpstreamCluster{
						UpstreamCluster: &serverexperimentationv1.SingleCluster{
							Name: "test-cluster",
						},
					},
					DownstreamType: &serverexperimentationv1.UpstreamEnforcing_DownstreamCluster{
						DownstreamCluster: &serverexperimentationv1.SingleCluster{
							Name: "test-cluster",
						},
					},
				},
			},
		},
	}

	a, err := ptypes.MarshalAny(&config)
	assert.NoError(t, err)

	experiment, err := storer.CreateExperiment(context.Background(), a, &now, &now)
	assert.NoError(t, err)

	return experiment
}

type testCriteria struct {
	evaluation bool

	sync.Mutex
}

func (t *testCriteria) ShouldTerminate(experiment *experimentationv1.Experiment, experimentConfig *ptypes.DynamicAny) (string, error) {
	t.Lock()
	defer t.Unlock()

	if t.evaluation {
		return "terminated", nil
	}

	return "", nil
}

func (t *testCriteria) update(evaluation bool) {
	t.Lock()
	defer t.Unlock()

	t.evaluation = evaluation
}

func awaitGaugeValue(ctx context.Context, t *testing.T, testScope tally.TestScope, name string, value float64) {
	t.Helper()

	checkTicker := time.NewTicker(100 * time.Millisecond)
	for {
		select {
		case <-checkTicker.C:
			g, ok := testScope.Snapshot().Gauges()[name+"+"]
			if ok && g.Value() == value {
				return
			}
		case <-ctx.Done():
			t.Errorf(context.Canceled.Error())
			return
		}
	}
}
