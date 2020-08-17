package authn

// <!-- START clutchdoc -->
// description: Registers login and callback endpoints for OAuth 2 flows.
// <!-- END clutchdoc -->

import (
	"context"
	"errors"
	"fmt"

	"github.com/golang/protobuf/ptypes/any"
	"github.com/uber-go/tally"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

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

	p, ok := svc.(authn.Provider)
	if !ok {
		return nil, errors.New("authn service was not the correct type")
	}

	return &mod{
		authnv1: &api{svc: p},
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
	svc authn.Provider
}

func (a *api) Login(ctx context.Context, request *authnv1.LoginRequest) (*authnv1.LoginResponse, error) {
	state, err := a.svc.GetStateNonce(request.RedirectUrl)
	if err != nil {
		return nil, err
	}
	authURL, err := a.svc.GetAuthCodeURL(ctx, state)
	if err != nil {
		return nil, err
	}

	return &authnv1.LoginResponse{
		AuthUrl: authURL,
	}, nil
}

func (a *api) Callback(ctx context.Context, request *authnv1.CallbackRequest) (*authnv1.CallbackResponse, error) {
	if request.Error != "" {
		return nil, fmt.Errorf("%s: %s", request.Error, request.ErrorDescription)
	}

	redirectURL, err := a.svc.ValidateStateNonce(request.State)
	if err != nil {
		return nil, err
	}

	token, err := a.svc.Exchange(ctx, request.Code)
	if err != nil {
		return nil, err
	}

	// set the cookie header
	// redirect back to original location
	md := metadata.New(map[string]string{
		"Location":         redirectURL,
		"Location-Status":  "303",
		"Set-Cookie-Token": token,
	})

	if err := grpc.SendHeader(ctx, md); err != nil {
		return nil, err
	}

	return &authnv1.CallbackResponse{
		Token: token,
	}, nil
}
