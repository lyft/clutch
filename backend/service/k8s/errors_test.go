package k8s

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	k8sv1 "github.com/lyft/clutch/backend/api/k8s/v1"
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
func TestConvertErrorOKEdgeCase(t *testing.T) {
	orig := &k8serrors.StatusError{ErrStatus: metav1.Status{
		Code: http.StatusOK,
	}}
	assert.Equal(t, orig, ConvertError(orig))
}

func TestAdditionalRichErrorFields(t *testing.T) {
	ks := metav1.Status{
		Status:  metav1.StatusFailure,
		Message: "whoopsie",
		Reason:  metav1.StatusReasonNotFound,
		Details: &metav1.StatusDetails{
			Name:  "name",
			Group: "grp",
			Kind:  "kind",
			UID:   "uid",
			Causes: []metav1.StatusCause{{
				Type:    "causeType",
				Message: "causeMsg",
				Field:   "causeField",
			}},
			RetryAfterSeconds: 100,
		},
		Code: http.StatusNotFound,
	}

	err := ConvertError(&k8serrors.StatusError{ErrStatus: ks})

	s, ok := status.FromError(err)
	assert.True(t, ok)
	assert.NotNil(t, s)
	assert.Len(t, s.Details(), 1)

	ret := s.Details()[0].(*k8sv1.Status)
	assert.EqualValues(t, ks.Message, ret.Message)
	assert.EqualValues(t, ks.Status, ret.Status)
	assert.EqualValues(t, ks.Reason, ret.Reason)
	assert.EqualValues(t, ks.Code, ret.Code)
	assert.EqualValues(t, ks.Details.Name, ret.Details.Name)
	assert.EqualValues(t, ks.Details.Group, ret.Details.Group)
	assert.EqualValues(t, ks.Details.Kind, ret.Details.Kind)
	assert.EqualValues(t, ks.Details.UID, ret.Details.Uid)
	assert.Len(t, ret.Details.Causes, len(ks.Details.Causes))
	assert.EqualValues(t, ks.Details.Causes[0].Message, ret.Details.Causes[0].Message)
	assert.EqualValues(t, ks.Details.Causes[0].Type, ret.Details.Causes[0].Type)
	assert.EqualValues(t, ks.Details.Causes[0].Field, ret.Details.Causes[0].Field)
}
