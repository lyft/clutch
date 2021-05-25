package authnmock

import (
	"context"
	"time"

	"golang.org/x/oauth2"

	authnv1 "github.com/lyft/clutch/backend/api/authn/v1"
)

type MockIssuer struct {
}

func (MockIssuer) CreateToken(ctx context.Context, subject string, tokenType authnv1.CreateTokenRequest_TokenType, expiry *time.Duration) (token *oauth2.Token, err error) {
	if expiry == nil {
		return &oauth2.Token{
			AccessToken: authnv1.CreateTokenRequest_TokenType_name[int32(tokenType)] + "_" + subject + "_token-without-expiry",
		}, nil
	}

	return &oauth2.Token{
		AccessToken: authnv1.CreateTokenRequest_TokenType_name[int32(tokenType)] + "_" + subject + "_token-with-expiry",
		Expiry:      time.Now().Add(*expiry),
	}, nil
}
