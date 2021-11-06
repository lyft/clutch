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
		logger.Error("context error", zap.String("err", err.Error()))
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
		// Increment retry counter.
		config.retryCounter.Inc(1)

		err := fn()
		if err != nil {
			logger.Error("retry failed", zap.Uint("attempt #", n+1), zap.String("err", err.Error()))
			errorLog = append(errorLog, err)

			// Avoid waiting if this is the last attempt.
			if n == config.maxRetries - 1 {
				break
			}

			delayTime := config.backoff(n, config)

			select {
			case <-time.After(delayTime):
			case <-ctx.Done():
				return ctx.Err()
			}
		} else {
			return nil
		}
		n++
	}
	return errorLog
}
