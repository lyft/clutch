package rtds

// <!-- START clutchdoc -->
// description: Runtime Discovery Service (RTDS) implementation that delivers chaos experiment values to subscribed Envoys.
// <!-- END clutchdoc -->

import (
	"context"
	"errors"
	"sync/atomic"
	"time"

	gcpV2 "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	gcpCoreV2 "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	gcpCoreV3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	gcpDiscoveryV2 "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v2"
	gcpDiscoveryV3 "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	gcpRuntimeServiceV3 "github.com/envoyproxy/go-control-plane/envoy/service/runtime/v3"
	gcpCacheV2 "github.com/envoyproxy/go-control-plane/pkg/cache/v2"
	gcpCacheV3 "github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	gcpServerV2 "github.com/envoyproxy/go-control-plane/pkg/server/v2"
	gcpServerV3 "github.com/envoyproxy/go-control-plane/pkg/server/v3"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/uber-go/tally"
	"go.uber.org/zap"
	rpc_status "google.golang.org/genproto/googleapis/rpc/status"

	rtdsconfigv1 "github.com/lyft/clutch/backend/api/config/module/chaos/experimentation/rtds/v1"
	"github.com/lyft/clutch/backend/module"
	"github.com/lyft/clutch/backend/service"
	"github.com/lyft/clutch/backend/service/chaos/experimentation/experimentstore"
)

const Name = "clutch.module.chaos.experimentation.rtds"

// Server serves RTDS
type Server struct {
	ctx context.Context

	// Experiment store
	storer experimentstore.Storer

	// RTDS built-in cache for V2 xDS
	snapshotCacheV2 gcpCacheV2.SnapshotCache

	// RTDS built-in cache for V3 xDS
	snapshotCacheV3 gcpCacheV3.SnapshotCache

	// duration of cache refresh in seconds
	cacheRefreshInterval time.Duration

	// Name of the RTDS layer in Envoy config i.e. envoy.yaml
	rtdsLayerName string

	// Runtime prefix for ingress faults
	ingressPrefix string

	// Runtime prefix for egress faults
	egressPrefix string

	// The TTL to set for RTDS resources.
	resourceTTL *time.Duration

	// The heartbeat interval to use for TTL'd resources.
	heartbeatInterval *time.Duration

	rtdsScope tally.Scope

	logger *zap.SugaredLogger
}

// ClusterHash implements NodeHash interface
type ClusterHashV2 struct{}

// ID is an override method to use Cluster instead of a Node
func (ClusterHashV2) ID(node *gcpCoreV2.Node) string {
	if node == nil {
		return ""
	}
	return node.Cluster
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
	config := &rtdsconfigv1.Config{}
	if err := ptypes.UnmarshalAny(cfg, config); err != nil {
		return nil, err
	}

	cacheRefreshInterval, err := ptypes.Duration(config.GetCacheRefreshInterval())
	if err != nil {
		return nil, errors.New("error parsing duration")
	}

	store, ok := service.Registry[experimentstore.Name]
	if !ok {
		return nil, errors.New("could not find experiment store service")
	}

	storer, ok := store.(experimentstore.Storer)
	if !ok {
		return nil, errors.New("service was not the correct type")
	}

	var heartbeatInterval *time.Duration
	var resourceTTL *time.Duration
	if config.ResourceTtl != nil {
		d, err := ptypes.Duration(config.ResourceTtl)
		if err != nil {
			return nil, err
		}
		resourceTTL = &d

		if config.HeartbeatInterval != nil {
			d, err := ptypes.Duration(config.HeartbeatInterval)
			if err != nil {
				return nil, err
			}
			heartbeatInterval = &d
		}
	}

	ctx := context.Background()
	var cacheV3 gcpCacheV3.SnapshotCache
	if heartbeatInterval != nil {
		cacheV3 = gcpCacheV3.NewSnapshotCacheWithHeartbeating(ctx, false, ClusterHashV3{}, logger.Sugar(), *heartbeatInterval)
	} else {
		cacheV3 = gcpCacheV3.NewSnapshotCache(false, ClusterHashV3{}, logger.Sugar())
	}
	gcpCacheV2 := gcpCacheV2.NewSnapshotCache(false, ClusterHashV2{}, logger.Sugar())

	return &Server{
		ctx:                  ctx,
		storer:               storer,
		snapshotCacheV2:      gcpCacheV2,
		snapshotCacheV3:      cacheV3,
		cacheRefreshInterval: cacheRefreshInterval,
		rtdsLayerName:        config.GetRtdsLayerName(),
		ingressPrefix:        config.GetIngressFaultRuntimePrefix(),
		egressPrefix:         config.GetEgressFaultRuntimePrefix(),
		resourceTTL:          resourceTTL,
		heartbeatInterval:    heartbeatInterval,
		rtdsScope:            scope,
		logger:               logger.Sugar(),
	}, nil
}

type serverStats struct {
	totalStreams         tally.Gauge
	totalResourcesServed tally.Counter
	totalErrorsReceived  tally.Counter
}

func (s *Server) newScopedStats(subScope string) serverStats {
	scope := s.rtdsScope.SubScope(subScope)
	return serverStats{
		totalStreams:         scope.Gauge("totalStreams"),
		totalResourcesServed: scope.Counter("totalResourcesServed"),
		totalErrorsReceived:  scope.Counter("totalErrorsReceived"),
	}
}

func (s *Server) Register(r module.Registrar) error {
	PeriodicallyRefreshCache(s)
	xdsServerV2 := gcpServerV2.NewServer(s.ctx, s.snapshotCacheV2, &callbacksV2{callbacksBase{s.newScopedStats("v2"), s.logger, 0}})
	xdsServerV3 := gcpServerV3.NewServer(s.ctx, s.snapshotCacheV3, &callbacksV3{callbacksBase{s.newScopedStats("v3"),
		s.logger, 0}})
	gcpRuntimeServiceV3.RegisterRuntimeDiscoveryServiceServer(r.GRPCServer(), xdsServerV3)
	gcpDiscoveryV2.RegisterRuntimeDiscoveryServiceServer(r.GRPCServer(), xdsServerV2)
	return nil
}

type callbacksBase struct {
	serverStats serverStats
	logger      *zap.SugaredLogger
	numStreams  int32
}

type callbacksV3 struct {
	callbacksBase
}

func (c *callbacksBase) OnStreamOpen(_ context.Context, streamID int64, typeURL string) error {
	c.logger.Debugw("RTDS onStreamOpen", "streamID", streamID, "typeURL", typeURL)
	numStreams := atomic.AddInt32(&c.numStreams, 1)
	c.serverStats.totalStreams.Update(float64(numStreams))
	return nil
}

func (c *callbacksBase) OnStreamClosed(streamID int64) {
	c.logger.Debugw("RTDS onStreamClosed", "streamID", streamID)
	numStreams := atomic.AddInt32(&c.numStreams, -1)
	c.serverStats.totalStreams.Update(float64(numStreams))
}

func (c *callbacksBase) onStreamRequest(streamID int64, cluster string, errorDetail *rpc_status.Status) {
	c.logger.Debugw("RTDS OnStreamRequest", "streamID", streamID, "cluster", cluster)
	if errorDetail != nil {
		c.serverStats.totalErrorsReceived.Inc(1)
		c.logger.Errorw("RTDS Error Request", "error", errorDetail.GetDetails())
	}
}

func (c *callbacksBase) onStreamResponse(streamID int64, cluster string, version string) {
	c.serverStats.totalResourcesServed.Inc(1)
	c.logger.Debugw("RTDS OnStreamResponse", "streamID", streamID, "cluster", cluster, "version", version)
}

func (c *callbacksV3) OnStreamRequest(streamID int64, request *gcpDiscoveryV3.DiscoveryRequest) error {
	c.onStreamRequest(streamID, request.Node.Cluster, request.ErrorDetail)
	return nil
}

func (c *callbacksV3) OnStreamResponse(streamID int64, request *gcpDiscoveryV3.DiscoveryRequest, response *gcpDiscoveryV3.DiscoveryResponse) {
	c.onStreamResponse(streamID, request.Node.Cluster, request.VersionInfo)
}

func (c *callbacksV3) OnFetchRequest(context.Context, *gcpDiscoveryV3.DiscoveryRequest) error {
	c.logger.Debugw("RTDS OnFetchRequest")
	return nil
}

func (c *callbacksV3) OnFetchResponse(*gcpDiscoveryV3.DiscoveryRequest, *gcpDiscoveryV3.DiscoveryResponse) {
	c.logger.Debugw("RTDS OnFetchResponse")
}

type callbacksV2 struct {
	callbacksBase
}

func (c *callbacksV2) OnStreamRequest(streamID int64, request *gcpV2.DiscoveryRequest) error {
	c.onStreamRequest(streamID, request.Node.Cluster, request.ErrorDetail)
	return nil
}

func (c *callbacksV2) OnStreamResponse(streamID int64, request *gcpV2.DiscoveryRequest, response *gcpV2.DiscoveryResponse) {
	c.onStreamResponse(streamID, request.Node.Cluster, request.VersionInfo)
}

func (c *callbacksV2) OnFetchRequest(context.Context, *gcpV2.DiscoveryRequest) error {
	c.logger.Debugw("RTDS OnFetchRequest")
	return nil
}

func (c *callbacksV2) OnFetchResponse(*gcpV2.DiscoveryRequest, *gcpV2.DiscoveryResponse) {
	c.logger.Debugw("RTDS OnFetchResponse")
}
