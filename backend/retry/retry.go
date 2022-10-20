package retry

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/uber-go/tally/v4"
	"go.uber.org/zap"
)

// RetryableFunc represents a function that is retryable.
type RetryableFunc func() error

// MultiError is a slice of errors ocurred while retrying.
type MultiError []error

// Error implements the error interface.
//
// Returns all the errors occurred during retry formatted as a string separated by
// newlines.
func (m MultiError) Error() string {
	errorLog := make([]string, len(m))
	for i, e := range m {
		errorLog[i] = fmt.Sprintf("Retry #%d: %s", i+1, e.Error())
	}

	return fmt.Sprintf("All retries failed:\n%s", strings.Join(errorLog, "\n"))
}

/*
Do retries a given function with delays between each attempts as defined by the caller
or their corresponding defaults.

err := retry.Do(

		ctx,
		logger,
		scope,
		func() error { return nil },
		Backoff(ExponentialBackoff),
	)
*/
func Do(ctx context.Context, logger *zap.Logger, scope tally.Scope, fn RetryableFunc, opts ...Option) error {
	if err := ctx.Err(); err != nil {
		logger.Error("context error", zap.Error(err))
		return err
	}

	var errorLog MultiError

	config := defaultConfig(scope)
	// apply user defined options.
	for _, opt := range opts {
		opt(config)
	}

	for n := uint(0); n < config.maxRetries; n++ {
		err := fn()
		if err == nil {
			return nil
		}

		logger.Error("attempt failed", zap.Uint("attempt", n+1), zap.Error(err))
		errorLog = append(errorLog, err)

		// Avoid waiting if this is the last attempt.
		if n == config.maxRetries-1 {
			break
		}

		// Increment the retry counter. First attempt is not
		// considered as a retry.
		config.retryCounter.Inc(1)

		delayTime := config.backoff(n, config)

		select {
		case <-time.After(delayTime):
			// Time for next retry.
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return errorLog
}
