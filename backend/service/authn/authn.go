package authn

// <!-- START clutchdoc -->
// description: Produces tokens for the configured OIDC provider.
// <!-- END clutchdoc -->

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/coreos/go-oidc"
	"github.com/dgrijalva/jwt-go"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/google/uuid"
	"github.com/uber-go/tally"
	"go.uber.org/zap"
	"golang.org/x/oauth2"

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

func New(cfg *any.Any, logger *zap.Logger, scope tally.Scope) (service.Service, error) {
	config := &authnv1.Config{}
	if err := ptypes.UnmarshalAny(cfg, config); err != nil {
		return nil, err
	}
	return NewProvider(config)
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
	Exchange(ctx context.Context, code string) (token string, err error)
}

type OIDCProvider struct {
	provider *oidc.Provider
	verifier *oidc.IDTokenVerifier
	oauth2   *oauth2.Config

	httpClient *http.Client

	sessionSecret string

	claimsFromOIDCToken ClaimsFromOIDCTokenFunc
}

func WithClaimsFromOIDCTokenFunc(p *OIDCProvider, fn ClaimsFromOIDCTokenFunc) *OIDCProvider {
	ret := *p
	ret.claimsFromOIDCToken = fn
	return &ret
}

func (p *OIDCProvider) GetAuthCodeURL(ctx context.Context, state string) (string, error) {
	opts := []oauth2.AuthCodeOption{oauth2.AccessTypeOffline}
	return p.oauth2.AuthCodeURL(state, opts...), nil
}

func (p *OIDCProvider) ValidateStateNonce(state string) (string, error) {
	claims := &stateClaims{}
	_, err := jwt.ParseWithClaims(state, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(p.sessionSecret), nil
	})
	if err != nil {
		return "", err
	}
	if err := claims.Valid(); err != nil {
		return "", err
	}
	return claims.RedirectURL, nil
}

func (p *OIDCProvider) GetStateNonce(redirectURL string) (string, error) {
	u, err := url.Parse(redirectURL)
	if err != nil {
		return "", err
	}
	if u.Host != "" {
		return "", errors.New("only relative redirects are supported")
	}
	dest := u.Path
	if !strings.HasPrefix(dest, "/") {
		dest = fmt.Sprintf("/%s", dest)
	}

	claims := &stateClaims{
		StandardClaims: &jwt.StandardClaims{
			Subject:   uuid.New().String(), // UUID serves as CSRF token.
			ExpiresAt: time.Now().Add(time.Minute * 5).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		RedirectURL: dest,
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(p.sessionSecret))
}

type stateClaims struct {
	*jwt.StandardClaims
	RedirectURL string `json:"redirect"`
}

func (p *OIDCProvider) Exchange(ctx context.Context, code string) (string, error) {
	// Exchange.
	ctx = oidc.ClientContext(ctx, p.httpClient)
	token, err := p.oauth2.Exchange(ctx, code)
	if err != nil {
		return "", err
	}

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return "", errors.New("'id_token' was not present in oauth token")
	}

	// Verify.
	idToken, err := p.verifier.Verify(ctx, rawIDToken)
	if err != nil {
		return "", err
	}

	// Issue token with claims.
	claims, err := p.claimsFromOIDCToken(ctx, idToken)
	if err != nil {
		return "", err
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(p.sessionSecret))
}

// Intermediate claims object for the ID token. Based on what scopes were requested.
type idClaims struct {
	Email string `json:"email"`
}

func oidcTokenToStandardClaims(t *oidc.IDToken) *jwt.StandardClaims {
	return &jwt.StandardClaims{
		ExpiresAt: t.Expiry.Unix(),
		IssuedAt:  t.IssuedAt.Unix(),
		Issuer:    t.Issuer,
		Subject:   t.Subject,
	}
}

type ClaimsFromOIDCTokenFunc func(ctx context.Context, t *oidc.IDToken) (*Claims, error)

// Extract claims from an OIDC token and return Clutch's standard claims object. This could be configurable at a later
// date to support subjects with IDs other than email (e.g. GitHub ID).
func DefaultClaimsFromOIDCToken(ctx context.Context, t *oidc.IDToken) (*Claims, error) {
	idc := &idClaims{}
	if err := t.Claims(idc); err != nil {
		return nil, err
	}
	if idc.Email == "" {
		return nil, errors.New("claims did not deserialize with desired fields")
	}

	sc := oidcTokenToStandardClaims(t)
	sc.Subject = idc.Email
	return &Claims{
		StandardClaims: sc,
		Groups:         []string{"12345"},
	}, nil
}

func (p *OIDCProvider) Verify(ctx context.Context, rawToken string) (*Claims, error) {
	claims := &Claims{}
	_, err := jwt.ParseWithClaims(rawToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(p.sessionSecret), nil
	})
	if err != nil {
		return nil, err
	}

	if err := claims.Valid(); err != nil {
		return nil, err
	}

	return claims, nil
}

// Evaluates what flows the provider claims to support.
type oidcProviderClaims struct {
	GrantTypesSupported []string `json:"grant_types_supported"`
}

func (pc *oidcProviderClaims) Check(grantType string) error {
	for _, gt := range pc.GrantTypesSupported {
		if gt == grantType {
			return nil
		}
	}
	return fmt.Errorf("grant type '%s' not supported by provider. supported: %v", grantType, pc.GrantTypesSupported)
}

func NewProvider(config *authnv1.Config) (Provider, error) {
	c := config.GetOidc()

	httpClient := &http.Client{}
	ctx := oidc.ClientContext(context.Background(), httpClient)
	provider, err := oidc.NewProvider(ctx, c.Issuer)
	if err != nil {
		return nil, err
	}

	verifier := provider.Verifier(&oidc.Config{
		ClientID: c.ClientId,
	})

	oc := &oauth2.Config{
		ClientID:     c.ClientId,
		ClientSecret: c.ClientSecret,
		Endpoint:     provider.Endpoint(),
		RedirectURL:  c.RedirectUrl,
		Scopes:       scopes,
	}

	// Verify the provider implements the same flow we do.
	pClaims := &oidcProviderClaims{}
	if err := provider.Claims(&pClaims); err != nil {
		return nil, err
	}
	if err := pClaims.Check("authorization_code"); err != nil {
		return nil, err
	}

	return &OIDCProvider{
		provider:            provider,
		verifier:            verifier,
		oauth2:              oc,
		httpClient:          httpClient,
		sessionSecret:       config.SessionSecret,
		claimsFromOIDCToken: DefaultClaimsFromOIDCToken,
	}, nil
}
