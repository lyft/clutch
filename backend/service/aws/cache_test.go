package aws

import (
	"context"
	"testing"
	"time"

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

	// Not currently used by AWS caching but its a required parameter
	cacheTtl := time.Minute * 1

	buff, err := c.StartTopologyCaching(context.Background(), cacheTtl)
	assert.NoError(t, err)
	assert.NotNil(t, buff)

	// Asserts that the topologyLock has already been aquired
	_, err = c.StartTopologyCaching(context.Background(), cacheTtl)
	assert.Error(t, err)
}

func TestStartTickerForCacheResource(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	c := client{}

	// This lets us track the number of invocations to processFakeResource()
	numProcessFakeResourceCalls := 0

	processFakeResource := func(ctx context.Context, client *regionalClient) {
		numProcessFakeResourceCalls++
	}

	go func() {
		// Give the ticker time to tick at least once
		time.Sleep(time.Second)
		cancel()
	}()

	c.startTickerForCacheResource(ctx, time.Millisecond, &regionalClient{}, processFakeResource)

	// This asserts that we setup a ticker, it was invoked at least once, and the cancel context was successful.
	assert.Greater(t, numProcessFakeResourceCalls, 0)
}
