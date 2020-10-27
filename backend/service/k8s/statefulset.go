package k8s

import (
	"context"
	"encoding/json"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/strategicpatch"
	"k8s.io/client-go/util/retry"

	k8sapiv1 "github.com/lyft/clutch/backend/api/k8s/v1"
)

func (s *svc) DescribeStatefulSet(ctx context.Context, clientset, cluster, namespace, name string) (*k8sapiv1.StatefulSet, error) {
	cs, err := s.manager.GetK8sClientset(clientset, cluster, namespace)
	if err != nil {
		return nil, err
	}

	statefulSets, err := cs.AppsV1().StatefulSets(cs.Namespace()).List(ctx, metav1.ListOptions{
		FieldSelector: "metadata.name=" + name,
	})
	if err != nil {
		return nil, err
	}

	if len(statefulSets.Items) == 1 {
		return ProtoForStatefulSet(cs.Cluster(), &statefulSets.Items[0]), nil
	} else if len(statefulSets.Items) > 1 {
		return nil, fmt.Errorf("Located multiple StatefulSets")
	}

	return nil, fmt.Errorf("Unable to locate StatefulSet")
}

func ProtoForStatefulSet(cluster string, statefulSet *appsv1.StatefulSet) *k8sapiv1.StatefulSet {
	clusterName := statefulSet.ClusterName
	if clusterName == "" {
		clusterName = cluster
	}
	return &k8sapiv1.StatefulSet{
		Cluster:     clusterName,
		Namespace:   statefulSet.Namespace,
		Name:        statefulSet.Name,
		Labels:      statefulSet.Labels,
		Annotations: statefulSet.Annotations,
	}
}

func (s *svc) UpdateStatefulSet(ctx context.Context, clientset, cluster, namespace, name string, fields *k8sapiv1.UpdateStatefulSetRequest_Fields) error {
	cs, err := s.manager.GetK8sClientset(clientset, cluster, namespace)
	if err != nil {
		return err
	}

	getOpts := metav1.GetOptions{}
	oldStatefulSet, err := cs.AppsV1().StatefulSets(cs.Namespace()).Get(ctx, name, getOpts)
	if err != nil {
		return err
	}

	newStatefulSet := oldStatefulSet.DeepCopy()
	mergeStatefulSetLabelsAndAnnotations(newStatefulSet, fields)

	patchBytes, err := GenerateStrategicPatch(oldStatefulSet, newStatefulSet, appsv1.StatefulSet{})
	if err != nil {
		return err
	}

	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		_, err := cs.AppsV1().StatefulSets(cs.Namespace()).Patch(ctx, oldStatefulSet.Name, types.StrategicMergePatchType, patchBytes, metav1.PatchOptions{})
		return err
	})
	return retryErr
}

func mergeStatefulSetLabelsAndAnnotations(statefulSet *appsv1.StatefulSet, fields *k8sapiv1.UpdateStatefulSetRequest_Fields) {
	if len(fields.Labels) > 0 {
		if statefulSet.Labels == nil {
			statefulSet.Labels = make(map[string]string)
		}
		for k, v := range fields.Labels {
			statefulSet.Labels[k] = v

			if statefulSet.Spec.Template.ObjectMeta.Labels == nil {
				statefulSet.Spec.Template.ObjectMeta.Labels = make(map[string]string)
			}

			statefulSet.Spec.Template.ObjectMeta.Labels[k] = v
		}
	}

	if len(fields.Annotations) > 0 {
		if statefulSet.Annotations == nil {
			statefulSet.Annotations = make(map[string]string)
		}
		for k, v := range fields.Annotations {
			statefulSet.Annotations[k] = v

			if statefulSet.Spec.Template.ObjectMeta.Annotations == nil {
				statefulSet.Spec.Template.ObjectMeta.Annotations = make(map[string]string)
			}

			statefulSet.Spec.Template.ObjectMeta.Annotations[k] = v
		}
	}
}

func generateStatefulSetStrategicPatch(oldStatefulSet, newStatefulSet *appsv1.StatefulSet) ([]byte, error) {
	old, err := json.Marshal(oldStatefulSet)
	if err != nil {
		return nil, err
	}

	new, err := json.Marshal(newStatefulSet)
	if err != nil {
		return nil, err
	}

	patchBytes, err := strategicpatch.CreateTwoWayMergePatch(old, new, appsv1.StatefulSet{})
	if err != nil {
		return nil, err
	}

	return patchBytes, nil
}
