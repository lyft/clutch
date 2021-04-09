package authn

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"golang.org/x/oauth2"

	authnmodulev1 "github.com/lyft/clutch/backend/api/authn/v1"
	authnv1 "github.com/lyft/clutch/backend/api/config/service/authn/v1"
)

// Default scopes, used if no scopes are provided in the configuration.
// Compatible with Okta offline access, a holdover from previous defaults.
var defaultScopes = []string{oidc.ScopeOpenID, oidc.ScopeOfflineAccess, "email"}

const clutchProvider = "clutch"

type OIDCProvider struct {
	provider *oidc.Provider
	verifier *oidc.IDTokenVerifier
	oauth2   *oauth2.Config

	httpClient *http.Client

	sessionSecret string

	tokenStorage  Storage
	providerAlias string

	claimsFromOIDCToken ClaimsFromOIDCTokenFunc

	enableServiceTokenCreation bool
}

// Clutch's state token claims used during the exchange.
type stateClaims struct {
	*jwt.StandardClaims
	RedirectURL string `json:"redirect"`
}

// Intermediate claims object for the ID token. Based on what scopes were requested.
type idClaims struct {
	Email string `json:"email"`
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
	if u.Scheme != "" || u.Host != "" {
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

func (p *OIDCProvider) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	// Exchange.
	ctx = oidc.ClientContext(ctx, p.httpClient)

	// offline_access is used to request issuance of a refresh_token. Some providers may request it as a scope though.
	// Also it may need to be configurable in the future depending on the requirements of other providers or users.
	token, err := p.oauth2.Exchange(ctx, code, oauth2.AccessTypeOffline)
	if err != nil {
		return nil, err
	}

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, errors.New("'id_token' was not present in oauth token")
	}

	// Verify. This is superfluous since the token was just issued but it can't hurt.
	idToken, err := p.verifier.Verify(ctx, rawIDToken)
	if err != nil {
		return nil, err
	}

	// Issue token with claims.
	claims, err := p.claimsFromOIDCToken(ctx, idToken)
	if err != nil {
		return nil, err
	}

	if p.tokenStorage != nil {
		err := p.tokenStorage.Store(ctx, claims.Subject, p.providerAlias, token)
		if err != nil {
			return nil, err
		}
	}

	return p.signNewToken(ctx, claims)
}

func (p *OIDCProvider) CreateToken(ctx context.Context, subject string, tokenType authnmodulev1.CreateTokenRequest_TokenType, expiry *time.Duration) (*oauth2.Token, error) {
	if !p.enableServiceTokenCreation {
		return nil, errors.New("not configured to allow service token creation")
	}

	var prefixedSubject string
	switch tokenType {
	case authnmodulev1.CreateTokenRequest_SERVICE:
		prefixedSubject = "service:" + subject
	default:
		return nil, errors.New("invalid token type")
	}

	issuedAt := time.Now()
	var expiresAt int64
	if expiry != nil {
		expiresAt = issuedAt.Add(*expiry).Unix()
	}

	claims := &Claims{
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: expiresAt,
			IssuedAt:  issuedAt.Unix(),
			Issuer:    clutchProvider,
			Subject:   prefixedSubject,
		},
	}

	return p.signNewToken(ctx, claims)
}

func (p *OIDCProvider) signNewToken(ctx context.Context, claims *Claims) (*oauth2.Token, error) {
	// Sign and issue token.
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(p.sessionSecret))
	if err != nil {
		return nil, err
	}
	t := &oauth2.Token{
		AccessToken:  accessToken,
		Expiry:       time.Unix(claims.ExpiresAt, 0),
		RefreshToken: "", // TODO: implement refresh_token flow with stateful sessions.
		TokenType:    "Bearer",
	}

	if p.tokenStorage != nil {
		err := p.tokenStorage.Store(ctx, claims.Subject, clutchProvider, t)
		if err != nil {
			return nil, err
		}
	}

	return t, nil
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
		Groups:         []string{""},
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

	// If the token doesn't exist in the token storage anymore, it must have been revoked.
	// Fail verification in this case.
	// TODO(perf): Cache the lookup result in-memory for min(60, timeToExpiry) to prevent
	// hitting the DB on each request. This should also cache whether we didn't find a token.
	if p.tokenStorage != nil {
		_, err := p.tokenStorage.Read(ctx, claims.Subject, clutchProvider)
		if err != nil {
			return nil, err
		}
	}

	return claims, nil
}

func NewOIDCProvider(ctx context.Context, config *authnv1.Config, tokenStorage Storage) (Provider, error) {
	c := config.GetOidc()

	// Allows injection of test client. If client not present then add the default.
	if v := ctx.Value(oauth2.HTTPClient); v == nil {
		ctx = oidc.ClientContext(ctx, &http.Client{})
	}

	u, err := url.Parse(c.Issuer)
	if err != nil {
		return nil, err
	}
	alias := u.Hostname()

	provider, err := oidc.NewProvider(ctx, c.Issuer)
	if err != nil {
		return nil, err
	}

	verifier := provider.Verifier(&oidc.Config{
		ClientID: c.ClientId,
	})

	scopes := c.Scopes
	if len(scopes) == 0 {
		scopes = defaultScopes
	}

	oc := &oauth2.Config{
		ClientID:     c.ClientId,
		ClientSecret: c.ClientSecret,
		Endpoint:     provider.Endpoint(),
		RedirectURL:  c.RedirectUrl,
		Scopes:       scopes,
	}

	// Verify the provider implements the same flow we do.
	pClaims := &oidcProviderClaims{}
	if err := provider.Claims(pClaims); err != nil {
		return nil, err
	}
	if err := pClaims.Check("authorization_code"); err != nil {
		return nil, err
	}

	p := &OIDCProvider{
		providerAlias:              alias,
		provider:                   provider,
		verifier:                   verifier,
		oauth2:                     oc,
		httpClient:                 ctx.Value(oauth2.HTTPClient).(*http.Client),
		sessionSecret:              config.SessionSecret,
		claimsFromOIDCToken:        DefaultClaimsFromOIDCToken,
		tokenStorage:               tokenStorage,
		enableServiceTokenCreation: tokenStorage != nil && config.EnableServiceTokenCreation,
	}

	return p, nil
}

func oidcTokenToStandardClaims(t *oidc.IDToken) *jwt.StandardClaims {
	return &jwt.StandardClaims{
		ExpiresAt: t.Expiry.Unix(),
		IssuedAt:  t.IssuedAt.Unix(),
		Issuer:    t.Issuer,
		Subject:   t.Subject,
	}
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
