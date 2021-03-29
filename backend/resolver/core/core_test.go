package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uber-go/tally"
	"go.uber.org/zap/zaptest"

	"github.com/lyft/clutch/backend/mock/service/topologymock"
	"github.com/lyft/clutch/backend/service"
)

func TestNewAwsResolver(t *testing.T) {
	service.Registry["clutch.service.topology"] = topologymock.New()

	log := zaptest.NewLogger(t)
	scope := tally.NewTestScope("", nil)

	aws, err := New(nil, log, scope)
	assert.NoError(t, err)
	assert.NotNil(t, aws)

	// Test successful construction without the topology service
	delete(service.Registry, "clutch.service.topology")
	coreNoTopology, err2 := New(nil, log, scope)
	assert.NoError(t, err2)
	assert.NotNil(t, coreNoTopology)
}
