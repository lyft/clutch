package k8s

import (
	"fmt"
	"strings"
)

type mismatchedAnnotation struct {
	Annotation    string
	ExpectedValue string
	CurrentValue  string
}

type AnnotationsPreconditionMismatchError struct {
	// List of annotations whose values were not as expected
	MismatchedAnnotations []*mismatchedAnnotation
}

func (p *AnnotationsPreconditionMismatchError) Error() string {
	var ret strings.Builder
	fmt.Fprintf(&ret, "annotation(s) values mismatched in precondition check")
	for _, mismatchedAnnotation := range p.MismatchedAnnotations {
		fmt.Fprintf(&ret,
			"\nannotation: %s; expectedValue: %s; currentValue: %s",
			mismatchedAnnotation.Annotation,
			mismatchedAnnotation.ExpectedValue,
			mismatchedAnnotation.CurrentValue)
	}
	return ret.String()
}
