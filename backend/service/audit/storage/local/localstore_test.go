package local

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/uber-go/tally/v4"
	"go.uber.org/zap/zaptest"

	auditv1 "github.com/lyft/clutch/backend/api/audit/v1"
	auditconfigv1 "github.com/lyft/clutch/backend/api/config/service/audit/v1"
	"github.com/lyft/clutch/backend/service/audit/storage"
)

func TestLocalstorage(t *testing.T) {
	cfg := &auditconfigv1.Config{
		StorageProvider: &auditconfigv1.Config_InMemory{InMemory: true},
	}
	log := zaptest.NewLogger(t)
	scope := tally.NewTestScope("", nil)

	storage, err := New(cfg, log, scope)
	assert.Nil(t, err)

	// No unsent events.
	unsent, err := storage.UnsentEvents(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 0, len(unsent))

	// Assert that later elements will have a larger ID.
	first, err := storage.WriteRequestEvent(context.Background(), &auditv1.RequestEvent{})
	assert.NoError(t, err)
	second, err := storage.WriteRequestEvent(context.Background(), &auditv1.RequestEvent{})
	assert.NoError(t, err)
	assert.Greater(t, second, first)

	// Assert updating is okay.
	diff := &auditv1.RequestEvent{
		Username: "foobar",
	}
	err = storage.UpdateRequestEvent(context.Background(), first, diff)
	assert.NoError(t, err)

	events, err := storage.ReadEvents(context.Background(), time.Time{}, nil, nil)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(events))
	assert.Equal(t, diff, events[0].GetEvent())

	event, err := storage.ReadEvent(context.Background(), 0)
	assert.NoError(t, err)
	assert.Equal(t, diff, event.GetEvent())

	// Assert the expected events are unsent.
	unsent, err = storage.UnsentEvents(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 2, len(unsent))

	// Assert storage is cleared after flushing.
	unsent, err = storage.UnsentEvents(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 0, len(unsent))

	// Assert updating will return an error.
	err = storage.UpdateRequestEvent(context.Background(), first, &auditv1.RequestEvent{})
	assert.Equal(t, errors.New("cannot update event because cannot find by id: 0"), err)
}

func TestReadEvents(t *testing.T) {
	testCases := []struct {
		id            string
		endTimeOffset int
		options       *storage.ReadOptions
		expectedCount int
	}{
		{
			id: "returns subset of events if limit is set",
			options: &storage.ReadOptions{
				Limit: 1,
			},
			expectedCount: 1,
		},
		{
			id: "returns subset of events if offset is set",
			options: &storage.ReadOptions{
				Offset: 1,
			},
			expectedCount: 1,
		},
		{
			id: "returns no events if offset is greater than length",
			options: &storage.ReadOptions{
				Offset: 10,
			},
			expectedCount: 0,
		},
		{
			id: "returns all events if offset is negative",
			options: &storage.ReadOptions{
				Offset: -10,
			},
			expectedCount: 2,
		},
		{
			id: "returns all events if limit is greater than length",
			options: &storage.ReadOptions{
				Limit: 100,
			},
			expectedCount: 2,
		},
		{
			id: "returns all events if limit is 0",
			options: &storage.ReadOptions{
				Limit: 0,
			},
			expectedCount: 2,
		},
	}
	for _, tt := range testCases {
		tt := tt
		t.Run(tt.id, func(t *testing.T) {
			cfg := &auditconfigv1.Config{
				StorageProvider: &auditconfigv1.Config_InMemory{InMemory: true},
			}
			log := zaptest.NewLogger(t)
			scope := tally.NewTestScope("", nil)

			storage, err := New(cfg, log, scope)
			assert.Nil(t, err)

			_, err = storage.WriteRequestEvent(context.Background(), &auditv1.RequestEvent{})
			assert.NoError(t, err)
			_, err = storage.WriteRequestEvent(context.Background(), &auditv1.RequestEvent{})
			assert.NoError(t, err)
			time.Sleep(100 * time.Millisecond)

			events, err := storage.ReadEvents(context.Background(), time.Now().Add(-5*time.Minute), nil, tt.options)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCount, len(events))
		})
	}
}

func TestReadEvent(t *testing.T) {
	testCases := []struct {
		id          string
		eventId     int
		expectError bool
	}{
		{
			id:      "returns event for valid event ID",
			eventId: 0,
		},
		{
			id:          "returns error for invalid event ID",
			eventId:     100,
			expectError: true,
		},
	}
	for _, tt := range testCases {
		tt := tt
		t.Run(tt.id, func(t *testing.T) {
			cfg := &auditconfigv1.Config{
				StorageProvider: &auditconfigv1.Config_InMemory{InMemory: true},
			}
			log := zaptest.NewLogger(t)
			scope := tally.NewTestScope("", nil)

			storage, err := New(cfg, log, scope)
			assert.Nil(t, err)

			e, err := storage.WriteRequestEvent(context.Background(), &auditv1.RequestEvent{})
			assert.NoError(t, err)

			event, err := storage.ReadEvent(context.Background(), int64(tt.eventId))
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, e, event.Id)
			}
		})
	}
}
