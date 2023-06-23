// helper methods for accessing k8s cluster metadata in object labels
// add detail as to why this needs to exist.

package k8s

import (
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
)

const clusterClutchNameLabel = "cluster.clutch.sh/name"

// (mikecutalo) - debate this return sig
func GetKubeCluster(obj any) string {
	ObjectMeta, err := meta.Accessor(obj)
	if err != nil {
		return ""
	}

	labels := ObjectMeta.GetLabels()
	if cluster, ok := labels[clusterClutchNameLabel]; ok {
		return cluster
	}

	return ""
}

// Applies the name of the cluster to a kube object
// ClusterName is still not set in kube v1.20 so we are setting this manually.
// https://github.com/kubernetes/apimachinery/blob/2456ebdaba229616fab2161a615148884b46644b/pkg/apis/meta/v1/types.go#L266-L270
func ApplyClusterMetadata(cluster string, obj runtime.Object) error {
	objMeta, err := meta.Accessor(obj)
	if err != nil {
		return err
	}

	objMeta.SetLabels(labels.Merge(objMeta.GetLabels(), labels.Set{
		clusterClutchNameLabel: cluster,
	}))

	return nil
}
