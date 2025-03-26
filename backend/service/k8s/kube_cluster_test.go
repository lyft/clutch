package k8s

import (
	"testing"

	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func TestGetKubeClusterName(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		obj    any
		expect string
	}{
		{
			name: "should find cluster",
			obj: &metav1.ObjectMeta{
				Labels: map[string]string{
					clutchLabelClusterName: "meow",
				},
			},
			expect: "meow",
		},
		{
			name: "cluster is not found",
			obj: &metav1.ObjectMeta{
				Labels: map[string]string{},
			},
			expect: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, GetKubeClusterName(tt.obj))
		})
	}
}

func TestApplyClusterLabels(t *testing.T) {
	t.Parallel()

	tests := []struct {
		id          string
		cluster     string
		obj         runtime.Object
		shouldError bool
	}{
		{
			id:          "pod",
			cluster:     "meow",
			obj:         &corev1.Pod{},
			shouldError: false,
		},
		{
			id:          "deployment",
			cluster:     "cat",
			obj:         &appsv1.Deployment{},
			shouldError: false,
		},
		{
			id:          "should error",
			cluster:     "",
			obj:         nil,
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.id, func(t *testing.T) {
			err := ApplyClusterLabels(tt.cluster, tt.obj)
			if !tt.shouldError {
				assert.NoError(t, err)
				assert.Equal(t, tt.cluster, GetKubeClusterName(tt.obj))
			} else {
				assert.Error(t, err)
			}
		})
	}
}
