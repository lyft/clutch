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

func TestErrorToStatus(t *testing.T) {
	var err error
	err = k8serrors.NewUnauthorized("nice try")
	err = ConvertError(err)
	spb, ok := status.FromError(err)
	assert.True(t, ok)
	assert.NotNil(t, spb)
	assert.Equal(t, "nice try", spb.Message())
	assert.EqualValues(t, codes.Unauthenticated, spb.Code())
}

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
