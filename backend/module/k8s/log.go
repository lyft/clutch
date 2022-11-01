package k8s

import (
	"context"
	k8sapiv1 "github.com/lyft/clutch/backend/api/k8s/v1"
)

func (a *k8sAPI) GetLogs(ctx context.Context, req *k8sapiv1.GetLogsRequest) (*k8sapiv1.GetLogsResponse, error) {
	err := a.k8s.GetLogs(ctx, req.Clientset, req.Cluster, req.Namespace, req.Name)
	if err != nil {
		return nil, err
	}
	return &k8sapiv1.GetLogsResponse{}, nil
}
