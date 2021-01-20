package k8s

import (
	"context"

	k8sapiv1 "github.com/lyft/clutch/backend/api/k8s/v1"
)

func (a *k8sAPI) ListConfigMaps(ctx context.Context, req *k8sapiv1.ListConfigMapsRequest) (*k8sapiv1.ListConfigMapsResponse, error) {
	configMaps, err := a.k8s.ListConfigMaps(ctx, req.Clientset, req.Cluster, req.Namespace, req.Options)
	if err != nil {
		return nil, err
	}
	return &k8sapiv1.ListConfigMapsResponse{ConfigMaps: configMaps}, nil
}
