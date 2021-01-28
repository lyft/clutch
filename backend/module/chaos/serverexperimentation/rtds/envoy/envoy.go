package envoy

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	bootstrap "github.com/envoyproxy/go-control-plane/envoy/config/bootstrap/v3"
	cluster "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	listener "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	hcm "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/http_connection_manager/v3"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"gopkg.in/yaml.v3"
)

const baseConfig = `
node:
  id: test
  cluster: test-cluster
admin:
  access_log_path: /dev/null
  address:
    socket_address: { address: 0.0.0.0, port_value: 9901 }
static_resources:
  listeners:
  - name: ingress
    address:
      socket_address: { address: 0.0.0.0, port_value: 10000 }
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

// EnvoyHandle is a handle to the Envoy instance under test, providing a startup check and provides
// utilities for interacting with the instance.
type EnvoyHandle struct{}

// NewEnvoyHandle creates a new handle for the Envoy under test after waiting for it to initialize.
func NewEnvoyHandle() (*EnvoyHandle, error) {
	timeout := time.NewTimer(5 * time.Second)

	for range time.NewTicker(100 * time.Millisecond).C {
		select {
		case <-timeout.C:
			return nil, errors.New("timed out waiting for Envoy to start up")
		default:
		}

		_, err := net.Dial("tcp", "envoy:10000")
		if err == nil {
			break
		}
	}

	return &EnvoyHandle{}, nil
}

// MakeSimpleCall issues a basic GET request to the Envoy under test, with the downstream cluster
// set to test-cluster.
func (e *EnvoyHandle) MakeSimpleCall() (int, error) {
	client := &http.Client{}

	r, err := http.NewRequest("GET", "http://envoy:10000", nil)
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

// EnvoyConfig provides a configuration builder that mirrors the upstream Envoy ConfigHelper:
// a base configuration is used which can be modified by a series of modifiers to create the
// final configuration.
type EnvoyConfig struct {
	configModifiers []func(*bootstrap.Bootstrap) *bootstrap.Bootstrap

	finalConfig *bootstrap.Bootstrap
}

// Generate generates the final configuration by applying all the config modifiers.
func (e *EnvoyConfig) Generate() (string, error) {
	if e.finalConfig == nil {
		b := &bootstrap.Bootstrap{}

		err := unmarshalYaml(b, baseConfig)
		if err != nil {
			return "", err
		}

		for _, m := range e.configModifiers {
			b = m(b)
		}

		e.finalConfig = b
	}

	m := jsonpb.Marshaler{}
	out := bytes.NewBuffer([]byte{})
	err := m.Marshal(out, e.finalConfig)

	return out.String(), err
}

// NewEnvoyConfig creates a new Envoy config builder, using a sensible default configuration that allows
// the provided EnvoyHandle to interact with the underlying Envoy instance.
func NewEnvoyConfig() *EnvoyConfig {
	return &EnvoyConfig{}
}

// AddRuntimeLayer adds a single runtime layer to the bootstrap.
func (e *EnvoyConfig) AddRuntimeLayer(input string) error {
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
func (e *EnvoyConfig) AddCluster(input string) error {
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
func (e *EnvoyConfig) AddHTTPFilter(input string) error {
	filter := &hcm.HttpFilter{}
	err := unmarshalYaml(filter, input)
	if err != nil {
		return err
	}

	e.configModifiers = append(e.configModifiers, func(b *bootstrap.Bootstrap) *bootstrap.Bootstrap {
		h := &hcm.HttpConnectionManager{}
		// TODO(snowp): Have config modifiers return (Bootstrap, error) so we can propagate errors.
		_ = b.StaticResources.Listeners[0].FilterChains[0].Filters[0].GetTypedConfig().UnmarshalTo(h)

		h.HttpFilters = append([]*hcm.HttpFilter{filter}, h.HttpFilters...)

		a, _ := ptypes.MarshalAny(h)

		b.StaticResources.Listeners[0].FilterChains[0].Filters[0].ConfigType = &listener.Filter_TypedConfig{
			TypedConfig: a,
		}

		return b
	})

	return nil
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
