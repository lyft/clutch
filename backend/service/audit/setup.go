package audit

// <!-- START clutchdoc -->
// description: Persists events emitted from the audit middleware. Other components can also use this service directly for advanced auditing capabilities.
// <!-- END clutchdoc -->

import (
	"context"
	"database/sql"
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
	"github.com/lyft/clutch/backend/service/auditsink"
	"github.com/lyft/clutch/backend/service/db/postgres"
)

const Name = "clutch.service.audit"

func New(cfg *any.Any, logger *zap.Logger, scope tally.Scope) (service.Service, error) {
	config := &auditconfigv1.Config{}
	if err := ptypes.UnmarshalAny(cfg, config); err != nil {
		return nil, err
	}

	db, ok := service.Registry[config.DbProvider]
	if !ok {
		return nil, fmt.Errorf("no database registered for saving audit events")
	}

	// TODO(maybe): Expand to more DB providers including non-SQL.
	sqlDB, ok := db.(postgres.Client)
	if !ok {
		return nil, fmt.Errorf("database in registry does not implement required interface")
	}

	// Ensure we have a non-nil set of Filters.
	filter := config.Filter
	if filter == nil {
		filter = &auditconfigv1.Filter{}
	}

	c := &client{
		logger: logger,
		scope:  scope,

		filter: filter,
		db:     sqlDB.DB(),
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

type client struct {
	logger *zap.Logger
	scope  tally.Scope

	filter    *auditconfigv1.Filter
	db        *sql.DB
	marshaler *jsonpb.Marshaler

	sinks []auditsink.Sink
}

// This should be called via `go` in order to avoid blocking main exectuion.
func (c *client) poll(interval time.Duration) {
	readAndFanout := func() {
		ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(1*time.Second))
		defer cancel()

		// TODO(maybe): Backpressure on continued failure.
		events, err := c.UnsentEvents(ctx)
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

func (c *client) Filter(event *auditv1.Event) bool {
	return auditsink.Filter(c.filter, event)
}

func (c *client) filterRequest(req *auditv1.RequestEvent) bool {
	return c.Filter(
		&auditv1.Event{EventType: &auditv1.Event_Event{Event: req}},
	)
}
