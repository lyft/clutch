package module

import (
	"context"

	"github.com/uber-go/tally"
	"go.uber.org/zap"

	"github.com/golang/protobuf/ptypes/any"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
)

type GatewayRegisterAPIHandlerFunc func(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error

type Registrar interface {
	GRPCServer() *grpc.Server
	RegisterJSONGateway(GatewayRegisterAPIHandlerFunc) error
}

type Module interface {
	Register(Registrar) error
}

type Factory map[string]func(*any.Any, *zap.Logger, tally.Scope) (Module, error)
