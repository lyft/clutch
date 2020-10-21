package sql

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/uber-go/tally"
	"go.uber.org/zap"
	rpcstatus "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/protobuf/encoding/protojson"

	apiv1 "github.com/lyft/clutch/backend/api/api/v1"
	auditv1 "github.com/lyft/clutch/backend/api/audit/v1"
	auditconfigv1 "github.com/lyft/clutch/backend/api/config/service/audit/v1"
	"github.com/lyft/clutch/backend/gateway/log"
	"github.com/lyft/clutch/backend/service"
	"github.com/lyft/clutch/backend/service/audit/storage"
	"github.com/lyft/clutch/backend/service/db/postgres"
)

type client struct {
	logger *zap.Logger
	scope  tally.Scope

	db *sql.DB
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
		WHERE sent = FALSE
		RETURNING id, occurred_at, details
	`

	return c.query(ctx, unsentEventsQuery)
}

func (c *client) ReadEvents(ctx context.Context, start time.Time, end *time.Time) ([]*auditv1.Event, error) {
	const readEventsRangeStatement = `
		SELECT id, occurred_at, details FROM audit_events
		WHERE occurred_at BETWEEN $1::timestamp AND $2::timestamp
		ORDER BY id
	`

	if end == nil {
		return c.query(ctx, readEventsRangeStatement, start, time.Now())
	}
	return c.query(ctx, readEventsRangeStatement, start, *end)
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

		occurred, err := ptypes.TimestampProto(row.OccurredAt)
		if err != nil {
			c.logger.Error("error in parsing db result's timestamp", zap.Error(err))
			return nil, err
		}

		proto := &auditv1.Event{
			OccurredAt: occurred,
			EventType: &auditv1.Event_Event{
				Event: requestEventProto(c.logger, row),
			},
		}
		events = append(events, proto)
	}

	return events, nil
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
		logger.Warn("unmarshallable object for RequestBody", zap.Error(err))
	}

	respBody, err := apiBodyProto(e.Details.ResponseBody)
	if err != nil {
		logger.Warn("unmarshallable object for ResponseBody", zap.Error(err))
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
func convertAPIBody(body *any.Any) (json.RawMessage, error) {
	// possible for Any to have a nil value
	if reflect.ValueOf(body).IsNil() {
		return nil, nil
	}

	b, err := protojson.Marshal(body)
	if err != nil {
		return nil, err
	}
	return json.RawMessage(b), nil
}

// Decodes JSON to proto Any message
func apiBodyProto(details json.RawMessage) (*any.Any, error) {
	// possible for json.RawMessage to be nil
	if details == nil {
		return nil, nil
	}

	body := &any.Any{}

	err := protojson.Unmarshal(details, body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
