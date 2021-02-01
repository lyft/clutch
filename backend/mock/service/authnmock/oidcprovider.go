package authnmock

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"math"
	"net"
	"net/http"
	"net/http/httptest"

	"github.com/dgrijalva/jwt-go"
	"gopkg.in/square/go-jose.v2"
)

const openIDConfiguration = `{"issuer":"http://foo.example.com","authorization_endpoint":"http://foo.example.com/oauth2/v1/authorize","token_endpoint":"http://foo.example.com/oauth2/v1/token","userinfo_endpoint":"http://foo.example.com/oauth2/v1/userinfo","registration_endpoint":"http://foo.example.com/oauth2/v1/clients","jwks_uri":"http://foo.example.com/oauth2/v1/keys","response_types_supported":["code","id_token","code id_token","code token","id_token token","code id_token token"],"response_modes_supported":["query","fragment","form_post","okta_post_message"],"grant_types_supported":["authorization_code","implicit","refresh_token","password"],"subject_types_supported":["public"],"id_token_signing_alg_values_supported":["RS256"],"scopes_supported":["openid","email","profile","address","phone","offline_access","groups"],"token_endpoint_auth_methods_supported":["client_secret_basic","client_secret_post","client_secret_jwt","private_key_jwt","none"],"claims_supported":["iss","ver","sub","aud","iat","exp","jti","auth_time","amr","idp","nonce","name","nickname","preferred_username","given_name","middle_name","family_name","email","email_verified","profile","zoneinfo","locale","address","phone_number","picture","website","gender","birthdate","updated_at","at_hash","c_hash"],"code_challenge_methods_supported":["S256"],"introspection_endpoint":"http://foo.example.com/oauth2/v1/introspect","introspection_endpoint_auth_methods_supported":["client_secret_basic","client_secret_post","client_secret_jwt","private_key_jwt","none"],"revocation_endpoint":"http://foo.example.com/oauth2/v1/revoke","revocation_endpoint_auth_methods_supported":["client_secret_basic","client_secret_post","client_secret_jwt","private_key_jwt","none"],"end_session_endpoint":"http://foo.example.com/oauth2/v1/logout","request_parameter_supported":true,"request_object_signing_alg_values_supported":["HS256","HS384","HS512","RS256","RS384","RS512","ES256","ES384","ES512"]}`

func NewMockOIDCProviderServer(email string) *MockOIDCProviderServer {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}

	m := &MockOIDCProviderServer{
		key:   key,
		email: email,
	}
	m.srv = httptest.NewServer(http.HandlerFunc(m.handle))
	m.client = m.srv.Client()
	m.client.Transport = &http.Transport{DialContext: func(_ context.Context, network, _ string) (net.Conn, error) {
		// Redirect all requests to the httptest server regardless of host.
		return net.Dial(network, m.srv.Listener.Addr().String())
	}}

	return m
}

type MockOIDCProviderServer struct {
	key    *rsa.PrivateKey
	srv    *httptest.Server
	client *http.Client

	email string
}

type testIdTokenClaims struct {
	*jwt.StandardClaims
	Email string `json:"email"`
}

func (m *MockOIDCProviderServer) Close() {
	m.srv.Close()
}

func (m *MockOIDCProviderServer) Client() *http.Client {
	return m.client
}

func (m *MockOIDCProviderServer) handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	switch r.URL.Path {
	case "/.well-known/openid-configuration":
		fmt.Fprintln(w, openIDConfiguration)
	case "/oauth2/v1/token":
		claims := &testIdTokenClaims{
			StandardClaims: &jwt.StandardClaims{
				Issuer:    "http://foo.example.com",
				Audience:  "my_client_id",
				ExpiresAt: math.MaxInt32,
			},
			Email: m.email,
		}

		tok, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(m.key)
		if err != nil {
			panic(err)
		}

		fmt.Fprintf(w, `{"token_type":"bearer","access_token":"AAAAAAAAAAAA", "id_token": "%s"}`, tok)
	case "/oauth2/v1/keys":
		jwk := jose.JSONWebKey{KeyID: "foo", Key: m.key.Public()}
		jwks := jose.JSONWebKeySet{Keys: []jose.JSONWebKey{jwk}}
		jks, _ := json.Marshal(jwks)
		_, _ = w.Write(jks)
	default:
		panic(fmt.Sprintf("mock received unknown URL '%s'", r.URL))
	}
}
