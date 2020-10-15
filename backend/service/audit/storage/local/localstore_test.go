package local

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/uber-go/tally"
	"go.uber.org/zap/zaptest"

	auditv1 "github.com/lyft/clutch/backend/api/audit/v1"
	auditconfigv1 "github.com/lyft/clutch/backend/api/config/service/audit/v1"
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

	events, err := storage.ReadEvents(context.Background(), time.Time{}, nil)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(events))
	assert.Equal(t, diff, events[0].GetEvent())

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
