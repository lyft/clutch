package moduletest

import (
	"context"
	"fmt"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"

	"github.com/lyft/clutch/backend/module"
)

// Similar to httptest, this is a helper package to allow module authors to verify modules register correctly.
// Note: this is only for testing module behavior. The actual RPC implementation should be tested directly against the API server.

type TestRegistrar struct {
	jsonCount int

	grpcServer *grpc.Server
	mux        *runtime.ServeMux
}

func NewRegisterChecker() *TestRegistrar {
	return &TestRegistrar{
		grpcServer: grpc.NewServer(),
		mux:        runtime.NewServeMux(),
	}
}

func (r *TestRegistrar) JSONRegistered() bool {
	return r.jsonCount >= 1
}

func (r *TestRegistrar) GRPCRegistered() bool {
	return len(r.grpcServer.GetServiceInfo()) >= 1
}

func (r *TestRegistrar) GRPCServer() *grpc.Server {
	return r.grpcServer
}

func (r *TestRegistrar) RegisterJSONGateway(handlerFunc module.GatewayRegisterAPIHandlerFunc) error {
	r.jsonCount++
	if r.jsonCount != len(r.grpcServer.GetServiceInfo()) {
		panic("RegisterJSONGateway called more than gRPC or no gRPC registration found")
	}
	if err := handlerFunc(context.TODO(), r.mux, nil); err != nil {
		// Panic in case error was ignored by caller.
		panic(err)
	}
	return nil
}

func (r *TestRegistrar) HasAPI(name string) error {
	services := r.grpcServer.GetServiceInfo()
	if _, ok := services[name]; !ok {
		keys := make([]string, 0, len(services))
		for key := range services {
			keys = append(keys, key)
		}
		return fmt.Errorf("service '%s' not found in %d service(s): %+v", name, len(keys), keys)
	}
	return nil
}
