package statuserr

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
)

func TestAllProtoCodesMatch(t *testing.T) {
	cases := []struct {
		name   string
		code   codes.Code
		expect bool
		status []*status.Status
	}{
		{
			name:   "empty",
			code:   codes.NotFound,
			expect: false,
		},
		{
			name:   "singleMatch",
			code:   codes.Unavailable,
			status: []*status.Status{{Code: int32(codes.Unavailable)}},
			expect: true,
		},
		{
			name:   "multipleMatch",
			code:   codes.Unauthenticated,
			status: []*status.Status{{Code: int32(codes.Unauthenticated)}, {Code: int32(codes.Unauthenticated)}},
			expect: true,
		},
		{
			name:   "singleNoMatch",
			code:   codes.InvalidArgument,
			status: []*status.Status{{Code: int32(codes.Unauthenticated)}},
			expect: false,
		},
		{
			name:   "multipleNoMatch",
			code:   codes.InvalidArgument,
			status: []*status.Status{{Code: int32(codes.Unauthenticated)}, {Code: int32(codes.Unimplemented)}},
			expect: false,
		},
		{
			name:   "multipleSomeMatch",
			code:   codes.InvalidArgument,
			status: []*status.Status{{Code: int32(codes.InvalidArgument)}, {Code: int32(codes.Unimplemented)}},
			expect: false,
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			assert.Equal(t, c.expect, AllProtoCodesMatch(c.code, c.status...))
		})
	}
}
