package topology

import (
	"database/sql"
	"encoding/json"
	"errors"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/uber-go/tally"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/structpb"

	topologyv1cfg "github.com/lyft/clutch/backend/api/config/service/topology/v1"
	topologyv1 "github.com/lyft/clutch/backend/api/topology/v1"
	"github.com/lyft/clutch/backend/service"
	pgservice "github.com/lyft/clutch/backend/service/db/postgres"
)

const Name = "clutch.service.topology"

type Service interface {
	GetTopology(ctx context.Context) error
	Search(ctx context.Context, search *topologyv1.SearchRequest) ([]*topologyv1.Resource, string, error)
	Autocomplete(ctx context.Context, typeURL, search string, limit uint64) ([]*topologyv1.Resource, error)
}

type client struct {
	config *topologyv1cfg.Config

	db    *sql.DB
	log   *zap.Logger
	scope tally.Scope

	cacheTTL        time.Duration
	bulkInsertMutex sync.RWMutex
}

// CacheableTopology is implemented by a service that wishes to enable the topology API feature set
//
// By implementing this interface the topology service will automatically setup all services which implement it.
// Automatically ingesting Resource objects via the `GetTopologyObjectChannel()` function.
// This enables users to make use of the Topology APIs with these new Topology Resources.
type CacheableTopology interface {
	CacheEnabled() bool

	// Notably the cache TTL is provided, this information can be used to ensure
	// all active resources are added to the cache more often than the cache TTL.
	StartTopologyCaching(ctx context.Context, ttl time.Duration) (<-chan *topologyv1.UpdateCacheRequest, error)
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
		config:          topologyConfig,
		db:              dbClient.DB(),
		log:             logger,
		scope:           scope,
		bulkInsertMutex: sync.RWMutex{},
	}

	if topologyConfig.Cache == nil {
		c.log.Info("Topology caching is disabled")
		return c, nil
	}

	// If TTL is not set default to two hours
	c.cacheTTL = time.Hour * 2
	if c.config.Cache.Ttl != nil {
		c.cacheTTL = c.config.Cache.Ttl.AsDuration()
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

func (c *client) Search(ctx context.Context, req *topologyv1.SearchRequest) ([]*topologyv1.Resource, string, error) {
	query, nextPageToken, err := paginatedQueryBuilder(
		req.Filter,
		req.Sort,
		req.PageToken,
		req.Limit,
	)
	if err != nil {
		return nil, "", err
	}

	rows, err := query.RunWith(c.db).Query()
	if err != nil {
		return nil, "", err
	}
	defer rows.Close()

	results := []*topologyv1.Resource{}
	for rows.Next() {
		var id string
		var data []byte
		var metadata []byte

		if err := rows.Scan(&id, &data, &metadata); err != nil {
			c.log.Error("Error scanning row", zap.Error(err))
			return nil, "", err
		}

		var dataAny any.Any
		if err := protojson.Unmarshal(data, &dataAny); err != nil {
			c.log.Error("Error unmarshaling data field", zap.Error(err))
			return nil, "", err
		}

		var metadataMap map[string]*structpb.Value
		if err := json.Unmarshal(metadata, &metadataMap); err != nil {
			c.log.Error("Error unmarshaling metadata", zap.Error(err))
			return nil, "", err
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
		return nil, "", err
	}

	// if their are no results then set the next page to nil
	nextPageTokenStr := ""
	if len(results) > 0 {
		nextPageTokenStr = strconv.FormatUint(nextPageToken, 10)
	}

	return results, nextPageTokenStr, nil
}

func (c *client) Autocomplete(ctx context.Context, typeURL, search string, limit uint64) ([]*topologyv1.Resource, error) {
	searchRequest := &topologyv1.SearchRequest{
		PageToken: "0",
		Limit:     limit,
		Sort: &topologyv1.SearchRequest_Sort{
			Direction: topologyv1.SearchRequest_Sort_ASCENDING,
			Field:     "column.id",
		},
		Filter: &topologyv1.SearchRequest_Filter{
			TypeUrl: typeURL,
			Search: &topologyv1.SearchRequest_Filter_Search{
				Field: "column.id",
				Text:  search,
			},
		},
	}

	results, _, err := c.Search(ctx, searchRequest)
	if err != nil {
		return nil, err
	}

	return results, err
}
