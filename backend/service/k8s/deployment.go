package k8s

import (
	"context"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/retry"

	k8sapiv1 "github.com/lyft/clutch/backend/api/k8s/v1"
)

func (s *svc) DescribeDeployment(ctx context.Context, clientset, cluster, namespace, name string) (*k8sapiv1.Deployment, error) {
	cs, err := s.manager.GetK8sClientset(ctx, clientset, cluster, namespace)
	if err != nil {
		return nil, err
	}

	deployments, err := cs.AppsV1().Deployments(cs.Namespace()).List(ctx, metav1.ListOptions{
		FieldSelector: "metadata.name=" + name,
	})
	if err != nil {
		return nil, err
	}

	if len(deployments.Items) == 1 {
		return ProtoForDeployment(cs.Cluster(), &deployments.Items[0]), nil
	} else if len(deployments.Items) > 1 {
		return nil, status.Error(codes.FailedPrecondition, "located multiple deployments")
	}

	return nil, status.Error(codes.NotFound, "unable to locate specified deployment")
}

func (s *svc) ListDeployments(ctx context.Context, clientset, cluster, namespace string, listOptions *k8sapiv1.ListOptions) ([]*k8sapiv1.Deployment, error) {
	cs, err := s.manager.GetK8sClientset(ctx, clientset, cluster, namespace)
	if err != nil {
		return nil, err
	}

	opts, err := ApplyListOptions(listOptions)
	if err != nil {
		return nil, err
	}

	deploymentList, err := cs.AppsV1().Deployments(cs.Namespace()).List(ctx, opts)
	if err != nil {
		return nil, err
	}

	var deployments []*k8sapiv1.Deployment
	for _, d := range deploymentList.Items {
		deployment := d
		deployments = append(deployments, ProtoForDeployment(cs.Cluster(), &deployment))
	}

	return deployments, nil
}

func ProtoForDeployment(cluster string, deployment *appsv1.Deployment) *k8sapiv1.Deployment {
	clusterName := deployment.ClusterName
	if clusterName == "" {
		clusterName = cluster
	}
	k8sDeployment := &k8sapiv1.Deployment{
		Cluster:          clusterName,
		Namespace:        deployment.Namespace,
		Name:             deployment.Name,
		Labels:           deployment.Labels,
		Annotations:      deployment.Annotations,
		DeploymentStatus: ProtoForDeploymentStatus(deployment.Status),
	}

	if !deployment.CreationTimestamp.IsZero() {
		// Convert Unix Timestamp to milliseconds
		k8sDeployment.CreationTimeMillis = deployment.CreationTimestamp.UnixNano() / 1e6
	}

	return k8sDeployment
}

func ProtoForDeploymentStatus(deploymentStatus appsv1.DeploymentStatus) *k8sapiv1.Deployment_DeploymentStatus {
	var deploymentConditions []*k8sapiv1.Deployment_DeploymentStatus_Condition
	for _, cond := range deploymentStatus.Conditions {
		var deploymentConditionType k8sapiv1.Deployment_DeploymentStatus_Condition_Type
		// TODO: Is this the preferred way of converting from one enum to another?
		if cond.Type != "" {
			deploymentConditionType = k8sapiv1.Deployment_DeploymentStatus_Condition_Type(
				k8sapiv1.Deployment_DeploymentStatus_Condition_Type_value[strings.ToUpper(string(cond.Type))])
		}
		var condStatus k8sapiv1.Deployment_DeploymentStatus_Condition_ConditionStatus
		switch cond.Status {
		case v1.ConditionTrue:
			{
				condStatus = k8sapiv1.Deployment_DeploymentStatus_Condition_CONDITION_TRUE
			}
		case v1.ConditionFalse:
			{
				condStatus = k8sapiv1.Deployment_DeploymentStatus_Condition_CONDITION_FALSE
			}
		default:
			condStatus = k8sapiv1.Deployment_DeploymentStatus_Condition_CONDITION_UNKNOWN
		}

		newCond := &k8sapiv1.Deployment_DeploymentStatus_Condition{
			Type:            deploymentConditionType,
			ConditionStatus: condStatus,
			Reason:          cond.Reason,
			Message:         cond.Message,
		}
		deploymentConditions = append(deploymentConditions, newCond)
	}

	return &k8sapiv1.Deployment_DeploymentStatus{
		Replicas:             uint32(deploymentStatus.Replicas),
		UpdatedReplicas:      uint32(deploymentStatus.UpdatedReplicas),
		ReadyReplicas:        uint32(deploymentStatus.ReadyReplicas),
		AvailableReplicas:    uint32(deploymentStatus.AvailableReplicas),
		UnavailableReplicas:  uint32(deploymentStatus.UnavailableReplicas),
		DeploymentConditions: deploymentConditions,
	}
}

func (s *svc) UpdateDeployment(ctx context.Context, clientset, cluster, namespace, name string, fields *k8sapiv1.UpdateDeploymentRequest_Fields) error {
	cs, err := s.manager.GetK8sClientset(ctx, clientset, cluster, namespace)
	if err != nil {
		return err
	}

	getOpts := metav1.GetOptions{}
	oldDeployment, err := cs.AppsV1().Deployments(cs.Namespace()).Get(ctx, name, getOpts)
	if err != nil {
		return err
	}

	newDeployment := oldDeployment.DeepCopy()
	mergeDeploymentLabelsAndAnnotations(newDeployment, fields)

	patchBytes, err := GenerateStrategicPatch(oldDeployment, newDeployment, appsv1.Deployment{})
	if err != nil {
		return err
	}

	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		_, err := cs.AppsV1().Deployments(cs.Namespace()).Patch(ctx, oldDeployment.Name, types.StrategicMergePatchType, patchBytes, metav1.PatchOptions{})
		return err
	})
	return retryErr
}

func (s *svc) DeleteDeployment(ctx context.Context, clientset, cluster, namespace, name string) error {
	cs, err := s.manager.GetK8sClientset(ctx, clientset, cluster, namespace)
	if err != nil {
		return err
	}

	opts := metav1.DeleteOptions{}
	return cs.AppsV1().Deployments(cs.Namespace()).Delete(ctx, name, opts)
}

func mergeDeploymentLabelsAndAnnotations(deployment *appsv1.Deployment, fields *k8sapiv1.UpdateDeploymentRequest_Fields) {
	if len(fields.Labels) > 0 {
		deployment.Labels = labels.Merge(labels.Set(deployment.Labels), labels.Set(fields.Labels))
		deployment.Spec.Template.ObjectMeta.Labels = labels.Merge(labels.Set(deployment.Spec.Template.ObjectMeta.Labels), labels.Set(fields.Labels))
	}

	if len(fields.Annotations) > 0 {
		deployment.Annotations = labels.Merge(labels.Set(deployment.Annotations), labels.Set(fields.Annotations))
		deployment.Spec.Template.ObjectMeta.Annotations = labels.Merge(labels.Set(deployment.Spec.Template.ObjectMeta.Annotations), labels.Set(fields.Annotations))
	}
}
