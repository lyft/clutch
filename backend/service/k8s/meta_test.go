package k8s

import (
	"testing"

	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	k8sapiv1 "github.com/lyft/clutch/backend/api/k8s/v1"
)

func TestApplyListOptions(t *testing.T) {
	t.Parallel()

	var testCases = []struct {
		id                  string
		listOptions         *k8sapiv1.ListOptions
		expectedListOptions metav1.ListOptions
	}{
		{
			id:                  "noop",
			listOptions:         &k8sapiv1.ListOptions{},
			expectedListOptions: metav1.ListOptions{},
		},
		{
			id: "adding label selectors",
			listOptions: &k8sapiv1.ListOptions{
				Labels: map[string]string{
					"foo": "bar",
					"key": "value",
				},
			},
			expectedListOptions: metav1.ListOptions{
				LabelSelector: "foo=bar,key=value",
			},
		},
		{
			id: "map and selector string",
			listOptions: &k8sapiv1.ListOptions{
				Labels: map[string]string{
					"foo": "bar",
					"key": "value",
				},
				SupplementalSelectorString: "abc",
			},
			expectedListOptions: metav1.ListOptions{
				LabelSelector: "foo=bar,key=value,abc",
			},
		},
		{
			id: "using selector string1",
			listOptions: &k8sapiv1.ListOptions{
				SupplementalSelectorString: "abc",
			},
			expectedListOptions: metav1.ListOptions{
				LabelSelector: "abc",
			},
		},
		{
			id: "using selector string2",
			listOptions: &k8sapiv1.ListOptions{
				SupplementalSelectorString: "abc,def",
			},
			expectedListOptions: metav1.ListOptions{
				LabelSelector: "abc,def",
			},
		},
		{
			id: "using selector string3",
			listOptions: &k8sapiv1.ListOptions{
				Labels: map[string]string{
					"foo": "bar",
				},
				SupplementalSelectorString: "!abc",
			},
			expectedListOptions: metav1.ListOptions{
				LabelSelector: "foo=bar,!abc",
			},
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.id, func(t *testing.T) {
			t.Parallel()

			listOpts, err := ApplyListOptions(tt.listOptions)
			assert.Equal(t, listOpts, tt.expectedListOptions)
			assert.Nil(t, err)
		})
	}
}

func TestApplyClusterMetadata(t *testing.T) {
	t.Parallel()

	var tests = []struct {
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
		tt := tt
		t.Run(tt.id, func(t *testing.T) {
			err := ApplyClusterMetadata(tt.cluster, tt.obj)
			if !tt.shouldError {
				assert.NoError(t, err)

				objMeta, err := meta.Accessor(tt.obj)
				assert.NoError(t, err)
				assert.Equal(t, tt.cluster, objMeta.GetClusterName())
			} else {
				assert.Error(t, err)
			}
		})
	}
}
