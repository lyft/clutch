package k8s

import (
	"testing"

	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestGenerateStrategicPatch(t *testing.T) {
	t.Parallel()

	var strategicPatchTest = []struct {
		id       string
		old      *appsv1.Deployment
		new      *appsv1.Deployment
		expected string
	}{
		{
			id:       "no change",
			old:      nil,
			new:      nil,
			expected: `{}`,
		},
		{
			id: "Adding Annotations",
			old: &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"foo": "bar"},
				},
				Spec: appsv1.DeploymentSpec{
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{"foo": "bar"},
						},
					},
				},
			},
			new: &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      map[string]string{"foo": "bar"},
					Annotations: map[string]string{"baz": "quuz"},
				},
				Spec: appsv1.DeploymentSpec{
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels:      map[string]string{"foo": "bar"},
							Annotations: map[string]string{"baz": "quuz"},
						},
					},
				},
			},
			expected: `{"metadata":{"annotations":{"baz":"quuz"}},"spec":{"template":{"metadata":{"annotations":{"baz":"quuz"}}}}}`,
		},
		{
			id: "Adding both Labels and Annotations",
			old: &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"foo": "bar"},
				},
				Spec: appsv1.DeploymentSpec{
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{"foo": "bar"},
						},
					},
				},
			},
			new: &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      map[string]string{"label1": "label1val"},
					Annotations: map[string]string{"annotations1": "annotations1val"},
				},
				Spec: appsv1.DeploymentSpec{
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels:      map[string]string{"label1": "label1val"},
							Annotations: map[string]string{"annotations1": "annotations1val"},
						},
					},
				},
			},
			expected: `{"metadata":{"annotations":{"annotations1":"annotations1val"},"labels":{"foo":null,"label1":"label1val"}},"spec":{"template":{"metadata":{"annotations":{"annotations1":"annotations1val"},"labels":{"foo":null,"label1":"label1val"}}}}}`,
		},
	}
	for _, tt := range strategicPatchTest {
		tt := tt
		t.Run(tt.id, func(t *testing.T) {
			t.Parallel()

			actualPatchBytes, err := GenerateStrategicPatch(tt.old, tt.new, appsv1.Deployment{})
			assert.NoError(t, err)
			assert.Equal(t, []byte(tt.expected), actualPatchBytes)
		})
	}
}
