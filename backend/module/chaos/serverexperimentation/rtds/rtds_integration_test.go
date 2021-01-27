package rtds

import (
	"context"
	"fmt"
	"testing"
	"time"

	_ "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/fault/v3"
	"github.com/golang/protobuf/ptypes"
	"github.com/stretchr/testify/assert"

	serverexperimentation "github.com/lyft/clutch/backend/api/chaos/serverexperimentation/v1"
	"github.com/lyft/clutch/backend/module/chaos/serverexperimentation/rtds/envoy"
	rtds_testing "github.com/lyft/clutch/backend/module/chaos/serverexperimentation/rtds/testing"
)

func TestEnvoyFaults(t *testing.T) {
	ts := newTestServer(t, false)
	defer ts.stop()

	e := startEnvoy(t)
	defer e.stop()

	code, err := e.MakeSimpleCall()
	assert.NoError(t, err)
	assert.Equal(t, 503, code)

	createTestExperiment(t, ts.storer)

	err = awaitExpectedReturnValueForSimpleCall(t, e.Envoy, awaitReturnValueParams{
		timeout:        2 * time.Second,
		expectedStatus: 400,
	})
	assert.NoError(t, err, "did not see faults enabled")
}

func TestTTls(t *testing.T) {
	ts := newTestServer(t, true)

	e := startEnvoy(t)
	defer e.stop()

	code, err := e.MakeSimpleCall()
	assert.NoError(t, err)
	assert.Equal(t, 503, code)

	createTestExperiment(t, ts.storer)

	err = awaitExpectedReturnValueForSimpleCall(t, e.Envoy, awaitReturnValueParams{
		timeout:        2 * time.Second,
		expectedStatus: 400,
	})
	assert.NoError(t, err, "did not see faults enabled")

	// We know that faults have been applied, now try to kill the server and ensure that we eventually reset the faults.
	ts.stop()

	err = awaitExpectedReturnValueForSimpleCall(t, e.Envoy, awaitReturnValueParams{
		timeout:        10 * time.Second,
		expectedStatus: 503,
	})
	assert.NoError(t, err, "did not see faults reverted")
}

const faultFilterConfig = `
name: envoy.fault
typed_config:
  '@type': type.googleapis.com/envoy.extensions.filters.http.fault.v3.HTTPFault
  abort:
    http_status: 400
    percentage: {denominator: HUNDRED, numerator: 0}
  delay:
    fixed_delay: 0.001s
    percentage: {denominator: HUNDRED, numerator: 0}
`

const rtdsLayerConfig = `
name: rtds
rtds_layer:
  name: rtds
  rtds_config:
    resource_api_version: V3
    api_config_source:
      api_type: GRPC
      transport_api_version: V3
      grpc_services:
      - envoy_grpc: 
          cluster_name: clutchxds
`

const rtdsCluster = `
connect_timeout: 0.250s
http2_protocol_options: {}
lb_policy: ROUND_ROBIN
load_assignment:
  cluster_name: clutchxds
  endpoints:
  - lb_endpoints:
    - endpoint:
        address:
          socket_address: 
            address: "127.0.0.1"
            port_value: 9000
name: clutchxds
type: STRICT_DNS
`

type rtdsEnvoy struct {
	stop func()

	envoy.Envoy
}

func startEnvoy(t *testing.T) *rtdsEnvoy {
	e := envoy.NewEnvoy(t)
	err := e.AddHTTPFilter(faultFilterConfig)
	assert.NoError(t, err)

	err = e.AddCluster(rtdsCluster)
	assert.NoError(t, err)

	err = e.AddRuntimeLayer(rtdsLayerConfig)
	assert.NoError(t, err)

	err = e.Start()
	assert.NoError(t, err)

	return &rtdsEnvoy{
		stop: func() {
			e.Stop(t)
			e.AwaitShutdown()
		},
		Envoy: *e,
	}
}

func createTestExperiment(t *testing.T, storer *rtds_testing.SimpleStorer) {
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

func awaitExpectedReturnValueForSimpleCall(t *testing.T, e envoy.Envoy, params awaitReturnValueParams) error {
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
