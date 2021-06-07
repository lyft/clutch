package k8s

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/retry"

	k8sapiv1 "github.com/lyft/clutch/backend/api/k8s/v1"
)

func (s *svc) DescribeStatefulSet(ctx context.Context, clientset, cluster, namespace, name string) (*k8sapiv1.StatefulSet, error) {
	cs, err := s.manager.GetK8sClientset(ctx, clientset, cluster, namespace)
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
		return nil, status.Error(codes.FailedPrecondition, "located multiple stateful sets")
	}

	return nil, status.Error(codes.NotFound, "unable to locate specified stateful set")
}

func (s *svc) ListStatefulSets(ctx context.Context, clientset, cluster, namespace string, listOptions *k8sapiv1.ListOptions) ([]*k8sapiv1.StatefulSet, error) {
	cs, err := s.manager.GetK8sClientset(ctx, clientset, cluster, namespace)
	if err != nil {
		return nil, err
	}

	opts, err := ApplyListOptions(listOptions)
	if err != nil {
		return nil, err
	}

	statefulSetList, err := cs.AppsV1().StatefulSets(cs.Namespace()).List(ctx, opts)
	if err != nil {
		return nil, err
	}

	var statefulSets []*k8sapiv1.StatefulSet
	for _, d := range statefulSetList.Items {
		statefulSet := d
		statefulSets = append(statefulSets, ProtoForStatefulSet(cs.Cluster(), &statefulSet))
	}

	return statefulSets, nil
}

// ProtoForStatefulSet maps a Kubernetes Stateful Set object to a k8sapiv1 object
func ProtoForStatefulSet(cluster string, statefulSet *appsv1.StatefulSet) *k8sapiv1.StatefulSet {
	clusterName := statefulSet.ClusterName
	if clusterName == "" {
		clusterName = cluster
	}
	k8sStateful := &k8sapiv1.StatefulSet{
		Cluster:     clusterName,
		Namespace:   statefulSet.Namespace,
		Name:        statefulSet.Name,
		Labels:      statefulSet.Labels,
		Annotations: statefulSet.Annotations,
		Status:      ProtoForStatus(statefulSet.Status),
	}

	if !statefulSet.CreationTimestamp.IsZero() {
		// Convert Unix Timestamp to milliseconds
		k8sStateful.CreationTimeMillis = statefulSet.CreationTimestamp.UnixNano() / 1e6
	}

	return k8sStateful
}

func ProtoForStatus(status appsv1.StatefulSetStatus) *k8sapiv1.StatefulSet_Status {
	return &k8sapiv1.StatefulSet_Status{
		Replicas:        uint32(status.Replicas),
		UpdatedReplicas: uint32(status.UpdatedReplicas),
		ReadyReplicas:   uint32(status.ReadyReplicas),
	}
}

func (s *svc) UpdateStatefulSet(ctx context.Context, clientset, cluster, namespace, name string, fields *k8sapiv1.UpdateStatefulSetRequest_Fields) error {
	cs, err := s.manager.GetK8sClientset(ctx, clientset, cluster, namespace)
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
		statefulSet.Labels = labels.Merge(labels.Set(statefulSet.Labels), labels.Set(fields.Labels))
		statefulSet.Spec.Template.ObjectMeta.Labels = labels.Merge(labels.Set(statefulSet.Spec.Template.ObjectMeta.Labels), labels.Set(fields.Labels))
	}

	if len(fields.Annotations) > 0 {
		statefulSet.Annotations = labels.Merge(labels.Set(statefulSet.Annotations), labels.Set(fields.Annotations))
		statefulSet.Spec.Template.ObjectMeta.Annotations = labels.Merge(labels.Set(statefulSet.Spec.Template.ObjectMeta.Annotations), labels.Set(fields.Annotations))
	}
}

func (s *svc) DeleteStatefulSet(ctx context.Context, clientset, cluster, namespace, name string) error {
	cs, err := s.manager.GetK8sClientset(ctx, clientset, cluster, namespace)
	if err != nil {
		return err
	}

	opts := metav1.DeleteOptions{}
	return cs.AppsV1().StatefulSets(cs.Namespace()).Delete(ctx, name, opts)
}
