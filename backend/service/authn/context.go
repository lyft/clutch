package authn

import (
	"context"
	"errors"

	"github.com/dgrijalva/jwt-go"
)

type ClaimsContextKey struct{}

const AnonymousSubject = "system:anonymous"

func ContextWithAnonymousClaims(ctx context.Context) context.Context {
	return context.WithValue(ctx, ClaimsContextKey{}, &Claims{
		StandardClaims: &jwt.StandardClaims{Subject: AnonymousSubject},
	})
}

func ContextWithClaims(ctx context.Context, claims *Claims) context.Context {
	return context.WithValue(ctx, ClaimsContextKey{}, claims)
}

func ClaimsFromContext(ctx context.Context) (*Claims, error) {
	if ctx == nil {
		return nil, errors.New("no context found when evaluating claims")
	}
	c, ok := ctx.Value(ClaimsContextKey{}).(*Claims)
	if !ok {
		return nil, errors.New("claims in context were not present or not the correct type")
	}
	return c, nil
}
