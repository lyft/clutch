package modulemock

import (
	mod "github.com/lyft/clutch/backend/module"
	"google.golang.org/grpc"
)

type MockRegistrar struct {
	Server *grpc.Server
}

func (m *MockRegistrar) GRPCServer() *grpc.Server { return m.Server }

func (m *MockRegistrar) RegisterJSONGateway(handlerFunc mod.GatewayRegisterAPIHandlerFunc) error {
	return nil
}
