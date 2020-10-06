package audit

// <!-- START clutchdoc -->
// description: Persists events emitted from the audit middleware. Other components can also use this service directly for advanced auditing capabilities.
// <!-- END clutchdoc -->

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/uber-go/tally"
	"go.uber.org/zap"

	auditv1 "github.com/lyft/clutch/backend/api/audit/v1"
	auditconfigv1 "github.com/lyft/clutch/backend/api/config/service/audit/v1"
	"github.com/lyft/clutch/backend/service"
	"github.com/lyft/clutch/backend/service/audit/storage"
	"github.com/lyft/clutch/backend/service/audit/storage/local"
	"github.com/lyft/clutch/backend/service/audit/storage/sql"
	"github.com/lyft/clutch/backend/service/auditsink"
)

const Name = "clutch.service.audit"

func New(cfg *any.Any, logger *zap.Logger, scope tally.Scope) (service.Service, error) {
	config := &auditconfigv1.Config{}
	if err := ptypes.UnmarshalAny(cfg, config); err != nil {
		return nil, err
	}

	storage, err := getStorageImpl(config, logger, scope)
	if err != nil {
		return nil, err
	}

	c := &client{
		logger: logger,
		scope:  scope,

		storage: storage,
		marshaler: &jsonpb.Marshaler{
			// Use the names from the .proto.
			OrigName: true,
			// Render zero values (useful for successful status).
			EmitDefaults: true,
		},
	}

	for _, sinkName := range config.Sinks {
		sinkService, ok := service.Registry[sinkName]
		if !ok {
			return nil, fmt.Errorf(
				"listed sink '%s' is unregistered",
				sinkName,
			)
		}

		sink, ok := sinkService.(auditsink.Sink)
		if !ok {
			return nil, fmt.Errorf(
				"listed sink '%s' does not implement required interface",
				sinkName,
			)
		}

		c.sinks = append(c.sinks, sink)
	}

	// Start polling loop against database.
	go c.poll(time.Second * 10)

	return c, nil
}

func getStorageImpl(cfg *auditconfigv1.Config, logger *zap.Logger, scope tally.Scope) (storage.Storage, error) {
	switch cfg.GetService().(type) {
	case *auditconfigv1.Config_DbProvider:
		return sql.New(cfg, logger, scope)
	case *auditconfigv1.Config_LocalAuditing:
		return local.New(cfg, logger, scope)
	}

	return nil, errors.New("reached end of non-exhaustive type switch")
}

type client struct {
	logger *zap.Logger
	scope  tally.Scope

	storage storage.Storage
	filter  *auditconfigv1.Filter

	marshaler *jsonpb.Marshaler
	sinks     []auditsink.Sink
}

func (c *client) Filter(event *auditv1.Event) bool {
	return auditsink.Filter(c.filter, event)
}

func (c *client) filterRequest(req *auditv1.RequestEvent) bool {
	return c.Filter(
		&auditv1.Event{EventType: &auditv1.Event_Event{Event: req}},
	)
}

func (c *client) WriteRequestEvent(ctx context.Context, req *auditv1.RequestEvent) (int64, error) {
	if !c.filterRequest(req) {
		return -1, ErrFailedFilters
	}

	return c.storage.WriteRequestEvent(ctx, req)
}

func (c *client) UpdateRequestEvent(ctx context.Context, id int64, update *auditv1.RequestEvent) error {
	return c.storage.UpdateRequestEvent(ctx, id, update)
}

func (c *client) ReadEvents(ctx context.Context, start time.Time, end *time.Time) ([]*auditv1.Event, error) {
	return c.storage.ReadEvents(ctx, start, end)
}

// This should be called via `go` in order to avoid blocking main exectuion.
func (c *client) poll(interval time.Duration) {
	readAndFanout := func() {
		ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(1*time.Second))
		defer cancel()

		// TODO(maybe): Backpressure on continued failure.
		events, err := c.storage.UnsentEvents(ctx)
		if err != nil {
			return
		}

		for _, event := range events {
			for _, s := range c.sinks {
				if err := s.Write(event); err != nil {
					c.logger.Error(
						"error writing audit event to sink",
						zap.Error(err),
					)
				}
			}
		}
	}

	ticker := time.NewTicker(interval)
	for {
		<-ticker.C
		readAndFanout()
	}
}
