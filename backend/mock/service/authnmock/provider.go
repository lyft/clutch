package authnmock

import (
	"context"
	"errors"
	"strings"

	"github.com/lyft/clutch/backend/service/authn"
	"golang.org/x/oauth2"
)

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
		AccessToken:  "test-access",
	}, err
}
