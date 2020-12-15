package k8s

import (
	"context"

	k8sapiv1 "github.com/lyft/clutch/backend/api/k8s/v1"
)

func (a *k8sAPI) DescribeCronJob(ctx context.Context, req *k8sapiv1.DescribeCronJobRequest) (*k8sapiv1.DescribeCronJobResponse, error) {
	cronJob, err := a.k8s.DescribeCronJob(ctx, req.Clientset, req.Cluster, req.Namespace, req.Name)
	if err != nil {
		return nil, err
	}
	return &k8sapiv1.DescribeCronJobResponse{Cronjob: cronJob}, nil
}

func (a *k8sAPI) DeleteCronJob(ctx context.Context, req *k8sapiv1.DeleteCronJobRequest) (*k8sapiv1.DeleteCronJobResponse, error) {
	err := a.k8s.DeleteCronJob(ctx, req.Clientset, req.Cluster, req.Namespace, req.Name)
	if err != nil {
		return nil, err
	}
	return &k8sapiv1.DeleteCronJobResponse{}, nil
}
