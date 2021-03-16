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

func testStatefulSetClientset() k8s.Interface {
	statefulSet := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:        "testing-statefulSet-name",
			Namespace:   "testing-namespace",
			Labels:      map[string]string{"foo": "bar"},
			Annotations: map[string]string{"baz": "quuz"},
		},
	}

	return fake.NewSimpleClientset(statefulSet)
}

func TestListStatefulSets(t *testing.T) {
	cs := testStatefulSetClientset()
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
	list, err := s.ListStatefulSets(context.Background(), "foo", "core-testing", "testing-namespace", opts)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(list))
	// Not Found
	opts = &k8sapiv1.ListOptions{Labels: map[string]string{"unknown": "bar"}}
	list, err = s.ListStatefulSets(context.Background(), "foo", "core-testing", "testing-namespace", opts)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(list))
}

func TestUpdateStatefulSet(t *testing.T) {
	t.Parallel()
	cs := testStatefulSetClientset()
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
	err := s.UpdateStatefulSet(context.Background(), "", "", "", "", nil)
	assert.Error(t, err)

	err = s.UpdateStatefulSet(context.Background(), "foo", "core-testing", "testing-namespace", "testing-statefulSet-name", &k8sapiv1.UpdateStatefulSetRequest_Fields{})
	assert.NoError(t, err)
}

func TestMergeStatefulSetLabelsAndAnnotations(t *testing.T) {
	t.Parallel()

	var mergeLabelAnnotationsTestCases = []struct {
		id     string
		fields *k8sapiv1.UpdateStatefulSetRequest_Fields
		expect *appsv1.StatefulSet
	}{
		{
			id:     "no changes",
			fields: &k8sapiv1.UpdateStatefulSetRequest_Fields{},
			expect: &appsv1.StatefulSet{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      map[string]string{"foo": "bar"},
					Annotations: map[string]string{"baz": "quuz"},
				},
			},
		},
		{
			id: "Adding labels",
			fields: &k8sapiv1.UpdateStatefulSetRequest_Fields{
				Labels: map[string]string{"field1": "field1val"},
			},
			expect: &appsv1.StatefulSet{
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
			fields: &k8sapiv1.UpdateStatefulSetRequest_Fields{
				Annotations: map[string]string{"field1": "field1val"},
			},
			expect: &appsv1.StatefulSet{
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
			fields: &k8sapiv1.UpdateStatefulSetRequest_Fields{
				Labels:      map[string]string{"label1": "label1val"},
				Annotations: map[string]string{"annotation1": "annotation1val"},
			},
			expect: &appsv1.StatefulSet{
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

	for _, tt := range mergeLabelAnnotationsTestCases {
		tt := tt
		t.Run(tt.id, func(t *testing.T) {
			t.Parallel()

			statefulSet := &appsv1.StatefulSet{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      map[string]string{"foo": "bar"},
					Annotations: map[string]string{"baz": "quuz"},
				},
			}

			mergeStatefulSetLabelsAndAnnotations(statefulSet, tt.fields)
			assert.Equal(t, tt.expect.ObjectMeta, statefulSet.ObjectMeta)
		})
	}
}

func TestProtoForStatefulSetClusterName(t *testing.T) {
	t.Parallel()

	var statefulSetTestCases = []struct {
		id                  string
		inputClusterName    string
		expectedClusterName string
		statefulSet         *appsv1.StatefulSet
	}{
		{
			id:                  "clustername already set",
			inputClusterName:    "notprod",
			expectedClusterName: "production",
			statefulSet: &appsv1.StatefulSet{
				ObjectMeta: metav1.ObjectMeta{
					ClusterName: "production",
				},
			},
		},
		{
			id:                  "custername is not set",
			inputClusterName:    "staging",
			expectedClusterName: "staging",
			statefulSet: &appsv1.StatefulSet{
				ObjectMeta: metav1.ObjectMeta{
					ClusterName: "",
				},
			},
		},
	}

	for _, tt := range statefulSetTestCases {
		tt := tt
		t.Run(tt.id, func(t *testing.T) {
			t.Parallel()

			statefulSet := ProtoForStatefulSet(tt.inputClusterName, tt.statefulSet)
			assert.Equal(t, tt.expectedClusterName, statefulSet.Cluster)
		})
	}
}

func TestDeleteStatefulSet(t *testing.T) {
	cs := testStatefulSetClientset()
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
	err := s.DeleteStatefulSet(context.Background(), "foo", "core-testing", "testing-namespace", "abc")
	assert.Error(t, err)

	err = s.DeleteStatefulSet(context.Background(), "foo", "core-testing", "testing-namespace", "testing-statefulSet-name")
	assert.NoError(t, err)

	// Not found.
	_, err = s.DescribeStatefulSet(context.Background(), "foo", "core-testing", "testing-namespace", "testing-statefulSet-name")
	assert.Error(t, err)
}
