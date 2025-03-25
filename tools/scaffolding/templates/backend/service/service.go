package {{ .ServiceName }}

import (
	"context"

	"github.com/golang/protobuf/ptypes/any"
	"github.com/lyft/clutch/backend/service"
	"github.com/uber-go/tally/v4"
	"go.uber.org/zap"
)

const Name = "{{ .RepoOwner }}.service.{{ .ServiceName }}"

func New(cfg *any.Any, logger *zap.Logger, scope tally.Scope) (service.Service, error) {
	ret := &svc{
		scope:  scope,
		logger: logger,
	}
	return ret, nil
}

type svc struct {
	logger *zap.Logger
	scope  tally.Scope
}

type Client interface {
	Greet(ctx context.Context, name string) (*GreetResponse, error)
}

type GreetResponse struct {
	Message string
}

func (c *svc) Greet(ctx context.Context, name string) (*GreetResponse, error) {
	response := &GreetResponse{
		Message: "Hello, " + name,
	}
	return response, nil
}
