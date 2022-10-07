package k8s

import (
	"context"

	k8sapiv1 "github.com/lyft/clutch/backend/api/k8s/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (s *svc) DescribeNode(ctx context.Context, clientset, cluster, name string) (*k8sapiv1.Node, error) {
	// Node is a cluster level object, so namespace should be ""
	cs, err := s.manager.GetK8sClientset(ctx, clientset, cluster, metav1.NamespaceAll)
	if err != nil {
		return nil, err
	}

	nodes, err := cs.CoreV1().Nodes().List(ctx, metav1.ListOptions{
		FieldSelector: "metadata.name=" + name,
	})

	if err != nil {
		return nil, err
	}

	if len(nodes.Items) == 1 {
		return ProtoForNode(cs.Cluster(), &nodes.Items[0]), nil
	} else if len(nodes.Items) > 1 {
		return nil, status.Error(codes.FailedPrecondition, "located multiple nodes")
	}
	return nil, status.Error(codes.NotFound, "unable to locate node")
}

func ProtoForNode(cluster string, k8snode *corev1.Node) *k8sapiv1.Node {
	clusterName := k8snode.ClusterName
	if clusterName == "" {
		clusterName = cluster
	}
	return &k8sapiv1.Node{
		Cluster:       clusterName,
		Name:          k8snode.Name,
		Unschedulable: k8snode.Spec.Unschedulable,
	}
}

func (s *svc) UpdateNode(ctx context.Context, clientset, cluster, name string, unschedulable bool) error {
	// Node is a cluster level object, so namespace should be ""
	cs, err := s.manager.GetK8sClientset(ctx, clientset, cluster, metav1.NamespaceAll)
	if err != nil {
		return err
	}

	node, err := cs.CoreV1().Nodes().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	if node.Spec.Unschedulable == unschedulable {
		// node schedulabliliy is already in desired state
		return nil
	}

	node.Spec.Unschedulable = unschedulable

	_, err = cs.CoreV1().Nodes().Update(ctx, node, metav1.UpdateOptions{})
	return err
}
