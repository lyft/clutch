package k8s

import (
	"context"

	k8sapiv1 "github.com/lyft/clutch/backend/api/k8s/v1"
)

func (a *k8sAPI) ListEvents(ctx context.Context, req *k8sapiv1.ListEventsRequest) (*k8sapiv1.ListEventsResponse, error) {
	events, err := a.k8s.ListEvents(ctx, req.Clientset, req.Cluster, req.Namespace, req.ObjectName, req.Kind)

	if err != nil {
		return nil, err
	}
	return &k8sapiv1.ListEventsResponse{Events: events}, nil
}
