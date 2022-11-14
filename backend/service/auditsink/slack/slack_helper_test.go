package slack

import (
	"testing"

	"github.com/stretchr/testify/assert"

	configv1 "github.com/lyft/clutch/backend/api/config/service/auditsink/slack/v1"
)

func TestNewOverrideLookup(t *testing.T) {
	testCases := []struct {
		override      []*configv1.CustomMessage
		output        OverrideLookup
		expectedEmpty bool
	}{
		// no overrides
		{
			override:      []*configv1.CustomMessage{},
			output:        OverrideLookup{},
			expectedEmpty: true,
		},
		{
			override: []*configv1.CustomMessage{{FullMethod: "foo", Message: "bar"}},
			output: OverrideLookup{
				messages: map[string]*configv1.CustomMessage{
					"foo": {FullMethod: "foo", Message: "bar"},
				},
			},
		},
	}
	for _, test := range testCases {
		result := NewOverrideLookup(test.override)
		if test.expectedEmpty {
			assert.Empty(t, result)
		} else {
			assert.Equal(t, 1, len(result.messages))

			v, ok := result.messages["foo"]
			assert.True(t, ok)
			assert.Equal(t, "bar", v.Message)
		}
	}
}

func TestGetOverrideMessage(t *testing.T) {
	testCases := []struct {
		input        OverrideLookup
		service      string
		method       string
		expectedOk   bool
		expectedText string
	}{
		// OverrideLookup is empty
		{
			input:      OverrideLookup{},
			service:    "foo",
			method:     "bar",
			expectedOk: false,
		},
		// OverrideLookup map is empty
		{
			input:      OverrideLookup{messages: map[string]*configv1.CustomMessage{}},
			service:    "foo",
			method:     "bar",
			expectedOk: false,
		},
		// no match
		{
			input: OverrideLookup{
				messages: map[string]*configv1.CustomMessage{
					"/clutch.k8s.v1.K8sAPI/DescribePod": {FullMethod: "/clutch.k8s.v1.K8sAPI/DescribePod", Message: "foo"},
				},
			},
			service:    "clutch.k8s.v1.K8sAPI",
			method:     "ResizeHPA",
			expectedOk: false,
		},
		// match
		{
			input: OverrideLookup{
				messages: map[string]*configv1.CustomMessage{
					"/clutch.k8s.v1.K8sAPI/DescribePod": {FullMethod: "/clutch.k8s.v1.K8sAPI/DescribePod", Message: "foo"},
					"foo":                               {FullMethod: "foo", Message: "foo"},
				},
			},
			service:      "clutch.k8s.v1.K8sAPI",
			method:       "DescribePod",
			expectedOk:   true,
			expectedText: "foo",
		},
	}

	for _, test := range testCases {
		message, ok := test.input.GetOverrideMessage(test.service, test.method)
		if !test.expectedOk {
			assert.False(t, ok)
			assert.Empty(t, message)
		} else {
			assert.True(t, ok)
			assert.Equal(t, test.expectedText, message)
		}
	}
}
