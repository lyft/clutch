package proxy

import (
	"context"
	"errors"

	"github.com/golang/protobuf/ptypes/any"
	proxyv1cfg "github.com/lyft/clutch/backend/api/config/module/proxy/v1"
	proxyv1 "github.com/lyft/clutch/backend/api/proxy/v1"
	"github.com/lyft/clutch/backend/module"
	"github.com/uber-go/tally"
	"go.uber.org/zap"
)

const Name = "clutch.module.proxy"

func New(cfg *any.Any, log *zap.Logger, scope tally.Scope) (module.Module, error) {
	config := &proxyv1cfg.Config{}
	err := cfg.UnmarshalTo(config)
	if err != nil {
		return nil, err
	}

	m := &mod{
		logger: log,
		scope:  scope,
	}

	return m, nil
}

type mod struct {
	logger *zap.Logger
	scope  tally.Scope
}

func (m *mod) Register(r module.Registrar) error {
	proxyv1.RegisterProxyAPIServer(r.GRPCServer(), m)
	return r.RegisterJSONGateway(proxyv1.RegisterProxyAPIHandler)
}

func (m *mod) RequestProxy(context.Context, *proxyv1.RequestProxyRequest) (*proxyv1.RequestProxyResponse, error) {
	return nil, errors.New("not implemented")
}
