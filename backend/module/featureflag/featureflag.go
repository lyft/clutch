package featureflag

import (
	"context"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/uber-go/tally"
	"go.uber.org/zap"

	featureflagcfgv1 "github.com/lyft/clutch/backend/api/config/module/featureflag/v1"
	featureflagv1 "github.com/lyft/clutch/backend/api/featureflag/v1"
	"github.com/lyft/clutch/backend/module"
)

const (
	Name = "clutch.module.featureflag"
)

func New(cfg *any.Any, log *zap.Logger, scope tally.Scope) (module.Module, error) {
	config := &featureflagcfgv1.Config{}

	if cfg != nil {
		if err := ptypes.UnmarshalAny(cfg, config); err != nil {
			return nil, err
		}
	}

	return newModuleImpl(config.GetSimple())
}

type featureFlagIface interface {
	module.Module
	featureflagv1.FeatureFlagAPIServer
}

func newModuleImpl(simple *featureflagcfgv1.Simple) (featureFlagIface, error) {
	return &moduleImpl{simple: simple}, nil
}

type moduleImpl struct {
	simple *featureflagcfgv1.Simple
}

func (m *moduleImpl) Register(r module.Registrar) error {
	featureflagv1.RegisterFeatureFlagAPIServer(r.GRPCServer(), m)
	return r.RegisterJSONGateway(featureflagv1.RegisterFeatureFlagAPIHandler)
}

func (m *moduleImpl) GetFlags(ctx context.Context, req *featureflagv1.GetFlagsRequest) (*featureflagv1.GetFlagsResponse, error) {
	flags := make(map[string]*featureflagv1.Flag, len(m.simple.Flags))
	for i, flag := range m.simple.Flags {
		flags[i] = &featureflagv1.Flag{Type: &featureflagv1.Flag_BooleanValue{BooleanValue: flag}}
	}

	return &featureflagv1.GetFlagsResponse{Flags: flags}, nil
}
