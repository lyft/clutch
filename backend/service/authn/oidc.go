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
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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
	dest := u.RequestURI()
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

	return p.issueAndStoreToken(ctx, claims, true)
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

	return p.issueAndStoreToken(ctx, claims, false)
}

// Refresh the issuer token. If the provider token is not valid, refresh it. If any error occurs continue auth code flow.
func (p *OIDCProvider) RefreshToken(ctx context.Context, t *oauth2.Token) (*oauth2.Token, error) {
	// Extract claims from refresh token. They are also validated by the parser (i.e. to check for expiry).
	claims := &jwt.StandardClaims{}
	_, err := jwt.ParseWithClaims(t.RefreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(p.sessionSecret), nil
	})
	if err != nil {
		return nil, err
	}

	// Look up refresh token and verify it matches the one stored in the database.
	rt, err := p.tokenStorage.Read(ctx, claims.Subject, clutchProvider)
	if err != nil {
		return nil, err
	}
	if rt.RefreshToken != t.RefreshToken {
		return nil, errors.New("refresh token did not match")
	}

	// Verify provider token is still valid.
	pt, err := p.tokenStorage.Read(ctx, claims.Subject, p.providerAlias)
	if err != nil {
		return nil, err
	}

	// Attempt to refresh provider token if not valid.
	if !pt.Valid() {
		pt, err = p.oauth2.TokenSource(ctx, pt).Token()
		if err != nil {
			return nil, err
		}

		// Store refreshed provider token.
		if err := p.tokenStorage.Store(ctx, claims.Subject, p.providerAlias, pt); err != nil {
			return nil, err
		}
	}

	// Create a new token.
	rawIDToken, ok := pt.Extra("id_token").(string)
	if !ok {
		return nil, errors.New("'id_token' was not present in provider token")
	}

	// Verify. This is superfluous since the token was just issued but it can't hurt.
	idToken, err := p.verifier.Verify(ctx, rawIDToken)
	if err != nil {
		return nil, err
	}

	// Issue token with claims.
	newClaims, err := p.claimsFromOIDCToken(ctx, idToken)
	if err != nil {
		return nil, err
	}

	newToken, err := p.issueAndStoreToken(ctx, newClaims, true)
	if err != nil {
		return nil, err
	}

	// Return token.
	return newToken, nil
}

// Issues and stores a token based on the provided claims. If refresh is true and storage is enabled, a refresh
// token will be issued as well.
func (p *OIDCProvider) issueAndStoreToken(ctx context.Context, claims *Claims, refresh bool) (*oauth2.Token, error) {
	// Sign and issue token.
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(p.sessionSecret))
	if err != nil {
		return nil, err
	}
	t := &oauth2.Token{
		AccessToken: accessToken,
		Expiry:      time.Unix(claims.ExpiresAt, 0),
		TokenType:   "Bearer",
	}

	if p.tokenStorage != nil {
		if refresh {
			refreshClaims := &jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * 12).Unix(), // TODO: configurable refresh token lifetime
				Issuer:    claims.Issuer,
				Subject:   claims.Subject,
			}

			refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(p.sessionSecret))
			if err != nil {
				return nil, err
			}
			t.RefreshToken = refreshToken
		}

		if err := p.tokenStorage.Store(ctx, claims.Subject, clutchProvider, t); err != nil {
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

func (p *OIDCProvider) Read(ctx context.Context, userID, provider string) (*oauth2.Token, error) {
	if p.tokenStorage == nil {
		return nil, status.Error(codes.Internal, "token read attempted but storage is not configured")
	}

	if provider != p.providerAlias {
		return nil, status.Errorf(codes.InvalidArgument, "provider '%s' cannot read '%s' tokens", p.providerAlias, provider)
	}

	// Get token from storage.
	t, err := p.tokenStorage.Read(ctx, userID, provider)
	if err != nil {
		return nil, err
	}

	// Validate token and return it if valid.
	if t.Valid() {
		return t, nil
	}

	if t.RefreshToken == "" {
		return nil, status.Error(codes.Unauthenticated, "the token has expired and no refresh token is present")
	}

	// If invalid, attempt refresh.
	newToken, err := p.oauth2.TokenSource(ctx, t).Token()
	if err != nil {
		return nil, err
	}

	// Store new token if refresh succeeded.
	if err := p.tokenStorage.Store(ctx, userID, provider, newToken); err != nil {
		return nil, err
	}

	return newToken, nil
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
