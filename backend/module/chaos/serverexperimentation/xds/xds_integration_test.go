// +build integration_only

package xds

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	_ "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/fault/v3"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	wrapperspb "github.com/golang/protobuf/ptypes/wrappers"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	experimentationv1 "github.com/lyft/clutch/backend/api/chaos/experimentation/v1"
	serverexperimentation "github.com/lyft/clutch/backend/api/chaos/serverexperimentation/v1"
	xdsconfigv1 "github.com/lyft/clutch/backend/api/config/module/chaos/experimentation/xds/v1"
	"github.com/lyft/clutch/backend/api/config/service/chaos/experimentation/terminator/v1"
	"github.com/lyft/clutch/backend/module/chaos/serverexperimentation/xds/internal/xdstest"
	"github.com/lyft/clutch/backend/service/chaos/experimentation/experimentstore"
	"github.com/lyft/clutch/backend/internal/test/integration/helper/envoytest"
	"github.com/lyft/clutch/backend/mock/service/chaos/experimentation/experimentstoremock"
	"github.com/lyft/clutch/backend/service/chaos/experimentation/terminator"
)

// These tests are intended to be run with docker-compose to in order to set up a running Envoy instance
// to run assertions against.
func TestEnvoyFaults(t *testing.T) {
	xdsConfig := &xdsconfigv1.Config{
		RtdsLayerName:             "rtds",
		CacheRefreshInterval:      ptypes.DurationProto(time.Second),
		IngressFaultRuntimePrefix: "fault.http",
		EgressFaultRuntimePrefix:  "egress",
	}

	ts := xdstest.NewTestModuleServer(New, true, xdsConfig)
	defer ts.Stop()

	e, err := envoytest.NewEnvoyHandle()
	assert.NoError(t, err)

	e.EnsureControlPlaneConnectivity(envoytest.RuntimeStatPrefix)
	assert.NoError(t, err)

	code, err := e.MakeSimpleCall()
	assert.NoError(t, err)
	assert.Equal(t, 503, code)

	createTestExperiment(t, 400, ts.Storer)

	err = awaitExpectedReturnValueForSimpleCall(t, e, awaitReturnValueParams{
		timeout:        4 * time.Second,
		expectedStatus: 400,
	})
	assert.NoError(t, err, "did not see faults enabled")

	// We know that faults have been applied, now try to kill the server and ensure that we eventually reset the faults.
	// This verifies that we're properly setting TTLs and that Envoy will honor this.
	ts.Stop()

	err = awaitExpectedReturnValueForSimpleCall(t, e, awaitReturnValueParams{
		timeout:        10 * time.Second,
		expectedStatus: 503,
	})
	assert.NoError(t, err, "did not see faults reverted")
}

func TestEnvoyFaultsTimeBasedTermination(t *testing.T) {
	t.Skip("flaky")

	xdsConfig := &xdsconfigv1.Config{
		RtdsLayerName:             "rtds",
		CacheRefreshInterval:      ptypes.DurationProto(time.Second),
		IngressFaultRuntimePrefix: "fault.http",
		EgressFaultRuntimePrefix:  "egress",
	}

	ts := xdstest.NewTestModuleServer(New, true, xdsConfig)
	defer ts.Stop()

	criteria := &testCriteria{}

	terminator.CriteriaFactories["type.googleapis.com/google.protobuf.StringValue"] = &testCriteriaFactory{testCriteria: criteria}

	anyString, err := ptypes.MarshalAny(&wrapperspb.StringValue{})
	assert.NoError(t, err)

	typedConfig := &terminatorv1.Config{
		PerConfigTypeConfiguration: map[string]*terminatorv1.Config_PerConfigTypeConfig{
			"type.googleapis.com/clutch.chaos.serverexperimentation.v1.HTTPFaultConfig": {
				TerminationCriteria: []*anypb.Any{anyString},
			},
		},
		OuterLoopInterval:          ptypes.DurationProto(time.Second),
		PerExperimentCheckInterval: ptypes.DurationProto(time.Second),
	}
	anyConfig, err := ptypes.MarshalAny(typedConfig)
	assert.NoError(t, err)

	terminator, err := terminator.NewMonitor(anyConfig, ts.Logger, ts.Scope)
	assert.NoError(t, err)

	// Cancel to ensure that the terminator doesn't leak into other tests.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	terminator.Run(ctx)

	e, err := envoytest.NewEnvoyHandle()
	assert.NoError(t, err)

	err = e.EnsureControlPlaneConnectivity(envoytest.RuntimeStatPrefix)
	assert.NoError(t, err)

	code, err := e.MakeSimpleCall()
	assert.NoError(t, err)
	assert.Equal(t, 503, code)

	createTestExperiment(t, 400, ts.Storer)

	err = awaitExpectedReturnValueForSimpleCall(t, e, awaitReturnValueParams{
		timeout:        4 * time.Second,
		expectedStatus: 400,
	})
	assert.NoError(t, err, "did not see faults enabled")

	criteria.start()

	// Since we've enabled a time based automatic termination, we expect to see faults get disabled on their own after some time.
	err = awaitExpectedReturnValueForSimpleCall(t, e, awaitReturnValueParams{
		timeout:        10 * time.Second,
		expectedStatus: 503,
	})
	assert.NoError(t, err, "did not see faults reverted")
}

func TestEnvoyECDSFaults(t *testing.T) {
	xdsConfig := &xdsconfigv1.Config{
		RtdsLayerName:             "rtds",
		CacheRefreshInterval:      ptypes.DurationProto(time.Second),
		IngressFaultRuntimePrefix: "fault.http",
		EgressFaultRuntimePrefix:  "egress",
		EcdsAllowList:             &xdsconfigv1.Config_ECDSAllowList{EnabledClusters: []string{"test-cluster"}},
	}

	ts := xdstest.NewTestModuleServer(New, true, xdsConfig)
	defer ts.Stop()

	e, err := envoytest.NewEnvoyHandle()
	assert.NoError(t, err)

	err = e.EnsureControlPlaneConnectivity(envoytest.EcdsStatPrefix)
	assert.NoError(t, err)

	code, err := e.MakeSimpleCall()
	assert.NoError(t, err)
	assert.Equal(t, 503, code)

	experiment := createTestExperiment(t, 404, ts.Storer)

	err = awaitExpectedReturnValueForSimpleCall(t, e, awaitReturnValueParams{
		timeout:        4 * time.Second,
		expectedStatus: 404,
	})
	assert.NoError(t, err, "did not see faults enabled")

	// TODO(kathan24): Test TTL by stopping the server instead of canceling the experiment. Currently, TTL is not not supported for ECDS in the upstream Envoy
	ts.Storer.CancelExperimentRun(context.Background(), experiment.RunId, "")

	err = awaitExpectedReturnValueForSimpleCall(t, e, awaitReturnValueParams{
		timeout:        10 * time.Second,
		expectedStatus: 503,
	})
	assert.NoError(t, err, "did not see faults reverted")
}

func createTestExperiment(t *testing.T, faultHttpStatus int, storer *experimentstoremock.SimpleStorer) *experimentationv1.Experiment {
	config := serverexperimentation.HTTPFaultConfig{
		Fault: &serverexperimentation.HTTPFaultConfig_AbortFault{
			AbortFault: &serverexperimentation.AbortFault{
				Percentage: &serverexperimentation.FaultPercentage{
					Percentage: 100,
				},
				AbortStatus: &serverexperimentation.FaultAbortStatus{
					HttpStatusCode: uint32(faultHttpStatus),
				},
			},
		},
		FaultTargeting: &serverexperimentation.FaultTargeting{
			Enforcer: &serverexperimentation.FaultTargeting_UpstreamEnforcing{
				UpstreamEnforcing: &serverexperimentation.UpstreamEnforcing{
					UpstreamType: &serverexperimentation.UpstreamEnforcing_UpstreamCluster{
						UpstreamCluster: &serverexperimentation.SingleCluster{
							Name: "test-cluster",
						},
					},
					DownstreamType: &serverexperimentation.UpstreamEnforcing_DownstreamCluster{
						DownstreamCluster: &serverexperimentation.SingleCluster{
							Name: "test-cluster",
						},
					},
				},
			},
		},
	}

	a, err := anypb.New(&config)
	assert.NoError(t, err)

	now := time.Date(2011, 0, 0, 0, 0, 0, 0, time.UTC)
	assert.NoError(t, err)
	future := now.Add(1 * time.Hour)
	futureTimestamp, err := ptypes.TimestampProto(future)
	assert.NoError(t, err)
	farFuture := future.Add(1 * time.Hour)
	farFutureTimestamp, err := ptypes.TimestampProto(farFuture)
	assert.NoError(t, err)

	d := &experimentationv1.CreateExperimentData{StartTime: futureTimestamp, EndTime: farFutureTimestamp, Config: a}
	s, err := experimentstore.NewExperimentSpecification(d, now)
	assert.NoError(t, err)
	experiment, err := storer.CreateExperiment(context.Background(), s)
	assert.NoError(t, err)

	return experiment
}

type awaitReturnValueParams struct {
	timeout        time.Duration
	expectedStatus int
}

func awaitExpectedReturnValueForSimpleCall(t *testing.T, e *envoytest.EnvoyHandle, params awaitReturnValueParams) error {
	timeout := time.NewTimer(params.timeout)

	for range time.NewTicker(100 * time.Millisecond).C {
		select {
		case <-timeout.C:
			return fmt.Errorf("did not receive expected status code before timeout")
		default:
		}

		code, err := e.MakeSimpleCall()
		if err == nil && code == params.expectedStatus {
			return nil
		}
	}

	return nil
}

type testCriteriaFactory struct {
	testCriteria *testCriteria
}

func (t *testCriteriaFactory) Create(*any.Any) (terminator.TerminationCriteria, error) {
	return t.testCriteria, nil
}

type testCriteria struct {
	startCheckingTime bool
	sync.Mutex
}

// We want the time check to be low but also avoid races, so this lets us prevent the criteria from activating until
// we know that we're seeing faults enabled.
func (t *testCriteria) start() {
	t.Lock()
	defer t.Unlock()

	t.startCheckingTime = true
}

func (t *testCriteria) ShouldTerminate(experiment *experimentationv1.Experiment, experimentConfig proto.Message) (string, error) {
	t.Lock()
	defer t.Unlock()

	// Sanity check that we get the expected configuration type.
	if _, ok := experimentConfig.(*serverexperimentation.HTTPFaultConfig); !ok {
		panic("received unexpected experiment config")
	}
	started, _ := ptypes.Timestamp(experiment.StartTime)
	if t.startCheckingTime && started.Add(1*time.Second).Before(time.Now()) {
		return "timed out", nil
	}

	return "", nil
}
