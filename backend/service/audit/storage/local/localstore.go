package local

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/uber-go/tally"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"

	auditv1 "github.com/lyft/clutch/backend/api/audit/v1"
	auditconfigv1 "github.com/lyft/clutch/backend/api/config/service/audit/v1"
	"github.com/lyft/clutch/backend/service/audit/storage"
)

// TODO(maybe): There are more performant ways to do this that are more
// complicated to read and write (by programmers, not computers). This is not
// currently on a high-read path, so the map-with-mutex approach "should" be
// fine.
type client struct {
	sync.RWMutex

	logger *zap.Logger
	scope  tally.Scope

	events map[int64]*auditv1.Event
}

func New(cfg *auditconfigv1.Config, logger *zap.Logger, scope tally.Scope) (storage.Storage, error) {
	c := &client{
		logger: logger,
		scope:  scope,
		events: make(map[int64]*auditv1.Event),
	}

	return c, nil
}

func (c *client) UnsentEvents(ctx context.Context) ([]*auditv1.Event, error) {
	c.Lock()
	defer c.Unlock()

	unsent := make([]*auditv1.Event, 0, len(c.events))
	for _, v := range c.events {
		unsent = append(unsent, v)
	}
	c.events = make(map[int64]*auditv1.Event)

	return unsent, nil
}

func (c *client) WriteRequestEvent(ctx context.Context, req *auditv1.RequestEvent) (int64, error) {
	c.Lock()
	defer c.Unlock()

	i := int64(len(c.events))
	c.events[i] = &auditv1.Event{
		OccurredAt: timestamppb.Now(),
		EventType:  &auditv1.Event_Event{Event: req},
	}
	return i, nil
}

func (c *client) UpdateRequestEvent(ctx context.Context, id int64, update *auditv1.RequestEvent) error {
	c.RLock()
	defer c.RUnlock()

	event, ok := c.events[id]
	if !ok {
		return fmt.Errorf("cannot update event because cannot find by id: %d", id)
	}

	proto.Merge(proto.MessageV1(event.GetEvent()), proto.MessageV1(update))
	return nil
}

// Does a full scan through and copies those with a timestamp that fits the bill.
func (c *client) ReadEvents(ctx context.Context, start time.Time, end *time.Time) ([]*auditv1.Event, error) {
	c.RLock()
	defer c.RUnlock()

	var stop time.Time
	if end == nil {
		stop = time.Now()
	} else {
		stop = *end
	}

	events := make([]*auditv1.Event, 0, len(c.events))
	for _, value := range c.events {
		t := value.OccurredAt.AsTime()
		if start.Before(t) && stop.After(t) {
			events = append(events, value)
		}
	}

	return events, nil
}
