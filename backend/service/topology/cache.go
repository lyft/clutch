package topology

import (
	"context"
	"database/sql"
	"hash/fnv"
	"time"

	"go.uber.org/zap"
)

const topologyCacheLockId = "topology:cache"

// Performs leader election by acquiring a postgres advisory lock.
// Once the lock is acquired topology caching is started.
func (c *client) acquireTopologyCacheLock(ctx context.Context) {
	// We create our own connection to use for acquiring the advisory lock
	// If the connection is severed for any reason the advisory lock will automatically unlock
	conn, err := c.db.Conn(ctx)
	if err != nil {
		c.log.Fatal("Unable to connect to the database", zap.Error(err))
	}
	defer conn.Close()

	ticker := time.NewTicker(time.Second * 10)
	// Infinitely try to acquire the advisory lock
	// Once the lock is acquired we start caching, this is a blocking operation
	for ; true; <-ticker.C {
		c.log.Debug("trying to acquire advisory lock")

		advisoryLockId, err := convertLockIdToAdvisoryLockId(topologyCacheLockId)
		if err != nil {
			c.log.Error("Unable to convert the advisory lock id", zap.Error(err))
			return
		}

		// TODO: We could in the future spread the load of the topology caching
		// across many clutch instances by having an a lock per service (e.g. AWS, k8s, etc)

		// Advisory locks only take a arbitrary bigint value
		// In this case 100 is the lock for topology caching
		if c.tryAdvisoryLock(ctx, conn, advisoryLockId) {
			c.log.Debug("acquired the advisory lock, starting to cache topology now...")
			c.startTopologyCache()
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

func (c *client) startTopologyCache() {
	c.log.Debug("implement me")
}

func convertLockIdToAdvisoryLockId(lockId string) (uint32, error) {
	h := fnv.New32a()
	_, err := h.Write([]byte(lockId))
	if err != nil {
		return 0, err
	}
	return h.Sum32(), nil
}
