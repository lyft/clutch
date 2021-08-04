package k8s

import (
	"context"
	"strings"

	"github.com/iancoleman/strcase"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	k8sapiv1 "github.com/lyft/clutch/backend/api/k8s/v1"
)

func (s *svc) ListEvents(ctx context.Context, clientset, cluster, namespace, object string, kind k8sapiv1.ObjectKind) ([]*k8sapiv1.Event, error) {
	cs, err := s.manager.GetK8sClientset(ctx, clientset, cluster, namespace)
	if err != nil {
		return nil, err
	}

	objKind := strcase.ToCamel(strings.ToLower(kind.String()))
	// returns the appropriate field selector based on the object involved
	fieldSelector := cs.CoreV1().Events(cs.Namespace()).GetFieldSelector(&object, &namespace, &objKind, nil)

	eventList, err := cs.CoreV1().Events(cs.Namespace()).List(ctx, metav1.ListOptions{FieldSelector: fieldSelector.String()})
	if err != nil {
		return nil, err
	}

	var events []*k8sapiv1.Event
	for i := range eventList.Items {
		events = append(events, ProtoForEvent(cs.Cluster(), &eventList.Items[i]))
	}

	return events, nil
}

func ProtoForEvent(cluster string, k8sEvent *corev1.Event) *k8sapiv1.Event {
	clusterName := k8sEvent.ClusterName
	if clusterName == "" {
		clusterName = cluster
	}
	return &k8sapiv1.Event{
		Cluster:            clusterName,
		Namespace:          k8sEvent.Namespace,
		Name:               k8sEvent.Name,
		Reason:             k8sEvent.Reason,
		Description:        k8sEvent.Message,
		InvolvedObjectName: k8sEvent.InvolvedObject.Name,
		Kind:               protoForObjectKind(k8sEvent.InvolvedObject.Kind),
	}
}

func protoForObjectKind(kind string) k8sapiv1.ObjectKind {
	// Look up value in generated enum map after ensuring consistent case with generated code.
	val, ok := k8sapiv1.ObjectKind_value[strcase.ToScreamingSnake(kind)]
	if !ok {
		return k8sapiv1.ObjectKind_UNKNOWN
	}
	return k8sapiv1.ObjectKind(val)
}
