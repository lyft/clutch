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

	// Retry count should be 2 as the first attempt is not counted
	// as a retry.
	retryAttempts := scope.Snapshot().Counters()["test.retry_attempts+"]
	assert.Equal(t, int64(2), retryAttempts.Value())
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

	// Retry count should be 0 as the first attempt succeeds.
	retryAttempts := scope.Snapshot().Counters()["test.retry_attempts+"]
	assert.Equal(t, int64(0), retryAttempts.Value())
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
	// Default backoff is 100ms and since we skip waiting on the last
	// attempt elapsed should be >=200ms instead of >=300ms.
	assert.True(t, elapsed >= 200*time.Millisecond, "3 times default retry should be >=200ms")
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
	assert.True(t, elapsed > 3*time.Millisecond, "4 times 100ms jitter backoff should be > 3ms")
	// elapsed should be < 400ms.
	assert.True(t, elapsed < 400*time.Millisecond, "4 times 100ms jitter backoff should be < 400ms")

	t.Run("user defined jitter", func(t *testing.T) {
		start = time.Now()
		err = Do(
			context.TODO(),
			log,
			scope,
			errFunc,
			Retries(4),
			Jitter(200*time.Millisecond),
			Backoff(JitterBackoff),
		)
		elapsed = time.Since(start)
		assert.Error(t, err)
		// elapsed should be > 6ms.
		assert.True(t, elapsed > 6*time.Millisecond, "4 times 200ms jitter backoff should be > 6ms")
		// elapsed should be < 800ms.
		assert.True(t, elapsed < 800*time.Millisecond, "4 times 200 ms jitter backoff should < 800ms")
	})

	t.Run("when user defined jitter is <= 0", func(t *testing.T) {
		start = time.Now()
		err = Do(
			context.TODO(),
			log,
			scope,
			errFunc,
			Retries(4),
			Jitter(-1),
			Backoff(JitterBackoff),
		)
		elapsed = time.Since(start)
		assert.Error(t, err)
		// elapsed should be > 3ms.
		assert.True(t, elapsed > 3*time.Millisecond, "4 times 100ms jitter backoff should be > 3ms")
		// elapsed should be < 400ms.
		assert.True(t, elapsed < 400*time.Millisecond, "4 times 100ms jitter backoff should be < 400ms")
	})
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
		ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
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
