package authn

// <!-- START clutchdoc -->
// description: Registers login and callback endpoints for OAuth 2 flows.
// <!-- END clutchdoc -->

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang/protobuf/ptypes/any"
	"github.com/uber-go/tally"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	authnv1 "github.com/lyft/clutch/backend/api/authn/v1"
	"github.com/lyft/clutch/backend/gateway/mux"
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

	p, ok := svc.(authn.Service)
	if !ok {
		return nil, errors.New("authn service was not the correct type")
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
	// Attempt refresh. If refresh fails for any reason continue the regular auth flow.
	md, _ := metadata.FromIncomingContext(ctx)
	if v := md.Get("grpcgateway-cookie"); len(v) > 0 {
		if refreshToken, err := mux.GetCookieValue(v, "refreshToken"); err == nil {
			if t, err := a.issuer.RefreshToken(ctx, &oauth2.Token{RefreshToken: refreshToken}); err == nil {
				md := metadata.New(map[string]string{
					"Location":         request.RedirectUrl,
					"Set-Cookie-Token": t.AccessToken,
					"Set-Cookie-Refresh-Token": t.RefreshToken,
				})
				if err := grpc.SendHeader(ctx, md); err != nil {
					return nil, err
				}

				return &authnv1.LoginResponse{
					Return: &authnv1.LoginResponse_Token_{
						Token: &authnv1.LoginResponse_Token{
							AccessToken:  t.AccessToken,
							RefreshToken: t.RefreshToken,
						},
					},
				}, nil

			}
		}
	}

	// Full login exchange flow.
	state, err := a.provider.GetStateNonce(request.RedirectUrl)
	if err != nil {
		return nil, err
	}
	authURL, err := a.provider.GetAuthCodeURL(ctx, state)
	if err != nil {
		return nil, err
	}

	if err := grpc.SendHeader(ctx, metadata.Pairs("Location", authURL)); err != nil {
		return nil, err
	}

	return &authnv1.LoginResponse{
		Return: &authnv1.LoginResponse_AuthUrl{AuthUrl: authURL},
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
		"Set-Cookie-Token": token.AccessToken,
	})

	if token.RefreshToken != "" {
		md.Set("Set-Cookie-Refresh-Token", token.RefreshToken)
	}

	if err := grpc.SendHeader(ctx, md); err != nil {
		return nil, err
	}

	return &authnv1.CallbackResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}, nil
}

func (a *api) CreateToken(ctx context.Context, request *authnv1.CreateTokenRequest) (*authnv1.CreateTokenResponse, error) {
	var expiry *time.Duration

	if request.Expiry != nil {
		convertedExpiry := request.Expiry.AsDuration()
		expiry = &convertedExpiry
	}

	token, err := a.issuer.CreateToken(ctx, request.Subject, request.TokenType, expiry)
	if err != nil {
		return nil, err
	}

	return &authnv1.CreateTokenResponse{
		AccessToken: token.AccessToken,
	}, nil
}
