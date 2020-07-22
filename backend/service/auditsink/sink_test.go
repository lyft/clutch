package auditsink

import (
	"testing"

	"github.com/stretchr/testify/assert"

	apiv1 "github.com/lyft/clutch/backend/api/api/v1"
	auditv1 "github.com/lyft/clutch/backend/api/audit/v1"
	configv1 "github.com/lyft/clutch/backend/api/config/service/audit/v1"
)

var filterTests = []struct {
	id       string
	filter   *configv1.Filter
	event    *auditv1.Event
	expected bool
}{
	{
		id:       "nil filter passes",
		expected: true,
	},
	{
		id: "service allowlist match passes",
		filter: &configv1.Filter{
			Rules: []*configv1.EventFilter{
				{
					Field: configv1.EventFilter_SERVICE,
					Value: &configv1.EventFilter_Text{Text: "service"},
				},
			},
		},
		event: &auditv1.Event{
			EventType: &auditv1.Event_Event{Event: &auditv1.RequestEvent{
				ServiceName: "service",
			}},
		},
		expected: true,
	},
	{
		id: "service denylist match fails",
		filter: &configv1.Filter{
			Denylist: true,
			Rules: []*configv1.EventFilter{
				{
					Field: configv1.EventFilter_SERVICE,
					Value: &configv1.EventFilter_Text{Text: "service"},
				},
			},
		},
		event: &auditv1.Event{
			EventType: &auditv1.Event_Event{Event: &auditv1.RequestEvent{
				ServiceName: "service",
			}},
		},
		expected: false,
	},
	{
		id: "method allowlist match passes",
		filter: &configv1.Filter{
			Rules: []*configv1.EventFilter{
				{
					Field: configv1.EventFilter_METHOD,
					Value: &configv1.EventFilter_Text{Text: "method"},
				},
			},
		},
		event: &auditv1.Event{
			EventType: &auditv1.Event_Event{Event: &auditv1.RequestEvent{
				MethodName: "method",
			}},
		},
		expected: true,
	},
	{
		id: "method denylist match fails",
		filter: &configv1.Filter{
			Denylist: true,
			Rules: []*configv1.EventFilter{
				{
					Field: configv1.EventFilter_METHOD,
					Value: &configv1.EventFilter_Text{Text: "method"},
				},
			},
		},
		event: &auditv1.Event{
			EventType: &auditv1.Event_Event{Event: &auditv1.RequestEvent{
				MethodName: "method",
			}},
		},
		expected: false,
	},
	{
		id: "type allowlist match passes",
		filter: &configv1.Filter{
			Rules: []*configv1.EventFilter{
				{
					Field: configv1.EventFilter_TYPE,
					Value: &configv1.EventFilter_Text{Text: "READ"},
				},
			},
		},
		event: &auditv1.Event{
			EventType: &auditv1.Event_Event{Event: &auditv1.RequestEvent{
				Type: apiv1.ActionType_READ,
			}},
		},
		expected: true,
	},
	{
		id: "type denylist match fails",
		filter: &configv1.Filter{
			Denylist: true,
			Rules: []*configv1.EventFilter{
				{
					Field: configv1.EventFilter_TYPE,
					Value: &configv1.EventFilter_Text{Text: "READ"},
				},
			},
		},
		event: &auditv1.Event{
			EventType: &auditv1.Event_Event{Event: &auditv1.RequestEvent{
				Type: apiv1.ActionType_READ,
			}},
		},
		expected: false,
	},
	{
		id: "beavior with multiple allowlist",
		filter: &configv1.Filter{
			Rules: []*configv1.EventFilter{
				{
					Field: configv1.EventFilter_TYPE,
					Value: &configv1.EventFilter_Text{Text: "CREATE"},
				},
				{
					Field: configv1.EventFilter_TYPE,
					Value: &configv1.EventFilter_Text{Text: "READ"},
				},
			},
		},
		event: &auditv1.Event{
			EventType: &auditv1.Event_Event{Event: &auditv1.RequestEvent{
				Type: apiv1.ActionType_READ,
			}},
		},
		expected: true,
	},
	{
		id: "beavior with multiple denylist",
		filter: &configv1.Filter{
			Denylist: true,
			Rules: []*configv1.EventFilter{
				{
					Field: configv1.EventFilter_TYPE,
					Value: &configv1.EventFilter_Text{Text: "CREATE"},
				},
				{
					Field: configv1.EventFilter_TYPE,
					Value: &configv1.EventFilter_Text{Text: "READ"},
				},
			},
		},
		event: &auditv1.Event{
			EventType: &auditv1.Event_Event{Event: &auditv1.RequestEvent{
				Type: apiv1.ActionType_READ,
			}},
		},
		expected: false,
	},
}

func TestFilter(t *testing.T) {
	t.Parallel()

	for _, tt := range filterTests {
		tt := tt
		t.Run(tt.id, func(t *testing.T) {
			t.Parallel()
			actual := Filter(tt.filter, tt.event)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
