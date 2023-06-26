package k8s

import (
	"context"
	"sort"
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

// For ease of use, we can sort the events in reverse chronological order
func sortEventsByTime(events []*k8sapiv1.Event) {
	sort.Slice(events, func(i, j int) bool {
		return events[i].CreationTimeMillis > events[j].CreationTimeMillis
	})
}

func convertTypeStringToEnum(str string) k8sapiv1.EventType {
	switch str {
	case "Normal":
		return k8sapiv1.EventType_NORMAL
	case "Warning":
		return k8sapiv1.EventType_WARNING
	case "Error":
		return k8sapiv1.EventType_ERROR
	default:
		return k8sapiv1.EventType_TYPE_UNSPECIFIED
	}
}

func ProtoForEvent(cluster string, k8sEvent *corev1.Event) *k8sapiv1.Event {
	// Note for timestamps - in k8s 1.25 LastTimestamp is deprecated in favor of
	// EventTime. However, some objects currently use EventTime, while others use LastTimestamp. Using the
	// CreationTime from the Metadata is an option as it refers to the creation of the object by the server.
	// See https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#ObjectMeta
	// See also https://github.com/kubernetes/kubernetes/issues/90482 for a discussion on different timestamps.
	return &k8sapiv1.Event{
		Cluster:            cluster,
		Namespace:          k8sEvent.Namespace,
		Name:               k8sEvent.Name,
		Reason:             k8sEvent.Reason,
		Description:        k8sEvent.Message,
		InvolvedObjectName: k8sEvent.InvolvedObject.Name,
		Kind:               protoForObjectKind(k8sEvent.InvolvedObject.Kind),
		CreationTimeMillis: k8sEvent.GetObjectMeta().GetCreationTimestamp().UnixMilli(),
		Type:               convertTypeStringToEnum(k8sEvent.Type),
	}
}

// Note in the case of a blank string being returned, that means there is no field selector,
// and all events will be returned.
// Note that chaining field selectors only results in AND not OR, thus to have multiples
// we must have multiple field selectors.
// See https://github.com/kubernetes/kubernetes/issues/32946 for more info
func convertTypesToFieldSelectors(types []k8sapiv1.EventType) []string {
	fs := []string{}
	setOfTypes := make(map[k8sapiv1.EventType]bool)
	for _, t := range types {
		setOfTypes[t] = true
	}

	if setOfTypes[k8sapiv1.EventType_NORMAL] {
		fs = append(fs, "type=Normal")
	}
	if setOfTypes[k8sapiv1.EventType_WARNING] {
		fs = append(fs, "type=Warning")
	}
	if setOfTypes[k8sapiv1.EventType_ERROR] {
		fs = append(fs, "type=Error")
	}

	return fs
}

func (s *svc) ListNamespaceEvents(ctx context.Context, clientset, cluster, namespace string, types []k8sapiv1.EventType) ([]*k8sapiv1.Event, error) {
	cs, err := s.manager.GetK8sClientset(ctx, clientset, cluster, namespace)
	if err != nil {
		return nil, err
	}

	totalEventList := []corev1.Event{}
	fs := convertTypesToFieldSelectors(types)
	// Note if field selector is blank it will return everything
	if len(fs) == 0 {
		eventList, err := cs.CoreV1().Events(cs.Namespace()).List(ctx, metav1.ListOptions{})
		if err != nil {
			return nil, err
		}
		totalEventList = append(totalEventList, eventList.Items...)
	}
	for _, f := range fs {
		eventList, err := cs.CoreV1().Events(cs.Namespace()).List(ctx, metav1.ListOptions{FieldSelector: f})
		if err != nil {
			return nil, err
		}
		totalEventList = append(totalEventList, eventList.Items...)
	}

	// in the future, could also potentially return full event object rather than the subset in
	// ProtoForEvent()
	var events []*k8sapiv1.Event
	for _, ev := range totalEventList {
		ev := ev
		events = append(events, ProtoForEvent(cs.Cluster(), &ev))
	}
	// Sort in reverse chronological order (by CreationTimeMillis),
	// this is needed because we might have multiple types of events
	sortEventsByTime(events)

	return events, nil
}

func protoForObjectKind(kind string) k8sapiv1.ObjectKind {
	// Look up value in generated enum map after ensuring consistent case with generated code.
	val, ok := k8sapiv1.ObjectKind_value[strcase.ToScreamingSnake(kind)]
	if !ok {
		return k8sapiv1.ObjectKind_UNKNOWN
	}
	return k8sapiv1.ObjectKind(val)
}
