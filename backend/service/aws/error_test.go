package aws

import (
	"errors"
	"net/http"
	"testing"

	awshttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
	"github.com/aws/smithy-go"
	smithyhttp "github.com/aws/smithy-go/transport/http"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/lyft/clutch/backend/middleware/errorintercept"
)

func TestImplementsInterceptorInterface(t *testing.T) {
	assert.Implements(t, (*errorintercept.Interceptor)(nil), (*client)(nil))
}

// Helper function to create ResponseError for testing.
func newResponseError(statusCode int, wrapped error) error {
	return &awshttp.ResponseError{
		ResponseError: &smithyhttp.ResponseError{
			Response: &smithyhttp.Response{
				Response: &http.Response{
					StatusCode: statusCode,
				},
			},
			Err: wrapped,
		},
		RequestID: "amzrequestid",
	}
}

func TestConvertError(t *testing.T) {
	re := newResponseError(401, &smithy.GenericAPIError{Code: "whoopsie", Message: "bad things happened"})

	e := ConvertError(re)

	s, ok := status.FromError(e)
	assert.True(t, ok)
	assert.Equal(t, codes.Unauthenticated, s.Code())
}

func TestConvertErrorNoEmbeddedAPIError(t *testing.T) {
	origErr := errors.New("something went wrong")

	re := newResponseError(404, origErr)
	e := ConvertError(re)

	s, ok := status.FromError(e)
	assert.True(t, ok)
	assert.Equal(t, codes.NotFound, s.Code())
	assert.Equal(t, s.Message(), re.Error())
}

func TestConvertErrorNotAWSError(t *testing.T) {
	e := errors.New("whoopsie")
	s := ConvertError(e)
	_, ok := status.FromError(s)
	assert.False(t, ok)
	assert.Equal(t, e, s)
}
