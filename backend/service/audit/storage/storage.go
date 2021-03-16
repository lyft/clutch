package storage

import (
	"context"
	"time"

	auditv1 "github.com/lyft/clutch/backend/api/audit/v1"
)

type Storage interface {
	// Used to get un-sent events.
	UnsentEvents(ctx context.Context) ([]*auditv1.Event, error)

	// Calls used by middleware to persist events during requests.
	WriteRequestEvent(ctx context.Context, req *auditv1.RequestEvent) (int64, error)
	UpdateRequestEvent(ctx context.Context, id int64, update *auditv1.RequestEvent) error

	// Used for services and modules to read past events within a timerange.
	// If end is nil, should search until the current time.
	ReadEvents(ctx context.Context, start time.Time, end *time.Time) ([]*auditv1.Event, error)
}
