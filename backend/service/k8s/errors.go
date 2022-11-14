package k8s

import (
	"fmt"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	k8sv1 "github.com/lyft/clutch/backend/api/k8s/v1"
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

// ConvertError takes an embedded K8s error and converts it to an error w/ embedded GRPC status containing
// all of the same information.
func ConvertError(e error) error {
	as, ok := e.(k8serrors.APIStatus)
	if !ok {
		return e
	}
	s := as.Status()

	c := service.CodeFromHTTPStatus(int(s.Code))
	if c == codes.OK {
		// Shouldn't happen, but I'm sure K8s APIs can behave in mysterious ways.
		// Just return the original error if so.
		return e
	}

	ret := status.New(c, s.Message)
	ret, _ = ret.WithDetails(k8sStatusToClutchK8sStatus(&s))

	return ret.Err()
}

// K8s protos panic when serializing with Go v2 proto APIs, so we have our own analogous type to
// convert the non-compliant struct.
func k8sStatusToClutchK8sStatus(s *v1.Status) *k8sv1.Status {
	ret := &k8sv1.Status{
		Status:  s.Status,
		Message: s.Message,
		Reason:  string(s.Reason),
		Code:    s.Code,
	}

	if s.Details != nil {
		ret.Details = &k8sv1.StatusDetails{
			Name:   s.Details.Name,
			Group:  s.Details.Group,
			Kind:   s.Details.Kind,
			Uid:    string(s.Details.UID),
			Causes: make([]*k8sv1.StatusCause, 0, len(s.Details.Causes)),
		}

		for _, c := range s.Details.Causes {
			ret.Details.Causes = append(ret.Details.Causes, &k8sv1.StatusCause{
				Message: c.Message,
				Field:   c.Field,
				Type:    string(c.Type),
			})
		}
	}
	return ret
}
