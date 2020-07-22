package echo

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/ptypes/any"
	"github.com/lyft/clutch/backend/module"
	"github.com/uber-go/tally"
	"go.uber.org/zap"

	echov1 "{{ .RepoProvider}}/{{ .RepoOwner }}/{{ .RepoName}}/backend/api/echo/v1"
)

const Name = "{{ .RepoOwner }}.module.echo"

func New(*any.Any, *zap.Logger, tally.Scope) (module.Module, error) {
	mod := &mod{
		api: newAPI(),
	}
	return mod, nil
}

type mod struct {
	api echov1.EchoAPIServer
}

func (m *mod) Register(r module.Registrar) error {
	echov1.RegisterEchoAPIServer(r.GRPCServer(), m.api)
	return r.RegisterJSONGateway(echov1.RegisterEchoAPIHandler)
}

func newAPI() echov1.EchoAPIServer {
	return &api{}
}

type api struct{}

func (a *api) SayHello(ctx context.Context, req *echov1.SayHelloRequest) (*echov1.SayHelloResponse, error) {
	return &echov1.SayHelloResponse{Message: fmt.Sprintf("Hello %s!", req.Name)}, nil
}
