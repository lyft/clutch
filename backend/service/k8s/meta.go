package k8s

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"

	k8sapiv1 "github.com/lyft/clutch/backend/api/k8s/v1"
)

// TODO (mcutalo): make this generic when adding additional `List` functionality
// TODO (cpuri): Added null check for safely passing in empty ListOptions, not sure if this is the best practice...
// for all k8s resources and support field selectors
func ApplyListOptions(listOpts *k8sapiv1.ListOptions) metav1.ListOptions {
	opts := metav1.ListOptions{}
	if listOpts != nil && len(listOpts.Labels) > 0 {
		opts.LabelSelector = labels.FormatLabels(listOpts.Labels)
	}
	return opts
}
