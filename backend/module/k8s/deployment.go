package k8s

import (
	"context"

	k8sapiv1 "github.com/lyft/clutch/backend/api/k8s/v1"
)

func (a *k8sAPI) UpdateDeployment(ctx context.Context, req *k8sapiv1.UpdateDeploymentRequest) (*k8sapiv1.UpdateDeploymentResponse, error) {
	err := a.k8s.UpdateDeployment(ctx, req.Clientset, req.Cluster, req.Namespace, req.Name, req.Fields)
	if err != nil {
		return nil, err
	}
	return &k8sapiv1.UpdateDeploymentResponse{}, nil
}
