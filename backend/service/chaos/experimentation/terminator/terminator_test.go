package terminator

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/uber-go/tally"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	experimentationv1 "github.com/lyft/clutch/backend/api/chaos/experimentation/v1"
	serverexperimentationv1 "github.com/lyft/clutch/backend/api/chaos/serverexperimentation/v1"
	terminatorv1 "github.com/lyft/clutch/backend/api/config/service/chaos/experimentation/terminator/v1"
	"github.com/lyft/clutch/backend/mock/service/chaos/experimentation/experimentstoremock"
	"github.com/lyft/clutch/backend/service"
	"github.com/lyft/clutch/backend/service/chaos/experimentation/experimentstore"
)

const testConfigType = "type.googleapis.com/clutch.chaos.serverexperimentation.v1.HTTPFaultConfig"

func TestConfigLoad(t *testing.T) {
	builtInCriteriaConfig := &terminatorv1.MaxTimeTerminationCriteria{
		MaxDuration: durationpb.New(time.Hour),
	}
	builtInCriteriaAnyConfig, err := anypb.New(builtInCriteriaConfig)
	assert.NoError(t, err)

	cfg := &terminatorv1.Config{
		PerConfigTypeConfiguration: map[string]*terminatorv1.Config_PerConfigTypeConfig{testConfigType: {TerminationCriteria: []*anypb.Any{builtInCriteriaAnyConfig}}},
		OuterLoopInterval:          durationpb.New(time.Second),
		PerExperimentCheckInterval: durationpb.New(time.Second),
	}

	any, err := anypb.New(cfg)
	assert.NoError(t, err)

	service.Registry[experimentstore.Name] = &experimentstoremock.MockStorer{}

	// Happy path testing.
	m, err := NewMonitor(any, zap.NewNop(), tally.NoopScope)
	assert.NoError(t, err)
	assert.Len(t, m.terminationCriteriaByTypeUrl, 1)
	assert.Len(t, m.terminationCriteriaByTypeUrl[testConfigType], 1)

	// If configured with a criteria we don't have registered we should error out.
	otherProto := &wrapperspb.StringValue{}
	otherProtoAny, err := anypb.New(otherProto)
	assert.NoError(t, err)

	cfg.PerConfigTypeConfiguration[testConfigType].TerminationCriteria = append(cfg.PerConfigTypeConfiguration[testConfigType].TerminationCriteria, otherProtoAny)
	any, err = anypb.New(cfg)
	assert.NoError(t, err)

	m, err = NewMonitor(any, zap.NewNop(), tally.NoopScope)
	assert.Nil(t, m)
	assert.EqualError(t, err, "terminator module configured with unknown criteria 'type.googleapis.com/google.protobuf.StringValue'")
}

func TestTerminator(t *testing.T) {
	l, err := zap.NewDevelopment()
	assert.NoError(t, err)

	testScope := tally.NewTestScope("", map[string]string{})

	store := &experimentstoremock.SimpleStorer{}
	criteria := &testCriteria{}
	monitor := Monitor{
		store:                        store,
		terminationCriteriaByTypeUrl: map[string][]TerminationCriteria{testConfigType: {criteria}},
		outerLoopInterval:            time.Millisecond,
		perExperimentCheckInterval:   time.Millisecond,
		log:                          l.Sugar(),
		activeMonitoringRoutines: trackingGauge{
			gauge: testScope.Gauge("active_routines"),
			value: 0,
		},
		criteriaEvaluationSuccessCount: testScope.Counter("criteria_success"),
		criteriaEvaluationFailureCount: testScope.Counter("criteria_failure"),
		terminationCount:               testScope.Counter("terminations"),
		marshallingErrorCount:          testScope.Counter("unpack_error"),
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	monitor.Run(ctx)

	createTestExperiment(t, 500, store)

	awaitGaugeValue(ctx, t, testScope, "active_routines", 1)

	es, err := store.GetExperiments(ctx, testConfigType, experimentationv1.GetExperimentsRequest_STATUS_RUNNING)
	assert.NoError(t, err)
	assert.Len(t, es, 1)

	criteria.update(true)

	// Ensure that we tear down the monitoring goroutines once we're done with an experiment.
	awaitGaugeValue(ctx, t, testScope, "active_routines", 0)
	assert.Equal(t, int64(1), testScope.Snapshot().Counters()["terminations+"].Value())

	es, err = store.GetExperiments(ctx, testConfigType, experimentationv1.GetExperimentsRequest_STATUS_RUNNING)
	assert.NoError(t, err)
	assert.Len(t, es, 0)
}

func createTestExperiment(t *testing.T, faultHttpStatus int, storer *experimentstoremock.SimpleStorer) *experimentationv1.Experiment {
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

	a, err := anypb.New(&config)
	assert.NoError(t, err)

	now := time.Now()
	endTime := time.Now().Add(5 * time.Minute)
	es := &experimentstore.ExperimentSpecification{StartTime: now, EndTime: &endTime, Config: a}
	experiment, err := storer.CreateExperiment(context.Background(), es)
	assert.NoError(t, err)

	return experiment
}

type testCriteria struct {
	evaluation bool

	sync.Mutex
}

func (t *testCriteria) ShouldTerminate(experiment *experimentationv1.Experiment, experimentConfig proto.Message) (string, error) {
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
