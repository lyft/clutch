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
		listOptions         *k8sapiv1.ListPodsOptions
		expectedListOptions metav1.ListOptions
	}{
		{
			id:                  "noop",
			listOptions:         &k8sapiv1.ListPodsOptions{},
			expectedListOptions: metav1.ListOptions{},
		},
		{
			id: "adding label selectors",
			listOptions: &k8sapiv1.ListPodsOptions{
				Labels: map[string]string{
					"foo": "bar",
					"key": "value",
				},
			},
			expectedListOptions: metav1.ListOptions{
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
