package k8s

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/client-go/kubernetes/scheme"

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

func (a *k8sAPI) CreateJob(ctx context.Context, req *k8sapiv1.CreateJobRequest) (*k8sapiv1.CreateJobResponse, error) {
	if req.JobConfig == nil {
		return nil, status.Error(codes.InvalidArgument, "unable to create job object, job configuration is missing")
	}
	// convert Job config into a runtime object and then type cast to a
	// batch job
	config := []byte(req.JobConfig.Value.GetStringValue())
	decode := scheme.Codecs.UniversalDeserializer().Decode
	obj, _, err := decode(config, nil, nil)
	if err != nil {
		return nil, err
	}
	job, ok := obj.(*batchv1.Job)
	if !ok {
		return nil, status.Error(codes.Internal, "unable to create Job object. Type assertion to Job failed")
	}
	result, err := a.k8s.CreateJob(ctx, req.Clientset, req.Cluster, req.Namespace, job)
	if err != nil {
		return nil, err
	}
	return &k8sapiv1.CreateJobResponse{Job: result}, nil
}
