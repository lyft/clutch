// {{ .ModuleName }}:
//
//   {{ .Description }}
package {{ .ModuleName }}

import (
	"context"
	"errors"

	"github.com/golang/protobuf/ptypes/any"
	"github.com/uber-go/tally/v4"
	"go.uber.org/zap"

	{{ .ModuleName }}{{ .ApiVersion }} "{{ .BackendGoPackage }}/api/{{ .ModuleName }}/{{ .ApiVersion }}"
	"github.com/lyft/clutch/backend/module"
{{ if .ServiceName }}
	"github.com/lyft/clutch/backend/service"
	{{ .ServiceName }}svc "{{ .BackendGoPackage }}/service/{{ .ServiceName }}" 
{{ end }}
)

const (
	Name = "{{ .RepoOwner }}.module.{{ .ModuleName }}"
)

func New(cfg *any.Any, logger *zap.Logger, scope tally.Scope) (module.Module, error) {
{{ if .ServiceName }}
	{{ .ServiceName }}Client, ok := service.Registry[{{ .ServiceName }}svc.Name]
	if !ok {
		return nil, errors.New("could not find service {{ .ServiceName }}")
	}

	c, ok := {{ .ServiceName }}Client.({{ .ServiceName }}svc.Client)
	if !ok {
		return nil, errors.New("{{ .ServiceName }}: service was not the correct type")
	}

	mod := &mod{
		client: c,
		logger: logger,
		scope:  scope,
	}
{{ else }}
	mod := &mod{
		logger: logger,
		scope:  scope,
	}
{{ end }}
	return mod, nil
}
{{ if .ServiceName }}
type mod struct {
	client {{ .ServiceName }}svc.Client
	logger *zap.Logger
	scope  tally.Scope
}
{{ else }}
type mod struct {
	logger *zap.Logger
	scope  tally.Scope
}
{{ end }}
func (m *mod) Register(r module.Registrar) error {
	{{ .ModuleName }}{{ .ApiVersion }}.Register{{ .ApiName }}APIServer(r.GRPCServer(), m)
	return r.RegisterJSONGateway({{ .ModuleName }}{{ .ApiVersion }}.Register{{ .ApiName }}APIHandler)
}
{{ if .ServiceName }}
func (m *mod) Greet(ctx context.Context, req *{{ .ModuleName }}{{ .ApiVersion }}.GreetRequest) (*{{ .ModuleName }}{{ .ApiVersion }}.GreetResponse, error) {
	response, err := m.client.Greet(ctx, req.Name)
	if err != nil {
		return nil, err
	}
	return &{{ .ModuleName }}{{ .ApiVersion }}.GreetResponse{
		Message: response.Message,
	}, nil
}
{{ else }}
func (m *mod) Greet(ctx context.Context, req *{{ .ModuleName }}{{ .ApiVersion }}.GreetRequest) (*{{ .ModuleName }}{{ .ApiVersion }}.GreetResponse, error) {
	return &{{ .ModuleName }}{{ .ApiVersion }}.GreetResponse{
		Message: "Hello, " + req.Name,
	}, nil
}
{{ end }}
