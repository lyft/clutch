package audit

// <!-- START clutchdoc -->
// description: Persists events emitted from the audit middleware. Other components can also use this service directly for advanced auditing capabilities.
// <!-- END clutchdoc -->

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/golang/protobuf/ptypes/any"
	"github.com/uber-go/tally"
	"go.uber.org/zap"
	"google.golang.org/protobuf/encoding/protojson"

	auditv1 "github.com/lyft/clutch/backend/api/audit/v1"
	auditconfigv1 "github.com/lyft/clutch/backend/api/config/service/audit/v1"
	"github.com/lyft/clutch/backend/gateway/log"
	"github.com/lyft/clutch/backend/service"
	"github.com/lyft/clutch/backend/service/audit/storage"
	"github.com/lyft/clutch/backend/service/audit/storage/local"
	"github.com/lyft/clutch/backend/service/audit/storage/sql"
	"github.com/lyft/clutch/backend/service/auditsink"
)

const (
	Name                 = "clutch.service.audit"
	auditEventSinkLockId = "audit:eventsink"
)

func New(cfg *any.Any, logger *zap.Logger, scope tally.Scope) (service.Service, error) {
	config := &auditconfigv1.Config{}
	if err := cfg.UnmarshalTo(config); err != nil {
		return nil, err
	}

	storageProvider, err := getStorageProvider(config, logger, scope)
	if err != nil {
		return nil, err
	}

	c := &client{
		logger: logger,
		scope:  scope,

		storage: storageProvider,
		marshaler: &protojson.MarshalOptions{
			// Use field names from the .proto rather than JSON camel case names.
			UseProtoNames: true,
			// Render zero values (useful for successful status).
			EmitUnpopulated: true,
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

		c.sinks = append(c.sinks, registeredSink{Sink: sink, name: sinkName})
	}

	// Start polling loop against database.
	go c.poll()

	return c, nil
}

func getStorageProvider(cfg *auditconfigv1.Config, logger *zap.Logger, scope tally.Scope) (storage.Storage, error) {
	switch cfg.GetStorageProvider().(type) {
	case *auditconfigv1.Config_DbProvider:
		return sql.New(cfg, logger, scope)
	case *auditconfigv1.Config_InMemory:
		return local.New(cfg, logger, scope)
	}

	return nil, errors.New("reached end of non-exhaustive type switch looking for audit storage")
}

type registeredSink struct {
	auditsink.Sink
	name string
}

type client struct {
	logger *zap.Logger
	scope  tally.Scope

	storage storage.Storage
	filter  *auditconfigv1.Filter

	marshaler *protojson.MarshalOptions
	sinks     []registeredSink
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
func (c *client) poll() {
	readAndFanout := func(ctx context.Context) {
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
						zap.String("sink", s.name),
						log.ProtoField("event", event),
						zap.Error(err),
					)
				}
			}
		}
	}

	lockID := convertLockToUint32(auditEventSinkLockId)

	// Generate random int with the floor of 10, to jitter this ticker.
	minTime := int64(10)
	ran, err := rand.Int(rand.Reader, big.NewInt(5))
	if err != nil {
		c.logger.Error("Unable to generate random int", zap.Error(err))
	}
	ticker := time.NewTicker(time.Second * time.Duration(ran.Int64()+minTime))
	for {
		<-ticker.C
		ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(30*time.Second))

		isLocked, err := c.storage.AttemptLock(ctx, lockID)
		if err != nil {
			c.logger.Error("Error while attempting to get lock", zap.Error(err))
		}

		if isLocked {
			_, err := c.storage.ReleaseLock(ctx, lockID)
			if err != nil {
				c.logger.Error("Error trying to release lock", zap.Error(err))
			}
			readAndFanout(ctx)
		}

		cancel()
	}
}

func convertLockToUint32(lockID string) uint32 {
	x := sha256.New().Sum([]byte(lockID))
	return binary.BigEndian.Uint32(x)
}
