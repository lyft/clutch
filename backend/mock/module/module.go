package modulemock

import (
	"google.golang.org/grpc"

	mod "github.com/lyft/clutch/backend/module"
)

type MockRegistrar struct {
	Server *grpc.Server
}

func (m *MockRegistrar) GRPCServer() *grpc.Server { return m.Server }

func (m *MockRegistrar) RegisterJSONGateway(handlerFunc mod.GatewayRegisterAPIHandlerFunc) error {
	return nil
}
