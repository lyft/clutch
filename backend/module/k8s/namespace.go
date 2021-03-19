package k8s

import (
	"context"

	k8sapiv1 "github.com/lyft/clutch/backend/api/k8s/v1"
)

func (a *k8sAPI) DescribeNamespace(ctx context.Context, req *k8sapiv1.DescribeNamespaceRequest) (*k8sapiv1.DescribeNamespaceResponse, error) {
	namespace, err := a.k8s.DescribeNamespace(ctx, req.Clientset, req.Cluster, req.Name)
	if err != nil {
		return nil, err
	}
	return &k8sapiv1.DescribeNamespaceResponse{Namespace: namespace}, nil
}
