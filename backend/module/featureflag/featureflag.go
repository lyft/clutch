package featureflag

import (
	"context"
	"github.com/golang/protobuf/ptypes/any"
	featureflagv1 "github.com/lyft/clutch/backend/api/featureflag/v1"
	"github.com/lyft/clutch/backend/module"
	"github.com/uber-go/tally"
	"go.uber.org/zap"
)

const (
	Name = "clutch.module.featureflag"
)

func New(*any.Any, *zap.Logger, tally.Scope) (module.Module, error) {
	return newModuleImpl()
}

type moduleIface interface {
	module.Module
	featureflagv1.FeatureFlagAPIServer
}

func newModuleImpl() (moduleIface, error) {
	return &moduleImpl{}, nil
}

type moduleImpl struct {

}

func (m *moduleImpl) Register(r module.Registrar) error {
	featureflagv1.RegisterFeatureFlagAPIServer(r.GRPCServer(), m)
	return r.RegisterJSONGateway(featureflagv1.RegisterFeatureFlagAPIHandler)
}

func (m *moduleImpl) GetFlags(ctx context.Context, req *featureflagv1.GetFlagsRequest) (*featureflagv1.GetFlagsResponse, error) {
	panic("implement me")
}

