package k8s

import (
	"context"

	"google.golang.org/protobuf/types/known/structpb"
	v1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"

	k8sapiv1 "github.com/lyft/clutch/backend/api/k8s/v1"
)

func (s *svc) DeleteJob(ctx context.Context, clientset, cluster, namespace, name string) error {
	cs, err := s.manager.GetK8sClientset(ctx, clientset, cluster, namespace)
	if err != nil {
		return err
	}

	opts := metav1.DeleteOptions{}

	return cs.BatchV1().Jobs(cs.Namespace()).Delete(ctx, name, opts)
}

func (s *svc) CreateJob(ctx context.Context, clientset, cluster, namespace string, jobConfig *structpb.Value) (*k8sapiv1.Job, error) {
	cs, err := s.manager.GetK8sClientset(ctx, clientset, cluster, namespace)
	if err != nil {
		return nil, err
	}

	// convert Job config into a runtime object and then type cast to a
	// batch job
	config := []byte(jobConfig.GetStringValue())
	decode := scheme.Codecs.UniversalDeserializer().Decode
	obj, _, err := decode(config, nil, nil)
	if err != nil {
		return nil, err
	}
	job := obj.(*v1.Job)
	opts := metav1.CreateOptions{}

	resultJob, err := cs.BatchV1().Jobs(cs.Namespace()).Create(ctx, job, opts)
	if err != nil {
		return nil, err
	}

	return protoForJob(cs.Cluster(), resultJob), nil
}

func (s *svc) ListJobs(ctx context.Context, clientset, cluster, namespace string, listOptions *k8sapiv1.ListOptions) ([]*k8sapiv1.Job, error) {
	cs, err := s.manager.GetK8sClientset(ctx, clientset, cluster, namespace)
	if err != nil {
		return nil, err
	}

	opts := ApplyListOptions(listOptions)

	jobList, err := cs.BatchV1().Jobs(cs.Namespace()).List(ctx, opts)
	if err != nil {
		return nil, err
	}

	var jobs []*k8sapiv1.Job
	for _, j := range jobList.Items {
		job := j
		jobs = append(jobs, protoForJob(cs.Cluster(), &job))
	}

	return jobs, nil
}

func protoForJob(cluster string, k8sJob *v1.Job) *k8sapiv1.Job {
	clusterName := k8sJob.ClusterName
	if clusterName == "" {
		clusterName = cluster
	}

	return &k8sapiv1.Job{
		Cluster:     clusterName,
		Namespace:   k8sJob.Namespace,
		Name:        k8sJob.Name,
		Labels:      k8sJob.Labels,
		Annotations: k8sJob.Annotations,
	}
}
