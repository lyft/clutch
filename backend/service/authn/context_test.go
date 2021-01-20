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

func TestContextWithAnonymousClaims(t *testing.T) {
	ctx := context.Background()
	ctx = ContextWithAnonymousClaims(ctx)
	cc, err := ClaimsFromContext(ctx)
	assert.NoError(t, err)
	assert.Equal(t, AnonymousSubject, cc.Subject)
}

func TestNilClaimsValueErrors(t *testing.T) {
	{
		cc, err := ClaimsFromContext(context.Background())
		assert.Nil(t, cc)
		assert.Error(t, err)
	}
	{
		cc, err := ClaimsFromContext(nil) // nolint:staticcheck
		assert.Nil(t, cc)
		assert.Error(t, err)
	}
}
