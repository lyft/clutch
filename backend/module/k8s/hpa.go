package k8s

import (
	"context"

	k8sapiv1 "github.com/lyft/clutch/backend/api/k8s/v1"
)

func (a *k8sAPI) ResizeHPA(ctx context.Context, req *k8sapiv1.ResizeHPARequest) (*k8sapiv1.ResizeHPAResponse, error) {
	err := a.k8s.ResizeHPA(ctx, req.Clientset, req.Cluster, req.Namespace, req.Name, req.Sizing)
	if err != nil {
		return nil, err
	}
	return &k8sapiv1.ResizeHPAResponse{}, nil
}

func (a *k8sAPI) DeleteHPA(ctx context.Context, req *k8sapiv1.DeleteHPARequest) (*k8sapiv1.DeleteHPAResponse, error) {
	err := a.k8s.DeleteHPA(ctx, req.Clientset, req.Cluster, req.Namespace, req.Name)
	if err != nil {
		return nil, err
	}
	return &k8sapiv1.DeleteHPAResponse{}, nil
}
