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

func testServiceClientset() k8s.Interface {
	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:        "testing-service-name",
			Namespace:   "testing-namespace",
			Labels:      map[string]string{"foo": "bar"},
			Annotations: map[string]string{"baz": "quuz"},
		},
	}

	return fake.NewSimpleClientset(svc)
}

func TestDeleteService(t *testing.T) {
	cs := testServiceClientset()
	s := &svc{
		manager: &managerImpl{
			clientsets: map[string]*ctxClientsetImpl{"foo": &ctxClientsetImpl{
				Interface: cs,
				namespace: "default",
				cluster:   "core-testing",
			}},
		},
	}

	// Not found.
	err := s.DeleteService(context.Background(), "foo", "core-testing", "testing-namespace", "abc")
	assert.Error(t, err)

	err = s.DeleteService(context.Background(), "foo", "core-testing", "testing-namespace", "testing-service-name")
	assert.NoError(t, err)

	// Not found.
	_, err = s.DescribeService(context.Background(), "foo", "core-testing", "testing-namespace", "testing-service-name")
	assert.Error(t, err)
}

func TestProtoForServiceClusterName(t *testing.T) {
	t.Parallel()

	var serviceTestCases = []struct {
		id                  string
		inputClusterName    string
		expectedClusterName string
		service             *corev1.Service
	}{
		{
			id:                  "clustername already set",
			inputClusterName:    "abc",
			expectedClusterName: "production",
			service: &corev1.Service{
				ObjectMeta: metav1.ObjectMeta{
					ClusterName: "production",
				},
			},
		},
		{
			id:                  "custername is not set",
			inputClusterName:    "staging",
			expectedClusterName: "staging",
			service: &corev1.Service{
				ObjectMeta: metav1.ObjectMeta{
					ClusterName: "",
				},
			},
		},
	}

	for _, tt := range serviceTestCases {
		tt := tt
		t.Run(tt.id, func(t *testing.T) {
			t.Parallel()

			service := ProtoForService(tt.inputClusterName, tt.service)
			assert.Equal(t, tt.expectedClusterName, service.Cluster)
		})
	}
}
