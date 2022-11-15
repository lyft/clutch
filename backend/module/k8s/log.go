package k8s

import (
	"context"

	k8sapiv1 "github.com/lyft/clutch/backend/api/k8s/v1"
)

func (a *k8sAPI) GetPodLogs(ctx context.Context, req *k8sapiv1.GetPodLogsRequest) (*k8sapiv1.GetPodLogsResponse, error) {
	return a.k8s.GetPodLogs(ctx, req.Clientset, req.Cluster, req.Namespace, req.Name, req.Options)
}
