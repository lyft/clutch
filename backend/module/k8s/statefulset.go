package k8s

import (
	"context"

	k8sapiv1 "github.com/lyft/clutch/backend/api/k8s/v1"
)

func (a *k8sAPI) UpdateStatefulSet(ctx context.Context, req *k8sapiv1.UpdateStatefulSetRequest) (*k8sapiv1.UpdateStatefulSetResponse, error) {
	err := a.k8s.UpdateStatefulSet(ctx, req.Clientset, req.Cluster, req.Namespace, req.Name, req.Fields)
	if err != nil {
		return nil, err
	}
	return &k8sapiv1.UpdateStatefulSetResponse{}, nil
}
