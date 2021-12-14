package k8s

import (
	"context"

	k8sapiv1 "github.com/lyft/clutch/backend/api/k8s/v1"
)

func (a *k8sAPI) DescribeNode(ctx context.Context, req *k8sapiv1.DescribeNodeRequest) (*k8sapiv1.DescribeNodeResponse, error) {
	node, err := a.k8s.DescribeNode(ctx, req.Clientset, req.Cluster, req.Name)
	if err != nil {
		return nil, err
	}
	return &k8sapiv1.DescribeNodeResponse{Node: node}, nil
}

func (a *k8sAPI) UpdateNode(ctx context.Context, req *k8sapiv1.UpdateNodeRequest) (*k8sapiv1.UpdateNodeResponse, error) {
	err := a.k8s.UpdateNode(ctx, req.Clientset, req.Cluster, req.Name, req.Unschedulable)
	if err != nil {
		return nil, err
	}
	return &k8sapiv1.UpdateNodeResponse{}, nil
}
