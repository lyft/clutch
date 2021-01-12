package authn

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"gopkg.in/square/go-jose.v2"
	"math"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"

	authnv1 "github.com/lyft/clutch/backend/api/config/service/authn/v1"
	apimock "github.com/lyft/clutch/backend/mock/api"
)

const openIDConfiguration = `{"issuer":"http://foo.example.com","authorization_endpoint":"http://foo.example.com/oauth2/v1/authorize","token_endpoint":"http://foo.example.com/oauth2/v1/token","userinfo_endpoint":"http://foo.example.com/oauth2/v1/userinfo","registration_endpoint":"http://foo.example.com/oauth2/v1/clients","jwks_uri":"http://foo.example.com/oauth2/v1/keys","response_types_supported":["code","id_token","code id_token","code token","id_token token","code id_token token"],"response_modes_supported":["query","fragment","form_post","okta_post_message"],"grant_types_supported":["authorization_code","implicit","refresh_token","password"],"subject_types_supported":["public"],"id_token_signing_alg_values_supported":["RS256"],"scopes_supported":["openid","email","profile","address","phone","offline_access","groups"],"token_endpoint_auth_methods_supported":["client_secret_basic","client_secret_post","client_secret_jwt","private_key_jwt","none"],"claims_supported":["iss","ver","sub","aud","iat","exp","jti","auth_time","amr","idp","nonce","name","nickname","preferred_username","given_name","middle_name","family_name","email","email_verified","profile","zoneinfo","locale","address","phone_number","picture","website","gender","birthdate","updated_at","at_hash","c_hash"],"code_challenge_methods_supported":["S256"],"introspection_endpoint":"http://foo.example.com/oauth2/v1/introspect","introspection_endpoint_auth_methods_supported":["client_secret_basic","client_secret_post","client_secret_jwt","private_key_jwt","none"],"revocation_endpoint":"http://foo.example.com/oauth2/v1/revoke","revocation_endpoint_auth_methods_supported":["client_secret_basic","client_secret_post","client_secret_jwt","private_key_jwt","none"],"end_session_endpoint":"http://foo.example.com/oauth2/v1/logout","request_parameter_supported":true,"request_object_signing_alg_values_supported":["HS256","HS384","HS512","RS256","RS384","RS512","ES256","ES384","ES512"]}`

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

type testIdTokenClaims struct {
	*jwt.StandardClaims
	Email string `json:"email"`
}

func TestNewProvider(t *testing.T) {
	cfg := &authnv1.Config{}
	apimock.FromYAML(`
session_secret: this_is_my_secret
oidc:
  issuer: http://foo.example.com
  client_id: my_client_id
  client_secret: my_client_secret
  redirect_url: "http://localhost:12000/v1/authn/callback"
`, cfg)

	key, _ := rsa.GenerateKey(rand.Reader, 2048)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("RECEIVED REQUEST", r.URL.Path)

		w.Header().Set("content-type", "application/json")
		switch r.URL.Path {
		case "/.well-known/openid-configuration":
			fmt.Fprintln(w, openIDConfiguration)
		case "/oauth2/v1/token":
			claims := &testIdTokenClaims{
				StandardClaims: &jwt.StandardClaims{
					Issuer: "http://foo.example.com",
					Audience: "my_client_id",
					ExpiresAt: math.MaxInt32,
				},
				Email: "user@example.com",
			}

			tok, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(key)
			if err != nil {
				panic(err)
			}

			fmt.Fprintln(w, fmt.Sprintf(`{"token_type":"bearer","access_token":"AAAAAAAAAAAA", "id_token": "%s"}`, tok))
		case "/oauth2/v1/keys":
			jwk := jose.JSONWebKey{KeyID: "foo", Key: key.Public()}
			jwks := jose.JSONWebKeySet{Keys: []jose.JSONWebKey{jwk}}
			jks, _ := json.Marshal(jwks)
			w.Write(jks)
		}
	}))
	defer ts.Close()

	ts.Client().Transport = &http.Transport{DialContext: func(_ context.Context, network, _ string) (net.Conn, error) {
		// Redirect all requests to the httptest server.
		return net.Dial(network, ts.Listener.Addr().String())
	}}

	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, ts.Client())
	p, err := NewOIDCProvider(ctx, cfg)
	assert.NoError(t, err)
	assert.NotNil(t, p)

	authURL, err := p.GetAuthCodeURL(context.Background(), "myState")
	assert.True(t, strings.HasPrefix(authURL, "http://foo.example.com/oauth2/v1/authorize"))
	assert.True(t, strings.Contains(authURL, "access_type=offline"))
	assert.True(t, strings.Contains(authURL, "state=myState"))

	token, err := p.Exchange(ctx, "aaa")
	assert.NoError(t, err)
	assert.NotNil(t, token)

	c, err := p.Verify(ctx, token.AccessToken)
	assert.NoError(t, err)
	assert.NotNil(t, c)

}
