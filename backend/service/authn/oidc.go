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

	authnv1 "github.com/lyft/clutch/backend/api/config/service/authn/v1"
)

type OIDCProvider struct {
	provider *oidc.Provider
	verifier *oidc.IDTokenVerifier
	oauth2   *oauth2.Config

	httpClient *http.Client

	sessionSecret string

	cryptographer *cryptographer

	claimsFromOIDCToken ClaimsFromOIDCTokenFunc
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
	token, err := p.oauth2.Exchange(ctx, code)
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

	// TODO: Store the encrypted refresh token in the database.
	if p.cryptographer != nil {
		if _, err := p.cryptographer.Encrypt([]byte(token.RefreshToken)); err != nil {
			return nil, err
		}
	}

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
	return t, err
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

	return claims, nil
}

func NewOIDCProvider(ctx context.Context, config *authnv1.Config) (Provider, error) {
	c := config.GetOidc()

	// Allows injection of test client. If client not present then add the default.
	if v := ctx.Value(oauth2.HTTPClient); v == nil {
		ctx = oidc.ClientContext(ctx, &http.Client{})
	}

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

	p := &OIDCProvider{
		provider:            provider,
		verifier:            verifier,
		oauth2:              oc,
		httpClient:          ctx.Value(oauth2.HTTPClient).(*http.Client),
		sessionSecret:       config.SessionSecret,
		claimsFromOIDCToken: DefaultClaimsFromOIDCToken,
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
