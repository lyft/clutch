package aws

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
	"golang.org/x/sync/semaphore"

	topologyv1 "github.com/lyft/clutch/backend/api/topology/v1"
)

func TestStartTopologyCache(t *testing.T) {
	log := zaptest.NewLogger(t)
	c := client{
		topologyObjectChan: make(chan *topologyv1.UpdateCacheRequest, topologyObjectChanBufferSize),
		topologyLock:       semaphore.NewWeighted(1),
		log:                log,
	}

	buff, err := c.StartTopologyCaching(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, buff)

	// Asserts that the topologyLock has already been aquired
	_, err = c.StartTopologyCaching(context.Background())
	assert.Error(t, err)
}
