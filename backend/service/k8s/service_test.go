package k8s

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8s "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"

	k8sv1 "github.com/lyft/clutch/backend/api/k8s/v1"
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

func TestDescribeService(t *testing.T) {
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

	_, err := s.DescribeService(context.Background(), "foo", "core-testing", "testing-namespace", "testing-service-name")
	assert.NoError(t, err)
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

func TestProtoForServiceType(t *testing.T) {
	assert.Equal(t, k8sv1.Service_CLUSTER_IP, protoForServiceType(corev1.ServiceTypeClusterIP))
	assert.Equal(t, k8sv1.Service_NODE_PORT, protoForServiceType(corev1.ServiceTypeNodePort))
	assert.Equal(t, k8sv1.Service_LOAD_BALANCER, protoForServiceType(corev1.ServiceTypeLoadBalancer))
	assert.Equal(t, k8sv1.Service_EXTERNAL_NAME, protoForServiceType(corev1.ServiceTypeExternalName))
}
