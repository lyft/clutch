package storage

import (
	"context"
	"time"

	auditv1 "github.com/lyft/clutch/backend/api/audit/v1"
)

type ReadOptions struct {
	Offset int64
	Limit  int64
}

type Storage interface {
	// Used to get un-sent events.
	UnsentEvents(ctx context.Context) ([]*auditv1.Event, error)

	// Calls used by middleware to persist events during requests.
	WriteRequestEvent(ctx context.Context, req *auditv1.RequestEvent) (int64, error)
	UpdateRequestEvent(ctx context.Context, id int64, update *auditv1.RequestEvent) error

	// Used for services and modules to read past events within a timerange.
	// If end is nil, should search until the current time.
	ReadEvents(ctx context.Context, start time.Time, end *time.Time, options *ReadOptions) ([]*auditv1.Event, error)

	ReadEvent(ctx context.Context, id int64) (*auditv1.Event, error)

	// Basic Locking functions to provide concurrency control.
	AttemptLock(ctx context.Context, lockID uint32) (bool, error)
	ReleaseLock(ctx context.Context, lockID uint32) (bool, error)
}
