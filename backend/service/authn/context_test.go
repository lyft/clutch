package authn

import (
	"context"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
)

func TestClaimsRoundTrip(t *testing.T) {
	ctx := context.Background()

	claims := &Claims{
		StandardClaims: &jwt.StandardClaims{Subject: "foo"},
	}

	newCtx := ContextWithClaims(ctx, claims)

	cc, err := ClaimsFromContext(newCtx)
	assert.NoError(t, err)
	assert.Equal(t, "foo", cc.Subject)
}
