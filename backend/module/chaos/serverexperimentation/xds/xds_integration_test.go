// +build integration_only

package xds

import (
	"context"
	"fmt"
	"testing"
	"time"

	_ "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/fault/v3"
	"github.com/golang/protobuf/ptypes"
	"github.com/stretchr/testify/assert"

	serverexperimentation "github.com/lyft/clutch/backend/api/chaos/serverexperimentation/v1"
	xds_testing "github.com/lyft/clutch/backend/module/chaos/serverexperimentation/xds/testing"
	"github.com/lyft/clutch/backend/test/envoy"
)

// These tests are intended to be run with docker-compose to in order to set up a running Envoy instance
// to run assertions against.
func TestEnvoyFaults(t *testing.T) {
	ts := xds_testing.NewTestServer(t, New, true)
	defer ts.Stop()

	e, err := testenvoy.NewEnvoyHandle()
	assert.NoError(t, err)

	code, err := e.MakeSimpleCall()
	assert.NoError(t, err)
	assert.Equal(t, 503, code)

	createTestExperiment(t, ts.Storer)

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

func createTestExperiment(t *testing.T, storer *xds_testing.SimpleStorer) {
	now := time.Now()
	config := serverexperimentation.HTTPFaultConfig{
		Fault: &serverexperimentation.HTTPFaultConfig_AbortFault{
			AbortFault: &serverexperimentation.AbortFault{
				Percentage: &serverexperimentation.FaultPercentage{
					Percentage: 100,
				},
				AbortStatus: &serverexperimentation.FaultAbortStatus{
					HttpStatusCode: 400,
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

	_, err = storer.CreateExperiment(context.Background(), a, &now, &now)
	assert.NoError(t, err)
}

type awaitReturnValueParams struct {
	timeout        time.Duration
	expectedStatus int
}

func awaitExpectedReturnValueForSimpleCall(t *testing.T, e *testenvoy.EnvoyHandle, params awaitReturnValueParams) error {
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
