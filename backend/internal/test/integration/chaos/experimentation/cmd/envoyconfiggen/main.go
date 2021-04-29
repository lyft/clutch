package main

import (
	"fmt"

	_ "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/fault/v3"

	testenvoy "github.com/lyft/clutch/backend/internal/test/integration/helper/envoytest"
)

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

	out, err := config.Generate()
	if err != nil {
		panic(err)
	}

	fmt.Print(out)
}
