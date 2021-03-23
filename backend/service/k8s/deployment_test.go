package k8s

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8s "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"

	k8sapiv1 "github.com/lyft/clutch/backend/api/k8s/v1"
)

func testDeploymentClientset() k8s.Interface {
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:        "testing-deployment-name",
			Namespace:   "testing-namespace",
			Labels:      map[string]string{"foo": "bar"},
			Annotations: map[string]string{"baz": "quuz"},
		},
	}

	return fake.NewSimpleClientset(deployment)
}

func TestDescribeDeployment(t *testing.T) {
	cs := testDeploymentClientset()
	s := &svc{
		manager: &managerImpl{
			clientsets: map[string]*ctxClientsetImpl{"foo": &ctxClientsetImpl{
				Interface: cs,
				namespace: "default",
				cluster:   "core-testing",
			}},
		},
	}

	d, err := s.DescribeDeployment(context.Background(), "foo", "core-testing", "testing-namespace", "testing-deployment-name")
	assert.NoError(t, err)
	assert.Equal(t, "testing-deployment-name", d.Name)
	// Not found.
	_, err = s.DescribeDeployment(context.Background(), "foo", "core-testing", "testing-namespace", "abc")
	assert.NoError(t, err)
}

func TestListDeployments(t *testing.T) {
	cs := testDeploymentClientset()
	s := &svc{
		manager: &managerImpl{
			clientsets: map[string]*ctxClientsetImpl{"foo": &ctxClientsetImpl{
				Interface: cs,
				namespace: "default",
				cluster:   "core-testing",
			}},
		},
	}

	opts := &k8sapiv1.ListOptions{Labels: map[string]string{"foo": "bar"}}
	list, err := s.ListDeployments(context.Background(), "foo", "core-testing", "testing-namespace", opts)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(list))
	// Not Found
	opts = &k8sapiv1.ListOptions{Labels: map[string]string{"unknown": "bar"}}
	list, err = s.ListDeployments(context.Background(), "foo", "core-testing", "testing-namespace", opts)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(list))
}

func TestUpdateDeployment(t *testing.T) {
	cs := testDeploymentClientset()
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
	err := s.UpdateDeployment(context.Background(), "", "", "", "", nil)
	assert.Error(t, err)

	err = s.UpdateDeployment(context.Background(), "foo", "core-testing", "testing-namespace", "testing-deployment-name", &k8sapiv1.UpdateDeploymentRequest_Fields{})
	assert.NoError(t, err)
}

func TestDeleteDeployment(t *testing.T) {
	cs := testDeploymentClientset()
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
	err := s.DeleteDeployment(context.Background(), "foo", "core-testing", "testing-namespace", "abc")
	assert.Error(t, err)

	err = s.DeleteDeployment(context.Background(), "foo", "core-testing", "testing-namespace", "testing-deployment-name")
	assert.NoError(t, err)

	// Not found.
	_, err = s.DescribeDeployment(context.Background(), "foo", "core-testing", "testing-namespace", "testing-deployment-name")
	assert.Error(t, err)
}

func TestMergeDeploymentLabelsAndAnnotations(t *testing.T) {
	t.Parallel()

	var mergeDeploymentLabelAnnotationsTestCases = []struct {
		id     string
		fields *k8sapiv1.UpdateDeploymentRequest_Fields
		expect *appsv1.Deployment
	}{
		{
			id:     "no changes",
			fields: &k8sapiv1.UpdateDeploymentRequest_Fields{},
			expect: &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      map[string]string{"foo": "bar"},
					Annotations: map[string]string{"baz": "quuz"},
				},
			},
		},
		{
			id: "Adding labels",
			fields: &k8sapiv1.UpdateDeploymentRequest_Fields{
				Labels: map[string]string{"field1": "field1val"},
			},
			expect: &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"foo":    "bar",
						"field1": "field1val",
					},
					Annotations: map[string]string{
						"baz": "quuz",
					},
				},
			},
		},
		{
			id: "Adding annotations",
			fields: &k8sapiv1.UpdateDeploymentRequest_Fields{
				Annotations: map[string]string{"field1": "field1val"},
			},
			expect: &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"foo": "bar",
					},
					Annotations: map[string]string{
						"baz":    "quuz",
						"field1": "field1val",
					},
				},
			},
		},
		{
			id: "Adding both labels and annotations",
			fields: &k8sapiv1.UpdateDeploymentRequest_Fields{
				Labels:      map[string]string{"label1": "label1val"},
				Annotations: map[string]string{"annotation1": "annotation1val"},
			},
			expect: &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"foo":    "bar",
						"label1": "label1val",
					},
					Annotations: map[string]string{
						"baz":         "quuz",
						"annotation1": "annotation1val",
					},
				},
			},
		},
	}

	for _, tt := range mergeDeploymentLabelAnnotationsTestCases {
		tt := tt
		t.Run(tt.id, func(t *testing.T) {
			t.Parallel()

			deployment := &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      map[string]string{"foo": "bar"},
					Annotations: map[string]string{"baz": "quuz"},
				},
			}

			mergeDeploymentLabelsAndAnnotations(deployment, tt.fields)
			assert.Equal(t, tt.expect.ObjectMeta, deployment.ObjectMeta)
		})
	}
}

func TestProtoForDeploymentClusterName(t *testing.T) {
	t.Parallel()

	var deploymentTestCases = []struct {
		id                  string
		inputClusterName    string
		expectedClusterName string
		deployment          *appsv1.Deployment
	}{
		{
			id:                  "clustername already set",
			inputClusterName:    "notprod",
			expectedClusterName: "production",
			deployment: &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					ClusterName: "production",
				},
			},
		},
		{
			id:                  "custername is not set",
			inputClusterName:    "staging",
			expectedClusterName: "staging",
			deployment: &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					ClusterName: "",
				},
			},
		},
	}

	for _, tt := range deploymentTestCases {
		tt := tt
		t.Run(tt.id, func(t *testing.T) {
			t.Parallel()

			deployment := ProtoForDeployment(tt.inputClusterName, tt.deployment)
			assert.Equal(t, tt.expectedClusterName, deployment.Cluster)
		})
	}
}

func TestProtoForDeploymentStatus(t *testing.T) {
	t.Parallel()

	var deploymentTestCases = []struct {
		id         string
		deployment *appsv1.Deployment
	}{
		{
			id: "foo",
			deployment: &appsv1.Deployment{
				Status: appsv1.DeploymentStatus{
					Replicas:            60,
					UpdatedReplicas:     60,
					AvailableReplicas:   10,
					UnavailableReplicas: 20,
					Conditions: []appsv1.DeploymentCondition{
						{
							Type:    appsv1.DeploymentAvailable,
							Status:  v1.ConditionTrue,
							Reason:  "reason",
							Message: "message",
						},
					},
				},
			},
		},
		{
			id: "oof",
			deployment: &appsv1.Deployment{
				Status: appsv1.DeploymentStatus{
					Replicas:            40,
					UpdatedReplicas:     314,
					AvailableReplicas:   3,
					UnavailableReplicas: 8,
					Conditions: []appsv1.DeploymentCondition{
						{
							Type:    appsv1.DeploymentReplicaFailure,
							Status:  v1.ConditionFalse,
							Reason:  "space",
							Message: "moon",
						},
					},
				},
			},
		},
	}

	for _, tt := range deploymentTestCases {
		tt := tt
		t.Run(tt.id, func(t *testing.T) {
			t.Parallel()

			deployment := ProtoForDeployment("", tt.deployment)
			assert.Equal(t, tt.deployment.Status.Replicas, int32(deployment.DeploymentStatus.Replicas))
			assert.Equal(t, tt.deployment.Status.UpdatedReplicas, int32(deployment.DeploymentStatus.UpdatedReplicas))
			assert.Equal(t, tt.deployment.Status.ReadyReplicas, int32(deployment.DeploymentStatus.ReadyReplicas))
			assert.Equal(t, tt.deployment.Status.AvailableReplicas, int32(deployment.DeploymentStatus.AvailableReplicas))
			assert.Equal(t, len(tt.deployment.Status.Conditions), len(deployment.DeploymentStatus.DeploymentConditions))
		})
	}
}
