package k8s

import (
	"context"
	"encoding/json"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/strategicpatch"
	"k8s.io/client-go/util/retry"

	k8sapiv1 "github.com/lyft/clutch/backend/api/k8s/v1"
)

func (s *svc) DescribeDeployment(ctx context.Context, clientset, cluster, namespace, name string) (*k8sapiv1.Deployment, error) {
	cs, err := s.manager.GetK8sClientset(clientset, cluster, namespace)
	if err != nil {
		return nil, err
	}

	deployments, err := cs.AppsV1().Deployments(cs.Namespace()).List(metav1.ListOptions{
		FieldSelector: "metadata.name=" + name,
	})
	if err != nil {
		return nil, err
	}

	if len(deployments.Items) == 1 {
		return ProtoForDeployment(cs.Cluster(), &deployments.Items[0]), nil
	} else if len(deployments.Items) > 1 {
		return nil, fmt.Errorf("Located multiple Deployments")
	}

	return nil, fmt.Errorf("Unable to locate Deployment")
}

func ProtoForDeployment(cluster string, deployment *appsv1.Deployment) *k8sapiv1.Deployment {
	clusterName := deployment.ClusterName
	if clusterName == "" {
		clusterName = cluster
	}
	return &k8sapiv1.Deployment{
		Cluster:     clusterName,
		Namespace:   deployment.Namespace,
		Name:        deployment.Name,
		Labels:      deployment.Labels,
		Annotations: deployment.Annotations,
	}
}

func (s *svc) UpdateDeployment(ctx context.Context, clientset, cluster, namespace, name string, fields *k8sapiv1.UpdateDeploymentRequest_Fields) error {
	cs, err := s.manager.GetK8sClientset(clientset, cluster, namespace)
	if err != nil {
		return err
	}

	getOpts := metav1.GetOptions{}
	oldDeployment, err := cs.AppsV1().Deployments(cs.Namespace()).Get(name, getOpts)
	if err != nil {
		return err
	}

	newDeployment := oldDeployment.DeepCopy()
	mergeLabelsAndAnnotations(newDeployment, fields)

	patchBytes, err := generateDeploymentStrategicPatch(oldDeployment, newDeployment)
	if err != nil {
		return err
	}

	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		_, err := cs.AppsV1().Deployments(cs.Namespace()).Patch(oldDeployment.Name, types.StrategicMergePatchType, patchBytes)
		return err
	})
	return retryErr
}

func mergeLabelsAndAnnotations(deployment *appsv1.Deployment, fields *k8sapiv1.UpdateDeploymentRequest_Fields) {
	if len(fields.Labels) > 0 {
		for k, v := range fields.Labels {
			deployment.Labels[k] = v

			if deployment.Spec.Template.ObjectMeta.Labels == nil {
				deployment.Spec.Template.ObjectMeta.Labels = make(map[string]string)
			}

			deployment.Spec.Template.ObjectMeta.Labels[k] = v
		}
	}

	if len(fields.Annotations) > 0 {
		for k, v := range fields.Annotations {
			deployment.Annotations[k] = v

			if deployment.Spec.Template.ObjectMeta.Annotations == nil {
				deployment.Spec.Template.ObjectMeta.Annotations = make(map[string]string)
			}

			deployment.Spec.Template.ObjectMeta.Annotations[k] = v
		}
	}
}

func generateDeploymentStrategicPatch(oldDeployment, newDeployment *appsv1.Deployment) ([]byte, error) {
	old, err := json.Marshal(oldDeployment)
	if err != nil {
		return nil, err
	}

	new, err := json.Marshal(newDeployment)
	if err != nil {
		return nil, err
	}

	patchBytes, err := strategicpatch.CreateTwoWayMergePatch(old, new, appsv1.Deployment{})
	if err != nil {
		return nil, err
	}

	return patchBytes, nil
}
