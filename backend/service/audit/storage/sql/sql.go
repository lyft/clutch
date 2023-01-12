package sql

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/uber-go/tally/v4"
	"go.uber.org/zap"
	rpcstatus "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	apiv1 "github.com/lyft/clutch/backend/api/api/v1"
	auditv1 "github.com/lyft/clutch/backend/api/audit/v1"
	auditconfigv1 "github.com/lyft/clutch/backend/api/config/service/audit/v1"
	"github.com/lyft/clutch/backend/gateway/log"
	"github.com/lyft/clutch/backend/service"
	"github.com/lyft/clutch/backend/service/audit/storage"
	"github.com/lyft/clutch/backend/service/db/postgres"
)

type client struct {
	logger           *zap.Logger
	scope            tally.Scope
	db               *sql.DB
	advisoryLockConn *sql.Conn
}

func New(cfg *auditconfigv1.Config, logger *zap.Logger, scope tally.Scope) (storage.Storage, error) {
	db, ok := service.Registry[cfg.GetDbProvider()]
	if !ok {
		return nil, fmt.Errorf("no database registered for saving audit events")
	}

	// TODO(maybe): Expand to more DB providers, including non-SQL.
	sqlDB, ok := db.(postgres.Client)
	if !ok {
		return nil, fmt.Errorf("database in registry does not implement required interface")
	}

	c := &client{
		logger: logger,
		scope:  scope,
		db:     sqlDB.DB(),
	}
	return c, nil
}

func (c *client) WriteRequestEvent(ctx context.Context, event *auditv1.RequestEvent) (int64, error) {
	if event == nil {
		return -1, errors.New("cannot write empty event to table")
	}

	reqBody, err := convertAPIBody(event.RequestMetadata.Body)
	if err != nil {
		return -1, err
	}

	dbEvent := &eventDetails{
		Username:         event.Username,
		Service:          event.ServiceName,
		Method:           event.MethodName,
		ActionType:       event.Type.String(),
		RequestResources: convertResources(event.Resources),
		RequestBody:      reqBody,
	}
	blob, err := json.Marshal(dbEvent)
	if err != nil {
		return -1, err
	}

	var id int64
	const writeEventStatement = `INSERT INTO audit_events (occurred_at, details) VALUES (NOW(), $1) RETURNING id`
	err = c.db.QueryRowContext(ctx, writeEventStatement, blob).Scan(&id)
	if err != nil {
		return -1, err
	}

	return id, nil
}

func (c *client) UpdateRequestEvent(ctx context.Context, id int64, update *auditv1.RequestEvent) error {
	respBody, err := convertAPIBody(update.ResponseMetadata.Body)
	if err != nil {
		return err
	}

	dbEvent := &eventDetails{
		Status: status{
			Code:    int(update.Status.Code),
			Message: update.Status.Message,
		},
		ResponseResources: convertResources(update.Resources),
		ResponseBody:      respBody,
	}
	blob, err := json.Marshal(dbEvent)
	if err != nil {
		return err
	}

	const updateEventStatement = `
		UPDATE audit_events
		SET details = details || $2::jsonb
		WHERE id = $1
    `
	if _, err = c.db.ExecContext(ctx, updateEventStatement, id, blob); err != nil {
		c.logger.Warn(
			"error updating audit row",
			zap.Int64("row_id", id),
			log.ProtoField("event", update),
			zap.Error(err),
		)
		return err
	}

	return nil
}

func (c *client) UnsentEvents(ctx context.Context) ([]*auditv1.Event, error) {
	const unsentEventsQuery = `
		UPDATE audit_events
		SET sent = TRUE
		WHERE id IN (SELECT id FROM audit_events WHERE sent = FALSE LIMIT 100)
		RETURNING id, occurred_at, details;
	`

	return c.query(ctx, unsentEventsQuery)
}

func (c *client) ReadEvents(ctx context.Context, start time.Time, end *time.Time, options *storage.ReadOptions) ([]*auditv1.Event, error) {
	readEventsRangeStatement := `
		SELECT id, occurred_at, details FROM audit_events
		WHERE occurred_at BETWEEN $1::timestamp AND $2::timestamp
		ORDER BY id
	`

	endTime := time.Now()
	if end != nil {
		endTime = *end
	}

	args := []interface{}{start, endTime}
	if options != nil {
		switch {
		case options.Limit != 0 && options.Offset != 0:
			readEventsRangeStatement = fmt.Sprintf(`%s LIMIT $3 OFFSET $4`, readEventsRangeStatement)
			args = append(args, options.Limit, options.Offset)
		case options.Limit != 0:
			readEventsRangeStatement = fmt.Sprintf(`%s LIMIT $3`, readEventsRangeStatement)
			args = append(args, options.Limit)
		case options.Offset != 0:
			readEventsRangeStatement = fmt.Sprintf(`%s OFFSET $4`, readEventsRangeStatement)
			args = append(args, options.Offset)
		}
	}

	return c.query(ctx, readEventsRangeStatement, args...)
}

func (c *client) ReadEvent(ctx context.Context, id int64) (*auditv1.Event, error) {
	const readEventsRangeStatement = `
		SELECT id, occurred_at, details FROM audit_events
		WHERE id = $1
	`

	events, err := c.query(ctx, readEventsRangeStatement, id)
	if err != nil {
		return nil, err
	}
	if len(events) == 0 {
		return nil, fmt.Errorf("cannot find event by id: %d", id)
	}
	return events[0], nil
}

func (c *client) query(ctx context.Context, query string, args ...interface{}) ([]*auditv1.Event, error) {
	rows, err := c.db.QueryContext(ctx, query, args...)
	if err != nil {
		c.logger.Error("error querying db", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var events []*auditv1.Event
	for rows.Next() {
		row := &event{
			Details: &eventDetails{},
		}
		var blob []byte
		if err := rows.Scan(&row.Id, &row.OccurredAt, &blob); err != nil {
			c.logger.Error("error scanning db results", zap.Error(err))
			return nil, err
		}

		if err := json.Unmarshal(blob, row.Details); err != nil {
			c.logger.Error("unmarshallable blob in db result", zap.Error(err))
			return nil, err
		}

		occurred := timestamppb.New(row.OccurredAt)
		if err := occurred.CheckValid(); err != nil {
			c.logger.Error("error in parsing db result's timestamp", zap.Error(err))
			return nil, err
		}

		proto := &auditv1.Event{
			// n.b. this is a safe casting since BIGSERIAL can only reach the max value for signed int64.
			Id:         int64(row.Id),
			OccurredAt: occurred,
			EventType: &auditv1.Event_Event{
				Event: requestEventProto(c.logger, row),
			},
		}
		events = append(events, proto)
	}

	return events, nil
}

// We create our own connection to use for acquiring the advisory lock.
// For advisory locks you must use the same session to issue the unlock
// If not the lock will stay locked until that connection is severed or
// the unlock request uses the same session in which the lock was created.
func (c *client) getAdvisoryConn() (*sql.Conn, error) {
	if c.advisoryLockConn != nil {
		return c.advisoryLockConn, nil
	}

	advisoryLockConn, err := c.db.Conn(context.Background())
	if err != nil {
		return nil, err
	}

	c.advisoryLockConn = advisoryLockConn
	return advisoryLockConn, nil
}

func (c *client) AttemptLock(ctx context.Context, lockID uint32) (bool, error) {
	conn, err := c.getAdvisoryConn()
	if err != nil {
		return false, err
	}

	var lock bool
	if err := conn.QueryRowContext(ctx, "SELECT pg_try_advisory_lock($1);", lockID).Scan(&lock); err != nil {
		c.logger.Error("Unable to query for a advisory lock", zap.Error(err))
		c.advisoryLockConn = nil
		return false, err
	}
	return lock, nil
}

func (c *client) ReleaseLock(ctx context.Context, lockID uint32) (bool, error) {
	conn, err := c.getAdvisoryConn()
	if err != nil {
		return false, err
	}

	var unlock bool
	if err := conn.QueryRowContext(ctx, "SELECT pg_advisory_unlock($1)", lockID).Scan(&unlock); err != nil {
		c.logger.Error("Unable to perform an advisory unlock", zap.Error(err))
		c.advisoryLockConn = nil
		return false, err
	}

	if conn != nil {
		conn.Close()
		c.advisoryLockConn = nil
	}

	return unlock, nil
}

type status struct {
	Code    int    `json:"code"`
	Message string `json:"message,omitempty"`
}

func (s *status) Status() *rpcstatus.Status {
	return &rpcstatus.Status{
		Code:    int32(s.Code),
		Message: s.Message,
	}
}

type resource struct {
	TypeUrl string `json:"type_url,omitempty"`
	Id      string `json:"id,omitempty"`
}

type eventDetails struct {
	Username          string          `json:"user_name,omitempty"`
	Service           string          `json:"service_name,omitempty"`
	Method            string          `json:"method_name,omitempty"`
	ActionType        string          `json:"type,omitempty"`
	Status            status          `json:"status,omitempty"`
	RequestResources  []*resource     `json:"request_resources,omitempty"`
	ResponseResources []*resource     `json:"response_resources,omitempty"`
	RequestBody       json.RawMessage `json:"request_body,omitempty"`
	ResponseBody      json.RawMessage `json:"response_body,omitempty"`
}

func (e *eventDetails) ResourcesProto() []*auditv1.Resource {
	resources := make([]*auditv1.Resource, 0, len(e.ResponseResources)+len(e.RequestResources))

	for _, resource := range e.RequestResources {
		proto := &auditv1.Resource{
			TypeUrl: resource.TypeUrl,
			Id:      resource.Id,
		}
		resources = append(resources, proto)
	}

	for _, resource := range e.ResponseResources {
		proto := &auditv1.Resource{
			TypeUrl: resource.TypeUrl,
			Id:      resource.Id,
		}
		resources = append(resources, proto)
	}

	return resources
}

type event struct {
	Id         uint64
	OccurredAt time.Time
	Details    *eventDetails
}

func requestEventProto(logger *zap.Logger, e *event) *auditv1.RequestEvent {
	reqBody, err := apiBodyProto(e.Details.RequestBody)
	if err != nil {
		logger.Error("unmarshallable object for RequestBody", zap.Error(err))
	}

	respBody, err := apiBodyProto(e.Details.ResponseBody)
	if err != nil {
		logger.Error("unmarshallable object for ResponseBody", zap.Error(err))
	}

	return &auditv1.RequestEvent{
		Username:    e.Details.Username,
		ServiceName: e.Details.Service,
		MethodName:  e.Details.Method,
		Type:        apiv1.ActionType(apiv1.ActionType_value[e.Details.ActionType]),
		Status:      e.Details.Status.Status(),
		Resources:   e.Details.ResourcesProto(),
		RequestMetadata: &auditv1.RequestMetadata{
			Body: reqBody,
		},
		ResponseMetadata: &auditv1.ResponseMetadata{
			Body: respBody,
		},
	}
}

func convertResources(proto []*auditv1.Resource) []*resource {
	resources := make([]*resource, 0, len(proto))
	for _, p := range proto {
		resources = append(resources, &resource{Id: p.Id, TypeUrl: p.TypeUrl})
	}
	return resources
}

// Encodes proto object in JSON format
func convertAPIBody(body *anypb.Any) (json.RawMessage, error) {
	// gracefully handles untyped and typed nils
	// returns empty object if input is untyped or type nil
	b, err := protojson.Marshal(body)
	if err != nil {
		return nil, err
	}
	return json.RawMessage(b), nil
}

// apiBodyProto reads the given []byte into the Any proto msg
func apiBodyProto(details json.RawMessage) (*anypb.Any, error) {
	// protojson.Marshal does not gracefully handle nil values
	if details == nil {
		return nil, nil
	}

	body := &anypb.Any{}

	err := protojson.Unmarshal(details, body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
