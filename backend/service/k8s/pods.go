package k8s

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	k8sapiv1 "github.com/lyft/clutch/backend/api/k8s/v1"
)

func (s *svc) DescribePod(ctx context.Context, clientset, cluster, namespace, name string) (*k8sapiv1.Pod, error) {
	cs, err := s.manager.GetK8sClientset(ctx, clientset, cluster, namespace)
	if err != nil {
		return nil, err
	}

	pods, err := cs.CoreV1().Pods(cs.Namespace()).List(ctx, metav1.ListOptions{
		FieldSelector: "metadata.name=" + name,
	})
	if err != nil {
		return nil, err
	}

	if len(pods.Items) == 1 {
		return podDescription(&pods.Items[0], cs.Cluster()), nil
	} else if len(pods.Items) > 1 {
		return nil, status.Error(codes.FailedPrecondition, "located multiple pods")
	}
	return nil, status.Error(codes.NotFound, "unable to locate specified pod")
}

func (s *svc) DeletePod(ctx context.Context, clientset, cluster, namespace, name string) error {
	cs, err := s.manager.GetK8sClientset(ctx, clientset, cluster, namespace)
	if err != nil {
		return err
	}

	var gracePeriod int64
	opts := metav1.DeleteOptions{}
	opts.GracePeriodSeconds = &gracePeriod

	return cs.CoreV1().Pods(cs.Namespace()).Delete(ctx, name, opts)
}

func (s *svc) ListPods(ctx context.Context, clientset, cluster, namespace string, listOpts *k8sapiv1.ListOptions) ([]*k8sapiv1.Pod, error) {
	cs, err := s.manager.GetK8sClientset(ctx, clientset, cluster, namespace)
	if err != nil {
		return nil, err
	}

	opts, err := ApplyListOptions(listOpts)
	if err != nil {
		return nil, err
	}

	podList, err := cs.CoreV1().Pods(cs.Namespace()).List(ctx, opts)
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
		return status.Error(codes.InvalidArgument, "update of pod labels not implemented")
	}

	// Ensure that the caller is not trying to delete an annotation and update it at the same time
	newAnnotations := objectMetaFields.GetAnnotations()
	for _, annotation := range removeObjectMetaFields.GetAnnotations() {
		_, annotationIsUpdated := newAnnotations[annotation]
		if annotationIsUpdated {
			return status.Errorf(codes.InvalidArgument, "annotation '%s' can't be updated and removed at once", annotation)
		}
	}

	cs, err := s.manager.GetK8sClientset(ctx, clientset, cluster, namespace)
	if err != nil {
		return err
	}

	pod, err := cs.CoreV1().Pods(cs.Namespace()).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	// Check that the current state of the pod matches with expectedObjectMetaFields.
	//
	// If there is a mismatch, checkExpectedObjectMetaFields() will return an error with the list of mismatches.
	err = s.checkExpectedObjectMetaFields(expectedObjectMetaFields, pod.GetObjectMeta())
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

	_, err = cs.CoreV1().Pods(cs.Namespace()).Update(ctx, pod, metav1.UpdateOptions{})
	return err
}

func (s *svc) checkExpectedObjectMetaFields(expectedObjectMetaFields *k8sapiv1.ExpectedObjectMetaFields, object metav1.Object) error {
	if len(expectedObjectMetaFields.Labels) > 0 {
		return status.Error(codes.InvalidArgument, "checking label expectations not implemented")
	}

	podAnnotations := object.GetAnnotations()
	var mismatchedAnnotations []*mismatchedAnnotation

	for expectedAnnotation, expectedValue := range expectedObjectMetaFields.GetAnnotations() {
		currentValue, annotationIsPresent := podAnnotations[expectedAnnotation]

		switch expectedValue.Kind.(type) {
		case *k8sapiv1.NullableString_Null:
			// Existance precondition not met
			if annotationIsPresent {
				mismatchedAnnotations = append(
					mismatchedAnnotations,
					&mismatchedAnnotation{
						Annotation:    expectedAnnotation,
						ExpectedValue: expectedValue.GetValue(),
						CurrentValue:  currentValue,
					},
				)
			}
		case *k8sapiv1.NullableString_Value:
			if !annotationIsPresent || expectedValue.GetValue() != currentValue {
				// Annotation values mismatched
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
	}

	if len(mismatchedAnnotations) == 0 {
		return nil
	}

	return &ExpectedObjectMetaFieldsCheckError{MismatchedAnnotations: mismatchedAnnotations}
}

func podDescription(k8spod *corev1.Pod, cluster string) *k8sapiv1.Pod {
	// TODO: There's a mismatch between the serialization of the timestamp here and what's expected
	// on the frontend.
	// var launch *timestamp.Timestamp
	// if converted, err := ptypes.TimestampProto(k8spod.Status.StartTime.Time); err == nil {
	// 	launch = converted
	// }
	clusterName := k8spod.ClusterName
	if clusterName == "" {
		clusterName = cluster
	}
	pod := &k8sapiv1.Pod{
		Cluster:    clusterName,
		Namespace:  k8spod.Namespace,
		Name:       k8spod.Name,
		Containers: makeContainers(k8spod.Status.ContainerStatuses),
		NodeIp:     k8spod.Status.HostIP,
		PodIp:      k8spod.Status.PodIP,
		State:      protoForPodState(k8spod.Status.Phase),
		//StartTime:   launch,
		Labels:         k8spod.Labels,
		Annotations:    k8spod.Annotations,
		StateReason:    k8spod.Status.Reason,
		PodConditions:  makeConditions(k8spod.Status.Conditions),
		InitContainers: makeContainers(k8spod.Status.InitContainerStatuses),
		Status:         getPodStatus(k8spod),
	}

	// This is here as a workaround since we can't use StartTime
	if k8spod.Status.StartTime != nil {
		// Note .Unix() gives seconds, not milliseconds
		pod.StartTimeMillis = k8spod.Status.StartTime.UnixNano() / 1e6
	}

	return pod
}

func makeConditions(conditions []corev1.PodCondition) []*k8sapiv1.PodCondition {
	podConditions := make([]*k8sapiv1.PodCondition, 0, len(conditions))
	for _, condition := range conditions {
		cond := &k8sapiv1.PodCondition{
			Type:   protoForConditionType(condition.Type),
			Status: protoForConditionStatus(condition.Status),
		}
		podConditions = append(podConditions, cond)
	}
	return podConditions
}

func makeContainers(statuses []corev1.ContainerStatus) []*k8sapiv1.Container {
	containers := make([]*k8sapiv1.Container, 0, len(statuses))
	for _, status := range statuses {
		container := &k8sapiv1.Container{
			Name:         status.Name,
			Image:        status.Image,
			State:        protoForContainerState(status.State),
			Ready:        status.Ready,
			RestartCount: status.RestartCount,
		}

		switch container.State {
		case k8sapiv1.Container_TERMINATED:
			container.StateDetails = protoForContainerStateTerminated(status.State.Terminated)
		case k8sapiv1.Container_RUNNING:
			container.StateDetails = protoForContainerStateRunning(status.State.Running)
		case k8sapiv1.Container_WAITING:
			container.StateDetails = protoForContainerStateWaiting(status.State.Waiting)
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

func protoForContainerStateWaiting(state *corev1.ContainerStateWaiting) *k8sapiv1.Container_StateWaiting {
	return &k8sapiv1.Container_StateWaiting{
		StateWaiting: &k8sapiv1.StateWaiting{
			Reason:  state.Reason,
			Message: state.Message,
		},
	}
}

func protoForContainerStateTerminated(state *corev1.ContainerStateTerminated) *k8sapiv1.Container_StateTerminated {
	return &k8sapiv1.Container_StateTerminated{
		StateTerminated: &k8sapiv1.StateTerminated{
			Reason:   state.Reason,
			Message:  state.Message,
			ExitCode: state.ExitCode,
			Signal:   state.Signal,
		},
	}
}

func protoForContainerStateRunning(state *corev1.ContainerStateRunning) *k8sapiv1.Container_StateRunning {
	return &k8sapiv1.Container_StateRunning{
		StateRunning: &k8sapiv1.StateRunning{
			// FE serialization currently does not support timestamp
			//StartTime: state.StartedAt,
		},
	}
}

func protoForConditionType(conditionType corev1.PodConditionType) k8sapiv1.PodCondition_Type {
	switch conditionType {
	case corev1.ContainersReady:
		return k8sapiv1.PodCondition_CONTAINERS_READY
	case corev1.PodInitialized:
		return k8sapiv1.PodCondition_INITIALIZED
	case corev1.PodReady:
		return k8sapiv1.PodCondition_READY
	case corev1.PodScheduled:
		return k8sapiv1.PodCondition_POD_SCHEDULED
	default:
		return k8sapiv1.PodCondition_TYPE_UNSPECIFIED
	}
}
func protoForConditionStatus(status corev1.ConditionStatus) k8sapiv1.PodCondition_Status {
	switch status {
	case corev1.ConditionTrue:
		return k8sapiv1.PodCondition_TRUE
	case corev1.ConditionFalse:
		return k8sapiv1.PodCondition_FALSE
	case corev1.ConditionUnknown:
		return k8sapiv1.PodCondition_UNKNOWN
	default:
		return k8sapiv1.PodCondition_STATUS_UNSPECIFIED
	}
}

func getPodStatus(pod *corev1.Pod) string {
	restarts := 0
	readyContainers := 0

	reason := string(pod.Status.Phase)
	if pod.Status.Reason != "" {
		reason = pod.Status.Reason
	}

	initializing := false
	for i := range pod.Status.InitContainerStatuses {
		container := pod.Status.InitContainerStatuses[i]
		restarts += int(container.RestartCount)
		switch {
		case container.State.Terminated != nil && container.State.Terminated.ExitCode == 0:
			continue
		case container.State.Terminated != nil:
			// initialization is failed
			if len(container.State.Terminated.Reason) == 0 {
				if container.State.Terminated.Signal != 0 {
					reason = fmt.Sprintf("Init: Signal %d", container.State.Terminated.Signal)
				} else {
					reason = fmt.Sprintf("Init: ExitCode %d", container.State.Terminated.ExitCode)
				}
			} else {
				reason = "Init:" + container.State.Terminated.Reason
			}
			initializing = true
		case container.State.Waiting != nil && len(container.State.Waiting.Reason) > 0 && container.State.Waiting.Reason != "PodInitializing":
			reason = fmt.Sprintf("Init: %s", container.State.Waiting.Reason)
			initializing = true
		default:
			reason = fmt.Sprintf("Init: %d/%d", i, len(pod.Spec.InitContainers))
			initializing = true
		}
		break
	}
	if !initializing {
		restarts = 0
		for i := len(pod.Status.ContainerStatuses) - 1; i >= 0; i-- {
			container := pod.Status.ContainerStatuses[i]

			restarts += int(container.RestartCount)
			switch {
			case container.State.Waiting != nil && container.State.Waiting.Reason != "":
				reason = container.State.Waiting.Reason
			case container.State.Terminated != nil && container.State.Terminated.Reason != "":
				reason = container.State.Terminated.Reason
			case container.State.Terminated != nil && container.State.Terminated.Reason == "":
				if container.State.Terminated.Signal != 0 {
					reason = fmt.Sprintf("Signal: %d", container.State.Terminated.Signal)
				} else {
					reason = fmt.Sprintf("ExitCode: %d", container.State.Terminated.ExitCode)
				}
			case container.Ready && container.State.Running != nil:
				readyContainers++
			}
		}
	}

	if pod.DeletionTimestamp != nil && pod.Status.Reason == "NodeLost" {
		reason = string(corev1.PodUnknown)
	} else if pod.DeletionTimestamp != nil {
		reason = "Terminating"
	}

	if len(reason) == 0 {
		reason = string(corev1.PodUnknown)
	}

	return reason
}
