package k8s

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	autoscalingv1 "k8s.io/api/autoscaling/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8s "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"

	k8sapiv1 "github.com/lyft/clutch/backend/api/k8s/v1"
)

var (
	newInt32 = func(n int32) *int32 { return &n }
)

func testHPAClientset() k8s.Interface {
	hpa := &autoscalingv1.HorizontalPodAutoscaler{
		ObjectMeta: metav1.ObjectMeta{
			Name:        "testing-hpa-name",
			Namespace:   "testing-namespace",
			Labels:      map[string]string{"foo": "bar"},
			Annotations: map[string]string{"baz": "quuz"},
		},
	}

	return fake.NewSimpleClientset(hpa)
}

func TestNormalizeChanges(t *testing.T) {
	t.Parallel()

	var applyScalingTestCases = []struct {
		id         string
		sizing     *k8sapiv1.ResizeHPARequest_Sizing
		resultSpec autoscalingv1.HorizontalPodAutoscalerSpec
	}{
		// HPA unchanged when struct is default constructed
		{id: "applying nil changes"},
		{
			id: "reflect simple changes",
			sizing: &k8sapiv1.ResizeHPARequest_Sizing{
				Min: 25,
				Max: 50,
			},
			resultSpec: autoscalingv1.HorizontalPodAutoscalerSpec{
				MinReplicas: newInt32(25),
				MaxReplicas: 50,
			},
		},
		// Min < Max always
		{
			id: "moving min above max",
			sizing: &k8sapiv1.ResizeHPARequest_Sizing{
				Min: 100,
				Max: 0,
			},
			resultSpec: autoscalingv1.HorizontalPodAutoscalerSpec{
				MinReplicas: newInt32(100),
				MaxReplicas: 100,
			},
		},
		{
			id: "omitting max",
			sizing: &k8sapiv1.ResizeHPARequest_Sizing{
				Min: 25,
			},
			resultSpec: autoscalingv1.HorizontalPodAutoscalerSpec{
				MinReplicas: newInt32(25),
				MaxReplicas: 25,
			},
		},
	}

	for _, tt := range applyScalingTestCases {
		tt := tt
		t.Run(tt.id, func(t *testing.T) {
			t.Parallel()

			hpa := autoscalingv1.HorizontalPodAutoscaler{
				Spec: autoscalingv1.HorizontalPodAutoscalerSpec{
					MinReplicas: newInt32(10),
					MaxReplicas: 50,
				},
			}
			normalizeHPAChanges(&hpa, tt.sizing)

			if tt.sizing == nil {
				tt.resultSpec = hpa.Spec
			}

			assert.Equal(t, *tt.resultSpec.MinReplicas, *hpa.Spec.MinReplicas)
			assert.Equal(t, tt.resultSpec.MaxReplicas, hpa.Spec.MaxReplicas)
			assert.LessOrEqual(t, *hpa.Spec.MinReplicas, hpa.Spec.MaxReplicas)
		})
	}
}

func TestResizeHPA(t *testing.T) {
	t.Parallel()

	cs := testHPAClientset()
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
	err := s.ResizeHPA(context.Background(), "", "", "", "", nil)
	assert.Error(t, err)

	err = s.ResizeHPA(context.Background(), "foo", "core-testing", "testing-namespace", "testing-hpa-name", nil)
	assert.NoError(t, err)
}

func TestProtoForHPAClusterName(t *testing.T) {
	t.Parallel()

	var hpaTestCases = []struct {
		id                  string
		inputClusterName    string
		expectedClusterName string
		hpa                 *autoscalingv1.HorizontalPodAutoscaler
	}{
		{
			id:                  "clustername already set",
			inputClusterName:    "notprod",
			expectedClusterName: "production",
			hpa: &autoscalingv1.HorizontalPodAutoscaler{
				ObjectMeta: metav1.ObjectMeta{
					ClusterName: "production",
				},
				Spec: autoscalingv1.HorizontalPodAutoscalerSpec{
					MinReplicas: newInt32(1),
					MaxReplicas: 2,
				},
			},
		},
		{
			id:                  "custername is not set",
			inputClusterName:    "staging",
			expectedClusterName: "staging",
			hpa: &autoscalingv1.HorizontalPodAutoscaler{
				ObjectMeta: metav1.ObjectMeta{
					ClusterName: "",
				},
				Spec: autoscalingv1.HorizontalPodAutoscalerSpec{
					MinReplicas: newInt32(1),
					MaxReplicas: 2,
				},
			},
		},
		{
			id:                  "foo",
			inputClusterName:    "staging",
			expectedClusterName: "staging",
			hpa: &autoscalingv1.HorizontalPodAutoscaler{
				ObjectMeta: metav1.ObjectMeta{
					ClusterName: "",
				},
				Spec: autoscalingv1.HorizontalPodAutoscalerSpec{
					MinReplicas:                    newInt32(1),
					MaxReplicas:                    69,
					TargetCPUUtilizationPercentage: newInt32(69),
				},
				Status: autoscalingv1.HorizontalPodAutoscalerStatus{
					CurrentReplicas:                 69,
					DesiredReplicas:                 69,
					CurrentCPUUtilizationPercentage: newInt32(69),
				},
			},
		},
		{
			id:                  "bar",
			inputClusterName:    "prod",
			expectedClusterName: "prod",
			hpa: &autoscalingv1.HorizontalPodAutoscaler{
				ObjectMeta: metav1.ObjectMeta{
					ClusterName: "",
				},
				Spec: autoscalingv1.HorizontalPodAutoscalerSpec{
					MinReplicas:                    newInt32(1),
					MaxReplicas:                    3,
					TargetCPUUtilizationPercentage: newInt32(42),
				},
				Status: autoscalingv1.HorizontalPodAutoscalerStatus{
					CurrentReplicas:                 69,
					DesiredReplicas:                 420,
					CurrentCPUUtilizationPercentage: newInt32(88),
				},
			},
		},
	}

	for _, tt := range hpaTestCases {
		tt := tt
		t.Run(tt.id, func(t *testing.T) {
			t.Parallel()

			hpa := ProtoForHPA(tt.inputClusterName, tt.hpa)
			assert.Equal(t, tt.expectedClusterName, hpa.Cluster)
			assert.Equal(t, *tt.hpa.Spec.MinReplicas, int32(hpa.Sizing.MinReplicas))
			assert.Equal(t, tt.hpa.Spec.MaxReplicas, int32(hpa.Sizing.MaxReplicas))

			if tt.hpa.Spec.TargetCPUUtilizationPercentage != nil {
				assert.Equal(t, *tt.hpa.Spec.TargetCPUUtilizationPercentage, hpa.TargetCpuUtilizationPercentage.Value)
			}
			if tt.hpa.Status.CurrentCPUUtilizationPercentage != nil {
				assert.Equal(t, *tt.hpa.Status.CurrentCPUUtilizationPercentage, hpa.CurrentCpuUtilizationPercentage.Value)
			}

			assert.Equal(t, tt.hpa.Status.CurrentReplicas, int32(hpa.Sizing.CurrentReplicas))
			assert.Equal(t, tt.hpa.Status.DesiredReplicas, int32(hpa.Sizing.DesiredReplicas))
		})
	}
}

func TestDeleteHPA(t *testing.T) {
	cs := testHPAClientset()
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
	err := s.DeleteHPA(context.Background(), "foo", "core-testing", "testing-namespace", "abc")
	assert.Error(t, err)

	err = s.DeleteHPA(context.Background(), "foo", "core-testing", "testing-namespace", "testing-hpa-name")
	assert.NoError(t, err)

	// Not found.
	_, err = s.DescribeHPA(context.Background(), "foo", "core-testing", "testing-namespace", "testing-hpa-name")
	assert.Error(t, err)
}
