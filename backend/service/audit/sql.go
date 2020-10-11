package audit

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"go.uber.org/zap"
	rpcstatus "google.golang.org/genproto/googleapis/rpc/status"

	apiv1 "github.com/lyft/clutch/backend/api/api/v1"
	auditv1 "github.com/lyft/clutch/backend/api/audit/v1"
)

func (c *client) WriteRequestEvent(ctx context.Context, event *auditv1.RequestEvent) (int64, error) {
	if event == nil {
		return -1, errors.New("cannot write empty event to table")
	}

	if !c.filterRequest(event) {
		return -1, ErrFailedFilters
	}

	dbEvent := &eventDetails{
		Username:         event.Username,
		Service:          event.ServiceName,
		Method:           event.MethodName,
		ActionType:       event.Type.String(),
		RequestResources: convertResources(event.Resources),
		RequestMetadata:  event.RequestMetadata,
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
	dbEvent := &eventDetails{
		Status: status{
			Code:    int(update.Status.Code),
			Message: update.Status.Message,
		},
		ResponseResources: convertResources(update.Resources),
		ResponseMetadata:  update.ResponseMetadata,
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
			zap.Any("event", update),
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
				Event: row.RequestEventProto(),
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
	Username          string      `json:"user_name,omitempty"`
	Service           string      `json:"service_name,omitempty"`
	Method            string      `json:"method_name,omitempty"`
	ActionType        string      `json:"type,omitempty"`
	Status            status      `json:"status,omitempty"`
	RequestResources  []*resource `json:"request_resources,omitempty"`
	ResponseResources []*resource `json:"response_resources,omitempty"`
	RequestMetadata   *any.Any    `json:"request_metadata,omitempty"`
	ResponseMetadata  *any.Any    `json:"response_metadata,omitempty"`
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

func (e *event) RequestEventProto() *auditv1.RequestEvent {
	return &auditv1.RequestEvent{
		Username:         e.Details.Username,
		ServiceName:      e.Details.Service,
		MethodName:       e.Details.Method,
		Type:             apiv1.ActionType(apiv1.ActionType_value[e.Details.ActionType]),
		Status:           e.Details.Status.Status(),
		Resources:        e.Details.ResourcesProto(),
		RequestMetadata:  e.Details.RequestMetadata,
		ResponseMetadata: e.Details.ResponseMetadata,
	}
}

func convertResources(proto []*auditv1.Resource) []*resource {
	resources := make([]*resource, 0, len(proto))
	for _, p := range proto {
		resources = append(resources, &resource{Id: p.Id, TypeUrl: p.TypeUrl})
	}
	return resources
}
