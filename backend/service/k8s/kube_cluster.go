// Clutch relies on cluster information throughout the application.
//
// Previously we relied on the `clusterName` field on objectmeta
// https://github.com/kubernetes/apimachinery/blob/2456ebdaba229616fab2161a615148884b46644b/pkg/apis/meta/v1/types.go#L266-L270
// This has since be depreacted as of https://github.com/kubernetes/kubernetes/commit/331525670b772eb8956b7f5204078c51c00aaef3
// and there was no replacement to this objectmeta field.
//
// To replace this we utilize our own label to denote which cluster the object belongs to,
// these helper functions are used to standardize the getting and setting of this field.
package k8s

import (
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
)

const clusterClutchNameLabel = "cluster.clutch.sh/name"

func GetKubeCluster(obj any) string {
	objectMeta, err := meta.Accessor(obj)
	if err != nil {
		// Callers are expected to handle nil cases
		return ""
	}

	labels := objectMeta.GetLabels()
	if cluster, ok := labels[clusterClutchNameLabel]; ok {
		return cluster
	}

	return ""
}

func ApplyClusterLabels(cluster string, obj runtime.Object) error {
	objMeta, err := meta.Accessor(obj)
	if err != nil {
		return err
	}

	objMeta.SetLabels(labels.Merge(objMeta.GetLabels(), labels.Set{
		clusterClutchNameLabel: cluster,
	}))

	return nil
}
