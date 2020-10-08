package topology

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/binary"
	"encoding/json"
	"time"

	"go.uber.org/zap"
	"google.golang.org/protobuf/encoding/protojson"

	topologyv1 "github.com/lyft/clutch/backend/api/topology/v1"
	"github.com/lyft/clutch/backend/service"
)

const topologyCacheLockId = "topology:cache"

// Performs leader election by acquiring a postgres advisory lock.
// Once the lock is acquired topology caching is started.
func (c *client) acquireTopologyCacheLock() {
	ctx, ctxCancelFunc := context.WithCancel(context.Background())

	// We create our own connection to use for acquiring the advisory lock
	// If the connection is severed for any reason the advisory lock will automatically unlock
	conn, err := c.db.Conn(ctx)
	if err != nil {
		c.log.Fatal("Unable to connect to the database", zap.Error(err))
	}
	defer conn.Close()

	advisoryLockId := convertLockIdToAdvisoryLockId(topologyCacheLockId)

	ticker := time.NewTicker(time.Second * 10)
	// Infinitely try to acquire the advisory lock
	// Once the lock is acquired we start caching, this is a blocking operation
	for ; true; <-ticker.C {
		c.log.Debug("trying to acquire advisory lock")

		// TODO: We could in the future spread the load of the topology caching
		// across many clutch instances by having an a lock per service (e.g. AWS, k8s, etc)

		// Advisory locks only take a arbitrary bigint value
		// In this case 100 is the lock for topology caching
		if c.tryAdvisoryLock(ctx, conn, advisoryLockId) {
			c.log.Debug("acquired the advisory lock, starting to cache topology now...")
			c.startTopologyCache(ctx)
		} else {
			ctxCancelFunc()
		}
	}

	ctxCancelFunc()
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
//
func (c *client) startTopologyCache(ctx context.Context) {
	for n, s := range service.Registry {
		if svc, ok := s.(CacheableTopology); ok {
			c.log.Debug("Processing Topology Objects for service", zap.String("service", n))
			go c.processTopologyObjectChannel(ctx, svc.GetTopologyObjectChannel(ctx))
		}
	}

	<-ctx.Done()
}

func (c *client) processTopologyObjectChannel(ctx context.Context, objs chan topologyv1.UpdateCacheRequest) {
	for obj := range objs {
		switch obj.Action {
		case topologyv1.UpdateCacheRequest_CREATE_OR_UPDATE:
			if err := c.setCache(ctx, obj.Resource); err != nil {
				c.log.Error("Error setting cache", zap.Error(err))
			}
		case topologyv1.UpdateCacheRequest_DELETE:
			if err := c.deleteCache(ctx, obj.Resource.Id); err != nil {
				c.log.Error("Error deleting cache", zap.Error(err))
			}
		}
	}
}

func (c *client) setCache(ctx context.Context, obj *topologyv1.Resource) error {
	const upsertQuery = `
		INSERT INTO topology_cache (id, resolver_type_url, data, metadata)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO UPDATE SET
			resolver_type_url = EXCLUDED.resolver_type_url,
			data = EXCLUDED.data,
			metadata = EXCLUDED.metadata,
			updated_at = NOW()
	`

	metadataJson, err := json.Marshal(obj.Metadata)
	if err != nil {
		c.scope.SubScope("cache").Counter("set.failure").Inc(1)
		return err
	}

	dataJson, err := protojson.Marshal(obj.Pb)
	if err != nil {
		c.scope.SubScope("cache").Counter("set.failure").Inc(1)
		return err
	}

	_, err = c.db.ExecContext(
		ctx,
		upsertQuery,
		obj.Id,
		obj.Pb.GetTypeUrl(),
		dataJson,
		metadataJson,
	)
	if err != nil {
		c.scope.SubScope("cache").Counter("set.failure").Inc(1)
		return err
	}

	c.scope.SubScope("cache").Counter("set.success").Inc(1)
	return nil
}

func (c *client) deleteCache(ctx context.Context, id string) error {
	const deleteQuery = `
		DELETE FROM topology_cache WHERE id = $1
	`

	_, err := c.db.ExecContext(ctx, deleteQuery, id)
	if err != nil {
		c.scope.SubScope("cache").Counter("delete.failure").Inc(1)
		return err
	}

	c.scope.SubScope("cache").Counter("delete.success").Inc(1)
	return nil
}
