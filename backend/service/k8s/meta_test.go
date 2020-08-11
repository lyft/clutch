package k8s

import (
	"testing"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

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
			id: "adding field selectors",
			listOptions: &k8sapiv1.ListOptions{
				FieldSelectors: "metadata.name=this-is-a-pod-name",
			},
			expectedListOptions: metav1.ListOptions{
				FieldSelector: "metadata.name=this-is-a-pod-name",
			},
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
			id: "adding both labels and field selectors",
			listOptions: &k8sapiv1.ListOptions{
				FieldSelectors: "metadata.name=this-is-a-pod-name",
				Labels: map[string]string{
					"foo": "bar",
					"key": "value",
				},
			},
			expectedListOptions: metav1.ListOptions{
				FieldSelector: "metadata.name=this-is-a-pod-name",
				LabelSelector: "foo=bar,key=value",
			},
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.id, func(t *testing.T) {
			t.Parallel()

			listOpts := ApplyListOptions(tt.listOptions)
			assert.Equal(t, listOpts, tt.expectedListOptions)
		})
	}
}
