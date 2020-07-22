package k8s

import (
	"context"

	k8sapiv1 "github.com/lyft/clutch/backend/api/k8s/v1"
)

func (a *k8sAPI) DescribePod(ctx context.Context, req *k8sapiv1.DescribePodRequest) (*k8sapiv1.DescribePodResponse, error) {
	pod, err := a.k8s.DescribePod(ctx, req.Clientset, req.Cluster, req.Namespace, req.Name)
	if err != nil {
		return nil, err
	}
	return &k8sapiv1.DescribePodResponse{Pod: pod}, nil
}

func (a *k8sAPI) DeletePod(ctx context.Context, req *k8sapiv1.DeletePodRequest) (*k8sapiv1.DeletePodResponse, error) {
	err := a.k8s.DeletePod(ctx, req.Clientset, req.Cluster, req.Namespace, req.Name)
	if err != nil {
		return nil, err
	}
	return &k8sapiv1.DeletePodResponse{}, nil
}
