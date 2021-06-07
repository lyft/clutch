package authn

import (
	"context"
	"crypto/rsa"
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
	"gopkg.in/square/go-jose.v2"

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
	assert.Empty(t, createdToken.RefreshToken)

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

	c, err = p.Verify(context.Background(), token.RefreshToken)
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

func createIDToken(key *rsa.PrivateKey) string {
	idToken := `{"iss":"http://foo.example.com","aud":"my_client_id","email":"user@example.com","exp":` + strconv.FormatInt(time.Now().Add(time.Hour).Unix(), 10) + `}`

	jwk := &jose.JSONWebKey{Key: key, Algorithm: "RS256"}
	signer, err := jose.NewSigner(jose.SigningKey{Algorithm: "RS256", Key: jwk}, nil)
	if err != nil {
		panic(err)
	}
	jws, err := signer.Sign([]byte(idToken))
	if err != nil {
		panic(err)
	}
	data, err := jws.CompactSerialize()
	if err != nil {
		panic(err)
	}

	return data
}

func TestIssuerRefresh(t *testing.T) {
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

	testcases := []struct {
		claims                  jwt.StandardClaims
		invalidateProviderToken bool
		storedTokenDontMatch    bool
		expectError             string
	}{
		{
			claims: jwt.StandardClaims{Subject: email},
		},
		{
			claims:                  jwt.StandardClaims{Subject: email},
			invalidateProviderToken: true,
		},
		{
			claims:      jwt.StandardClaims{Subject: email, ExpiresAt: time.Now().Add(-1 * time.Hour).Unix()},
			expectError: "token is expired",
		},
		{
			claims:               jwt.StandardClaims{Subject: email},
			storedTokenDontMatch: true,
			expectError:          "did not match",
		},
	}

	for idx, tc := range testcases {
		tc := tc
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			mockprovider := authnmock.NewMockOIDCProviderServer(email)
			defer mockprovider.Close()
			ctx := context.WithValue(context.Background(), oauth2.HTTPClient, mockprovider.Client())

			mockStorage := authnmock.NewMockStorage()
			p, err := NewOIDCProvider(ctx, cfg, mockStorage)
			assert.NoError(t, err)

			refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, tc.claims).SignedString([]byte("this_is_my_secret"))
			assert.NoError(t, err)

			providerToken := &oauth2.Token{AccessToken: "AAAA", RefreshToken: refreshToken}
			providerToken = providerToken.WithExtra(map[string]interface{}{"id_token": createIDToken(mockprovider.Key)})

			if tc.invalidateProviderToken {
				providerToken.Expiry = time.Now().Add(-1 * time.Hour)
			}

			if tc.storedTokenDontMatch {
				assert.NoError(t, mockStorage.Store(ctx, email, clutchProvider, &oauth2.Token{AccessToken: "AAAA", RefreshToken: refreshToken + "BADTOKEN"}))
			} else {
				assert.NoError(t, mockStorage.Store(ctx, email, clutchProvider, &oauth2.Token{AccessToken: "AAAA", RefreshToken: refreshToken}))
			}
			assert.NoError(t, mockStorage.Store(ctx, email, "foo.example.com", providerToken))

			tok, err := p.(*OIDCProvider).RefreshToken(ctx, &oauth2.Token{
				RefreshToken: refreshToken,
			})

			if tc.expectError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectError)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, tok)
			}

			if tc.invalidateProviderToken {
				assert.Equal(t, 1, mockprovider.TokenCount)
			} else {
				assert.Equal(t, 0, mockprovider.TokenCount)
			}
		})
	}
}
