// +build integration_only

package xds

import (
	"context"
	"fmt"
	"github.com/lyft/clutch/backend/mock/service/chaos/experimentation/experimentstoremock"
	"github.com/lyft/clutch/backend/module/chaos/serverexperimentation/xds/internal/xdstest"
	"testing"
	"time"

	_ "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/fault/v3"
	"github.com/golang/protobuf/ptypes"
	"github.com/stretchr/testify/assert"

	experimentationv1 "github.com/lyft/clutch/backend/api/chaos/experimentation/v1"
	serverexperimentation "github.com/lyft/clutch/backend/api/chaos/serverexperimentation/v1"
	"github.com/lyft/clutch/backend/internal/test/integration/helper/envoytest"
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

	code, err := e.MakeSimpleCall()
	assert.NoError(t, err)
	assert.Equal(t, 503, code)

	_ = createTestExperiment(t, 400, ts.Storer)

	err = awaitExpectedReturnValueForSimpleCall(t, e, awaitReturnValueParams{
		timeout:        2 * time.Second,
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

func createTestExperiment(t *testing.T, faultHttpStatus int, storer *experimentstoremock.SimpleStorer) {
	now := time.Now()
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

	a, err := ptypes.MarshalAny(&config)
	assert.NoError(t, err)

	experiment, err := storer.CreateExperiment(context.Background(), a, &now, &now)
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

func TestEnvoyECDSFaults(t *testing.T) {
	xdsConfig := &xdsconfigv1.Config{
		RtdsLayerName:             "rtds",
		CacheRefreshInterval:      ptypes.DurationProto(time.Second),
		IngressFaultRuntimePrefix: "fault.http",
		EgressFaultRuntimePrefix:  "egress",
		EcdsAllowList:             &xdsconfigv1.Config_ECDSAllowList{EnabledClusters: []string{"test-cluster"}},
	}

	ts := xds_testing.NewTestServer(t, xdsConfig, xds.New, true)
	defer ts.Stop()

	e, err := testenvoy.NewEnvoyHandle()
	assert.NoError(t, err)

	code, err := e.MakeSimpleCall()
	assert.NoError(t, err)
	assert.Equal(t, 503, code)

	experiment := createTestExperiment(t, 404, ts.Storer)

	err = awaitExpectedReturnValueForSimpleCall(t, e, awaitReturnValueParams{
		timeout:        2 * time.Second,
		expectedStatus: 404,
	})
	assert.NoError(t, err, "did not see faults enabled")

	// Not testing TTL by stopping the server as its currently not supported for ECDS in Upstream Envoy
	ts.Storer.CancelExperimentRun(context.Background(), experiment.Id)

	err = awaitExpectedReturnValueForSimpleCall(t, e, awaitReturnValueParams{
		timeout:        10 * time.Second,
		expectedStatus: 503,
	})
	assert.NoError(t, err, "did not see faults reverted")
}
