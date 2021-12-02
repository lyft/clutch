package k8s

import (
	"context"
	"errors"

	k8sapiv1 "github.com/lyft/clutch/backend/api/k8s/v1"
)

func (a *k8sAPI) DescribeNode(ctx context.Context, req *k8sapiv1.DescribeNodeRequest) (*k8sapiv1.DescribeNodeResponse, error) {
	return nil, errors.New("not implemented")
}

func (a *k8sAPI) UpdateNode(ctx context.Context, req *k8sapiv1.UpdateNodeRequest) (*k8sapiv1.UpdateNodeResponse, error) {
	return nil, errors.New("not implemented")
}
