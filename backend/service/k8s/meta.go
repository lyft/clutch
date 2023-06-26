package k8s

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"

	k8sapiv1 "github.com/lyft/clutch/backend/api/k8s/v1"
)

// TODO (mcutalo): make this generic when adding additional `List` functionality
// for all k8s resources and support field selectors
func ApplyListOptions(listOpts *k8sapiv1.ListOptions) (metav1.ListOptions, error) {
	opts := metav1.ListOptions{}
	if len(listOpts.Labels) > 0 {
		opts.LabelSelector = labels.FormatLabels(listOpts.Labels)
	}
	// use the selector string as an addition
	// Example: "!abc"
	// Another example: "a=b,!c,d!=e"
	if len(listOpts.SupplementalSelectorString) > 0 {
		// If we already got a string from the labels, we need to add a comma
		if len(opts.LabelSelector) > 0 {
			opts.LabelSelector += ","
		}
		opts.LabelSelector += listOpts.SupplementalSelectorString
	}
	// Parse() validates the selector string, and will give an err if it fails
	_, err := labels.Parse(opts.LabelSelector)
	return opts, err
}
