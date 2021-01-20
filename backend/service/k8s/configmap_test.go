package k8s

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"

	k8sv1 "github.com/lyft/clutch/backend/api/k8s/v1"
)

func testConfigMapClientset() *fake.Clientset {
	testConfigMaps := []runtime.Object{
		&v1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "testing-configmap-name",
				Namespace: "testing-namespace",
				Labels:    map[string]string{"foo": "bar"},
			},
		},
		&v1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "testing-configmap-name-1",
				Namespace: "testing-namespace",
				Labels:    map[string]string{"foo": "bar"},
			},
		},
		&v1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "testing-configmap-name-2",
				Namespace: "testing-namespace",
				Labels:    map[string]string{"foo1": "bar"},
			},
		},
	}

	return fake.NewSimpleClientset(testConfigMaps...)
}

func TestListConfigMaps(t *testing.T) {
	t.Parallel()

	cs := testConfigMapClientset()

	s := &svc{
		manager: &managerImpl{
			clientsets: map[string]*ctxClientsetImpl{"testing-clientset": &ctxClientsetImpl{
				Interface: cs,
				namespace: "testing-namespace",
				cluster:   "testing-cluster",
			}},
		},
	}

	// No matching ConfigMaps
	result, err := s.ListConfigMaps(
		context.Background(),
		"testing-clientset",
		"testing-cluster",
		"testing-namespace",
		&k8sv1.ListOptions{Labels: map[string]string{"unknown-label": "bar"}},
	)
	assert.NoError(t, err)
	assert.Empty(t, result)

	// Two matching ConfigMaps
	result, err = s.ListConfigMaps(
		context.Background(),
		"testing-clientset",
		"testing-cluster",
		"testing-namespace",
		&k8sv1.ListOptions{Labels: map[string]string{"foo": "bar"}},
	)
	assert.NoError(t, err)
	assert.Len(t, result, 2)

	// All ConfigMaps in the namespace
	result, err = s.ListConfigMaps(
		context.Background(),
		"testing-clientset",
		"testing-cluster",
		"testing-namespace",
		&k8sv1.ListOptions{},
	)
	assert.NoError(t, err)
	assert.Len(t, result, 3)
}

func TestProtoForConfigMap(t *testing.T) {
	t.Parallel()

	var configMapTestCases = []struct {
		id                  string
		inputClusterName    string
		expectedClusterName string
		configMap           *v1.ConfigMap
	}{
		{
			id:                  "clustername already set",
			inputClusterName:    "notprod",
			expectedClusterName: "production",
			configMap: &v1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					ClusterName: "production",
				},
			},
		},
		{
			id:                  "custername is not set",
			inputClusterName:    "staging",
			expectedClusterName: "staging",
			configMap: &v1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					ClusterName: "",
				},
			},
		},
	}

	for _, tt := range configMapTestCases {
		tt := tt
		t.Run(tt.id, func(t *testing.T) {
			t.Parallel()

			configMap := protoForConfigMap(tt.inputClusterName, tt.configMap)
			assert.Equal(t, tt.expectedClusterName, configMap.Cluster)
		})
	}
}
