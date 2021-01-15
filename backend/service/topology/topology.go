package topology

import (
	"database/sql"
	"encoding/json"
	"errors"
	"os"
	"os/signal"
	"syscall"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/uber-go/tally"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"google.golang.org/protobuf/encoding/protojson"

	topologyv1cfg "github.com/lyft/clutch/backend/api/config/service/topology/v1"
	topologyv1 "github.com/lyft/clutch/backend/api/topology/v1"
	"github.com/lyft/clutch/backend/service"
	pgservice "github.com/lyft/clutch/backend/service/db/postgres"
)

const Name = "clutch.service.topology"

type Service interface {
	GetTopology(ctx context.Context) error
	Search(ctx context.Context, search *topologyv1.SearchRequest) ([]*topologyv1.Resource, uint64, error)
}

type client struct {
	config *topologyv1cfg.Config

	db    *sql.DB
	log   *zap.Logger
	scope tally.Scope
}

// CacheableTopology is implemented by a service that wishes to enable the topology API feature set
//
// By implementing this interface the topology service will automatically setup all services which implement it.
// Automatically ingesting Resource objects via the `GetTopologyObjectChannel()` function.
// This enables users to make use of the Topology APIs with these new Topology Resources.
//
type CacheableTopology interface {
	CacheEnabled() bool
	StartTopologyCaching(ctx context.Context) (<-chan *topologyv1.UpdateCacheRequest, error)
}

func New(cfg *any.Any, logger *zap.Logger, scope tally.Scope) (service.Service, error) {
	topologyConfig := &topologyv1cfg.Config{}
	err := ptypes.UnmarshalAny(cfg, topologyConfig)
	if err != nil {
		return nil, err
	}

	p, ok := service.Registry[pgservice.Name]
	if !ok {
		return nil, errors.New("Please config the datastore [clutch.service.db.postgres] to use the topology service")
	}

	dbClient, ok := p.(pgservice.Client)
	if !ok {
		return nil, errors.New("Unable to get the datastore client")
	}

	c := &client{
		config: topologyConfig,
		db:     dbClient.DB(),
		log:    logger,
		scope:  scope,
	}

	if topologyConfig.Cache == nil {
		c.log.Info("Topology caching is disabled")
		return c, nil
	}

	ctx, ctxCancelFunc := context.WithCancel(context.Background())
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sigc
		c.log.Info("Caught shutdown signal, shutting down topology caching and releasing advisory lock")
		ctxCancelFunc()
	}()

	c.log.Info("Topology caching is enabled")
	go c.acquireTopologyCacheLock(ctx)

	return c, nil
}

func (c *client) GetTopology(ctx context.Context) error {
	return nil
}

func (c *client) Search(ctx context.Context, req *topologyv1.SearchRequest) ([]*topologyv1.Resource, uint64, error) {
	query, nextPageToken, err := paginatedQueryBuilder(
		req.Filter,
		req.Sort,
		req.PageToken,
		req.Limit,
	)
	if err != nil {
		return nil, 0, err
	}

	rows, err := query.RunWith(c.db).Query()
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	results := []*topologyv1.Resource{}
	for rows.Next() {
		var id string
		var data []byte
		var metadata []byte

		if err := rows.Scan(&id, &data, &metadata); err != nil {
			c.log.Error("Error scanning row", zap.Error(err))
			return nil, 0, err
		}

		var dataAny any.Any
		if err := protojson.Unmarshal(data, &dataAny); err != nil {
			c.log.Error("Error unmarshaling data field", zap.Error(err))
			return nil, 0, err
		}

		var metadataMap map[string]string
		if err := json.Unmarshal(metadata, &metadataMap); err != nil {
			c.log.Error("Error unmarshaling metadata", zap.Error(err))
			return nil, 0, err
		}

		results = append(results, &topologyv1.Resource{
			Id:       id,
			Pb:       &dataAny,
			Metadata: metadataMap,
		})
	}

	// Rows.Err will report the last error encountered by Rows.Scan.
	if err := rows.Err(); err != nil {
		c.log.Error("Error processing rows for topology search query", zap.Error(err))
		return nil, 0, err
	}

	return results, nextPageToken, nil
}
