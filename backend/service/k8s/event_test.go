package k8s

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	k8s "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"

	k8sapiv1 "github.com/lyft/clutch/backend/api/k8s/v1"
)

func testEventClientset() k8s.Interface {
	testEvents := []runtime.Object{
		&corev1.Event{
			ObjectMeta: metav1.ObjectMeta{
				Name:              "testing-event-name-1",
				Namespace:         "testing-namespace",
				CreationTimestamp: metav1.Time{Time: time.Unix(1, 0)},
			},
			InvolvedObject: corev1.ObjectReference{
				Kind:      "Pod",
				Namespace: "testing-namespace",
				Name:      "Pod1",
			},
			Reason:        "testing-reason-1",
			Message:       "testing-message-1",
			EventTime:     metav1.MicroTime{Time: time.Unix(1, 0)},
			LastTimestamp: metav1.Time{Time: time.Unix(1, 0)},
		},
		&corev1.Event{
			ObjectMeta: metav1.ObjectMeta{
				Name:              "testing-event-name-2",
				Namespace:         "testing-namespace",
				CreationTimestamp: metav1.Time{Time: time.Unix(2, 0)},
			},
			InvolvedObject: corev1.ObjectReference{
				Kind:      "Pod",
				Namespace: "testing-namespace",
				Name:      "Pod1",
			},
			Reason:        "testing-reason-2",
			Message:       "testing-message-2",
			EventTime:     metav1.MicroTime{Time: time.Unix(2, 0)},
			LastTimestamp: metav1.Time{Time: time.Unix(2, 0)},
		},
	}

	return fake.NewSimpleClientset(testEvents...)
}
func TestListEvents(t *testing.T) {
	cs := testEventClientset()
	s := &svc{
		manager: &managerImpl{
			clientsets: map[string]*ctxClientsetImpl{"foo": {
				Interface: cs,
				namespace: "default",
				cluster:   "core-testing",
			}},
		},
	}

	val, ok := k8sapiv1.ObjectKind_value[strings.ToUpper("Pod")]
	assert.Equal(t, true, ok)
	kind := k8sapiv1.ObjectKind(val)
	list, err := s.ListEvents(context.Background(), "foo", "core-testing", "testing-namespace", "Pod1", kind)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(list))
	for i, v := range list {
		reason := fmt.Sprintf("testing-reason-%d", i+1)
		assert.Equal(t, reason, v.Reason)
		assert.Equal(t, kind, v.Kind)
		timeVal := (int64)(i + 1)
		assert.Equal(t, time.Unix(timeVal, 0).UnixMilli(), v.CreationTimeMillis)
	}
}
