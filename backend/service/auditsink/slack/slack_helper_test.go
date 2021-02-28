package slack

import (
	"testing"

	"github.com/stretchr/testify/assert"

	auditv1 "github.com/lyft/clutch/backend/api/audit/v1"
	configv1 "github.com/lyft/clutch/backend/api/config/service/auditsink/slack/v1"
)

func TestGetSlackOverrideText(t *testing.T) {
	testCases := []struct {
		override     *configv1.Override
		event        *auditv1.RequestEvent
		expectedOk   bool
		expectedText string
	}{
		// no override specified
		{
			override:   &configv1.Override{},
			event:      &auditv1.RequestEvent{},
			expectedOk: false,
		},
		// the override's FullMethod is invalid
		{
			override: &configv1.Override{
				CustomSlackMessages: []*configv1.CustomSlackMessage{{FullMethod: "foo", Message: "foo"}},
			},
			event:      &auditv1.RequestEvent{},
			expectedOk: false,
		},
		// the override's FullMethod does not matche the service and method in the audit event
		{
			override: &configv1.Override{
				CustomSlackMessages: []*configv1.CustomSlackMessage{
					{FullMethod: "/clutch.k8s.v1.K8sAPI/DescribePod", Message: "foo"},
				},
			},
			event:      &auditv1.RequestEvent{ServiceName: "clutch.k8s.v1.K8sAPI", MethodName: "ResizeHPA"},
			expectedOk: false,
		},
		// the override's FullMethod matches the service and method in the audit event
		{
			override: &configv1.Override{
				CustomSlackMessages: []*configv1.CustomSlackMessage{
					{FullMethod: "/clutch.k8s.v1.K8sAPI/DescribePod", Message: "foo"},
				},
			},
			event:        &auditv1.RequestEvent{ServiceName: "clutch.k8s.v1.K8sAPI", MethodName: "DescribePod"},
			expectedOk:   true,
			expectedText: "foo",
		},
	}

	for _, test := range testCases {
		result, ok := GetSlackOverrideText(test.override, test.event)
		if !test.expectedOk {
			assert.False(t, ok)
			assert.Empty(t, result)
		} else {
			assert.True(t, ok)
			assert.Equal(t, test.expectedText, result)
		}
	}
}
