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

func (a *k8sAPI) ListPods(ctx context.Context, req *k8sapiv1.ListPodsRequest) (*k8sapiv1.ListPodsResponse, error) {
	pods, err := a.k8s.ListPods(ctx, req.Clientset, req.Cluster, req.Namespace, req.Options)
	if err != nil {
		return nil, err
	}
	return &k8sapiv1.ListPodsResponse{Pods: pods}, nil
}

func (a *k8sAPI) UpdatePod(ctx context.Context, req *k8sapiv1.UpdatePodRequest) (*k8sapiv1.UpdatePodResponse, error) {
	// TODO Introduce a list of allowed annotations/labels and verify that the request satisfies it
	err := a.k8s.UpdatePod(ctx, req.Clientset, req.Cluster, req.Namespace, req.Name, req.ExpectedObjectMetaFields, req.ObjectMetaFields, req.RemoveObjectMetaFields)
	if err != nil {
		return nil, err
	}
	return &k8sapiv1.UpdatePodResponse{}, nil
}
