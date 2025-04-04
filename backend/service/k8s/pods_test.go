package k8s

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
	k8stesting "k8s.io/client-go/testing"

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
				Name:      "testing-pod-name",
				Namespace: "testing-namespace",
				Labels: map[string]string{
					"foo":                  "bar",
					clutchLabelClusterName: "production",
				},
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
				Name:      "testing-pod-name-1",
				Namespace: "testing-namespace",
				Labels: map[string]string{
					"foo":                  "bar",
					clutchLabelClusterName: "staging",
				},
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

func testListFakeClientset(numPods int) *fake.Clientset {
	var fakeClient fake.Clientset
	fakeClient.AddReactor("list", "pods",
		func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
			pods := corev1.PodList{}

			for i := 0; i < numPods; i++ {
				pods.Items = append(pods.Items, corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Name:      fmt.Sprintf("testing-pod-name-%b", i),
						Namespace: "testing-namespace",
						Labels: map[string]string{
							clutchLabelClusterName: "staging",
						},
					},
				})
			}

			return true, &pods, nil
		})
	return &fakeClient
}

func TestDescribePod(t *testing.T) {
	t.Parallel()

	s := &svc{
		manager: &managerImpl{
			clientsets: map[string]*ctxClientsetImpl{
				"foo": {
					Interface: testListFakeClientset(3),
					namespace: "testing-namespace",
					cluster:   "core-testing",
				},
				"bar": {
					Interface: testListFakeClientset(1),
					namespace: "testing-namespace",
					cluster:   "core-testing",
				},
			},
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

	// Found more than 1 pod
	result, err = s.DescribePod(context.Background(),
		"foo",
		"",
		"testing-namespace",
		"testing-pod-name-1",
	)
	assert.Error(t, err)
	assert.Nil(t, result)

	// Found exactly 1 pod
	result, err = s.DescribePod(context.Background(),
		"bar",
		"",
		"testing-namespace",
		"testing-pod-name-0",
	)
	assert.NoError(t, err)
	assert.NotNil(t, result)
}

func TestDeletePod(t *testing.T) {
	t.Parallel()

	cs := testPodClientset()
	s := &svc{
		manager: &managerImpl{
			clientsets: map[string]*ctxClientsetImpl{"foo": {
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
			clientsets: map[string]*ctxClientsetImpl{"testing-clientset": {
				Interface: cs,
				namespace: "testing-namespace",
				cluster:   "testing-cluster",
			}},
		},
	}

	// Two matching pods
	result, err := s.ListPods(
		context.Background(),
		"testing-clientset",
		"testing-cluster",
		"testing-namespace",
		&k8sv1.ListOptions{Labels: map[string]string{"foo": "bar"}},
	)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestPodDescription(t *testing.T) {
	t.Parallel()

	podTestCases := []struct {
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
					Labels: map[string]string{
						clutchLabelClusterName: "production",
					},
				},
				Status: corev1.PodStatus{
					StartTime: &metav1.Time{},
					ContainerStatuses: []corev1.ContainerStatus{
						{Name: "container1"},
						{Name: "container2"},
						{Name: "container3"},
					},
					InitContainerStatuses: []corev1.ContainerStatus{
						{Name: "initcontainer1"},
						{Name: "initcontainer2"},
						{Name: "initcontainer3"},
					},
					Reason: "Evicted",
					Conditions: []corev1.PodCondition{
						{
							Type:   corev1.ContainersReady,
							Status: corev1.ConditionTrue,
						},
					},
				},
			},
		},
		{
			id:                  "custername is not set",
			inputClusterName:    "staging",
			expectedClusterName: "staging",
			pod: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						clutchLabelClusterName: "",
					},
				},
				Status: corev1.PodStatus{
					StartTime: &metav1.Time{},
					ContainerStatuses: []corev1.ContainerStatus{
						{Name: "container1"},
						{Name: "container2"},
						{Name: "container3"},
					},
					InitContainerStatuses: []corev1.ContainerStatus{
						{Name: "initcontainer1"},
						{Name: "initcontainer2"},
						{Name: "initcontainer3"},
					},
					Reason: "Evicted",
					Conditions: []corev1.PodCondition{
						{
							Type:   corev1.ContainersReady,
							Status: corev1.ConditionTrue,
						},
					},
				},
			},
		},
	}

	for _, tt := range podTestCases {
		t.Run(tt.id, func(t *testing.T) {
			t.Parallel()

			pod := podDescription(tt.pod, tt.inputClusterName)
			assert.Equal(t, tt.expectedClusterName, pod.Cluster)
			assert.Equal(t, tt.pod.Status.Reason, pod.StateReason)
			assert.Equal(t, tt.pod.Status.InitContainerStatuses[0].Name, pod.InitContainers[0].Name)
			assert.Equal(t, k8sv1.PodCondition_Type(1), pod.PodConditions[0].Type)
			assert.Equal(t, k8sv1.PodCondition_Status(1), pod.PodConditions[0].Status)
			assert.Equal(t, "Init: 0/0", pod.Status)
			assert.NotNil(t, pod.StartTimeMillis)
		})
	}
}

func TestUpdatePod(t *testing.T) {
	t.Parallel()

	pods := &corev1.PodList{}
	pods.Items = append(pods.Items, corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "testing-pod-name",
			Namespace: "testing-namespace",
			Labels: map[string]string{
				"foo":                  "bar",
				clutchLabelClusterName: "staging",
			},
			Annotations: map[string]string{"baz": "quuz"},
		},
	})

	var fakeClient fake.Clientset
	fakeClient.AddReactor("list", "pods",
		func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
			return true, pods, nil
		})
	fakeClient.AddReactor("get", "pods",
		func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
			name := action.(k8stesting.GetAction).GetName()

			if name != pods.Items[0].Name {
				return true, nil, fmt.Errorf("no pod found")
			}

			return true, &pods.Items[0], nil
		})

	s := &svc{
		manager: &managerImpl{
			clientsets: map[string]*ctxClientsetImpl{
				"testing-clientset": {
					Interface: &fakeClient,
					namespace: "testing-namespace",
					cluster:   "testing-cluster",
				},
			},
		},
	}

	// Pod not found.
	err := s.UpdatePod(context.Background(),
		"testing-clientset",
		"testing-cluster",
		"testing-namespace",
		"non-existent-pod-name",
		&k8sv1.ExpectedObjectMetaFields{},
		&k8sv1.ObjectMetaFields{Annotations: map[string]string{"new-annotation": "foo"}},
		&k8sv1.RemoveObjectMetaFields{},
	)
	assert.Error(t, err)

	// Returns an error when the precondition is not met
	err = s.UpdatePod(context.Background(),
		"testing-clientset",
		"testing-cluster",
		"testing-namespace",
		"testing-pod-name",
		&k8sv1.ExpectedObjectMetaFields{Annotations: map[string]*k8sv1.NullableString{"foo": {Kind: &k8sv1.NullableString_Value{Value: "non-matching-value"}}}},
		&k8sv1.ObjectMetaFields{Annotations: map[string]string{"new-annotation": "foo"}},
		&k8sv1.RemoveObjectMetaFields{},
	)
	assert.Error(t, err)

	// Successfully sets an annotation when the precondition is met
	err = s.UpdatePod(context.Background(),
		"testing-clientset",
		"testing-cluster",
		"testing-namespace",
		"testing-pod-name",
		&k8sv1.ExpectedObjectMetaFields{Annotations: map[string]*k8sv1.NullableString{"baz": {Kind: &k8sv1.NullableString_Value{Value: "quuz"}}}},
		&k8sv1.ObjectMetaFields{Annotations: map[string]string{"baz": "new-value"}},
		&k8sv1.RemoveObjectMetaFields{},
	)
	assert.NoError(t, err)

	// Successfully removes an annotation. This step also verifies that the previous step has properly updated the annotation.
	err = s.UpdatePod(context.Background(),
		"testing-clientset",
		"testing-cluster",
		"testing-namespace",
		"testing-pod-name",
		&k8sv1.ExpectedObjectMetaFields{Annotations: map[string]*k8sv1.NullableString{"baz": {Kind: &k8sv1.NullableString_Value{Value: "new-value"}}}},
		&k8sv1.ObjectMetaFields{},
		&k8sv1.RemoveObjectMetaFields{Annotations: []string{"baz"}},
	)
	assert.NoError(t, err)

	pod, err := s.DescribePod(context.Background(),
		"testing-clientset",
		"testing-cluster",
		"testing-namespace",
		"testing-pod-name",
	)
	assert.NoError(t, err)

	_, annotationPresent := pod.Annotations["baz"]
	assert.False(t, annotationPresent)

	// Checking that an annotation is not set works
	err = s.UpdatePod(context.Background(),
		"testing-clientset",
		"testing-cluster",
		"testing-namespace",
		"testing-pod-name",
		&k8sv1.ExpectedObjectMetaFields{Annotations: map[string]*k8sv1.NullableString{"baz": {Kind: &k8sv1.NullableString_Null{}}}},
		&k8sv1.ObjectMetaFields{Annotations: map[string]string{"baz": "new-value"}},
		&k8sv1.RemoveObjectMetaFields{},
	)
	assert.NoError(t, err)

	pod, err = s.DescribePod(context.Background(),
		"testing-clientset",
		"testing-cluster",
		"testing-namespace",
		"testing-pod-name",
	)
	assert.NoError(t, err)

	annotationValue, annotationPresent := pod.Annotations["baz"]
	assert.True(t, annotationPresent)
	assert.Equal(t, "new-value", annotationValue)
}

func TestMakeContainers(t *testing.T) {
	t.Parallel()

	podTestCases := []struct {
		id                 string
		expectedContainers []*k8sv1.Container
		statuses           []corev1.ContainerStatus
	}{
		{
			id: "cont1",
			expectedContainers: []*k8sv1.Container{
				{
					Name:         "bar",
					Image:        "baz",
					Ready:        false,
					RestartCount: 0,
					State:        k8sv1.Container_RUNNING,
					StateDetails: &k8sv1.Container_StateRunning{
						StateRunning: &k8sv1.StateRunning{},
					},
				},
				{
					Name:         "TheContainer",
					Image:        "foo",
					Ready:        true,
					RestartCount: 5,
					State:        k8sv1.Container_WAITING,
					StateDetails: &k8sv1.Container_StateWaiting{
						StateWaiting: &k8sv1.StateWaiting{},
					},
				},
			},
			statuses: []corev1.ContainerStatus{
				{
					Name:         "bar",
					Image:        "baz",
					Ready:        false,
					RestartCount: 0,
					State: corev1.ContainerState{
						Running: &corev1.ContainerStateRunning{},
					},
				},
				{
					Name:         "TheContainer",
					Image:        "foo",
					Ready:        true,
					RestartCount: 5,
					State: corev1.ContainerState{
						Waiting: &corev1.ContainerStateWaiting{},
					},
				},
			},
		},
		{
			id: "cont2",
			expectedContainers: []*k8sv1.Container{
				{
					Name:         "foo",
					Image:        "giraffe",
					Ready:        true,
					RestartCount: 1,
					State:        k8sv1.Container_TERMINATED,
					StateDetails: &k8sv1.Container_StateTerminated{
						StateTerminated: &k8sv1.StateTerminated{},
					},
				},
			},
			statuses: []corev1.ContainerStatus{
				{
					Name:         "foo",
					Image:        "giraffe",
					Ready:        true,
					RestartCount: 1,
					State: corev1.ContainerState{
						Terminated: &corev1.ContainerStateTerminated{},
					},
				},
			},
		},
	}

	for _, tt := range podTestCases {
		t.Run(tt.id, func(t *testing.T) {
			t.Parallel()
			containers := makeContainers(tt.statuses)
			assert.ElementsMatch(t, tt.expectedContainers, containers)
		})
	}
}

func TestPodStatus(t *testing.T) {
	t.Parallel()

	timeStamp := metav1.Now()
	podTestCases := []struct {
		id                string
		expectedPodStatus string
		pod               *corev1.Pod
	}{
		{
			id:                "Error Status",
			expectedPodStatus: "CreateContainerError",
			pod: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						clutchLabelClusterName: "production",
					},
				},
				Status: corev1.PodStatus{
					ContainerStatuses: []corev1.ContainerStatus{
						{
							Name:         "container1",
							Ready:        false,
							RestartCount: 0,
							State: corev1.ContainerState{
								Waiting: &corev1.ContainerStateWaiting{
									Reason: "CreateContainerError",
								},
							},
						},
					},
				},
			},
		},
		{
			id:                "Error Status",
			expectedPodStatus: "CrashLoopBackOff",
			pod: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						clutchLabelClusterName: "production",
					},
				},
				Status: corev1.PodStatus{
					ContainerStatuses: []corev1.ContainerStatus{
						{
							Name:         "container1",
							Ready:        false,
							RestartCount: 50,
							State: corev1.ContainerState{
								Waiting: &corev1.ContainerStateWaiting{
									Reason: "CrashLoopBackOff",
								},
							},
						},
					},
				},
			},
		},
		{
			id:                "Running Successfully",
			expectedPodStatus: "Running",
			pod: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						clutchLabelClusterName: "production",
					},
				},
				Status: corev1.PodStatus{
					Phase: "Running",
					ContainerStatuses: []corev1.ContainerStatus{
						{
							Name:         "container1",
							Ready:        true,
							RestartCount: 0,
							State: corev1.ContainerState{
								Running: &corev1.ContainerStateRunning{
									StartedAt: metav1.Now(),
								},
							},
						},
					},
				},
			},
		},
		{
			id:                "Terminating",
			expectedPodStatus: "Terminating",
			pod: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					DeletionTimestamp: &timeStamp,
				},
			},
		},
		{
			id:                "PodInitializing",
			expectedPodStatus: "Init: 0/0",
			pod: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						clutchLabelClusterName: "production",
					},
				},
				Status: corev1.PodStatus{
					InitContainerStatuses: []corev1.ContainerStatus{
						{
							Name:         "container1",
							Ready:        true,
							RestartCount: 0,
							State: corev1.ContainerState{
								Waiting: &corev1.ContainerStateWaiting{
									Reason: "PodInitializing",
								},
							},
						},
					},
				},
			},
		},
		{
			id:                "PodTerminating",
			expectedPodStatus: "Signal: 9",
			pod: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						clutchLabelClusterName: "production",
					},
				},
				Status: corev1.PodStatus{
					ContainerStatuses: []corev1.ContainerStatus{
						{
							Name:         "container1",
							Ready:        true,
							RestartCount: 0,
							State: corev1.ContainerState{
								Terminated: &corev1.ContainerStateTerminated{
									Signal: 9,
								},
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range podTestCases {
		t.Run(tt.id, func(t *testing.T) {
			t.Parallel()

			pod := podDescription(tt.pod, "staging")
			assert.Equal(t, tt.expectedPodStatus, pod.Status)
		})
	}
}
