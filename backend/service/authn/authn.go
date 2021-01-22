package authn

// <!-- START clutchdoc -->
// description: Produces tokens for the configured OIDC provider.
// <!-- END clutchdoc -->

import (
	"context"
	"fmt"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/dgrijalva/jwt-go"
	"github.com/uber-go/tally"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"google.golang.org/protobuf/types/known/anypb"

	authnv1 "github.com/lyft/clutch/backend/api/config/service/authn/v1"
	"github.com/lyft/clutch/backend/service"
)

// TODO(maybe): make scopes part of config?
// For more documentation on scopes see: https://developer.okta.com/docs/reference/api/oidc/#scopes
var scopes = []string{
	oidc.ScopeOpenID,
	oidc.ScopeOfflineAccess, // offline_access is used to request issuance of a refresh_token
	"email",
}

const Name = "clutch.service.authn"

func New(cfg *anypb.Any, logger *zap.Logger, scope tally.Scope) (service.Service, error) {
	config := &authnv1.Config{}
	if err := cfg.UnmarshalTo(config); err != nil {
		return nil, err
	}

	tokenStorage, err := newStorage(config.Storage)
	if err != nil {
		return nil, err
	}

	switch t := config.Type.(type) {
	case *authnv1.Config_Oidc:
		return NewOIDCProvider(context.Background(), config, tokenStorage)
	default:
		return nil, fmt.Errorf("authn provider type '%T' not implemented", t)
	}
}

// Standardized representation of a user's claims.
type Claims struct {
	*jwt.StandardClaims

	// Groups could be derived from the token or an external mapping.
	Groups []string `json:"grp,omitempty"`
}

type Provider interface {
	GetStateNonce(redirectURL string) (string, error)
	ValidateStateNonce(state string) (redirectURL string, err error)

	Verify(ctx context.Context, rawIDToken string) (*Claims, error)
	GetAuthCodeURL(ctx context.Context, state string) (string, error)
	Exchange(ctx context.Context, code string) (token *oauth2.Token, err error)
}
