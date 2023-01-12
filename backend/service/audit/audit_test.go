package audit

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/uber-go/tally/v4"
	"go.uber.org/zap/zaptest"
	"google.golang.org/protobuf/types/known/anypb"

	apiv1 "github.com/lyft/clutch/backend/api/api/v1"
	auditv1 "github.com/lyft/clutch/backend/api/audit/v1"
	auditconfigv1 "github.com/lyft/clutch/backend/api/config/service/audit/v1"
)

func TestNew(t *testing.T) {
	cfg, _ := anypb.New(&auditconfigv1.Config{
		StorageProvider: &auditconfigv1.Config_InMemory{InMemory: true},
		Filter:          &auditconfigv1.Filter{Denylist: true},
	})
	log := zaptest.NewLogger(t)
	scope := tally.NewTestScope("", nil)
	svc, err := New(cfg, log, scope)
	assert.Nil(t, err)

	// Assert conformance to public interface.
	_, ok := svc.(Auditor)
	assert.True(t, ok)

	as := svc.(*client)
	assert.True(t, as.filter.Denylist)
}

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
		{
			id: "filter out healthcheck via allowlist for something else",
			filter: &auditconfigv1.Filter{
				Rules: []*auditconfigv1.EventFilter{
					{
						Field: auditconfigv1.EventFilter_METHOD,
						Value: &auditconfigv1.EventFilter_Text{Text: "Readiness"},
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

func writeRequestEvents(c *client, name string, count int) error {
	for i := 0; i < count; i++ {
		_, err := c.WriteRequestEvent(context.Background(), &auditv1.RequestEvent{
			Username:    "CI",
			ServiceName: "Test Service",
			MethodName:  fmt.Sprintf("Run - %s #%d", name, count),
			Type:        apiv1.ActionType_READ,
		})
		if err != nil {
			return err
		}
		time.Sleep(100 * time.Millisecond)
	}
	return nil
}

func TestReadEvents(t *testing.T) {
	testCases := []struct {
		id            string
		endTimeOffset int
		options       *ReadOptions
		expectedCount int
	}{
		{
			id:            "returns all events until current time by default",
			expectedCount: 5,
		},
		{
			id:            "returns subset of events if end time set",
			endTimeOffset: -250,
			expectedCount: 3,
		},
		{
			id: "returns subset of events if limit is set",
			options: &ReadOptions{
				Limit: 1,
			},
			expectedCount: 1,
		},
		{
			id: "returns subset of events if offset is set",
			options: &ReadOptions{
				Offset: 3,
			},
			expectedCount: 2,
		},
		{
			id: "returns no events if offset is greater than length",
			options: &ReadOptions{
				Offset: 10,
			},
			expectedCount: 0,
		},
		{
			id: "returns all events if offset is negative",
			options: &ReadOptions{
				Offset: -10,
			},
			expectedCount: 5,
		},
		{
			id: "returns all events if limit is greater than length",
			options: &ReadOptions{
				Limit: 100,
			},
			expectedCount: 5,
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.id, func(t *testing.T) {
			cfg, _ := anypb.New(&auditconfigv1.Config{
				StorageProvider: &auditconfigv1.Config_InMemory{InMemory: true},
				Filter:          &auditconfigv1.Filter{Denylist: true},
			})
			log := zaptest.NewLogger(t)
			scope := tally.NewTestScope("", nil)
			svc, err := New(cfg, log, scope)
			assert.Nil(t, err)
			as := svc.(*client)
			err = writeRequestEvents(as, tt.id, 5)
			assert.NoError(t, err)
			var events []*auditv1.Event
			if tt.endTimeOffset != 0 {
				endTime := time.Now().Add(time.Duration(tt.endTimeOffset) * time.Millisecond)
				events, err = as.ReadEvents(context.Background(), time.Now().Add(-5*time.Minute), &endTime, tt.options)
			} else {
				events, err = as.ReadEvents(context.Background(), time.Now().Add(-5*time.Minute), nil, tt.options)
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCount, len(events))
		})
	}
}

func TestReadEvent(t *testing.T) {
	cfg, _ := anypb.New(&auditconfigv1.Config{
		StorageProvider: &auditconfigv1.Config_InMemory{InMemory: true},
		Filter:          &auditconfigv1.Filter{Denylist: true},
	})
	log := zaptest.NewLogger(t)
	scope := tally.NewTestScope("", nil)
	svc, err := New(cfg, log, scope)
	assert.Nil(t, err)
	as := svc.(*client)
	err = writeRequestEvents(as, "TestReadEvent", 5)
	assert.NoError(t, err)
	event, err := as.ReadEvent(context.Background(), 3)
	assert.NoError(t, err)
	assert.Equal(t, int64(3), event.Id)
}

func TestConvertLockToUint32(t *testing.T) {
	tests := []struct {
		id     string
		input  string
		expect uint32
	}{
		{
			id:     "key with chars",
			input:  "audit:eventsink",
			expect: 1635083369,
		},
		{
			id:     "key with special chars",
			input:  "*()#@&!*(#!@",
			expect: 707275043,
		},
		{
			id:     "key with numbers",
			input:  "123123123",
			expect: 825373489,
		},
	}

	for _, tt := range tests {
		id := convertLockToUint32(tt.input)
		assert.Equal(t, tt.expect, id)
	}
}
