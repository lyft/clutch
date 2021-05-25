package envoyadmin

// <!-- START clutchdoc -->
// description: Executes remote queries against the Envoy Proxy admin interface.
// <!-- END clutchdoc -->

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/golang/protobuf/ptypes/any"
	"github.com/uber-go/tally"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	envoyadminv1 "github.com/lyft/clutch/backend/api/config/service/envoyadmin/v1"
	envoytriagev1 "github.com/lyft/clutch/backend/api/envoytriage/v1"
	"github.com/lyft/clutch/backend/service"
)

const Name = "clutch.service.envoyadmin"

func New(cfg *any.Any, logger *zap.Logger, scope tally.Scope) (service.Service, error) {
	return NewWithHTTPClient(cfg, logger, scope, &http.Client{})
}

func NewWithHTTPClient(cfg *any.Any, logger *zap.Logger, scope tally.Scope, httpClient *http.Client) (service.Service, error) {
	config := &envoyadminv1.Config{}
	if err := cfg.UnmarshalTo(config); err != nil {
		return nil, err
	}

	return &client{
		defaultPort: config.DefaultRemotePort,
		httpClient:  httpClient,
	}, nil
}

type Client interface {
	// Get performs read-only operations concurrently and returns the results. If any of the operations fail,
	// an error is returned.
	Get(ctx context.Context, operation *envoytriagev1.ReadOperation) (*envoytriagev1.Result, error)
}

type client struct {
	defaultPort uint32
	httpClient  *http.Client
}

func makeRequest(ctx context.Context, cl *http.Client, baseURL, path string) ([]byte, error) {
	url := fmt.Sprintf("%s%s", baseURL, path)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := cl.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("received non-200 status '%d %s'", resp.StatusCode, resp.Status)
	}

	return ioutil.ReadAll(resp.Body)
}

func (c *client) Get(ctx context.Context, operation *envoytriagev1.ReadOperation) (*envoytriagev1.Result, error) {
	defer c.httpClient.CloseIdleConnections()

	// Use default port if one was not provided and construct address.
	port := operation.Address.Port
	if port == 0 {
		port = c.defaultPort
	}
	addr := fmt.Sprintf("%s:%d", operation.Address.Host, port)
	baseURL := fmt.Sprintf("http://%s", addr)

	// Make an empty result.
	result := &envoytriagev1.Result{
		Address: &envoytriagev1.Address{Host: operation.Address.Host, Port: port},
		Output:  &envoytriagev1.Result_Output{},
	}

	// Make concurrent requests using an error group.
	g, ctx := errgroup.WithContext(ctx)
	if operation.Include.Clusters {
		g.Go(func() error {
			resp, err := makeRequest(ctx, c.httpClient, baseURL, "/clusters?format=json")
			if err != nil {
				return err
			}

			v, err := clustersFromResponse(resp)
			result.Output.Clusters = v
			return err
		})
	}

	if operation.Include.ConfigDump {
		g.Go(func() error {
			resp, err := makeRequest(ctx, c.httpClient, baseURL, "/config_dump")
			if err != nil {
				return err
			}

			v, err := configDumpFromResponse(resp)
			result.Output.ConfigDump = v
			return err
		})
	}

	if operation.Include.Listeners {
		g.Go(func() error {
			resp, err := makeRequest(ctx, c.httpClient, baseURL, "/listeners?format=json")
			if err != nil {
				return err
			}

			v, err := listenersFromResponse(resp)
			result.Output.Listeners = v
			return err
		})
	}

	if operation.Include.Runtime {
		g.Go(func() error {
			resp, err := makeRequest(ctx, c.httpClient, baseURL, "/runtime")
			if err != nil {
				return err
			}

			v, err := runtimeFromResponse(resp)
			result.Output.Runtime = v
			return err
		})
	}

	if operation.Include.Stats {
		g.Go(func() error {
			resp, err := makeRequest(ctx, c.httpClient, baseURL, "/stats")
			if err != nil {
				return err
			}

			v, err := statsFromResponse(resp)
			result.Output.Stats = v
			return err
		})
	}

	// Always fetch server info so we can populate node metadata.
	g.Go(func() error {
		resp, err := makeRequest(ctx, c.httpClient, baseURL, "/server_info")
		if err != nil {
			return err
		}

		// Include server info.
		if operation.Include.ServerInfo {
			v, err := serverInfoFromResponse(resp)
			if err != nil {
				return err
			}
			result.Output.ServerInfo = v
		}

		// Include node metadata.
		v, err := nodeMetadataFromResponse(resp)
		result.NodeMetadata = v
		return err
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return result, nil
}
