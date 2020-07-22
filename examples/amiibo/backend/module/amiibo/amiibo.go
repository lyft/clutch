// +build ignore

package amiibo

import (
	"context"
	"errors"

	"github.com/golang/protobuf/ptypes/any"
	"github.com/uber-go/tally"
	"go.uber.org/zap"

	amiibov1 "github.com/lyft/clutch/backend/api/amiibo/v1"
	"github.com/lyft/clutch/backend/module"
	"github.com/lyft/clutch/backend/service"
	amiiboservice "github.com/lyft/clutch/backend/service/amiibo"
)

const Name = "clutch.module.amiibo"

func New(*any.Any, *zap.Logger, tally.Scope) (module.Module, error) {
	svc, ok := service.Registry["clutch.service.amiibo"]
	if !ok {
		return nil, errors.New("no amiibo service was registered")
	}

	client, ok := svc.(amiiboservice.Client)
	if !ok {
		return nil, errors.New("amiibo service in registry was the wrong type")
	}

	return &mod{client: client}, nil
}

type mod struct {
	client amiiboservice.Client
}

func (m *mod) Register(r module.Registrar) error {
	amiibov1.RegisterAmiiboAPIServer(r.GRPCServer(), m)
	return r.RegisterJSONGateway(amiibov1.RegisterAmiiboAPIHandler)
}

func (m *mod) GetAmiibo(ctx context.Context, request *amiibov1.GetAmiiboRequest) (*amiibov1.GetAmiiboResponse, error) {
	a, err := m.client.GetAmiibo(ctx, request.Name)
	if err != nil {
		return nil, err
	}
	return &amiibov1.GetAmiiboResponse{Amiibo: a}, nil
}
