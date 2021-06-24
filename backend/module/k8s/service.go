package k8s

import (
	"context"

	k8sapiv1 "github.com/lyft/clutch/backend/api/k8s/v1"
)

func (a *k8sAPI) DescribeService(ctx context.Context, req *k8sapiv1.DescribeServiceRequest) (*k8sapiv1.DescribeServiceResponse, error) {
	service, err := a.k8s.DescribeService(ctx, req.Clientset, req.Cluster, req.Namespace, req.Name)
	if err != nil {
		return nil, err
	}
	return &k8sapiv1.DescribeServiceResponse{Service: service}, nil
}

func (a *k8sAPI) DeleteService(ctx context.Context, req *k8sapiv1.DeleteServiceRequest) (*k8sapiv1.DeleteServiceResponse, error) {
	err := a.k8s.DeleteService(ctx, req.Clientset, req.Cluster, req.Namespace, req.Name)
	if err != nil {
		return nil, err
	}
	return &k8sapiv1.DeleteServiceResponse{}, nil
}

func (a *k8sAPI) ListServices(ctx context.Context, req *k8sapiv1.ListServicesRequest) (*k8sapiv1.ListServicesResponse, error) {
	services, err := a.k8s.ListServices(ctx, req.Clientset, req.Cluster, req.Namespace, req.Options)
	if err != nil {
		return nil, err
	}
	return &k8sapiv1.ListServicesResponse{Services: services}, nil
}
