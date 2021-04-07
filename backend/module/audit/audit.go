package audit

// <!-- START clutchdoc -->
// description: Returns audit event logs from the database.
// <!-- END clutchdoc -->

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang/protobuf/ptypes/any"
	"github.com/uber-go/tally"
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
	switch req.GetWindow().(type) {
	case *auditv1.GetEventsRequest_Range:
		timerange := req.GetRange()
		start := timerange.StartTime.AsTime()
		end := timerange.EndTime.AsTime()

		events, err := m.client.ReadEvents(ctx, start, &end)
		if err != nil {
			return nil, err
		}
		return &auditv1.GetEventsResponse{Events: events}, nil
	case *auditv1.GetEventsRequest_Since:
		if err := req.GetSince().CheckValid(); err != nil {
			return nil, fmt.Errorf("problem parsing duration: %w", err)
		}

		window := req.GetSince().AsDuration()
		start := time.Now().Add(-window)

		events, err := m.client.ReadEvents(ctx, start, nil)
		if err != nil {
			return nil, err
		}
		return &auditv1.GetEventsResponse{Events: events}, nil
	default:
		return nil, errors.New("no time window requested")
	}
}
