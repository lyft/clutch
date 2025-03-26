package k8s

import (
	"testing"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	k8sapiv1 "github.com/lyft/clutch/backend/api/k8s/v1"
)

func TestApplyListOptions(t *testing.T) {
	t.Parallel()

	testCases := []struct {
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
		t.Run(tt.id, func(t *testing.T) {
			t.Parallel()

			listOpts, err := ApplyListOptions(tt.listOptions)
			assert.Equal(t, listOpts, tt.expectedListOptions)
			assert.Nil(t, err)
		})
	}
}
