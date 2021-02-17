package k8s

import (
	"context"

	k8sapiv1 "github.com/lyft/clutch/backend/api/k8s/v1"
)

func (a *k8sAPI) DeleteJob(ctx context.Context, req *k8sapiv1.DeleteJobRequest) (*k8sapiv1.DeleteJobResponse, error) {
	err := a.k8s.DeleteJob(ctx, req.Clientset, req.Cluster, req.Namespace, req.Name)
	if err != nil {
		return nil, err
	}
	return &k8sapiv1.DeleteJobResponse{}, nil
}

func (a *k8sAPI) ListJobs(ctx context.Context, req *k8sapiv1.ListJobsRequest) (*k8sapiv1.ListJobsResponse, error) {
	jobs, err := a.k8s.ListJobs(ctx, req.Clientset, req.Cluster, req.Namespace, req.Options)
	if err != nil {
		return nil, err
	}
	return &k8sapiv1.ListJobsResponse{Jobs: jobs}, nil
}
