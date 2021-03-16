package xdstest

import (
	"net"

	"github.com/uber-go/tally"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/lyft/clutch/backend/module"
	"github.com/lyft/clutch/backend/module/moduletest"
)

type TestServer struct {
	registrar *moduletest.TestRegistrar
	Scope     tally.TestScope
	Logger    *zap.Logger
}

func New() *TestServer {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	server := &TestServer{
		registrar: moduletest.NewRegisterChecker(),
		Scope:     tally.NewTestScope("test", nil),
		Logger:    logger,
	}
	return server
}

func (t *TestServer) Register(m module.Module) {
	err := m.Register(t.registrar)
	if err != nil {
		panic(err)
	}
}

func (t *TestServer) Run() {
	//nolint:gosec
	l, err := net.Listen("tcp", "0.0.0.0:9000")
	if err != nil {
		panic(err)
	}
	go func() {
		err = t.registrar.GRPCServer().Serve(l)
		if err != nil {
			panic(err)
		}
	}()
}

func (t *TestServer) Stop() {
	t.registrar.GRPCServer().Stop()
}

func (t *TestServer) ClientConn() (*grpc.ClientConn, error) {
	return grpc.Dial("localhost:9000", grpc.WithInsecure())
}
