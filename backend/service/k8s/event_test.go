package k8s

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8s "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"

	k8sapiv1 "github.com/lyft/clutch/backend/api/k8s/v1"
	k8sv1 "github.com/lyft/clutch/backend/api/k8s/v1"
)

func testEventClientset() k8s.Interface {
	svc := &corev1.Event{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "testing-event-name",
			Namespace: "testing-namespace",
		},
		InvolvedObject: corev1.ObjectReference{
			Kind:      "Pod",
			Namespace: "testing-namespace",
			Name:      "Pod1",
		},
		Reason:  "testing-reason-1",
		Message: "testing-message-1",
	}

	return fake.NewSimpleClientset(svc)
}

func TestListEvents(t *testing.T) {
	cs := testEventClientset()
	s := &svc{
		manager: &managerImpl{
			clientsets: map[string]*ctxClientsetImpl{"foo": &ctxClientsetImpl{
				Interface: cs,
				namespace: "default",
				cluster:   "core-testing",
			}},
		},
	}

	val, ok := k8sapiv1.ObjectKind_value[strings.ToUpper("Pod")]
	assert.Equal(t, true, ok)
	kind := k8sv1.ObjectKind(val)
	list, err := s.ListEvents(context.Background(), "foo", "core-testing", "testing-namespace", "Pod1", kind)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(list))
	assert.Equal(t, "testing-reason-1", list[0].Reason)
	assert.Equal(t, kind, list[0].Kind)
}
