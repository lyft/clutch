package k8s

import (
	"context"

	k8sapiv1 "github.com/lyft/clutch/backend/api/k8s/v1"
)

func (a *k8sAPI) DescribeService(ctx context.Context, req *k8sapiv1.DescribeServiceRequest) (*k8sapiv1.DescribeServiceResponse, error) {
	svc, err := a.k8s.DescribeService(ctx, req.Clientset, req.Cluster, req.Namespace, req.Name)
	if err != nil {
		return nil, err
	}
	return &k8sapiv1.DescribeServiceResponse{Svc: svc}, nil
}

func (a *k8sAPI) ListServices(ctx context.Context, req *k8sapiv1.ListServicesRequest) (*k8sapiv1.ListServicesResponse, error) {
	svcList, err := a.k8s.ListServices(ctx, req.Clientset, req.Cluster, req.Namespace, req.Options)
	if err != nil {
		return nil, err
	}
	return &k8sapiv1.ListServicesResponse{Svc: svcList}, nil
}
