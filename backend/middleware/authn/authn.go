package authn

// <!-- START clutchdoc -->
// description: Validates authentication headers using the authn service and adds claims to the request context.
// <!-- END clutchdoc -->

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/golang/protobuf/ptypes/any"
	"github.com/uber-go/tally"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/lyft/clutch/backend/middleware"
	"github.com/lyft/clutch/backend/service"
	"github.com/lyft/clutch/backend/service/authn"
)

const Name = "clutch.middleware.authn"

// List of method patterns that should not be blocked by authn.
// TODO(maybe): convert this to an API annotation or make configurable on the middleware.
var allowlist = []string{
	"/clutch.authn.v1.AuthnAPI/*",
	"/clutch.healthcheck.v1.HealthcheckAPI/*",
}

func New(cfg *any.Any, logger *zap.Logger, scope tally.Scope) (middleware.Middleware, error) {
	svc, ok := service.Registry["clutch.service.authn"]
	if !ok {
		return nil, errors.New("unable to get authn service")
	}

	p, ok := svc.(authn.Provider)
	if !ok {
		return nil, errors.New("authn service was not the correct type")
	}

	m := &mid{
		provider: p,
	}

	return m, nil
}

type mid struct {
	provider authn.Provider
}

func (m *mid) UnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Check for auth.
		authenticatedCtx, authErr := m.authenticate(ctx)

		// Determine if it's on the allow list.
		checkRequired := true
		for _, allow := range allowlist {
			if middleware.MatchMethodOrResource(allow, info.FullMethod) {
				checkRequired = false
				break
			}
		}

		// Assert auth if required.
		if checkRequired {
			if authErr != nil {
				return nil, status.New(codes.Unauthenticated, authErr.Error()).Err()
			}
			return handler(authenticatedCtx, req)
		}

		// If auth not required, we still append claims for logging purposes or anonymously accessible APIs.
		if _, err := authn.ClaimsFromContext(authenticatedCtx); err != nil {
			// Anonymous claims if there weren't any authenticated claims.
			return handler(authn.ContextWithAnonymousClaims(ctx), req)
		}
		return handler(authenticatedCtx, req)
	}
}

// getCookieValue is the easiest way to parse a cookie string in a non-HTTP request context.
func getCookieValue(raw, key string) (string, error) {
	request := http.Request{Header: http.Header{}}
	request.Header.Add("Cookie", raw)
	c, err := request.Cookie(key)
	if err != nil {
		return "", err
	}
	return c.Value, nil
}

// getToken looks for the token in the authorization header or cookies.
func getToken(md metadata.MD) (string, error) {
	if tokens := md.Get("authorization"); len(tokens) > 0 {
		splitToken := strings.Split(tokens[0], "Token")
		if len(splitToken) != 2 {
			return "", errors.New("bad token format, expected Authorization: Token <token>")
		}
		return strings.TrimSpace(splitToken[1]), nil
	}

	v := md.Get("grpcgateway-cookie")
	if len(v) == 0 {
		return "", errors.New("token not present in authorization header or cookies")
	}

	return getCookieValue(v[0], "token")
}

func (m *mid) authenticate(ctx context.Context) (context.Context, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("no headers present on request")
	}

	token, err := getToken(md)
	if err != nil {
		return nil, err
	}

	// Verify token.
	claims, err := m.provider.Verify(ctx, token)
	if err != nil {
		return nil, err
	}

	// Append claims information to context.
	return authn.ContextWithClaims(ctx, claims), nil
}
