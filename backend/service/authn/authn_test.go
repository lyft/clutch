package authn

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
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

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("RECEIVED REQUEST", r.URL)
		if r.URL.Path == "/.well-known/openid-configuration" {
			fmt.Fprintln(w, openIDConfiguration)
		}
	}))
	defer ts.Close()

	ts.Client().Transport = &http.Transport{DialContext: func(_ context.Context, network, _ string) (net.Conn, error) {
		// Redirect all requests to the httptest server.
		return net.Dial(network, ts.Listener.Addr().String())
	}}

	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, ts.Client())
	p, err := NewProvider(ctx, cfg)
	assert.NoError(t, err)
	assert.NotNil(t, p)
}
