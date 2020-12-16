package k8s

import (
	"context"
	"fmt"

	v1beta1 "k8s.io/api/batch/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	k8sapiv1 "github.com/lyft/clutch/backend/api/k8s/v1"
)

func (s *svc) DescribeCronJob(ctx context.Context, clientset, cluster, namespace, name string) (*k8sapiv1.CronJob, error) {
	cs, err := s.manager.GetK8sClientset(clientset, cluster, namespace)
	if err != nil {
		return nil, err
	}

	cronJobs, err := cs.BatchV1beta1().CronJobs(cs.Namespace()).List(ctx, metav1.ListOptions{
		FieldSelector: "metadata.name=" + name,
	})
	if err != nil {
		return nil, err
	}

	if len(cronJobs.Items) == 1 {
		return ProtoForCronJob(cs.Cluster(), &cronJobs.Items[0]), nil
	} else if len(cronJobs.Items) > 1 {
		return nil, fmt.Errorf("Located multiple CronJobs")
	}
	return nil, fmt.Errorf("Unable to locate cronJob")
}

func (s *svc) DeleteCronJob(ctx context.Context, clientset, cluster, namespace, name string) error {
	cs, err := s.manager.GetK8sClientset(clientset, cluster, namespace)
	if err != nil {
		return err
	}

	opts := metav1.DeleteOptions{}
	return cs.BatchV1beta1().CronJobs(cs.Namespace()).Delete(ctx, name, opts)
}

func ProtoForCronJob(cluster string, k8scronJob *v1beta1.CronJob) *k8sapiv1.CronJob {
	clusterName := k8scronJob.ClusterName
	if clusterName == "" {
		clusterName = cluster
	}
	return &k8sapiv1.CronJob{
		Cluster:     clusterName,
		Namespace:   k8scronJob.Namespace,
		Name:        k8scronJob.Name,
		Schedule:    k8scronJob.Spec.Schedule,
		Labels:      k8scronJob.Labels,
		Annotations: k8scronJob.Annotations,
	}
}
