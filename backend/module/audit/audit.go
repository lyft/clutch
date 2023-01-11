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

	var page int
	if req.PageToken != "" {
		p, err := strconv.Atoi(req.PageToken)
		if err != nil || p < 0 {
			return nil, fmt.Errorf("invalid page token: %s", req.PageToken)
		}
		page = p
	}

	limit := int(req.Limit)
	// if a request has no limit we return all found events.
	if req.Limit == 0 {
		// in the event there is no request limit set limit to be 1 for purposes
		// of determining page offset later.
		limit = 1
	}

	startIdx := page * limit
	// n.b. the user defined limit is increased by 1 to help determine if there is
	// a subsequent page of information that should be denoted in the response
	additionalEventsNumber := int64(limit + 1)
	options := &audit.ReadOptions{Offset: int64(startIdx)}
	// if request limit is specified pass that value into the read options call
	if req.Limit != 0 {
		options.Limit = additionalEventsNumber
	}
	events, err := m.client.ReadEvents(ctx, start, end, options)
	if err != nil {
		return nil, err
	}

	// This logic is only used for purposes of pagination and should only be triggered if:
	// 1. there are more than 0 events
	// 2. if the request has a limit specified
	// 3. if the number of events is equal to the request limit + 1 meaning there are
	//    additional events to request on a subsequent page.
	if len(events) > 0 && req.Limit != 0 && len(events) == int(additionalEventsNumber) {
		resp.NextPageToken = strconv.FormatInt(int64(page+1), 10)
		events = events[:len(events)-1]
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
