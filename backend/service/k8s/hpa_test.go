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

	var (
		newInt32 = func(n int32) *int32 { return &n }
	)

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
