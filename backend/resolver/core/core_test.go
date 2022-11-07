package core

import (
	"testing"

	"github.com/lyft/clutch/backend/mock/service/topologymock"
	"github.com/lyft/clutch/backend/service"
	"github.com/stretchr/testify/assert"
	"github.com/uber-go/tally/v4"
	"go.uber.org/zap/zaptest"
)

func TestNewCoreResolver(t *testing.T) {
	service.Registry["clutch.service.topology"] = topologymock.New()

	log := zaptest.NewLogger(t)
	scope := tally.NewTestScope("", nil)

	core, err := New(nil, log, scope)
	assert.NoError(t, err)
	assert.NotNil(t, core)

	// Test successful construction without the topology service
	delete(service.Registry, "clutch.service.topology")
	coreNoTopology, err2 := New(nil, log, scope)
	assert.NoError(t, err2)
	assert.NotNil(t, coreNoTopology)
}
