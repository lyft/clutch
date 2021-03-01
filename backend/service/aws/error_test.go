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

func TestConvertError(t *testing.T) {
	re := &awshttp.ResponseError{
		ResponseError: &smithyhttp.ResponseError{
			Response: &smithyhttp.Response{
				Response: &http.Response{
					StatusCode: 401,
				},
			},
			Err: &smithy.GenericAPIError{Code: "whoopsie", Message: "bad things happened"},
		},
		RequestID: "amzrequestid",
	}

	e := ConvertError(re)

	s, ok := status.FromError(e)
	assert.True(t, ok)
	assert.Equal(t, codes.Unauthenticated, s.Code())
}

func TestConvertErrorNoEmbeddedAPIError(t *testing.T) {
	origErr := errors.New("something went wrong")
	re := &awshttp.ResponseError{
		ResponseError: &smithyhttp.ResponseError{
			Response: &smithyhttp.Response{
				Response: &http.Response{
					StatusCode: 401,
				},
			},
			Err: origErr,
		},
		RequestID: "amzrequestid",
	}

	e := ConvertError(re)

	s, ok := status.FromError(e)
	assert.True(t, ok)
	assert.Equal(t, codes.Unauthenticated, s.Code())
	assert.Equal(t, s.Message(), re.Error())
}

func TestConvertErrorNotAWSError(t *testing.T) {
	e := errors.New("whoopsie")
	s := ConvertError(e)
	_, ok := status.FromError(s)
	assert.False(t, ok)
	assert.Equal(t, e, s)
}
