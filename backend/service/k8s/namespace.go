package k8s

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	k8sapiv1 "github.com/lyft/clutch/backend/api/k8s/v1"
)

func (s *svc) DescribeNamespace(ctx context.Context, clientset, cluster, name string) (*k8sapiv1.Namespace, error) {
	cs, err := s.manager.GetK8sClientset(ctx, clientset, cluster, name)
	if err != nil {
		return nil, err
	}

	namespaces, err := cs.CoreV1().Namespaces().List(ctx, metav1.ListOptions{
		FieldSelector: "metadata.name=" + name,
	})
	if err != nil {
		return nil, err
	}

	if len(namespaces.Items) == 1 {
		return ProtoForNamespace(cs.Cluster(), &namespaces.Items[0]), nil
	} else if len(namespaces.Items) > 1 {
		return nil, fmt.Errorf("Located multiple Namespaces")
	}
	return nil, fmt.Errorf("Unable to locate namespaces")
}

func ProtoForNamespace(cluster string, k8snamespace *corev1.Namespace) *k8sapiv1.Namespace {
	clusterName := k8snamespace.ClusterName
	if clusterName == "" {
		clusterName = cluster
	}
	return &k8sapiv1.Namespace{
		Cluster:     clusterName,
		Name:        k8snamespace.Name,
		Labels:      k8snamespace.Labels,
		Annotations: k8snamespace.Annotations,
	}
}
