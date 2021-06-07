package project

import (
	"testing"

	"github.com/lyft/clutch/backend/mock/service/dbmock"
	"github.com/lyft/clutch/backend/service"
	"github.com/stretchr/testify/assert"
	"github.com/uber-go/tally"
	"go.uber.org/zap/zaptest"
	"google.golang.org/protobuf/types/known/anypb"
)

func TestNew(t *testing.T) {
	log := zaptest.NewLogger(t)
	scope := tally.NewTestScope("", nil)

	// Should throw error looking for postgres
	_, err := New(&anypb.Any{}, log, scope)
	assert.Error(t, err)

	// Should not error with postgres in the registry
	service.Registry["clutch.service.db.postgres"] = dbmock.NewMockDB()
	_, err = New(&anypb.Any{}, log, scope)
	assert.NoError(t, err)
}
