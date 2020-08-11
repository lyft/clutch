package k8s

import (
	"context"
	"errors"
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8slabels "k8s.io/apimachinery/pkg/labels"

	k8sapiv1 "github.com/lyft/clutch/backend/api/k8s/v1"
)

func (s *svc) DescribePod(ctx context.Context, clientset, cluster, namespace, name string) (*k8sapiv1.Pod, error) {
	cs, err := s.manager.GetK8sClientset(clientset, cluster, namespace)
	if err != nil {
		return nil, err
	}
	opts := metav1.GetOptions{}
	pod, err := cs.CoreV1().Pods(cs.Namespace()).Get(name, opts)
	if err != nil {
		return nil, err
	}
	return podDescription(pod, cs.Cluster()), nil
}

func (s *svc) DeletePod(ctx context.Context, clientset, cluster, namespace, name string) error {
	cs, err := s.manager.GetK8sClientset(clientset, cluster, namespace)
	if err != nil {
		return err
	}

	var gracePeriod int64
	opts := &metav1.DeleteOptions{}
	opts.GracePeriodSeconds = &gracePeriod

	return cs.CoreV1().Pods(cs.Namespace()).Delete(name, opts)
}

func (s *svc) ListPods(ctx context.Context, clientset, cluster, namespace string, listPodsOpts *k8sapiv1.ListPodsOptions) ([]*k8sapiv1.Pod, error) {
	cs, err := s.manager.GetK8sClientset(clientset, cluster, namespace)
	if err != nil {
		return nil, err
	}

	opts := metav1.ListOptions{}
	if len(listPodsOpts.Labels) > 0 {
		opts.LabelSelector = k8slabels.FormatLabels(listPodsOpts.Labels)
	}

	podList, err := cs.CoreV1().Pods(cs.Namespace()).List(opts)
	if err != nil {
		return nil, err
	}

	var pods []*k8sapiv1.Pod
	for _, p := range podList.Items {
		pod := p
		pods = append(pods, podDescription(&pod, cs.Cluster()))
	}

	return pods, nil
}

// Update pod fields if the current field values match with that's described by expectedObjectMetaFields
//
// TODO: add support for updating pod labels
func (s *svc) UpdatePod(ctx context.Context, clientset, cluster, namespace, name string, expectedObjectMetaFields *k8sapiv1.ExpectedObjectMetaFields, objectMetaFields *k8sapiv1.ObjectMetaFields, removeObjectMetaFields *k8sapiv1.RemoveObjectMetaFields) error {
	if len(objectMetaFields.GetLabels()) > 0 || len(removeObjectMetaFields.GetLabels()) > 0 {
		return errors.New("update of pod labels not implemented")
	}

	// Ensure that the caller is not trying to delete an annotation and update it at the same time
	newAnnotations := objectMetaFields.GetAnnotations()
	for _, annotation := range removeObjectMetaFields.GetAnnotations() {
		_, annotationIsUpdated := newAnnotations[annotation]
		if annotationIsUpdated {
			return fmt.Errorf("annotation '%s' can't be updated and removed at once", annotation)
		}
	}

	cs, err := s.manager.GetK8sClientset(clientset, cluster, namespace)
	if err != nil {
		return err
	}

	pod, err := cs.CoreV1().Pods(cs.Namespace()).Get(name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	// Check that the current state of the pod matches with expectedObjectMetaFields.
	//
	// If there is a mismatch, checkExpectedObjectMetaFields() will return an error with the list of mismatches.
	err = checkExpectedObjectMetaFields(expectedObjectMetaFields, pod)
	if err != nil {
		return err
	}

	// Update/add annotations
	podAnnotations := pod.GetAnnotations()
	for annotation, value := range newAnnotations {
		podAnnotations[annotation] = value
	}

	// Delete annotations to be removed
	for _, annotation := range removeObjectMetaFields.GetAnnotations() {
		delete(podAnnotations, annotation)
	}

	_, err = cs.CoreV1().Pods(cs.Namespace()).Update(pod)
	return err
}

func checkExpectedObjectMetaFields(expectedObjectMetaFields *k8sapiv1.ExpectedObjectMetaFields, pod *corev1.Pod) error {
	if len(expectedObjectMetaFields.Labels) > 0 {
		return errors.New("checking for pod labels not implemented")
	}

	podAnnotations := pod.GetAnnotations()
	var mismatchedAnnotations []*mismatchedAnnotation

	for expectedAnnotation, expectedValue := range expectedObjectMetaFields.GetAnnotations() {
		// "" is a valid annotation value, so nil is used to indicate that the
		// annotation shouldn't be set
		annotationShouldBePresent := expectedValue != nil
		currentValue, annotationIsPresent := podAnnotations[expectedAnnotation]

		// Existance precondition not met
		if annotationShouldBePresent != annotationIsPresent {
			mismatchedAnnotations = append(
				mismatchedAnnotations,
				&mismatchedAnnotation{
					Annotation:    expectedAnnotation,
					ExpectedValue: expectedValue.GetValue(),
					CurrentValue:  currentValue,
				},
			)

			continue
		}

		// Annotation values mismatched
		if expectedValue.GetValue() != currentValue {
			mismatchedAnnotations = append(
				mismatchedAnnotations,
				&mismatchedAnnotation{
					Annotation:    expectedAnnotation,
					ExpectedValue: expectedValue.GetValue(),
					CurrentValue:  currentValue,
				},
			)
		}
	}

	if len(mismatchedAnnotations) == 0 {
		return nil
	}

	return &ExpectedObjectMetaFieldsCheckError{MismatchedAnnotations: mismatchedAnnotations}
}

func podDescription(k8spod *corev1.Pod, cluster string) *k8sapiv1.Pod {
	// TODO: There's a mismatch between the serialization of the timestamp here and what's expected
	// on the frontend.
	//var launch *timestamp.Timestamp
	//if converted, err := ptypes.TimestampProto(k8spod.Status.StartTime.Time); err == nil {
	//	launch = converted
	//}
	clusterName := k8spod.ClusterName
	if clusterName == "" {
		clusterName = cluster
	}
	return &k8sapiv1.Pod{
		Cluster:    clusterName,
		Namespace:  k8spod.Namespace,
		Name:       k8spod.Name,
		Containers: makeContainers(k8spod.Status.ContainerStatuses),
		NodeIp:     k8spod.Status.HostIP,
		PodIp:      k8spod.Status.PodIP,
		State:      protoForPodState(k8spod.Status.Phase),
		//StartTime:   launch,
		Labels:      k8spod.Labels,
		Annotations: k8spod.Annotations,
	}
}

func makeContainers(statuses []corev1.ContainerStatus) []*k8sapiv1.Container {
	containers := make([]*k8sapiv1.Container, 0, len(statuses))
	for _, status := range statuses {
		container := &k8sapiv1.Container{
			Name:  status.Name,
			Image: status.Image,
			State: protoForContainerState(status.State),
			Ready: status.Ready,
		}
		containers = append(containers, container)
	}
	return containers
}

func protoForPodState(state corev1.PodPhase) k8sapiv1.Pod_State {
	// Look up value in generated enum map after ensuring consistent case with generated code.
	val, ok := k8sapiv1.Pod_State_value[strings.ToUpper(string(state))]
	if !ok {
		return k8sapiv1.Pod_UNKNOWN
	}

	return k8sapiv1.Pod_State(val)
}

func protoForContainerState(state corev1.ContainerState) k8sapiv1.Container_State {
	switch {
	case state.Terminated != nil:
		return k8sapiv1.Container_TERMINATED
	case state.Running != nil:
		return k8sapiv1.Container_RUNNING
	case state.Waiting != nil:
		return k8sapiv1.Container_WAITING
	default:
		return k8sapiv1.Container_UNKNOWN
	}
}
