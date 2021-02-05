package k8s

import (
	"context"

	k8sapiv1 "github.com/lyft/clutch/backend/api/k8s/v1"
)

func (a *k8sAPI) DescribeStatefulSet(ctx context.Context, req *k8sapiv1.DescribeStatefulSetRequest) (*k8sapiv1.DescribeStatefulSetResponse, error) {
	set, err := a.k8s.DescribeStatefulSet(ctx, req.Clientset, req.Cluster, req.Namespace, req.Name)
	if err != nil {
		return nil, err
	}
	return &k8sapiv1.DescribeStatefulSetResponse{StatefulSet: set}, nil
}

func (a *k8sAPI) UpdateStatefulSet(ctx context.Context, req *k8sapiv1.UpdateStatefulSetRequest) (*k8sapiv1.UpdateStatefulSetResponse, error) {
	err := a.k8s.UpdateStatefulSet(ctx, req.Clientset, req.Cluster, req.Namespace, req.Name, req.Fields)
	if err != nil {
		return nil, err
	}
	return &k8sapiv1.UpdateStatefulSetResponse{}, nil
}

func (a *k8sAPI) DeleteStatefulSet(ctx context.Context, req *k8sapiv1.DeleteStatefulSetRequest) (*k8sapiv1.DeleteStatefulSetResponse, error) {
	err := a.k8s.DeleteStatefulSet(ctx, req.Clientset, req.Cluster, req.Namespace, req.Name)
	if err != nil {
		return nil, err
	}
	return &k8sapiv1.DeleteStatefulSetResponse{}, nil
}
