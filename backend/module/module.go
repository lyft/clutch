package module

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/uber-go/tally/v4"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/anypb"
)

type GatewayRegisterAPIHandlerFunc func(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error

type Registrar interface {
	GRPCServer() *grpc.Server
	RegisterJSONGateway(GatewayRegisterAPIHandlerFunc) error
}

type Module interface {
	Register(Registrar) error
}

type Factory map[string]func(*anypb.Any, *zap.Logger, tally.Scope) (Module, error)
