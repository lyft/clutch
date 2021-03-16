package aws

import (
	"errors"
	"fmt"

	awshttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
	"github.com/aws/smithy-go"
	"google.golang.org/grpc/status"

	"github.com/lyft/clutch/backend/service"
)

// Use error handling techniques outlined in https://aws.github.io/aws-sdk-go-v2/docs/handling-errors/
// to extract embedded HTTP status information and error messages from AWS SDK errors and create the
// corresponding gRPC status.
func ConvertError(e error) error {
	// Extract HTTP status details.
	var re *awshttp.ResponseError
	if !errors.As(e, &re) {
		return e
	}
	code := service.CodeFromHTTPStatus(re.Response.StatusCode)

	// Extract AWS API error details.
	var msg string
	var ae smithy.APIError
	if errors.As(re.Err, &ae) {
		// e.g. InvalidAccessKeyId: The key you provided does not exist in the database.
		msg = fmt.Sprintf("%s: %s", ae.ErrorCode(), ae.ErrorMessage())
	} else {
		// Use the original error message if APIError not found, though it should never happen.
		msg = e.Error()
	}

	// TODO: unpack remaining error information, e.g. AWS request ID into error details.
	return status.Error(code, msg)
}
