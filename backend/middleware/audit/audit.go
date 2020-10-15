package audit

// <!-- START clutchdoc -->
// description: Emits audit events from requests and responses based on annotations.
// <!-- END clutchdoc -->

import (
	"context"
	"errors"
	"fmt"

	"github.com/golang/protobuf/descriptor"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/uber-go/tally"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	auditv1 "github.com/lyft/clutch/backend/api/audit/v1"
	"github.com/lyft/clutch/backend/gateway/log"
	"github.com/lyft/clutch/backend/gateway/meta"
	"github.com/lyft/clutch/backend/middleware"
	"github.com/lyft/clutch/backend/service"
	auditservice "github.com/lyft/clutch/backend/service/audit"
	"github.com/lyft/clutch/backend/service/authn"
)

const Name = "clutch.middleware.audit"

func New(_ *any.Any, logger *zap.Logger, scope tally.Scope) (middleware.Middleware, error) {
	svc, ok := service.Registry[auditservice.Name]
	if !ok {
		return nil, fmt.Errorf("no audit svc with path '%s' registered for middleware", auditservice.Name)
	}

	auditService, ok := svc.(auditservice.Auditor)
	if !ok {
		return nil, errors.New("service in registry does not implement required interface")
	}

	return &mid{logger: logger, scope: scope, audit: auditService}, nil
}

type mid struct {
	logger *zap.Logger
	scope  tally.Scope
	audit  auditservice.Auditor
}

type auditEntryContextKey struct{}

func (m *mid) UnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		event, err := m.eventFromRequest(ctx, req, info)
		if err != nil {
			return nil, err
		}

		id, err := m.audit.WriteRequestEvent(ctx, event)
		if err != nil && !errors.Is(err, auditservice.ErrFailedFilters) {
			return nil, fmt.Errorf("could not make call %s because failed to audit: %w", info.FullMethod, err)
		}

		ctx = context.WithValue(ctx, auditEntryContextKey{}, id)
		resp, err := handler(ctx, req)

		if id != -1 {
			update, err := m.eventFromResponse(resp, err)
			if err != nil {
				return nil, err
			}

			if auditErr := m.audit.UpdateRequestEvent(ctx, id, update); auditErr != nil {
				m.logger.Warn("error updating audit event",
					zap.Int64("auditId", id),
					log.ProtoField("updateEvent", update),
				)
			}
		}

		return resp, err
	}
}

func (m *mid) eventFromRequest(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo) (*auditv1.RequestEvent, error) {
	svc, method, ok := middleware.SplitFullMethod(info.FullMethod)
	if !ok {
		m.logger.Warn("could not parse gRPC method", zap.String("fullMethod", info.FullMethod))
	}

	username := "UNKNOWN"
	if claims, err := authn.ClaimsFromContext(ctx); err == nil {
		username = claims.Subject
	}

	reqBody, err := meta.APIBody(req)
	if err != nil {
		return nil, err
	}

	return &auditv1.RequestEvent{
		Username:    username,
		ServiceName: svc,
		MethodName:  method,
		Type:        meta.GetAction(info.FullMethod),
		Resources:   meta.ResourceNames(req.(descriptor.Message)),
		RequestMetadata: &auditv1.RequestMetadata{
			Body: reqBody,
		},
	}, nil
}

func (m *mid) eventFromResponse(resp interface{}, err error) (*auditv1.RequestEvent, error) {
	s := status.Convert(err)
	if s == nil {
		s = status.New(codes.OK, "")
	}

	respBody, err := meta.APIBody(resp)
	if err != nil {
		return nil, err
	}

	return &auditv1.RequestEvent{
		Status:    s.Proto(),
		Resources: meta.ResourceNames(resp.(descriptor.Message)),
		ResponseMetadata: &auditv1.ResponseMetadata{
			Body: respBody,
		},
	}, nil
}
