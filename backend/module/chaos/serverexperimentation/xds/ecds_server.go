package xds

import (
	"context"
	"errors"

	discovery "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	extensionconfigservice "github.com/envoyproxy/go-control-plane/envoy/service/extension/v3"
	"github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	"github.com/envoyproxy/go-control-plane/pkg/resource/v3"
	gcpServerV3 "github.com/envoyproxy/go-control-plane/pkg/server/v3"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ECDSServer interface {
	extensionconfigservice.ExtensionConfigDiscoveryServiceServer
}

type ecdsServer struct {
	server gcpServerV3.Server
}

func NewECDSServer(ctx context.Context, config cache.Cache, callbacks gcpServerV3.Callbacks) ECDSServer {
	return &ecdsServer{
		server: gcpServerV3.NewServer(ctx, config, callbacks),
	}
}

func (e ecdsServer) StreamExtensionConfigs(stream extensionconfigservice.ExtensionConfigDiscoveryService_StreamExtensionConfigsServer) error {
	return e.server.StreamHandler(stream, resource.ExtensionConfigType)
}

func (e ecdsServer) DeltaExtensionConfigs(extensionconfigservice.ExtensionConfigDiscoveryService_DeltaExtensionConfigsServer) error {
	return errors.New("not implemented")
}

func (e ecdsServer) FetchExtensionConfigs(ctx context.Context, req *discovery.DiscoveryRequest) (*discovery.DiscoveryResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.Unavailable, "empty request")
	}
	req.TypeUrl = resource.ExtensionConfigType
	return e.server.Fetch(ctx, req)
}
