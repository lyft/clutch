package k8s

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
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
