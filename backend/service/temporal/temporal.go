package temporal

import (
	"fmt"

	"github.com/uber-go/tally"
	"go.temporal.io/sdk/client"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/anypb"

	temporalv1 "github.com/lyft/clutch/backend/api/config/service/temporal/v1"
	"github.com/lyft/clutch/backend/service"
)

// <!-- START clutchdoc -->
// description: Workflow client for temporal.io.
// <!-- END clutchdoc -->

const Name = "clutch.service.temporal"

func New(cfg *anypb.Any, logger *zap.Logger, scope tally.Scope) (service.Service, error) {
	config := &temporalv1.Config{}
	if err := cfg.UnmarshalTo(config); err != nil {
		return nil, err
	}

	tc, err := client.NewClient(client.Options{
		HostPort:           fmt.Sprintf("%s:%d", config.Host, config.Port),
		Namespace:          "",
		Logger:             nil,
		MetricsHandler:     &metricsHandler{Scope: scope},
		Identity:           "",
		DataConverter:      nil,
		ContextPropagators: nil,
		ConnectionOptions:  client.ConnectionOptions{},
		HeadersProvider:    nil,
		TrafficController:  nil,
		Interceptors:       nil,
	})
	if err != nil {
		return nil, err
	}

	return nil, fmt.Errorf("implement me")
}
