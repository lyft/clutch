package stats

import (
	"context"
	"testing"
	"time"

	"github.com/uber-go/tally"

)

func TestNewRuntimeStatsCancelTicker(t *testing.T) {
	scope := tally.NewTestScope("", nil)
	ctx, cancel := context.WithCancel(context.Background())

	runtimeCollector := NewRuntimeStats(scope, &gatewayv1.Stats_GoRuntimeStats{})
	go runtimeCollector.Collect(ctx)

	// Give time for the ticker to spin up
	time.Sleep(time.Millisecond * 50)
	cancel()
}
