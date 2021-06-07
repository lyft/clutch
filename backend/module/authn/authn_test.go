package authn

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/durationpb"

	authnv1 "github.com/lyft/clutch/backend/api/authn/v1"
	"github.com/lyft/clutch/backend/gateway/meta"
	"github.com/lyft/clutch/backend/mock/grpcmock"
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

	transportStream := &grpcmock.MockServerTransportStream{}
	ctx := grpc.NewContextWithServerTransportStream(context.Background(), transportStream)

	response, err := api.Login(ctx, &authnv1.LoginRequest{
		RedirectUrl: "foo.com",
	})
	assert.NoError(t, err)

	assert.Equal(t, "https://auth.com/?param=nonce-foo.com", response.GetAuthUrl())
}

func TestAPILoginRefresh(t *testing.T) {
	api := api{
		provider: MockProvider{},
		issuer:   authnmock.MockIssuer{AllowRefresh: true},
		logger:   zaptest.NewLogger(t, zaptest.Level(zap.WarnLevel)),
	}

	transportStream := &grpcmock.MockServerTransportStream{}
	ctx := grpc.NewContextWithServerTransportStream(context.Background(), transportStream)
	ctx = metadata.NewIncomingContext(ctx, metadata.Pairs("grpcgateway-cookie", "token=zoom;refreshToken=zip"))

	response, err := api.Login(ctx, &authnv1.LoginRequest{
		RedirectUrl: "foo.com",
	})
	assert.NoError(t, err)

	assert.Equal(t, "newAccess", response.GetToken().AccessToken)
	assert.Equal(t, "refreshed", response.GetToken().RefreshToken)
}

func TestAPILoginRefreshFailsFlowContinues(t *testing.T) {
	api := api{
		provider: MockProvider{},
		issuer:   authnmock.MockIssuer{AllowRefresh: false},
		logger:   zaptest.NewLogger(t, zaptest.Level(zap.WarnLevel)),
	}

	transportStream := &grpcmock.MockServerTransportStream{}
	ctx := grpc.NewContextWithServerTransportStream(context.Background(), transportStream)
	ctx = metadata.NewIncomingContext(ctx, metadata.Pairs("grpcgateway-cookie", "token=zoom;refreshToken=zip"))

	response, err := api.Login(ctx, &authnv1.LoginRequest{
		RedirectUrl: "foo.com",
	})
	assert.NoError(t, err)

	assert.Equal(t, "https://auth.com/?param=nonce-foo.com", response.GetAuthUrl())
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

	transportStream := &grpcmock.MockServerTransportStream{}
	ctx := grpc.NewContextWithServerTransportStream(context.Background(), transportStream)
	response, err = api.Callback(ctx, &authnv1.CallbackRequest{
		State: "nonce-foo.com",
		Code:  "whatever",
	})
	assert.NoError(t, err)
	assert.Equal(t, response.AccessToken, "test-access")
	assert.Equal(t, response.RefreshToken, "")
	assert.Len(t, transportStream.SentHeaders, 1)
	assert.EqualValues(t, metadata.MD{"location": []string{"foo.com"}, "set-cookie-token": []string{"test-access"}}, transportStream.SentHeaders[0])
}

func TestAPICallbackWithRefresh(t *testing.T) {
	api := api{
		provider: MockProvider{IssueRefresh: true},
		issuer:   authnmock.MockIssuer{},
	}

	response, err := api.Callback(context.Background(), &authnv1.CallbackRequest{
		Error:            "error",
		ErrorDescription: "description",
	})
	assert.Nil(t, response)
	assert.Error(t, err, "error: description")

	transportStream := &grpcmock.MockServerTransportStream{}
	ctx := grpc.NewContextWithServerTransportStream(context.Background(), transportStream)
	response, err = api.Callback(ctx, &authnv1.CallbackRequest{
		State: "nonce-foo.com",
		Code:  "whatever",
	})
	assert.NoError(t, err)
	assert.Equal(t, response.AccessToken, "test-access")
	assert.Equal(t, response.RefreshToken, "test-refresh")
	assert.Len(t, transportStream.SentHeaders, 1)

	assert.EqualValues(t,
		metadata.MD{"location": []string{"foo.com"}, "set-cookie-token": []string{"test-access"}, "set-cookie-refresh-token": []string{"test-refresh"}},
		transportStream.SentHeaders[0])
}

func TestAPICreateToken(t *testing.T) {
	api := api{
		provider: MockProvider{},
		issuer:   authnmock.MockIssuer{},
	}

	response, err := api.CreateToken(context.Background(), &authnv1.CreateTokenRequest{
		Subject:   "name",
		TokenType: authnv1.CreateTokenRequest_SERVICE,
	})
	assert.NoError(t, err)
	assert.Equal(t, "SERVICE_name_token-without-expiry", response.AccessToken)

	response, err = api.CreateToken(context.Background(), &authnv1.CreateTokenRequest{
		Subject:   "name",
		Expiry:    &durationpb.Duration{Seconds: 1},
		TokenType: authnv1.CreateTokenRequest_SERVICE,
	})
	assert.NoError(t, err)
	assert.Equal(t, "SERVICE_name_token-with-expiry", response.AccessToken)
}

// TODO(snowp): Ideally this should be in the authmocks package, but this introduces a circular dep
// due to the dependency on authn.Claims.
type MockProvider struct {
	MockClaims   authn.Claims
	IssueRefresh bool
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

func (m MockProvider) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	t := &oauth2.Token{
		AccessToken: "test-access",
	}

	if m.IssueRefresh {
		t.RefreshToken = "test-refresh"
	}
	return t, nil
}
