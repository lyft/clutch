package authn

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"

	authnmodulev1 "github.com/lyft/clutch/backend/api/authn/v1"
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
enable_service_token_creation: true
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
	createdToken, err := p.(Issuer).CreateToken(ctx, "some subject", authnmodulev1.CreateTokenRequest_SERVICE, &expiry)
	assert.NoError(t, err)
	assert.NotNil(t, createdToken)

	// The token should have been recorded in the database.
	assert.Len(t, mockStorage.Tokens[clutchProvider], 1)
	assert.Equal(t, mockStorage.Tokens[clutchProvider]["service:some subject"], createdToken)

	claims, err := p.Verify(ctx, createdToken.AccessToken)
	assert.NoError(t, err)

	assert.NotZero(t, claims.StandardClaims.IssuedAt)
	assert.Equal(t, claims.StandardClaims.ExpiresAt, time.Unix(claims.StandardClaims.IssuedAt, 0).Add(5*time.Hour).Unix())
	assert.Equal(t, claims.StandardClaims.Subject, "service:some subject")

	// Create a token without expiry.
	createdToken, err = p.(Issuer).CreateToken(ctx, "some subject", authnmodulev1.CreateTokenRequest_SERVICE, nil)
	assert.NoError(t, err)
	assert.NotNil(t, createdToken)

	// The token should have been recorded in the database.
	assert.Len(t, mockStorage.Tokens[clutchProvider], 1)
	assert.Equal(t, mockStorage.Tokens[clutchProvider]["service:some subject"], createdToken)

	claims, err = p.Verify(ctx, createdToken.AccessToken)
	assert.NoError(t, err)

	assert.NotZero(t, claims.StandardClaims.IssuedAt)
	assert.Equal(t, claims.StandardClaims.ExpiresAt, int64(0))
	assert.Equal(t, claims.StandardClaims.Subject, "service:some subject")

	// If we don't have a configured store we should be unable to create a token.
	p, err = NewOIDCProvider(ctx, cfg, nil)
	assert.NoError(t, err)
	assert.NotNil(t, p)

	createdToken, err = p.(Issuer).CreateToken(ctx, "some subject", authnmodulev1.CreateTokenRequest_SERVICE, nil)
	assert.Nil(t, createdToken)
	assert.Error(t, err)

	disabledServiceTokenConfig := &authnv1.Config{}
	apimock.FromYAML(`
session_secret: this_is_my_secret
enable_service_token_creation: false
oidc:
  issuer: http://foo.example.com
  client_id: my_client_id
  client_secret: my_client_secret
  redirect_url: "http://localhost:12000/v1/authn/callback"
  scopes:
  - openid
  - email
`, disabledServiceTokenConfig)

	// We have a configured store but token creation is not explicitly enabled.
	p, err = NewOIDCProvider(ctx, disabledServiceTokenConfig, mockStorage)
	assert.NoError(t, err)
	assert.NotNil(t, p)

	createdToken, err = p.(Issuer).CreateToken(ctx, "some subject", authnmodulev1.CreateTokenRequest_SERVICE, nil)
	assert.Nil(t, createdToken)
	assert.Error(t, err)
}

func TestReadThrough(t *testing.T) {
	cfg := &authnv1.Config{}
	apimock.FromYAML(`
session_secret: this_is_my_secret
enable_service_token_creation: false
oidc:
  issuer: http://foo.example.com
  client_id: my_client_id
  client_secret: my_client_secret
  redirect_url: "http://localhost:12000/v1/authn/callback"
  scopes:
  - openid
  - email
`, cfg)

	mockStorage := authnmock.NewMockStorage()

	mockprovider := authnmock.NewMockOIDCProviderServer("foo@example.com")
	defer mockprovider.Close()
	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, mockprovider.Client())

	p, err := NewOIDCProvider(ctx, cfg, mockStorage)
	assert.NoError(t, err)
	assert.NotNil(t, p)

	tr, ok := p.(*OIDCProvider)
	assert.True(t, ok)
	assert.NotNil(t, tr)

	assert.NoError(t, mockStorage.Store(context.Background(), "foo@example.com", tr.providerAlias, &oauth2.Token{AccessToken: "AAAA", RefreshToken: "RRRR"}))
	assert.NoError(t, mockStorage.Store(context.Background(), "fooExpired@example.com", tr.providerAlias, &oauth2.Token{AccessToken: "AAAA", RefreshToken: "RRRR", Expiry: time.Unix(1, 0)}))

	// Valid token returned.
	readToken, err := tr.Read(context.Background(), "foo@example.com", tr.providerAlias)
	assert.NoError(t, err)
	assert.Equal(t, "AAAA", readToken.AccessToken)

	// Now test an expired token refreshes.
	readToken, err = tr.Read(ctx, "fooExpired@example.com", tr.providerAlias)
	assert.NoError(t, err)
	assert.Equal(t, "AAAAAAAAAAAA", readToken.AccessToken)
	assert.Equal(t, "REFRESH", readToken.RefreshToken)
	storedToken, _ := mockStorage.Read(context.Background(), "fooExpired@example.com", tr.providerAlias)
	assert.Equal(t, "REFRESH", storedToken.RefreshToken)
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

	// Revoke all the tokens for clutch, leaving the provider token. This verifies that we explicitly
	// check for the existence of the clutch issued token in the database.
	delete(mockStorage.Tokens, clutchProvider)

	// Verification should now fail.
	c, err = p.Verify(context.Background(), token.AccessToken)
	assert.Error(t, err)
	assert.Nil(t, c)
}
