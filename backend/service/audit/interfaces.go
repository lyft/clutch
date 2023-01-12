package audit

import (
	"context"
	"errors"
	"time"

	auditv1 "github.com/lyft/clutch/backend/api/audit/v1"
)

var ErrFailedFilters = errors.New("event did not pass auditor's filters")

type ReadOptions struct {
	Offset int64
	Limit  int64
}

// Required functions to save/share events processed by Clutch.
type Auditor interface {
	// Check if an event passes the configured filters for this store. True if it should save the event.
	Filter(event *auditv1.Event) bool

	// Calls used by middleware to persist events during requests.
	WriteRequestEvent(ctx context.Context, req *auditv1.RequestEvent) (int64, error)
	UpdateRequestEvent(ctx context.Context, id int64, update *auditv1.RequestEvent) error

	// Used for services and modules to read past events within a timerange.
	// If end is nil, should search until the current time.
	ReadEvents(ctx context.Context, start time.Time, end *time.Time, options *ReadOptions) ([]*auditv1.Event, error)

	// Used for services and modules to read a specific event.
	ReadEvent(ctx context.Context, id int64) (*auditv1.Event, error)
}
