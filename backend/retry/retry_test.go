package retry

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/uber-go/tally"
	"go.uber.org/zap/zaptest"
)

var errFunc = func() error { return errors.New("failed") }

func TestDoAllFail(t *testing.T) {
	t.Parallel()

	log := zaptest.NewLogger(t)
	scope := tally.NewTestScope("test", map[string]string{})
	err := Do(
		context.TODO(),
		log,
		scope,
		errFunc,
		Retries(3),
	)
	assert.Error(t, err)

	// Error format should be same.
	expectedErrorFormat := `All retries failed:
Retry #1: failed
Retry #2: failed
Retry #3: failed`
	assert.Equal(t, expectedErrorFormat, err.Error())

	// Retry count should be 3.
	retryAttempts := scope.Snapshot().Counters()["test.retry_attempts+"]
	assert.Equal(t, int64(3), retryAttempts.Value())
}

func TestDoSuccess(t *testing.T) {
	t.Parallel()

	log := zaptest.NewLogger(t)
	scope := tally.NewTestScope("test", map[string]string{})
	err := Do(
		context.TODO(),
		log,
		scope,
		func() error { return nil },
		Retries(3),
	)
	assert.NoError(t, err)

	// Retry count should be 1.
	retryAttempts := scope.Snapshot().Counters()["test.retry_attempts+"]
	assert.Equal(t, int64(1), retryAttempts.Value())
}

func TestDoWithDefaultConfig(t *testing.T) {
	t.Parallel()

	log := zaptest.NewLogger(t)
	scope := tally.NewTestScope("", nil)
	start := time.Now()
	err := Do(
		context.TODO(),
		log,
		scope,
		errFunc,
	)
	elapsed := time.Since(start)
	assert.Error(t, err)
	// Default backoff is one second and since we skip waiting on the last
	// attempt elapsed should be >=2s instead of >=3s.
	assert.True(t, elapsed >= 2*time.Second, "3 times default retry should be >=2s")
}

func TestDoWithFixedBackoff(t *testing.T) {
	t.Parallel()

	log := zaptest.NewLogger(t)
	scope := tally.NewTestScope("", nil)
	start := time.Now()
	err := Do(
		context.TODO(),
		log,
		scope,
		errFunc,
		Delay(time.Second*2),
		Backoff(FixedBackoff),
	)
	elapsed := time.Since(start)
	assert.Error(t, err)
	// elapsed should be >= 4s.
	assert.True(t, elapsed >= 4*time.Second, "3 times fixed backoff with 2s delay should be >=4s")
}

func TestDoWithExponentialBackoff(t *testing.T) {
	t.Parallel()

	log := zaptest.NewLogger(t)
	scope := tally.NewTestScope("", nil)
	start := time.Now()
	err := Do(
		context.TODO(),
		log,
		scope,
		errFunc,
		Retries(4),
		Backoff(ExponentialBackoff),
	)
	elapsed := time.Since(start)
	assert.Error(t, err)
	// elapsed should be >= 7s.
	assert.True(t, elapsed >= 7*time.Second, "4 times exponential backoff should be >= 7s")
}

func TestDoWithJitterBackoff(t *testing.T) {
	t.Parallel()

	log := zaptest.NewLogger(t)
	scope := tally.NewTestScope("", nil)
	start := time.Now()
	err := Do(
		context.TODO(),
		log,
		scope,
		errFunc,
		Retries(4),
		Backoff(JitterBackoff),
	)
	elapsed := time.Since(start)
	assert.Error(t, err)
	// elapsed should be > 3ms.
	assert.True(t, elapsed > 3*time.Millisecond)
	// elapsed should be < 400ms.
	assert.True(t, elapsed < 400*time.Millisecond)
}

func TestDoWithFailedContext(t *testing.T) {
	t.Parallel()

	log := zaptest.NewLogger(t)
	scope := tally.NewTestScope("", nil)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	// When context is canceled before retry attempts.
	cancel()

	err := Do(
		ctx,
		log,
		scope,
		errFunc,
	)
	assert.Error(t, err)
	assert.Equal(t, context.Canceled, err)

	t.Run("when context deadline exceeds", func(t *testing.T) {
		ctx, cancel = context.WithTimeout(context.Background(), 3 * time.Second)
		defer cancel()

		err := Do(
			ctx,
			log,
			scope,
			errFunc,
			Retries(4),
			Backoff(ExponentialBackoff),
		)
		assert.Error(t, err)
		assert.Equal(t, context.DeadlineExceeded, err)
	})
}
