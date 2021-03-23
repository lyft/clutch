package experiment_terminator

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/stretchr/testify/assert"

	experimentationv1 "github.com/lyft/clutch/backend/api/chaos/experimentation/v1"
	serverexperimentationv1 "github.com/lyft/clutch/backend/api/chaos/serverexperimentation/v1"
	"github.com/lyft/clutch/backend/mock/service/chaos/experimentation/experimentstoremock"
)

type testCriteria struct {
	evaluation bool

	sync.Mutex
}

func (t *testCriteria) ShouldTerminate(experimentStarted time.Time, config interface{}) error {
	t.Lock()
	defer t.Unlock()

	if t.evaluation {
		return errors.New("terminated")
	}

	return nil
}

func (t *testCriteria) update(evaluation bool) {
	t.Lock()
	defer t.Unlock()

	t.evaluation = evaluation
}

func TestTerminator(t *testing.T) {
	store := &experimentstoremock.SimpleStorer{}
	criteria := &testCriteria{}
	terminator := Terminator{
		Store:                      store,
		EnabledConfigTypes:         []string{"type.googleapis.com/clutch.chaos.serverexperimentation.v1.HTTPFaultConfig"},
		Criterias:                  []TerminationCriteria{criteria},
		OuterLoopInterval:          1,
		PerExperimentCheckInterval: 1,
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	terminator.Run(ctx)

	createTestExperiment(t, 500, store)

	time.Sleep(5 * time.Second)

	es, err := store.GetExperiments(ctx, []string{"type.googleapis.com/clutch.chaos.serverexperimentation.v1.HTTPFaultConfig"}, experimentationv1.GetExperimentsRequest_STATUS_RUNNING)
	assert.NoError(t, err)
	assert.Len(t, es, 1)

	criteria.update(true)

	time.Sleep(2 * time.Second)

	es, err = store.GetExperiments(ctx, []string{"type.googleapis.com/clutch.chaos.serverexperimentation.v1.HTTPFaultConfig"}, experimentationv1.GetExperimentsRequest_STATUS_RUNNING)
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
