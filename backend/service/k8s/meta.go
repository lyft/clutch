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
func ApplyListOptions(listOpts *k8sapiv1.ListOptions) metav1.ListOptions {
	opts := metav1.ListOptions{}
	if len(listOpts.Labels) > 0 {
		opts.LabelSelector = labels.FormatLabels(listOpts.Labels)
	}

	return opts
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
