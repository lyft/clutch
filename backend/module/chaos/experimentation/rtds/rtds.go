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
	gcpCore "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	gcpDiscovery "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v2"
	gcpCache "github.com/envoyproxy/go-control-plane/pkg/cache"
	gcpServer "github.com/envoyproxy/go-control-plane/pkg/server"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/uber-go/tally"
	"go.uber.org/zap"

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
	experimentStore experimentstore.ExperimentStore

	// RTDS built-in cache
	snapshotCache gcpCache.SnapshotCache

	// duration of cache refresh in seconds
	cacheRefreshInterval time.Duration

	// Name of the RTDS layer in Envoy config i.e. envoy.yaml
	rtdsLayerName string

	// Total number of open streams
	totalStreams tally.Gauge

	// Total runtime resources served
	totalResourcesServed tally.Counter

	logger *zap.SugaredLogger
}

// ClusterHash implements NodeHash interface
type ClusterHash struct{}

// ID is an override method to use Cluster instead of a Node
func (ClusterHash) ID(node *gcpCore.Node) string {
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
	rtdsLayerName := config.GetRtdsLayerName()

	store, ok := service.Registry[experimentstore.Name]
	if !ok {
		return nil, errors.New("could not find service")
	}

	experimentStore, ok := store.(experimentstore.ExperimentStore)
	if !ok {
		return nil, errors.New("service was not the correct type")
	}

	initializeGCPCache := gcpCache.NewSnapshotCache(false, ClusterHash{}, logger.Sugar())
	rtdsScope := scope.SubScope("rtds")

	return &Server{
		ctx:                  context.Background(),
		experimentStore:      experimentStore,
		snapshotCache:        initializeGCPCache,
		cacheRefreshInterval: cacheRefreshInterval,
		rtdsLayerName:        rtdsLayerName,
		totalStreams:         rtdsScope.Gauge("totalStreams"),
		totalResourcesServed: rtdsScope.Counter("totalResourcesServed"),
		logger:               logger.Sugar(),
	}, nil
}

func (s *Server) Register(r module.Registrar) error {
	PeriodicallyRefreshCache(s)
	xdsServer := gcpServer.NewServer(s.ctx, s.snapshotCache, &callbacks{s.totalStreams,
		s.totalResourcesServed, s.logger, 0})
	gcpDiscovery.RegisterRuntimeDiscoveryServiceServer(r.GRPCServer(), xdsServer)
	return nil
}

type callbacks struct {
	totalStreams         tally.Gauge
	totalResourcesServed tally.Counter
	logger               *zap.SugaredLogger
	numStreams           int32
}

func (c *callbacks) OnStreamOpen(_ context.Context, streamID int64, typeURL string) error {
	c.logger.Debugw("RTDS onStreamOpen", "streamID", streamID, "typeURL", typeURL)
	numStreams := atomic.AddInt32(&c.numStreams, 1)
	c.totalStreams.Update(float64(numStreams))
	return nil
}

func (c *callbacks) OnStreamClosed(streamID int64) {
	c.logger.Debugw("RTDS onStreamClosed", "streamID", streamID)
	numStreams := atomic.AddInt32(&c.numStreams, -1)
	c.totalStreams.Update(float64(numStreams))
}

func (c *callbacks) OnStreamRequest(streamID int64, request *gcpV2.DiscoveryRequest) error {
	c.logger.Debugw("RTDS OnStreamRequest", "streamID", streamID, "cluster", request.Node.Cluster)
	return nil
}

func (c *callbacks) OnStreamResponse(streamID int64, request *gcpV2.DiscoveryRequest, response *gcpV2.DiscoveryResponse) {
	c.totalResourcesServed.Inc(1)
	c.logger.Debugw("RTDS OnStreamResponse", "streamID", streamID, "cluster", request.Node.Cluster, "version", response.VersionInfo)
}

func (c *callbacks) OnFetchRequest(context.Context, *gcpV2.DiscoveryRequest) error {
	c.logger.Debugw("RTDS OnFetchRequest")
	return nil
}

func (c *callbacks) OnFetchResponse(*gcpV2.DiscoveryRequest, *gcpV2.DiscoveryResponse) {
	c.logger.Debugw("RTDS OnFetchResponse")
}
