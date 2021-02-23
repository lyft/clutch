package k8s

import (
	"fmt"
	"strings"

	"google.golang.org/grpc/status"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"

	"github.com/lyft/clutch/backend/service"
)

type mismatchedAnnotation struct {
	Annotation    string
	ExpectedValue string
	CurrentValue  string
}

type ExpectedObjectMetaFieldsCheckError struct {
	// List of annotations whose values were not as expected
	MismatchedAnnotations []*mismatchedAnnotation
}

func (p *ExpectedObjectMetaFieldsCheckError) Error() string {
	var ret strings.Builder
	fmt.Fprintf(&ret, "annotation(s) values mismatched in expected object metadata fields check")
	for _, mismatchedAnnotation := range p.MismatchedAnnotations {
		fmt.Fprintf(&ret,
			"\nannotation: %s; expectedValue: %s; currentValue: %s",
			mismatchedAnnotation.Annotation,
			mismatchedAnnotation.ExpectedValue,
			mismatchedAnnotation.CurrentValue)
	}
	return ret.String()
}

func ConvertError(e error) error {
	as, ok := e.(k8serrors.APIStatus)
	if !ok {
		return e
	}
	s := as.Status()

	c := service.CodeFromHTTPStatus(int(s.Code))
	ret := status.New(c, s.Message)

	return ret.Err()
}
