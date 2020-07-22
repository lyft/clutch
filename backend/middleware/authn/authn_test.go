package authn

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
)

func TestGetCookieValue(t *testing.T) {
	cookie := "foo=bar;baz=bang;foo=qux"
	res, err := getCookieValue(cookie, "foo")
	assert.NoError(t, err)
	assert.Equal(t, "bar", res)

	res, err = getCookieValue(cookie, "baz")
	assert.NoError(t, err)
	assert.Equal(t, "bang", res)
}

func TestGetToken(t *testing.T) {
	tokenVal := "quux"

	tests := []struct {
		md  metadata.MD
		err bool
	}{
		{md: metadata.Pairs("authorization", "Token "+tokenVal)},
		{md: metadata.Pairs("Authorization", "Token "+tokenVal)},
		{md: metadata.Pairs("grpcgateway-cookie", "foo=bar;token="+tokenVal)},
		{md: metadata.Pairs("GRPCGateway-Cookie", "foo=bar;token="+tokenVal)},
		{md: metadata.Pairs("Authorization", tokenVal), err: true},
		{md: metadata.Pairs(), err: true},
	}

	for idx, tt := range tests {
		tt := tt
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			t.Parallel()

			result, err := getToken(tt.md)
			if tt.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tokenVal, result)
			}
		})
	}
}
