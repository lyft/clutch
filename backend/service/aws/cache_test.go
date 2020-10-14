package aws

import (
	"context"
	"testing"

	topologyv1 "github.com/lyft/clutch/backend/api/topology/v1"
	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/semaphore"
)

func TestStartTopologyCache(t *testing.T) {
	c := client{
		topologyObjectChan: make(chan *topologyv1.UpdateCacheRequest, topologyObjectChanBufferSize),
		topologyLock:       semaphore.NewWeighted(1),
	}

	buff, err := c.StartTopologyCaching(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, buff)

	// Asserts that the topologyLock has already been aquired
	_, err = c.StartTopologyCaching(context.Background())
	assert.Error(t, err)
}
