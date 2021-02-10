package k8s

import (
	"context"

	k8sapiv1 "github.com/lyft/clutch/backend/api/k8s/v1"
)

func (a *k8sAPI) DescribeDeployment(ctx context.Context, req *k8sapiv1.DescribeDeploymentRequest) (*k8sapiv1.DescribeDeploymentResponse, error) {
	deployment, err := a.k8s.DescribeDeployment(ctx, req.Clientset, req.Cluster, req.Namespace, req.Name)
	if err != nil {
		return nil, err
	}
	return &k8sapiv1.DescribeDeploymentResponse{Deployment: deployment}, nil
}

func (a *k8sAPI) ListDeployments(ctx context.Context, req *k8sapiv1.ListDeploymentsRequest) (*k8sapiv1.ListDeploymentsResponse, error) {
	deployments, err := a.k8s.ListDeployments(ctx, req.Clientset, req.Cluster, req.Namespace, req.Options)
	if err != nil {
		return nil, err
	}
	return &k8sapiv1.ListDeploymentsResponse{Deployments: deployments}, nil
}

func (a *k8sAPI) UpdateDeployment(ctx context.Context, req *k8sapiv1.UpdateDeploymentRequest) (*k8sapiv1.UpdateDeploymentResponse, error) {
	err := a.k8s.UpdateDeployment(ctx, req.Clientset, req.Cluster, req.Namespace, req.Name, req.Fields)
	if err != nil {
		return nil, err
	}
	return &k8sapiv1.UpdateDeploymentResponse{}, nil
}

func (a *k8sAPI) DeleteDeployment(ctx context.Context, req *k8sapiv1.DeleteDeploymentRequest) (*k8sapiv1.DeleteDeploymentResponse, error) {
	err := a.k8s.DeleteDeployment(ctx, req.Clientset, req.Cluster, req.Namespace, req.Name)
	if err != nil {
		return nil, err
	}
	return &k8sapiv1.DeleteDeploymentResponse{}, nil
}
