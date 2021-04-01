package authnmock

import (
	"context"

	"time"

	"golang.org/x/oauth2"
)

type MockIssuer struct{
}

func (MockIssuer) CreateToken(ctx context.Context, subject string, expiry *time.Duration) (token *oauth2.Token, err error) {
	if expiry == nil {
		return &oauth2.Token{
			AccessToken: subject+"_token-without-expiry",
		}, nil
	}

	return &oauth2.Token{
		AccessToken: subject+"_token-with-expiry",
		Expiry:      time.Now().Add(*expiry),
	}, nil
}
