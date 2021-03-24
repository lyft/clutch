package k8s

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
	autoscalingv1 "k8s.io/api/autoscaling/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/retry"

	k8sapiv1 "github.com/lyft/clutch/backend/api/k8s/v1"
)

func (s *svc) DescribeHPA(ctx context.Context, clientset, cluster, namespace, name string) (*k8sapiv1.HPA, error) {
	cs, err := s.manager.GetK8sClientset(ctx, clientset, cluster, namespace)
	if err != nil {
		return nil, err
	}

	hpas, err := cs.AutoscalingV1().HorizontalPodAutoscalers(cs.Namespace()).List(ctx, metav1.ListOptions{
		FieldSelector: "metadata.name=" + name,
	})
	if err != nil {
		return nil, err
	}

	if len(hpas.Items) == 1 {
		return ProtoForHPA(cs.Cluster(), &hpas.Items[0]), nil
	} else if len(hpas.Items) > 1 {
		return nil, status.Error(codes.FailedPrecondition, "located multiple HPAs")
	}

	return nil, status.Error(codes.NotFound, "unable to locate specified HPA")
}

func ProtoForHPA(cluster string, autoscaler *autoscalingv1.HorizontalPodAutoscaler) *k8sapiv1.HPA {
	clusterName := autoscaler.ClusterName
	if clusterName == "" {
		clusterName = cluster
	}
	ret := &k8sapiv1.HPA{
		Cluster:   clusterName,
		Namespace: autoscaler.Namespace,
		Name:      autoscaler.Name,
		Sizing: &k8sapiv1.HPA_Sizing{
			MinReplicas:     uint32(*autoscaler.Spec.MinReplicas),
			MaxReplicas:     uint32(autoscaler.Spec.MaxReplicas),
			CurrentReplicas: uint32(autoscaler.Status.CurrentReplicas),
			DesiredReplicas: uint32(autoscaler.Status.DesiredReplicas),
		},
		Labels:      autoscaler.Labels,
		Annotations: autoscaler.Annotations,
	}

	if autoscaler.Spec.TargetCPUUtilizationPercentage != nil {
		ret.TargetCpuUtilizationPercentage = &wrapperspb.Int32Value{Value: *autoscaler.Spec.TargetCPUUtilizationPercentage}
	}

	if autoscaler.Status.CurrentCPUUtilizationPercentage != nil {
		ret.CurrentCpuUtilizationPercentage = &wrapperspb.Int32Value{Value: *autoscaler.Status.CurrentCPUUtilizationPercentage}
	}
	return ret
}

func (s *svc) ResizeHPA(ctx context.Context, clientset, cluster, namespace, name string, sizing *k8sapiv1.ResizeHPARequest_Sizing) error {
	cs, err := s.manager.GetK8sClientset(ctx, clientset, cluster, namespace)
	if err != nil {
		return err
	}

	opts := metav1.GetOptions{}
	hpa, err := cs.AutoscalingV1().HorizontalPodAutoscalers(cs.Namespace()).Get(ctx, name, opts)
	if err != nil {
		return err
	}

	normalizeHPAChanges(hpa, sizing)

	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		_, err := cs.AutoscalingV1().HorizontalPodAutoscalers(cs.Namespace()).Update(ctx, hpa, metav1.UpdateOptions{})
		return err
	})
	return retryErr
}

func normalizeHPAChanges(hpa *autoscalingv1.HorizontalPodAutoscaler, sizing *k8sapiv1.ResizeHPARequest_Sizing) {
	if sizing == nil {
		return
	}

	min := int32(sizing.Min)
	hpa.Spec.MinReplicas = &min
	hpa.Spec.MaxReplicas = int32(sizing.Max)

	if *hpa.Spec.MinReplicas > hpa.Spec.MaxReplicas {
		hpa.Spec.MaxReplicas = *hpa.Spec.MinReplicas
	}
}

func (s *svc) DeleteHPA(ctx context.Context, clientset, cluster, namespace, name string) error {
	cs, err := s.manager.GetK8sClientset(ctx, clientset, cluster, namespace)
	if err != nil {
		return err
	}

	opts := metav1.DeleteOptions{}
	return cs.AutoscalingV1().HorizontalPodAutoscalers(cs.Namespace()).Delete(ctx, name, opts)
}
