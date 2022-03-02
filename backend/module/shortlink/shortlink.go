package shortlink

import (
	"context"
	"errors"

	"github.com/golang/protobuf/ptypes/any"
	"github.com/uber-go/tally/v4"
	"go.uber.org/zap"

	shortlinkv1 "github.com/lyft/clutch/backend/api/shortlink/v1"
	"github.com/lyft/clutch/backend/module"
	"github.com/lyft/clutch/backend/service"
	"github.com/lyft/clutch/backend/service/shortlink"
)

const (
	Name = "clutch.module.shortlink"
)

func New(*any.Any, *zap.Logger, tally.Scope) (module.Module, error) {
	client, ok := service.Registry["clutch.service.shortlink"]
	if !ok {
		return nil, errors.New("could not find shortlink service")
	}

	svc, ok := client.(shortlink.Service)
	if !ok {
		return nil, errors.New("service was not the correct type")
	}

	mod := &mod{
		shortlink: newShortlinkAPI(svc),
	}

	return mod, nil
}

type mod struct {
	shortlink shortlinkv1.ShortlinkAPIServer
}

func (m *mod) Register(r module.Registrar) error {
	shortlinkv1.RegisterShortlinkAPIServer(r.GRPCServer(), m.shortlink)
	return r.RegisterJSONGateway(shortlinkv1.RegisterShortlinkAPIHandler)
}

type shortlinkAPI struct {
	shortlink shortlink.Service
}

func newShortlinkAPI(svc shortlink.Service) shortlinkv1.ShortlinkAPIServer {
	return &shortlinkAPI{
		shortlink: svc,
	}
}

func (s *shortlinkAPI) Create(ctx context.Context, req *shortlinkv1.CreateRequest) (*shortlinkv1.CreateResponse, error) {
	return nil, errors.New("not implemented")
}

func (s *shortlinkAPI) Get(ctx context.Context, req *shortlinkv1.GetRequest) (*shortlinkv1.GetResponse, error) {
	return nil, errors.New("not implemented")
}
