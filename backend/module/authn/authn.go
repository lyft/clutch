package authn

// <!-- START clutchdoc -->
// description: Registers login and callback endpoints for OAuth 2 flows.
// <!-- END clutchdoc -->

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/uber-go/tally"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	authnv1 "github.com/lyft/clutch/backend/api/authn/v1"
	"github.com/lyft/clutch/backend/module"
	"github.com/lyft/clutch/backend/service"
	"github.com/lyft/clutch/backend/service/authn"
)

const Name = "clutch.module.authn"

func New(*any.Any, *zap.Logger, tally.Scope) (module.Module, error) {
	svc, ok := service.Registry["clutch.service.authn"]
	if !ok {
		return nil, errors.New("unable to get authn service")
	}

	p, ok := svc.(authn.IssuerProvider)
	if !ok {
		return nil, errors.New("authn service was not the correct type (is not a IssuerProvider)")
	}

	return &mod{
		authnv1: &api{provider: p, issuer: p},
	}, nil
}

type mod struct {
	authnv1 authnv1.AuthnAPIServer
}

func (m *mod) Register(r module.Registrar) error {
	authnv1.RegisterAuthnAPIServer(r.GRPCServer(), m.authnv1)
	return r.RegisterJSONGateway(authnv1.RegisterAuthnAPIHandler)
}

type api struct {
	provider authn.Provider
	issuer   authn.Issuer
}

func (a *api) Login(ctx context.Context, request *authnv1.LoginRequest) (*authnv1.LoginResponse, error) {
	state, err := a.provider.GetStateNonce(request.RedirectUrl)
	if err != nil {
		return nil, err
	}
	authURL, err := a.provider.GetAuthCodeURL(ctx, state)
	if err != nil {
		return nil, err
	}

	if request.RedirectUrl == "" {
		md := metadata.Pairs("Location", authURL)
		if err := grpc.SendHeader(ctx, md); err != nil {
			return nil, err
		}
	}

	return &authnv1.LoginResponse{
		AuthUrl: authURL,
	}, nil
}

func (a *api) Callback(ctx context.Context, request *authnv1.CallbackRequest) (*authnv1.CallbackResponse, error) {
	if request.Error != "" {
		return nil, fmt.Errorf("%s: %s", request.Error, request.ErrorDescription)
	}

	redirectURL, err := a.provider.ValidateStateNonce(request.State)
	if err != nil {
		return nil, err
	}

	token, err := a.provider.Exchange(ctx, request.Code)
	if err != nil {
		return nil, err
	}

	// set the cookie header
	// redirect back to original location
	md := metadata.New(map[string]string{
		"Location":         redirectURL,
		"Location-Status":  "303",
		"Set-Cookie-Token": token.AccessToken,
	})

	if err := grpc.SendHeader(ctx, md); err != nil {
		return nil, err
	}

	return &authnv1.CallbackResponse{
		AccessToken: token.AccessToken,
	}, nil
}

func (a *api) CreateToken(ctx context.Context, request *authnv1.CreateTokenRequest) (*authnv1.CreateTokenResponse, error) {
	var prefixedSubject string
	switch request.TokenType {
	case authnv1.CreateTokenRequest_SERVICE:
		prefixedSubject = "service:" + request.Subject
	default:
		return nil, status.Error(codes.InvalidArgument, "invalid token type")
	}

	var expiry *time.Duration

	if request.Expiry != nil {
		convertedExpiry, err := ptypes.Duration(request.Expiry)
		if err != nil {
			return nil, err
		}
		expiry = &convertedExpiry
	}

	token, err := a.issuer.CreateToken(ctx, prefixedSubject, expiry)
	if err != nil {
		return nil, err
	}

	return &authnv1.CreateTokenResponse{
		AccessToken: token.AccessToken,
	}, nil
}
