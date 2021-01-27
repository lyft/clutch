package envoy

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"testing"
	"time"

	bootstrap "github.com/envoyproxy/go-control-plane/envoy/config/bootstrap/v3"
	listener "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	hcm "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/http_connection_manager/v3"
	cluster "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

const baseConfig = `
node:
  id: test
  cluster: test-cluster
admin:
  access_log_path: /dev/null
  address:
    socket_address: { address: 127.0.0.1, port_value: 9901 }
static_resources:
  listeners:
  - name: ingress
    address:
      socket_address: { address: 127.0.0.1, port_value: 10000 }
    filter_chains:
    - filters:
      - name: envoy.filters.network.http_connection_manager
        typed_config:
          "@type": type.googleapis.com/envoy.extensions.filters.network.http_connection_manager.v3.HttpConnectionManager
          stat_prefix: ingress_http
          codec_type: AUTO
          route_config:
            name: local_route
            virtual_hosts:
            - name: local_service
              domains: ["*"]
              routes:
              - match: { prefix: "/" }
                route: { cluster: cluster_0 }
          http_filters:
          - name: envoy.filters.http.router
  clusters:
  - name: cluster_0
    connect_timeout: 0.25s
    type: STATIC
    lb_policy: ROUND_ROBIN
    load_assignment:
      cluster_name: cluster_0
      endpoints:
      - lb_endpoints:
        - endpoint:
            address:
              socket_address:
                address: 127.0.0.1
                port_value: 1234
`

// Envoy provides lifetime management of and utility functions for an Envoy instance under test.
// This makes use of a similar configuration system as Envoy integration tests, where a base
// config is modified by a series of modifiers before arriving at the final bootstrap config.
// TODO(snowp): If we ever need to support multiple Envoy instances under test, use some port selection
// logic to avoid port collisions.
type Envoy struct {
	doneCh chan error
	configModifiers []func(*bootstrap.Bootstrap) *bootstrap.Bootstrap

	finalConfig *bootstrap.Bootstrap
}

// Start generates the final configuration and starts the Envoy proccess.
func (e *Envoy) Start() error {
	if e.finalConfig == nil {
		b := &bootstrap.Bootstrap{}

		err := unmarshalYaml(b, baseConfig)
		if err != nil {
			return err
		}

		fmt.Println(len(e.configModifiers))
		for _, m := range e.configModifiers {
			b = m(b)
		}

		e.finalConfig = b
	}

	m := jsonpb.Marshaler{}
	out := bytes.NewBuffer([]byte{})
	err := m.Marshal(out, e.finalConfig)
	if err != nil {
		return err
	}

	startupCh := make(chan error)

	go func(startupCh chan error) {
		timeout := time.NewTimer(5 * time.Second)

		for range time.NewTicker(100 * time.Millisecond).C {
			select {
			case <- timeout.C:
					startupCh <- errors.New("timed out waiting for Envoy to start up")
					close(startupCh)
					return
			default:
			}

			_, err := net.Dial("tcp", "localhost:10000")
			if err == nil {
				close(startupCh)
				return
			}
		}
	}(startupCh)

	go func(doneCh chan error) {
		defer close(doneCh)
		c := exec.Command("/usr/local/bin/getenvoy", "run", "standard:1.17.0", "--", "--config-yaml", string(out.Bytes()))
		c.Stderr = os.Stderr
		c.Stdout = os.Stdout
		err := c.Run()
		doneCh <- err
	}(e.doneCh)

	err = <- startupCh

	return err
}

// Stop issues a stop signal to the underlying Envoy.
func (e *Envoy) Stop(t *testing.T) {
	_, err := http.Post("http://localhost:9901/quitquitquit", "text/plain", nil)
	assert.NoError(t, err)
}

// MakeSimpleCall issues a basic GET request to the Envoy under test, with the downstream cluster
// set to test-cluster.
func (e *Envoy) MakeSimpleCall() (int, error) {
	client := &http.Client{
	}

	r, err := http.NewRequest("GET", "http://localhost:10000", nil)
	if err != nil {
		return 0, err
	}

	r.Header.Add("x-envoy-downstream-service-cluster", "test-cluster")
	resp, err := client.Do(r)
	if err != nil {
		return 0, err
	}

	return resp.StatusCode, nil
}

// AddRuntimeLayer adds a single runtime layer to the bootstrap.
func (e *Envoy) AddRuntimeLayer(input string) error {
	runtimeLayer := &bootstrap.RuntimeLayer{}

	err := unmarshalYaml(runtimeLayer, input)
	if err != nil {
		return err
	}

	e.configModifiers = append(e.configModifiers, func(b *bootstrap.Bootstrap) *bootstrap.Bootstrap {
		if b.LayeredRuntime == nil {
			b.LayeredRuntime = &bootstrap.LayeredRuntime{}
		}
		b.LayeredRuntime.Layers = append(b.LayeredRuntime.Layers, runtimeLayer)
		
		return b
	})

	return nil
}

// AddCluster adds a cluster to the list of static clusters.
func (e *Envoy) AddCluster(input string) error {
	cluster := &cluster.Cluster{}

	err := unmarshalYaml(cluster, input)
	if err != nil {
		return err
	}

	e.configModifiers = append(e.configModifiers, func(b *bootstrap.Bootstrap) *bootstrap.Bootstrap {
		b.StaticResources.Clusters = append(b.StaticResources.Clusters, cluster)
		
		return b
	})

	return nil
}

// AddHTTPFilter adds a HTTP filter in front of the list of HTTP filters for the default listener.
func (e *Envoy) AddHTTPFilter(input string) error {
	filter := &hcm.HttpFilter{}
	err := unmarshalYaml(filter, input)
	if err != nil {
		return err
	}

	e.configModifiers = append(e.configModifiers, func(b *bootstrap.Bootstrap) *bootstrap.Bootstrap {
		h := &hcm.HttpConnectionManager{}
		b.StaticResources.Listeners[0].FilterChains[0].Filters[0].GetTypedConfig().UnmarshalTo(h)

		h.HttpFilters = append([]*hcm.HttpFilter{filter}, h.HttpFilters...)

		a, _ := ptypes.MarshalAny(h)

		b.StaticResources.Listeners[0].FilterChains[0].Filters[0].ConfigType = &listener.Filter_TypedConfig{
			TypedConfig: a,
		}

		return b
	})

	return nil
}

// AwaitShutdown blocks until the Envoy instance has terminated.
func (e *Envoy) AwaitShutdown() {
	<- e.doneCh
}

// NewEnvoy creates a new Envoy instance to run under test.
func NewEnvoy(t *testing.T) *Envoy {
	// We use of getenvoy here to make things really simple, but it might be nicer
	// to use something that supports more specific Envoy versions.
	if _, err := os.Stat("/usr/local/bin/getenvoy"); os.IsNotExist(err) {
		t.Fatal("getenvoy not installed at /usr/local/bin/getenvoy")
		return nil
	}

	envoy := &Envoy{
		doneCh: make(chan error),
	}

	return envoy
}

// Helper to convert an input yaml into a typed Protobuf message.
func unmarshalYaml(m proto.Message, input string) error {
		intermediate := map[string]interface{}{}

		err := yaml.Unmarshal([]byte(input), intermediate)
		if err != nil {
			return fmt.Errorf("error unmarshaling yaml: %s", err)
		}

		asJSON, err := json.Marshal(intermediate)
		if err != nil {
			return fmt.Errorf("error marshaling json: %s", err)
		}

		u := jsonpb.Unmarshaler{}
		err = u.Unmarshal(bytes.NewReader(asJSON), m)
		if err != nil {
			return fmt.Errorf("error unmarshaling json (%s) to proto: %s", string(asJSON), err)
		}

		return nil
}
