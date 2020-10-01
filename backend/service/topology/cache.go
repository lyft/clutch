package topology

import (
	"context"
	"database/sql"
	"time"

	"go.uber.org/zap"
)

// Perfroms leader election by acquiring a postgres advisory lock.
// Once the lock is acquired topology caching is started.
func (c *client) acquireTopologyCacheLock() {
	// We create our own connection to use for acquiring the advisory lock
	// If the connection is severed for any reason the advisory lock will automatically unlock
	conn, err := c.db.Conn(context.Background())
	if err != nil {
		c.log.Fatal("Unable to connect to the database", zap.Error(err))
	}

	// Infinitely try to acquire the advisory lock
	// Once the lock is acquired we start caching, this is a blocking operation
	for {
		c.log.Debug("trying to acquire advisory lock")

		// Advisory locks only take a arbitrary bigint value
		// In this case 100 is the lock for topology caching
		if tryAdvisoryLock(conn, 100) {
			c.log.Debug("acquired the advisory lock, starting to cache topology now...")
			c.startTopologyCache()
		}

		time.Sleep(time.Second * 10)
	}
}

func tryAdvisoryLock(conn *sql.Conn, lockId int) bool {
	var lock bool

	// Notably we do not use `pg_advisory_lock` as the behavior of this function stack locks.
	// For each lock invocation requires the same number of unlocks to release the advisory lock.
	// https://www.postgresql.org/docs/12/functions-admin.html#FUNCTIONS-ADVISORY-LOCKS-TABL
	_ = conn.QueryRowContext(context.Background(), "SELECT pg_try_advisory_lock($1);", lockId).Scan(&lock)
	return lock
}

func (c *client) startTopologyCache() {
	c.log.Debug("implement me")
}
