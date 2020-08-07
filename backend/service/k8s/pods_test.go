package k8s

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"

	k8sv1 "github.com/lyft/clutch/backend/api/k8s/v1"
)

func TestProtoForPodState(t *testing.T) {
	assert.Equal(t, k8sv1.Pod_PENDING, protoForPodState(corev1.PodPending))
	assert.Equal(t, k8sv1.Pod_RUNNING, protoForPodState(corev1.PodRunning))
	assert.Equal(t, k8sv1.Pod_SUCCEEDED, protoForPodState(corev1.PodSucceeded))
	assert.Equal(t, k8sv1.Pod_FAILED, protoForPodState(corev1.PodFailed))
	assert.Equal(t, k8sv1.Pod_UNKNOWN, protoForPodState(corev1.PodUnknown))
}

func TestProtoForContainerState(t *testing.T) {
	assert.Equal(t, k8sv1.Container_RUNNING, protoForContainerState(corev1.ContainerState{Running: &corev1.ContainerStateRunning{}}))
	assert.Equal(t, k8sv1.Container_WAITING, protoForContainerState(corev1.ContainerState{Waiting: &corev1.ContainerStateWaiting{}}))
	assert.Equal(t, k8sv1.Container_TERMINATED, protoForContainerState(corev1.ContainerState{Terminated: &corev1.ContainerStateTerminated{}}))
}

func testPodClientset() *fake.Clientset {
	testPods := []runtime.Object{
		&corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:        "testing-pod-name",
				Namespace:   "testing-namespace",
				ClusterName: "production",
				Labels:      map[string]string{"foo": "bar"},
				Annotations: map[string]string{"baz": "quuz"},
			},
			Status: corev1.PodStatus{
				StartTime: &metav1.Time{},
				ContainerStatuses: []corev1.ContainerStatus{
					{Name: "container1"},
					{Name: "container2"},
					{Name: "container3"},
				},
			},
		},
		&corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:        "testing-pod-name-1",
				Namespace:   "testing-namespace",
				ClusterName: "staging",
				Labels:      map[string]string{"foo": "bar"},
				Annotations: map[string]string{"baz": "quuz"},
			},
			Status: corev1.PodStatus{
				StartTime: &metav1.Time{},
				ContainerStatuses: []corev1.ContainerStatus{
					{Name: "container1"},
					{Name: "container2"},
					{Name: "container3"},
				},
			},
		},
		&corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:        "testing-pod-name-2",
				Namespace:   "testing-namespace",
				Labels:      map[string]string{"foo1": "bar"},
				Annotations: map[string]string{"baz": "quuz"},
			},
			Status: corev1.PodStatus{
				StartTime: &metav1.Time{},
				ContainerStatuses: []corev1.ContainerStatus{
					{Name: "container1"},
					{Name: "container2"},
					{Name: "container3"},
				},
			},
		},
	}

	return fake.NewSimpleClientset(testPods...)
}

func TestDescribePod(t *testing.T) {
	t.Parallel()

	cs := testPodClientset()
	s := &svc{
		manager: &managerImpl{
			clientsets: map[string]*ctxClientsetImpl{"foo": {
				Interface: cs,
				namespace: "testing-namespace",
				cluster:   "core-testing",
			}},
		},
	}
	// Not found.
	result, err := s.DescribePod(
		context.Background(),
		"",
		"",
		"",
		"",
	)
	assert.Error(t, err)
	assert.Nil(t, result)

	result, err = s.DescribePod(context.Background(),
		"foo",
		"",
		"testing-namespace",
		"testing-pod-name",
	)
	assert.NoError(t, err)
	assert.NotNil(t, result)

	// Assert that the pod is found using the clientset's namespace if none is provided.
	result, err = s.DescribePod(context.Background(),
		"foo",
		"",
		"",
		"testing-pod-name",
	)
	assert.NoError(t, err)
	assert.NotNil(t, result)
}

func TestDeletePod(t *testing.T) {
	t.Parallel()

	cs := testPodClientset()
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
	err := s.DeletePod(context.Background(),
		"foo",
		"",
		"",
		"",
	)
	assert.Error(t, err)

	err = s.DeletePod(context.Background(),
		"foo",
		"",
		"testing-namespace",
		"testing-pod-name",
	)
	assert.NoError(t, err)
}

func TestListPods(t *testing.T) {
	t.Parallel()

	cs := testPodClientset()
	s := &svc{
		manager: &managerImpl{
			clientsets: map[string]*ctxClientsetImpl{"testing-clientset": &ctxClientsetImpl{
				Interface: cs,
				namespace: "testing-namespace",
				cluster:   "testing-cluster",
			}},
		},
	}
	// Clientset not found
	result, err := s.ListPods(
		context.Background(),
		"unknown-clientset",
		"testing-cluster",
		"testing-namespace",
		&k8sv1.ListPodsOptions{},
	)
	assert.Error(t, err)
	assert.Nil(t, result)

	// No matching pods
	result, err = s.ListPods(
		context.Background(),
		"testing-clientset",
		"testing-cluster",
		"testing-namespace",
		&k8sv1.ListPodsOptions{Labels: map[string]string{"unknown-annotation": "bar"}},
	)
	assert.NoError(t, err)
	assert.Empty(t, result)

	// Two matching pods
	result, err = s.ListPods(
		context.Background(),
		"testing-clientset",
		"testing-cluster",
		"testing-namespace",
		&k8sv1.ListPodsOptions{Labels: map[string]string{"foo": "bar"}},
	)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestPodDescriptionClusterName(t *testing.T) {
	t.Parallel()

	var podTestCases = []struct {
		id                  string
		inputClusterName    string
		expectedClusterName string
		pod                 *corev1.Pod
	}{
		{
			id:                  "clustername already set",
			inputClusterName:    "notprod",
			expectedClusterName: "production",
			pod: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					ClusterName: "production",
				},
			},
		},
		{
			id:                  "custername is not set",
			inputClusterName:    "staging",
			expectedClusterName: "staging",
			pod: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					ClusterName: "",
				},
			},
		},
	}

	for _, tt := range podTestCases {
		tt := tt
		t.Run(tt.id, func(t *testing.T) {
			t.Parallel()

			pod := podDescription(tt.pod, tt.inputClusterName)
			assert.Equal(t, tt.expectedClusterName, pod.Cluster)
		})
	}
}
