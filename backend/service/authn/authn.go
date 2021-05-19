package authn

// <!-- START clutchdoc -->
// description: Produces tokens for the configured OIDC provider.
// <!-- END clutchdoc -->

import (
	"context"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/uber-go/tally"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"google.golang.org/protobuf/types/known/anypb"

	authnmodulev1 "github.com/lyft/clutch/backend/api/authn/v1"
	authnv1 "github.com/lyft/clutch/backend/api/config/service/authn/v1"
	"github.com/lyft/clutch/backend/service"
)

const Name = "clutch.service.authn"

// AlwaysAllowedMethods is a list of method patterns that are always open and not blocked by authn or authz.
// TODO(maybe): convert this to an API annotation or make configurable on the middleware that use the list.
var AlwaysAllowedMethods = []string{
	"/clutch.authn.v1.AuthnAPI/Callback",
	"/clutch.authn.v1.AuthnAPI/Login",
	"/clutch.healthcheck.v1.HealthcheckAPI/*",
}

func New(cfg *anypb.Any, logger *zap.Logger, scope tally.Scope) (service.Service, error) {
	config := &authnv1.Config{}
	if err := cfg.UnmarshalTo(config); err != nil {
		return nil, err
	}

	tokenStorage, _ := service.Registry[StorageName].(Storage) // Ignoring 'ok', nil allowed.

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

type Issuer interface {
	// CreateToken creates a new OAuth2 for the provided subject with the provided expiration. If expiry is nil,
	// the token will never expire.
	CreateToken(ctx context.Context, subject string, tokenType authnmodulev1.CreateTokenRequest_TokenType, expiry *time.Duration) (token *oauth2.Token, err error)
	RefreshToken(ctx context.Context, token *oauth2.Token) (*oauth2.Token, error)
}

type Service interface {
	Issuer
	Provider
	TokenReader // Read calls are proxied through the IssuerProvider so the token can be refreshed if needed.
}

type TokenReader interface {
	Read(ctx context.Context, userID, provider string) (*oauth2.Token, error)
}

type TokenStorer interface {
	Store(ctx context.Context, userID, provider string, token *oauth2.Token) error
}

type Storage interface {
	TokenReader
	TokenStorer
}
