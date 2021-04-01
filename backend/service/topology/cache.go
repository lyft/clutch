package topology

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
	"google.golang.org/protobuf/encoding/protojson"

	topologyv1 "github.com/lyft/clutch/backend/api/topology/v1"
	"github.com/lyft/clutch/backend/service"
)

const topologyCacheLockId = "topology:cache"

// Performs leader election by acquiring a postgres advisory lock.
// Once the lock is acquired topology caching is started.
func (c *client) acquireTopologyCacheLock(ctx context.Context) {
	advisoryLockId := convertLockIdToAdvisoryLockId(topologyCacheLockId)
	ticker := time.NewTicker(time.Second * 10)

	// Infinitely try to acquire the advisory lock
	// Once the lock is acquired we start caching, this is a blocking operation
	for {
		c.log.Info("trying to acquire advisory lock")

		// We create our own connection to use for acquiring the advisory lock
		// If the connection is severed for any reason the advisory lock will automatically unlock
		conn, err := c.db.Conn(ctx)
		switch err {
		case nil:
			// TODO: We could in the future spread the load of the topology caching
			// across many clutch instances by having an a lock per service (e.g. AWS, k8s, etc)
			if c.tryAdvisoryLock(ctx, conn, advisoryLockId) {
				c.log.Info("acquired the advisory lock, starting to cache topology now...")
				go c.expireCache(ctx)
				c.startTopologyCache(ctx)
			}
		default:
			c.log.Error("lost connection to database, trying to reconnect...", zap.Error(err))
		}
		// Closes the db connection which will also release the advisory lock
		conn.Close()

		select {
		case <-ticker.C:
			continue
		case <-ctx.Done():
			ticker.Stop()
			return
		}
	}
}

// TODO: The advisory locking logic can be decomposed into its own service (e.g. "global locking service").
// Which can be used generically for anything that needs distributed locking functionality.
func (c *client) tryAdvisoryLock(ctx context.Context, conn *sql.Conn, lockId uint32) bool {
	var lock bool

	// Notably we do not use `pg_advisory_lock` as the behavior of this function stack locks.
	// For each lock invocation requires the same number of unlocks to release the advisory lock.
	// https://www.postgresql.org/docs/12/functions-admin.html#FUNCTIONS-ADVISORY-LOCKS-TABL
	if err := conn.QueryRowContext(ctx, "SELECT pg_try_advisory_lock($1);", lockId).Scan(&lock); err != nil {
		c.log.Error("Unable to query for a advisory lock", zap.Error(err))
	}

	return lock
}

func convertLockIdToAdvisoryLockId(lockID string) uint32 {
	x := sha256.New().Sum([]byte(lockID))
	return binary.BigEndian.Uint32(x)
}

// This will check all services that are currently registered for the given clutch configuration
// If any of the services implement the CacheableTopology interface we will start consuming
// topology objects until the context has been cancelled.
func (c *client) startTopologyCache(ctx context.Context) {
	for name, s := range service.Registry {
		if svc, ok := s.(CacheableTopology); ok {
			if svc.CacheEnabled() {
				c.log.Info("Processing Topology Objects for service", zap.String("service", name))
				topologyChannel, err := svc.StartTopologyCaching(ctx, c.cacheTTL)
				if err != nil {
					c.log.Error("Unable to start topology caching", zap.String("service", name), zap.Error(err))
					continue
				}
				go c.processTopologyObjectChannel(ctx, topologyChannel, name)
			}
		}
	}

	<-ctx.Done()
}

// Processes all topology cache update requests.
// This supports insert batching which defaults to 1, and can be configured by settting BatchInsertSize.
// Notably a flush will happen every 30 seconds, ensuring all items make it into cache even if the
// configured BatchInsertSize is not met.
func (c *client) processTopologyObjectChannel(ctx context.Context, objs <-chan *topologyv1.UpdateCacheRequest, service string) {
	var batchInsert []*topologyv1.Resource
	queueDepth := c.scope.Tagged(map[string]string{
		"service": service,
	}).SubScope("cache").Gauge("object_channel.queue_depth")

	for {
		select {
		case obj, ok := <-objs:
			if !ok {
				return
			}

			queueDepth.Update(float64(len(objs)))

			switch obj.Action {
			case topologyv1.UpdateCacheRequest_CREATE_OR_UPDATE:
				batchInsert = append(batchInsert, obj.Resource)
			case topologyv1.UpdateCacheRequest_DELETE:
				if err := c.deleteCache(ctx, obj.Resource.Id, obj.Resource.Pb.TypeUrl); err != nil {
					c.log.Error("Error deleting cache", zap.Error(err))
				}
			default:
				c.log.Warn("UpdateCacheRequest action is not implemented", zap.String("action", obj.Action.String()))
			}

			if len(batchInsert) >= c.batchInsertSize {
				if err := c.setCache(ctx, batchInsert); err != nil {
					c.log.Error("Error setting cache", zap.Error(err))
				}
				batchInsert = batchInsert[:0]
			}
		case <-time.After(c.batchInsertFlush):
			if len(batchInsert) > 0 {
				if err := c.setCache(ctx, batchInsert); err != nil {
					c.log.Error("Error setting cache", zap.Error(err))
				}
				batchInsert = batchInsert[:0]
			}
		}
	}
}

func (c *client) prepareBulkCacheInsert(obj []*topologyv1.Resource) ([]interface{}, string) {
	queryParams := make([]string, 0, len(obj))
	// Total object length x 4 as 4 is the number of columns we need to populate.
	queryArgs := make([]interface{}, 0, len(obj)*4)

	for i, o := range obj {
		// To parameterize all of the inputs for the bulk insert we need to increment the variable name.
		// EG: A bulk insert of two items should look like so "($1, $2, $3, $4),($5, $6, $7, $8)"
		queryParams = append(queryParams, fmt.Sprintf("($%d, $%d, $%d, $%d)", i*4+1, i*4+2, i*4+3, i*4+4))

		metadataJson, err := json.Marshal(o.Metadata)
		if err != nil {
			c.scope.SubScope("cache").Counter("set.failure").Inc(1)
			return nil, ""
		}

		dataJson, err := protojson.Marshal(o.Pb)
		if err != nil {
			c.scope.SubScope("cache").Counter("set.failure").Inc(1)
			return nil, ""
		}

		// append args in column order
		queryArgs = append(queryArgs, o.Id, o.Pb.GetTypeUrl(), dataJson, metadataJson)
	}

	return queryArgs, strings.Join(queryParams, ",")
}

func (c *client) setCache(ctx context.Context, obj []*topologyv1.Resource) error {
	args, queryParams := c.prepareBulkCacheInsert(obj)
	// upsertQuery := "INSERT INTO topology_cache"
	upsertQuery := fmt.Sprintf(`
		INSERT INTO topology_cache (id, resolver_type_url, data, metadata)
		VALUES %s
		ON CONFLICT (id, resolver_type_url) DO UPDATE SET
			resolver_type_url = EXCLUDED.resolver_type_url,
			data = EXCLUDED.data,
			metadata = EXCLUDED.metadata,
			updated_at = NOW()
	`, queryParams)

	_, err := c.db.ExecContext(
		ctx,
		upsertQuery,
		args...,
	)
	if err != nil {
		c.scope.SubScope("cache").Counter("set.failure").Inc(int64(len(obj)))
		return err
	}

	c.scope.SubScope("cache").Counter("set.success").Inc(int64(len(obj)))
	return nil
}

func (c *client) deleteCache(ctx context.Context, id, type_url string) error {
	const deleteQuery = `
		DELETE FROM topology_cache WHERE id = $1 and resolver_type_url = $2
	`

	_, err := c.db.ExecContext(ctx, deleteQuery, id, type_url)
	if err != nil {
		c.scope.SubScope("cache").Counter("delete.failure").Inc(1)
		return err
	}

	c.scope.SubScope("cache").Counter("delete.success").Inc(1)
	return nil
}

func (c *client) expireCache(ctx context.Context) {
	const expireQuery = `
		DELETE FROM topology_cache WHERE updated_at <= NOW() - CAST ( $1 AS INTERVAL );
	`

	ticker := time.NewTicker(time.Minute * 20)
	for {
		result, err := c.db.ExecContext(ctx, expireQuery, strconv.Itoa(int(c.cacheTTL.Seconds())))
		if err != nil {
			c.scope.SubScope("cache").Counter("expire.failure").Inc(1)
			c.log.Error("unable to expire cache", zap.Error(err))
			continue
		}

		numOfItemsRemoved, err := result.RowsAffected()
		if err != nil {
			c.scope.SubScope("cache").Counter("expire.failure").Inc(1)
			c.log.Error("unable to get rows removed from cache expiry query", zap.Error(err))
		} else {
			c.scope.SubScope("cache").Counter("expire.success").Inc(1)
			c.log.Info("successfully removed expired cache", zap.Int64("count", numOfItemsRemoved))
		}

		select {
		case <-ticker.C:
			continue
		case <-ctx.Done():
			ticker.Stop()
			return
		}
	}
}
