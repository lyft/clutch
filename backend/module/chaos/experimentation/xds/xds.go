package xds

// <!-- START clutchdoc -->
// description: Chaos Experimentation Framework - Envoy Discovery Service (xDS) implementation that delivers chaos experiment values to subscribed Envoys.
// <!-- END clutchdoc -->

import (
	"context"
	"errors"
	"sync/atomic"

	gcpCoreV3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	gcpDiscoveryV3 "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	gcpExtencionServiceV3 "github.com/envoyproxy/go-control-plane/envoy/service/extension/v3"
	gcpRuntimeServiceV3 "github.com/envoyproxy/go-control-plane/envoy/service/runtime/v3"
	gcpCacheV3 "github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	gcpServerV3 "github.com/envoyproxy/go-control-plane/pkg/server/v3"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/uber-go/tally"
	"go.uber.org/zap"
	rpc_status "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/protobuf/proto"

	xdsconfigv1 "github.com/lyft/clutch/backend/api/config/module/chaos/experimentation/xds/v1"
	"github.com/lyft/clutch/backend/module"
	"github.com/lyft/clutch/backend/service"
	"github.com/lyft/clutch/backend/service/chaos/experimentation/experimentstore"
)

const Name = "clutch.module.chaos.experimentation.xds"

func TypeUrl(message proto.Message) string {
	return "type.googleapis.com/" + string(message.ProtoReflect().Descriptor().FullName())
}

// Server serves xDS
type Server struct {
	ctx context.Context

	poller *Poller

	ecdsConfig *ECDSConfig

	scope  tally.Scope
	logger *zap.SugaredLogger
}

// ClusterHash implements NodeHash interface
type ClusterHashV3 struct{}

// ID is an override method to use Cluster instead of a Node
func (ClusterHashV3) ID(node *gcpCoreV3.Node) string {
	if node == nil {
		return ""
	}
	return node.Cluster
}

func New(cfg *any.Any, logger *zap.Logger, scope tally.Scope) (module.Module, error) {
	config := &xdsconfigv1.Config{}
	if err := cfg.UnmarshalTo(config); err != nil {
		return nil, err
	}

	store, ok := service.Registry[experimentstore.Name]
	if !ok {
		return nil, errors.New("could not find experiment store service")
	}

	storer, ok := store.(experimentstore.Storer)
	if !ok {
		return nil, errors.New("service was not the correct type")
	}

	enabledECDSClusters := make(map[string]struct{})
	for _, cluster := range config.GetEcdsAllowList().GetEnabledClusters() {
		enabledECDSClusters[cluster] = struct{}{}
	}

	ecdsConfig := &ECDSConfig{
		enabledClusters: enabledECDSClusters,
		ecdsResourceMap: &SafeEcdsResourceMap{
			requestedResourcesMap: make(map[string]map[string]struct{}),
		},
	}

	ctx := context.Background()
	cacheV3 := gcpCacheV3.NewSnapshotCacheWithHeartbeating(ctx, false, ClusterHashV3{}, logger.Sugar(), config.HeartbeatInterval.AsDuration())

	poller := Poller{
		storer: storer,
		cache:  cacheV3,

		resourceTtl:          config.ResourceTtl.AsDuration(),
		cacheRefreshInterval: config.CacheRefreshInterval.AsDuration(),

		rtdsConfig: &RTDSConfig{
			layerName: config.RtdsLayerName,
		},
		ecdsConfig: ecdsConfig,

		rtdsGeneratorsByTypeUrl: RTDSGeneratorsByTypeUrl,
		ecdsGeneratorsByTypeUrl: ECDSGeneratorsByTypeUrl,

		rtdsResourceGenerationFailureCount:        scope.Counter("rtds_resource_generation_failure"),
		ecdsResourceGenerationFailureCount:        scope.Counter("ecds_resource_generation_failure"),
		ecdsDefaultResourceGenerationFailureCount: scope.Counter("ecds_default_resource_generation_failure"),
		setCacheSnapshotSuccessCount:              scope.Counter("set_snapshot_success"),
		setCacheSnapshotFailureCount:              scope.Counter("set_snapshot_failure"),
		activeFaultsGauge:                         scope.Gauge("active_faults"),

		logger: logger.Sugar(),
	}

	return &Server{
		ctx:        context.Background(),
		poller:     &poller,
		ecdsConfig: ecdsConfig,
		scope:      scope,

		logger: logger.Sugar(),
	}, nil
}

type serverStats struct {
	totalStreams         tally.Gauge
	totalResourcesServed tally.Counter
	totalErrorsReceived  tally.Counter
}

func (s *Server) newScopedStats(subScope string) serverStats {
	scope := s.scope.SubScope(subScope)
	return serverStats{
		totalStreams:         scope.Gauge("totalStreams"),
		totalResourcesServed: scope.Counter("totalResourcesServed"),
		totalErrorsReceived:  scope.Counter("totalErrorsReceived"),
	}
}

func (s *Server) Register(r module.Registrar) error {
	ctx := context.Background()
	s.poller.Start(ctx)
	// RTDS V3 Server
	rtdsServer := gcpServerV3.NewServer(s.ctx, s.poller.cache, &rtdsCallbacks{callbacksBase{s.newScopedStats("rtds"),
		s.logger, 0}})
	gcpRuntimeServiceV3.RegisterRuntimeDiscoveryServiceServer(r.GRPCServer(), rtdsServer)

	ecdsServer := NewECDSServer(s.ctx, s.poller.cache, &ecdsCallbacks{callbacksBase{s.newScopedStats("ecds"), s.logger, 0}, s.ecdsConfig.ecdsResourceMap})
	gcpExtencionServiceV3.RegisterExtensionConfigDiscoveryServiceServer(r.GRPCServer(), ecdsServer)
	return nil
}

type callbacksBase struct {
	serverStats serverStats
	logger      *zap.SugaredLogger
	numStreams  int32
}

func (c *callbacksBase) onStreamOpen(_ context.Context) error {
	numStreams := atomic.AddInt32(&c.numStreams, 1)
	c.serverStats.totalStreams.Update(float64(numStreams))
	return nil
}

func (c *callbacksBase) onStreamClosed(streamID int64) {
	numStreams := atomic.AddInt32(&c.numStreams, -1)
	c.serverStats.totalStreams.Update(float64(numStreams))
}

func (c *callbacksBase) onStreamRequest(streamID int64, cluster string, errorDetail *rpc_status.Status) {
	if errorDetail != nil {
		c.serverStats.totalErrorsReceived.Inc(1)
		c.logger.Errorw("xDS Error Request", "streamID", streamID, "cluster", cluster, "error", errorDetail.GetDetails())
	}
}

func (c *callbacksBase) onStreamResponse(_ int64, _ string, _ string) {
	c.serverStats.totalResourcesServed.Inc(1)
}

// RTDS Callbacks
type rtdsCallbacks struct {
	callbacksBase
}

func (c *rtdsCallbacks) OnStreamOpen(ctx context.Context, streamID int64, typeURL string) error {
	c.logger.Debugw("RTDS onStreamOpen", "streamID", streamID, "typeURL", typeURL)
	return c.onStreamOpen(ctx)
}

func (c *rtdsCallbacks) OnStreamClosed(streamID int64) {
	c.logger.Debugw("RTDS onStreamClosed", "streamID", streamID)
	c.onStreamClosed(streamID)
}

func (c *rtdsCallbacks) OnStreamRequest(streamID int64, req *gcpDiscoveryV3.DiscoveryRequest) error {
	c.logger.Debugw("RTDS OnStreamRequest", "streamID", streamID, "cluster", req.Node.Cluster)
	c.onStreamRequest(streamID, req.Node.Cluster, req.ErrorDetail)

	return nil
}

func (c *rtdsCallbacks) OnStreamResponse(streamID int64, request *gcpDiscoveryV3.DiscoveryRequest, response *gcpDiscoveryV3.DiscoveryResponse) {
	c.logger.Debugw("RTDS OnStreamResponse", "streamID", streamID, "cluster", request.Node.Cluster, "version", request.VersionInfo)
	c.onStreamResponse(streamID, request.Node.Cluster, request.VersionInfo)
}

func (c *rtdsCallbacks) OnFetchRequest(context.Context, *gcpDiscoveryV3.DiscoveryRequest) error {
	c.logger.Debugw("RTDS OnFetchRequest")
	return nil
}

func (c *rtdsCallbacks) OnFetchResponse(*gcpDiscoveryV3.DiscoveryRequest, *gcpDiscoveryV3.DiscoveryResponse) {
	c.logger.Debugw("RTDS OnFetchResponse")
}

// ECDS Callbacks
type ecdsCallbacks struct {
	callbacksBase

	// Track all the seen ECDS resources globally across all the streams. This allows us to query all the requested
	// resources and present a default value for all the ones that don't have a specific value.
	// This allows us to set a default value for all the dynamic ECDS resources for all clusters, relying on go-control-plane
	// to only respond with ones actually requested by the client.
	safeECDSResources *SafeEcdsResourceMap
}

func (c *ecdsCallbacks) OnStreamOpen(ctx context.Context, streamID int64, typeURL string) error {
	c.logger.Debugw("ECDS onStreamOpen", "streamID", streamID, "typeURL", typeURL)
	return c.onStreamOpen(ctx)
}

func (c *ecdsCallbacks) OnStreamClosed(streamID int64) {
	c.logger.Debugw("ECDS onStreamClosed", "streamID", streamID)
	c.onStreamClosed(streamID)
}

func (c *ecdsCallbacks) OnStreamRequest(streamID int64, req *gcpDiscoveryV3.DiscoveryRequest) error {
	c.safeECDSResources.setResourcesForCluster(req.Node.Cluster, req.ResourceNames)

	c.logger.Debugw("ECDS OnStreamRequest", "streamID", streamID, "cluster", req.Node.Cluster, "resources", req.ResourceNames)
	c.onStreamRequest(streamID, req.Node.Cluster, req.ErrorDetail)

	return nil
}

func (c *ecdsCallbacks) OnStreamResponse(streamID int64, request *gcpDiscoveryV3.DiscoveryRequest, response *gcpDiscoveryV3.DiscoveryResponse) {
	c.logger.Debugw("ECDS OnStreamResponse", "streamID", streamID, "cluster", request.Node.Cluster, "version", request.VersionInfo)
	c.onStreamResponse(streamID, request.Node.Cluster, request.VersionInfo)
}

func (c *ecdsCallbacks) OnFetchRequest(context.Context, *gcpDiscoveryV3.DiscoveryRequest) error {
	c.logger.Debugw("ECDS OnFetchRequest")
	return nil
}

func (c *ecdsCallbacks) OnFetchResponse(*gcpDiscoveryV3.DiscoveryRequest, *gcpDiscoveryV3.DiscoveryResponse) {
	c.logger.Debugw("ECDS OnFetchResponse")
}
