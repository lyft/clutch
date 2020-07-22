package audit

import (
	"testing"

	"github.com/stretchr/testify/assert"

	auditv1 "github.com/lyft/clutch/backend/api/audit/v1"
	auditconfigv1 "github.com/lyft/clutch/backend/api/config/service/audit/v1"
)

func TestFilter(t *testing.T) {
	testCases := []struct {
		id       string
		filter   *auditconfigv1.Filter
		event    *auditv1.Event
		expected bool
	}{
		{
			id: "filter out healthcheck",
			filter: &auditconfigv1.Filter{
				Denylist: true,
				Rules: []*auditconfigv1.EventFilter{
					{
						Field: auditconfigv1.EventFilter_METHOD,
						Value: &auditconfigv1.EventFilter_Text{Text: "Healthcheck"},
					},
				},
			},
			event: &auditv1.Event{
				EventType: &auditv1.Event_Event{
					Event: &auditv1.RequestEvent{
						MethodName: "Healthcheck",
					}}},
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.id, func(t *testing.T) {
			c := &client{filter: tt.filter}
			actual := c.Filter(tt.event)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
