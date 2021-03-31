package authn

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"

	authnv1 "github.com/lyft/clutch/backend/api/config/service/authn/v1"
	apimock "github.com/lyft/clutch/backend/mock/api"
	"github.com/lyft/clutch/backend/mock/service/authnmock"
)

func TestStateNonceRoundTrip(t *testing.T) {
	p := &OIDCProvider{
		sessionSecret: "this-is-my-secret",
	}

	url := "/foo"

	state, err := p.GetStateNonce(url)
	assert.NoError(t, err)
	assert.NotEmpty(t, state)

	outUrl, err := p.ValidateStateNonce(state)
	assert.NoError(t, err)
	assert.Equal(t, url, outUrl)

	// Check that the same fails if not signed correctly.
	p2 := &OIDCProvider{
		sessionSecret: "this-is-a-different-secret",
	}
	_, err = p2.ValidateStateNonce(state)
	assert.Error(t, err)
}

func TestStateNoncePrependsLeadingSlash(t *testing.T) {
	p := &OIDCProvider{
		sessionSecret: "this-is-my-secret",
	}
	s, err := p.GetStateNonce("dest/foo")
	assert.NoError(t, err)
	u, _ := p.ValidateStateNonce(s)
	assert.Equal(t, "/dest/foo", u)
}

func TestStateNonceRejections(t *testing.T) {
	p := &OIDCProvider{
		sessionSecret: "this-is-my-secret",
	}

	tests := []string{
		"unix:///tmp/redis.sock",
		"https://google.com",
		"http://google.com",
		"123%45%6",
	}

	for _, tt := range tests {
		tt := tt // pin!
		t.Run(tt, func(t *testing.T) {
			_, err := p.GetStateNonce(tt)
			assert.Error(t, err)
		})
	}
}

func TestNewOIDCProvider(t *testing.T) {
	cfg := &authnv1.Config{}
	apimock.FromYAML(`
session_secret: this_is_my_secret
oidc:
  issuer: http://foo.example.com
  client_id: my_client_id
  client_secret: my_client_secret
  redirect_url: "http://localhost:12000/v1/authn/callback"
  scopes:
  - openid
  - email
`, cfg)

	email := "user@example.com"

	mockprovider := authnmock.NewMockOIDCProviderServer(email)
	defer mockprovider.Close()

	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, mockprovider.Client())
	p, err := NewOIDCProvider(ctx, cfg, nil)
	assert.NoError(t, err)
	assert.NotNil(t, p)

	authURL, err := p.GetAuthCodeURL(context.Background(), "myState")
	assert.NoError(t, err)
	assert.True(t, strings.HasPrefix(authURL, "http://foo.example.com/oauth2/v1/authorize"))
	assert.True(t, strings.Contains(authURL, "access_type=offline"))
	assert.True(t, strings.Contains(authURL, "state=myState"))

	token, err := p.Exchange(context.Background(), "aaa")
	assert.NoError(t, err)
	assert.NotNil(t, token)

	c, err := p.Verify(context.Background(), token.AccessToken)
	assert.NoError(t, err)
	assert.NotNil(t, c)
	assert.Equal(t, email, c.Subject)
}

func TestCreateNewToken(t *testing.T) {
	cfg := &authnv1.Config{}
	apimock.FromYAML(`
session_secret: this_is_my_secret
oidc:
  issuer: http://foo.example.com
  client_id: my_client_id
  client_secret: my_client_secret
  redirect_url: "http://localhost:12000/v1/authn/callback"
  scopes:
  - openid
  - email
`, cfg)

	email := "user@example.com"

	mockprovider := authnmock.NewMockOIDCProviderServer(email)
	defer mockprovider.Close()

	mockStorage := authnmock.NewMockStorage()

	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, mockprovider.Client())

	p, err := NewOIDCProvider(ctx, cfg, mockStorage)
	assert.NoError(t, err)
	assert.NotNil(t, p)

	expiry := 5 * time.Hour
	createdToken, err := p.(Issuer).CreateToken(ctx, "some subject", &expiry)
	assert.NoError(t, err)
	assert.NotNil(t, createdToken)

	// The token should have been recorded in the database.
	assert.Len(t, mockStorage.Tokens[clutchProvider], 1)
	assert.Equal(t, mockStorage.Tokens[clutchProvider]["some subject"], createdToken)

	claims, err := p.Verify(ctx, createdToken.AccessToken)
	assert.NoError(t, err)

	assert.NotZero(t, claims.StandardClaims.IssuedAt)
	assert.Equal(t, claims.StandardClaims.ExpiresAt, time.Unix(claims.StandardClaims.IssuedAt, 0).Add(5*time.Hour).Unix())
	assert.Equal(t, claims.StandardClaims.Subject, "some subject")

	// Create a token without expiry.
	createdToken, err = p.(Issuer).CreateToken(ctx, "some subject", nil)
	assert.NoError(t, err)
	assert.NotNil(t, createdToken)

	// The token should have been recorded in the database.
	assert.Len(t, mockStorage.Tokens[clutchProvider], 1)
	assert.Equal(t, mockStorage.Tokens[clutchProvider]["some subject"], createdToken)

	claims, err = p.Verify(ctx, createdToken.AccessToken)
	assert.NoError(t, err)

	assert.NotZero(t, claims.StandardClaims.IssuedAt)
	assert.Equal(t, claims.StandardClaims.ExpiresAt, int64(0))
	assert.Equal(t, claims.StandardClaims.Subject, "some subject")
}

func TestTokenRevocationFlow(t *testing.T) {
	cfg := &authnv1.Config{}
	apimock.FromYAML(`
session_secret: this_is_my_secret
oidc:
  issuer: http://foo.example.com
  client_id: my_client_id
  client_secret: my_client_secret
  redirect_url: "http://localhost:12000/v1/authn/callback"
  scopes:
  - openid
  - email
`, cfg)

	email := "user@example.com"

	mockprovider := authnmock.NewMockOIDCProviderServer(email)
	defer mockprovider.Close()

	mockStorage := authnmock.NewMockStorage()

	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, mockprovider.Client())

	p, err := NewOIDCProvider(ctx, cfg, mockStorage)
	assert.NoError(t, err)
	assert.NotNil(t, p)

	authURL, err := p.GetAuthCodeURL(context.Background(), "myState")
	assert.NoError(t, err)
	assert.True(t, strings.HasPrefix(authURL, "http://foo.example.com/oauth2/v1/authorize"))
	assert.True(t, strings.Contains(authURL, "access_type=offline"))
	assert.True(t, strings.Contains(authURL, "state=myState"))

	token, err := p.Exchange(context.Background(), "aaa")
	assert.NoError(t, err)
	assert.NotNil(t, token)

	// Check the store to make sure we recorded the provider token.
	providerToken, err := mockStorage.Read(context.Background(), "user@example.com", "foo.example.com")
	assert.NoError(t, err)
	assert.NotNil(t, providerToken)

	// Check the store to make sure we recorded the clutch issued token.
	clutchToken, err := mockStorage.Read(context.Background(), "user@example.com", clutchProvider)
	assert.NoError(t, err)
	assert.NotNil(t, clutchToken)

	assert.NotEqual(t, providerToken, clutchToken)

	c, err := p.Verify(context.Background(), token.AccessToken)
	assert.NoError(t, err)
	assert.NotNil(t, c)
	assert.Equal(t, email, c.Subject)

	// TODO(snowp): Delete clutch issued token and verify that we are no longer valid.
}
