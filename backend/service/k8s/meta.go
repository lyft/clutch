package k8s

import (
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"

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

// Applies the name of the cluster to a kube object
// ClusterName is still not set in kube v1.20 so we are setting this manually.
// https://github.com/kubernetes/apimachinery/blob/2456ebdaba229616fab2161a615148884b46644b/pkg/apis/meta/v1/types.go#L266-L270
func ApplyClusterMetadata(cluster string, obj runtime.Object) error {
	objMeta, err := meta.Accessor(obj)
	if err != nil {
		return err
	}

	objMeta.SetClusterName(cluster)
	return nil
}
