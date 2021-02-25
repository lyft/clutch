package k8s

import (
	"context"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/golang/protobuf/ptypes/wrappers"
	v1beta1 "k8s.io/api/batch/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	k8sapiv1 "github.com/lyft/clutch/backend/api/k8s/v1"
)

func (s *svc) DescribeCronJob(ctx context.Context, clientset, cluster, namespace, name string) (*k8sapiv1.CronJob, error) {
	cs, err := s.manager.GetK8sClientset(ctx, clientset, cluster, namespace)
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
		return nil, status.Error(codes.FailedPrecondition, "located multiple cron jobs")
	}
	return nil, status.Error(codes.NotFound, "unable to locate specified cron job")
}

func (s *svc) ListCronJobs(ctx context.Context, clientset, cluster, namespace string, listOptions *k8sapiv1.ListOptions) ([]*k8sapiv1.CronJob, error) {
	cs, err := s.manager.GetK8sClientset(ctx, clientset, cluster, namespace)
	if err != nil {
		return nil, err
	}

	opts := ApplyListOptions(listOptions)
	cronJobList, err := cs.BatchV1beta1().CronJobs(cs.Namespace()).List(ctx, opts)
	if err != nil {
		return nil, err
	}

	var cronJobs []*k8sapiv1.CronJob
	for _, d := range cronJobList.Items {
		cronJob := d
		cronJobs = append(cronJobs, ProtoForCronJob(cs.Cluster(), &cronJob))
	}

	return cronJobs, nil
}

func (s *svc) DeleteCronJob(ctx context.Context, clientset, cluster, namespace, name string) error {
	cs, err := s.manager.GetK8sClientset(ctx, clientset, cluster, namespace)
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
	// Required fields
	ret := &k8sapiv1.CronJob{
		Cluster:     clusterName,
		Namespace:   k8scronJob.Namespace,
		Name:        k8scronJob.Name,
		Schedule:    k8scronJob.Spec.Schedule,
		Labels:      k8scronJob.Labels,
		Annotations: k8scronJob.Annotations,
	}

	// Update optional fields
	if k8scronJob.Spec.Suspend != nil {
		ret.Suspend = *k8scronJob.Spec.Suspend
	}
	if k8scronJob.Spec.ConcurrencyPolicy != "" {
		ret.ConcurrencyPolicy = k8sapiv1.CronJob_ConcurrencyPolicy(
			k8sapiv1.CronJob_ConcurrencyPolicy_value[strings.ToUpper(string(k8scronJob.Spec.ConcurrencyPolicy))])
	}
	if k8scronJob.Status.Active != nil {
		ret.NumActiveJobs = int32(len(k8scronJob.Status.Active))
	}
	if k8scronJob.Spec.StartingDeadlineSeconds != nil {
		ret.StartingDeadlineSeconds = &wrappers.Int64Value{Value: *k8scronJob.Spec.StartingDeadlineSeconds}
	}
	return ret
}
