package audit

// <!-- START clutchdoc -->
// description: Returns audit event logs from the database.
// <!-- END clutchdoc -->

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang/protobuf/ptypes/any"
	"github.com/uber-go/tally/v4"
	"go.uber.org/zap"

	auditv1 "github.com/lyft/clutch/backend/api/audit/v1"
	"github.com/lyft/clutch/backend/module"
	"github.com/lyft/clutch/backend/service"
	"github.com/lyft/clutch/backend/service/audit"
)

const Name = "clutch.module.audit"

func New(*any.Any, *zap.Logger, tally.Scope) (module.Module, error) {
	auditClient, ok := service.Registry["clutch.service.audit"]
	if !ok {
		return nil, errors.New("could not find service")
	}

	c, ok := auditClient.(audit.Auditor)
	if !ok {
		return nil, errors.New("service was not the correct type")
	}

	mod := &mod{
		client: c,
	}

	return mod, nil
}

type mod struct {
	client audit.Auditor
}

func (m *mod) Register(r module.Registrar) error {
	auditv1.RegisterAuditAPIServer(r.GRPCServer(), m)
	return r.RegisterJSONGateway(auditv1.RegisterAuditAPIHandler)
}

func (m *mod) GetEvents(ctx context.Context, req *auditv1.GetEventsRequest) (*auditv1.GetEventsResponse, error) {
	resp := &auditv1.GetEventsResponse{}

	var start time.Time
	var end *time.Time

	switch req.GetWindow().(type) {
	case *auditv1.GetEventsRequest_Range:
		timerange := req.GetRange()
		start = timerange.StartTime.AsTime()
		endTime := timerange.EndTime.AsTime()
		end = &endTime
	case *auditv1.GetEventsRequest_Since:
		if err := req.GetSince().CheckValid(); err != nil {
			return nil, fmt.Errorf("problem parsing duration: %w", err)
		}

		window := req.GetSince().AsDuration()
		start = time.Now().Add(-window)
	default:
		return nil, errors.New("no time window requested")
	}

	eventCount, err := m.client.CountEvents(ctx, start, end)
	if err != nil {
		return nil, err
	}

	var options *audit.ReadOptions
	if req.PageToken != "" {
		startIdx := 0
		page, err := strconv.Atoi(req.PageToken)
		if err != nil {
			return nil, fmt.Errorf("invalid page token: %s", req.PageToken)
		}
		limit := int(req.Limit)
		if req.Limit == 0 {
			limit = 10
		}

		startIdx = page * limit
		if startIdx > int(eventCount) {
			startIdx = 0
		}

		endIdx := startIdx + limit
		if endIdx < int(eventCount) {
			resp.NextPageToken = strconv.FormatInt(int64(page+1), 10)
		}

		options = &audit.ReadOptions{Offset: int64(startIdx), Limit: int64(limit)}
	}

	events, err := m.client.ReadEvents(ctx, start, end, options)
	if err != nil {
		return nil, err
	}

	resp.Events = events

	return resp, nil
}

func (m *mod) GetEvent(ctx context.Context, req *auditv1.GetEventRequest) (*auditv1.GetEventResponse, error) {
	event, err := m.client.ReadEvent(ctx, req.EventId)
	if err != nil {
		return nil, err
	}

	resp := &auditv1.GetEventResponse{
		Event: event,
	}
	return resp, nil
}
