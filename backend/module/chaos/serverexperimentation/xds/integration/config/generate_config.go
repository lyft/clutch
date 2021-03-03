package main

import (
	"fmt"

	_ "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/fault/v3"

	testenvoy "github.com/lyft/clutch/backend/test/envoy"
)

const faultFilterConfig = `
name: envoy.fault
typed_config:
  '@type': type.googleapis.com/envoy.extensions.filters.http.fault.v3.HTTPFault
  abort:
    http_status: 503
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

const ecdsFilterConfig = `
name: envoy.extension_config
config_discovery:
  config_source:
    api_config_source: 
      api_type: GRPC
      grpc_services: 
        - envoy_grpc: 
            cluster_name: clutchxds
      transport_api_version: V3
    initial_fetch_timeout: 5s
    resource_api_version: V3
  default_config: 
    '@type': type.googleapis.com/envoy.extensions.filters.http.fault.v3.HTTPFault
    abort: 
      http_status: 503
      percentage: {denominator: HUNDRED, numerator: 0}
    abort_http_status_runtime: ecds_runtime_override_do_not_use.http.abort.http_status
    abort_percent_runtime: ecds_runtime_override_do_not_use.http.abort.abort_percent
    delay: 
      fixed_delay: 0.001s
      percentage: {denominator: HUNDRED, numerator: 0}
    delay_duration_runtime: ecds_runtime_override_do_not_use.http.delay.fixed_duration_ms
    delay_percent_runtime: ecds_runtime_override_do_not_use.http.delay.percentage
  apply_default_config_without_warming: false
  type_urls: 
    - type.googleapis.com/envoy.extensions.filters.http.fault.v3.HTTPFault
`

const xdsCluster = `
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
            address: "test_runner"
            port_value: 9000
name: clutchxds
type: STRICT_DNS
`

func main() {
	config := testenvoy.NewEnvoyConfig()

	err := config.AddCluster(xdsCluster)
	if err != nil {
		panic(err)
	}

	err = config.AddHTTPFilter(ecdsFilterConfig)
	if err != nil {
		panic(err)
	}

	err = config.AddHTTPFilter(faultFilterConfig)
	if err != nil {
		panic(err)
	}

	err = config.AddRuntimeLayer(rtdsLayerConfig)
	if err != nil {
		panic(err)
	}

	out, err := config.Generate()
	if err != nil {
		panic(err)
	}

	fmt.Print(out)
}
