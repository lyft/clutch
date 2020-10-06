package local

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uber-go/tally"
	"go.uber.org/zap/zaptest"

	auditv1 "github.com/lyft/clutch/backend/api/audit/v1"
	auditconfigv1 "github.com/lyft/clutch/backend/api/config/service/audit/v1"
)

func TestLocalstorage(t *testing.T) {
	cfg := &auditconfigv1.Config{
		Service: &auditconfigv1.Config_LocalAuditing{LocalAuditing: true},
	}
	log := zaptest.NewLogger(t)
	scope := tally.NewTestScope("", nil)

	storage, err := New(cfg, log, scope)
	assert.Nil(t, err)

	// No unsent events.
	unsent, err := storage.UnsentEvents(context.Background())
	assert.Nil(t, err)
	assert.Equal(t, 0, len(unsent))

	// Assert that later elements will have a larger ID.
	first, err := storage.WriteRequestEvent(context.Background(), &auditv1.RequestEvent{})
	assert.Nil(t, err)
	second, err := storage.WriteRequestEvent(context.Background(), &auditv1.RequestEvent{})
	assert.Nil(t, err)
	assert.Greater(t, second, first)

	// Assert updating is okay.
	err = storage.UpdateRequestEvent(context.Background(), first, &auditv1.RequestEvent{})
	assert.Nil(t, err)

	// Assert the expected events are unsent.
	unsent, err = storage.UnsentEvents(context.Background())
	assert.Nil(t, err)
	assert.Equal(t, 2, len(unsent))

	// Assert storage is cleared after flushing.
	unsent, err = storage.UnsentEvents(context.Background())
	assert.Nil(t, err)
	assert.Equal(t, 0, len(unsent))

	// Assert updating will return an error.
	err = storage.UpdateRequestEvent(context.Background(), first, &auditv1.RequestEvent{})
	assert.Equal(t, errors.New("cannot update event because cannot find by id: 0"), err)
}
