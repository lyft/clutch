package authn

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/durationpb"

	authnv1 "github.com/lyft/clutch/backend/api/authn/v1"
	"github.com/lyft/clutch/backend/gateway/meta"
	"github.com/lyft/clutch/backend/mock/service/authnmock"
	"github.com/lyft/clutch/backend/service/authn"
)

func TestAllRequestResponseAreRedacted(t *testing.T) {
	obj := []proto.Message{
		&authnv1.CallbackRequest{},
		&authnv1.CallbackResponse{},
		&authnv1.LoginRequest{},
		&authnv1.LoginResponse{},
		&authnv1.CreateTokenRequest{},
		&authnv1.CreateTokenResponse{},
	}

	for _, o := range obj {
		assert.True(t, meta.IsRedacted(o))
	}
}

func TestAPILogin(t *testing.T) {
	api := api{
		provider: MockProvider{},
		issuer:   authnmock.MockIssuer{},
	}

	response, err := api.Login(context.Background(), &authnv1.LoginRequest{
		RedirectUrl: "foo.com",
	})
	assert.NoError(t, err)

	assert.Equal(t, "https://auth.com/?param=nonce-foo.com", response.AuthUrl)
}

func TestAPICallback(t *testing.T) {
	api := api{
		provider: MockProvider{},
		issuer:   authnmock.MockIssuer{},
	}

	response, err := api.Callback(context.Background(), &authnv1.CallbackRequest{
		Error:            "error",
		ErrorDescription: "description",
	})
	assert.Nil(t, response)
	assert.Error(t, err, "error: description")

	transportStream := &authnmock.MockServerTransportStream{}
	ctx := grpc.NewContextWithServerTransportStream(context.Background(), transportStream)
	response, err = api.Callback(ctx, &authnv1.CallbackRequest{
		State: "nonce-foo.com",
		Code:  "whatever",
	})
	assert.NoError(t, err)
	assert.Equal(t, response.AccessToken, "test-access")
	assert.Len(t, transportStream.SentHeaders, 1)
	assert.Equal(t, metadata.MD{"location": []string{"foo.com"}, "location-status": []string{"303"}, "set-cookie-token": []string{"test-access"}}, transportStream.SentHeaders[0])
}

func TestAPICreateToken(t *testing.T) {
	api := api{
		provider: MockProvider{},
		issuer:   authnmock.MockIssuer{},
	}

	response, err := api.CreateToken(context.Background(), &authnv1.CreateTokenRequest{
		Subject: "name",
	})
	assert.Nil(t, response)
	assert.Error(t, err, "invalid token type")

	response, err = api.CreateToken(context.Background(), &authnv1.CreateTokenRequest{
		Subject:   "name",
		TokenType: authnv1.CreateTokenRequest_SERVICE,
	})
	assert.NoError(t, err)
	assert.Equal(t, response.AccessToken, "service:name_token-without-expiry")

	response, err = api.CreateToken(context.Background(), &authnv1.CreateTokenRequest{
		Subject:   "name",
		Expiry:    &durationpb.Duration{Seconds: 1},
		TokenType: authnv1.CreateTokenRequest_SERVICE,
	})
	assert.NoError(t, err)
	assert.Equal(t, response.AccessToken, "service:name_token-with-expiry")
}

// TODO(snowp): Ideally this should be in the authmocks package, but this introduces a circular dep
// due to the dependency on authn.Claims.
type MockProvider struct {
	MockClaims authn.Claims
}

func (MockProvider) GetStateNonce(redirectURL string) (string, error) {
	return strings.Join([]string{"nonce", redirectURL}, "-"), nil
}

func (MockProvider) ValidateStateNonce(state string) (redirectURL string, err error) {
	parts := strings.Split(state, "-")
	if len(parts) != 2 {
		return "", errors.New("invalid nonce formace")
	}

	return parts[1], nil
}

func (m MockProvider) Verify(ctx context.Context, rawIDToken string) (*authn.Claims, error) {
	return &m.MockClaims, nil
}

func (MockProvider) GetAuthCodeURL(ctx context.Context, state string) (string, error) {
	return "https://auth.com/?param=" + state, nil
}

func (MockProvider) Exchange(ctx context.Context, code string) (token *oauth2.Token, err error) {
	return &oauth2.Token{
		AccessToken: "test-access",
	}, err
}
