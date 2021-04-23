// +build integration_only

package serverexperimentation

import (
	"context"
	"fmt"
	"github.com/lyft/clutch/backend/module/chaos/experimentation/xds"
	"github.com/lyft/clutch/backend/service"
	"github.com/lyft/clutch/backend/service/chaos/experimentation/experimentstore"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	experimentationv1 "github.com/lyft/clutch/backend/api/chaos/experimentation/v1"
	serverexperimentationv1 "github.com/lyft/clutch/backend/api/chaos/serverexperimentation/v1"
	xdsconfigv1 "github.com/lyft/clutch/backend/api/config/module/chaos/experimentation/xds/v1"
	serverexperimentationconfigv1 "github.com/lyft/clutch/backend/api/config/module/chaos/serverexperimentation/v1"
	"github.com/lyft/clutch/backend/internal/test/integration/helper/envoytest"
	xdstest "github.com/lyft/clutch/backend/internal/test/integration/helper/testserver"
	"github.com/lyft/clutch/backend/mock/service/chaos/experimentation/experimentstoremock"
)

// These tests are intended to be run with docker-compose to in order to set up a running Envoy instance
// to run assertions against.
func TestEnvoyRTDSFaults(t *testing.T) {
	xdsConfig := &xdsconfigv1.Config{
		RtdsLayerName: "rtds",
	}
	seConfig := &serverexperimentationconfigv1.Config{
		IngressFaultRuntimePrefix: "fault.http",
	}

	ts, err := newServerExperimentationTestModuleServer(seConfig, xdsConfig)
	assert.NoError(t, err)
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

func TestEnvoyECDSFaults(t *testing.T) {
	xdsConfig := &xdsconfigv1.Config{
		RtdsLayerName:        "rtds",
		CacheRefreshInterval: durationpb.New(time.Second),
		EcdsAllowList:        &xdsconfigv1.Config_ECDSAllowList{EnabledClusters: []string{"test-cluster"}},
	}
	seConfig := &serverexperimentationconfigv1.Config{}

	ts, err := newServerExperimentationTestModuleServer(seConfig, xdsConfig)
	assert.NoError(t, err)
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
	ts.Storer.CancelExperimentRun(context.Background(), experiment.Run.Id, "")

	err = awaitExpectedReturnValueForSimpleCall(t, e, awaitReturnValueParams{
		timeout:        10 * time.Second,
		expectedStatus: 503,
	})
	assert.NoError(t, err, "did not see faults reverted")
}

// Helper class for initializing a test RTDS server.
type TestModuleServer struct {
	*xdstest.TestServer

	Storer *experimentstoremock.SimpleStorer
}

func newServerExperimentationTestModuleServer(seConfig *serverexperimentationconfigv1.Config, xdsConfig *xdsconfigv1.Config) (*TestModuleServer, error) {
	xdsConfig.ResourceTtl = &durationpb.Duration{
		Seconds: 2,
	}
	xdsConfig.HeartbeatInterval = &durationpb.Duration{
		Seconds: 1,
	}
	xdsConfig.CacheRefreshInterval = &durationpb.Duration{
		Seconds: 1,
	}

	anySeConfig, err := anypb.New(seConfig)
	if err != nil {
		return nil, err
	}

	anyXdsConfig, err := anypb.New(xdsConfig)
	if err != nil {
		return nil, err
	}

	ms := &TestModuleServer{
		TestServer: xdstest.New(),
		Storer:     &experimentstoremock.SimpleStorer{},
	}
	service.Registry[experimentstore.Name] = ms.Storer

	_, err = New(anySeConfig, ms.Logger, ms.Scope)
	if err != nil {
		return nil, err
	}

	xdsm, err := xds.New(anyXdsConfig, ms.Logger, ms.Scope)
	if err != nil {
		return nil, err
	}

	ms.Register(xdsm)

	err = ms.Run()
	if err != nil {
		return nil, err
	}

	return ms, nil
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

func createTestExperiment(t *testing.T, faultHttpStatus int, storer *experimentstoremock.SimpleStorer) *experimentstore.Experiment {
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

	now := time.Date(2011, 0, 0, 0, 0, 0, 0, time.UTC)
	future := now.Add(1 * time.Hour)
	farFuture := future.Add(1 * time.Hour)

	d := &experimentationv1.CreateExperimentData{StartTime: timestamppb.New(future), EndTime: timestamppb.New(farFuture), Config: a}
	s, err := experimentstore.NewExperimentSpecification(d, now)
	assert.NoError(t, err)
	experiment, err := storer.CreateExperiment(context.Background(), s)
	assert.NoError(t, err)

	return experiment
}