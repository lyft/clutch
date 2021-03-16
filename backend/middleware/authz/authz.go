package authz

// <!-- START clutchdoc -->
// description: Extracts resource information from requests and passes it to the authz service to allow or deny access.
// <!-- END clutchdoc -->

import (
	"context"
	"errors"

	"github.com/uber-go/tally"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	authzv1 "github.com/lyft/clutch/backend/api/authz/v1"
	"github.com/lyft/clutch/backend/gateway/meta"
	"github.com/lyft/clutch/backend/middleware"
	"github.com/lyft/clutch/backend/service"
	"github.com/lyft/clutch/backend/service/authn"
	"github.com/lyft/clutch/backend/service/authz"
)

const Name = "clutch.middleware.authz"

// List of method patterns that should never go through the authz engine.
// TODO(maybe): convert this to an API annotation or make configurable on the authz middleware.
var allowlist = []string{
	"/clutch.authn.v1.AuthnAPI/*",
	"/clutch.healthcheck.v1.HealthcheckAPI/*",
}

func New(cfg *anypb.Any, logger *zap.Logger, scope tally.Scope) (middleware.Middleware, error) {
	svc, ok := service.Registry["clutch.service.authz"]
	if !ok {
		return nil, errors.New("unable to get authz service")
	}

	zclient, ok := svc.(authz.Client)
	if !ok {
		return nil, errors.New("authz service was not the correct type")
	}

	return &mid{client: zclient}, nil
}

type mid struct {
	client authz.Client
}

func (m *mid) evaluate(ctx context.Context, r *authzv1.CheckRequest) error {
	resp, err := m.client.Check(ctx, r)
	if err != nil {
		return err
	}
	if resp.Decision != authzv1.Decision_ALLOW {
		s := status.New(codes.PermissionDenied, "permission denied by authz")
		s, _ = s.WithDetails(r)
		return s.Err()
	}
	return nil
}

func (m *mid) UnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		// Never interfere with allowlisted flows.
		for _, allow := range allowlist {
			if middleware.MatchMethodOrResource(allow, info.FullMethod) {
				return handler(ctx, req)
			}
		}

		claims, err := authn.ClaimsFromContext(ctx)
		if err != nil {
			return nil, err
		}

		actionType := meta.GetAction(info.FullMethod)
		resources := meta.ResourceNames(req.(proto.Message))

		subject := &authzv1.Subject{
			User:   claims.Subject,
			Groups: claims.Groups,
		}

		if len(resources) == 0 {
			check := &authzv1.CheckRequest{
				Subject:    subject,
				Method:     info.FullMethod,
				ActionType: actionType,
			}
			if err := m.evaluate(ctx, check); err != nil {
				return nil, err
			}
		}

		for _, resource := range resources {
			check := &authzv1.CheckRequest{
				Subject:    subject,
				Method:     info.FullMethod,
				ActionType: actionType,
				Resource:   resource.Id,
			}

			if err := m.evaluate(ctx, check); err != nil {
				return nil, err
			}
		}

		return handler(ctx, req)
	}
}
