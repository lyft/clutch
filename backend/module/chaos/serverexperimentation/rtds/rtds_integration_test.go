package rtds

import (
	"context"
	"testing"
	"time"

	_ "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/fault/v3"
	"github.com/golang/protobuf/ptypes"
	"github.com/stretchr/testify/assert"

	serverexperimentation "github.com/lyft/clutch/backend/api/chaos/serverexperimentation/v1"
	"github.com/lyft/clutch/backend/module/chaos/serverexperimentation/rtds/envoy"
)

const FaultFilterConfig = `
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

const RtdsLayerConfig = `
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

const RtdsCluster = `
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

func TestEnvoyFaults(t *testing.T) {
	ts := newTestServer(t, false)
	defer ts.stop()

	e := envoy.NewEnvoy(t)
	err := e.AddHTTPFilter(FaultFilterConfig)
	assert.NoError(t, err)

	err = e.AddCluster(RtdsCluster)
	assert.NoError(t, err)

	err = e.AddRuntimeLayer(RtdsLayerConfig)
	assert.NoError(t, err)

	err = e.Start()
	assert.NoError(t, err)
	defer e.AwaitShutdown()
	defer e.Stop(t)

	code, err := e.MakeSimpleCall()
	assert.NoError(t, err)
	assert.Equal(t, 503, code)

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

	_, err = ts.storer.CreateExperiment(context.Background(), a, &now, &now)
	assert.NoError(t, err)

	timeout := time.NewTimer(2 * time.Second)

	for range time.NewTicker(100 * time.Millisecond).C {
		select {
		case <-timeout.C:
			t.Errorf("timed out waiting for faults to take effect")
			break
		default:
		}

		code, err = e.MakeSimpleCall()
		if err == nil && code == 400 {
			break
		}
	}
}
