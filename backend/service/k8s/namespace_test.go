package k8s

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8s "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
)

func testNamespaceClientset() k8s.Interface {
	svc := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:        "testing-namespace",
			Labels:      map[string]string{"foo": "bar"},
			Annotations: map[string]string{"baz": "quuz"},
		},
	}

	return fake.NewSimpleClientset(svc)
}

func TestDescribeNamespace(t *testing.T) {
	cs := testNamespaceClientset()
	s := &svc{
		manager: &managerImpl{
			clientsets: map[string]*ctxClientsetImpl{"foo": &ctxClientsetImpl{
				Interface: cs,
				namespace: "testing-namespace",
				cluster:   "core-testing",
			}},
		},
	}

	namespace, err := s.DescribeNamespace(context.Background(), "foo", "core-testing", "testing-namespace")
	assert.NoError(t, err)
	assert.Equal(t, "testing-namespace", namespace.Name)
}
