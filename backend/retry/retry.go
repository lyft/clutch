package retry

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/uber-go/tally"
	"go.uber.org/zap"
)

// RetryableFunc represents a function that is retryable.
type RetryableFunc func() error

// MultiError is a slice of errors ocurred while retrying.
type MultiError []error

// Error implements error interface.
//
// Returns all the errors ocurred during retry formatted as a string.
func (m MultiError) Error() string {
	errorLog := make([]string, len(m))
	for i, e := range m {
		errorLog[i] = fmt.Sprintf("Retry #%d: %s", i+1, e.Error())
	}

	return fmt.Sprintf("All retries failed:\n%s", strings.Join(errorLog, "\n"))
}

/*
Do retries a given function with delays between each attempts as defined by the caller.

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

	n := uint(0)
	var errorLog MultiError

	config := defaultConfig(scope)
	// apply user defined options.
	for _, opt := range opts {
		opt(config)
	}

	for n < config.maxRetries {
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

		n++
	}
	return errorLog
}
