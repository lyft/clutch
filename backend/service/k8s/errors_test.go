package k8s

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"

	"github.com/lyft/clutch/backend/middleware/errorintercept"
)

func TestImplementsInterceptorInterface(t *testing.T) {
	assert.Implements(t, (*errorintercept.Interceptor)(nil), (*svc)(nil))
}

func TestConvertError(t *testing.T) {
	testCases := []struct {
		err          error
		expectStatus bool
		expectedCode codes.Code
	}{
		{
			err: errors.New("plain error"),
		},
		{
			err:          k8serrors.NewBadRequest("bad request"),
			expectStatus: true,
			expectedCode: codes.FailedPrecondition,
		},
		{
			err:          k8serrors.NewTimeoutError("timeout", 3),
			expectStatus: true,
			expectedCode: codes.DeadlineExceeded,
		},
		{
			err:          k8serrors.NewResourceExpired("expired"),
			expectStatus: true,
			expectedCode: codes.NotFound,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.err.Error(), func(t *testing.T) {
			result := ConvertError(tc.err)

			s, ok := status.FromError(result)
			assert.Equal(t, tc.expectStatus, ok)
			if tc.expectStatus {
				assert.Equal(t, tc.expectedCode, s.Code())
				assert.Contains(t, result.Error(), s.Message())
			} else {
				assert.Equal(t, tc.err, result)
			}
		})
	}
}
