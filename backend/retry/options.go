package retry

import (
	"math/rand"
	"time"

	"github.com/uber-go/tally"
)

// BackoffStrategy is used to determine the delay between retries.
type BackoffStrategy func(uint, *Config) time.Duration

// Config defines the behavior of retries.
type Config struct {
	maxRetries   uint
	maxDelay     time.Duration
	maxJitter    time.Duration
	backoff      BackoffStrategy
	retryCounter tally.Counter
}

// defaultConfig returns retry defaults.
func defaultConfig(scope tally.Scope) *Config {
	return &Config{
		maxRetries:   uint(3),
		maxDelay:     1 * time.Second,
		maxJitter:    100 * time.Millisecond,
		backoff:      DefaultBackoff,
		retryCounter: scope.Counter("retry_attempts"),
	}
}

// Option is an option for retry.
type Option func(*Config)

// Retries sets the number of retries.
func Retries(count uint) Option {
	return func(c *Config) {
		c.maxRetries = count
	}
}

// Backoff sets the backoff mechanism used while retrying.
func Backoff(backoff BackoffStrategy) Option {
	return func(c *Config) {
		c.backoff = backoff
	}
}

// Delay sets the delay between each retries.
func Delay(delay time.Duration) Option {
	return func(c *Config) {
		c.maxDelay = delay
	}
}

// DefaultBackoff always returns 1 second.
func DefaultBackoff(_ uint, c *Config) time.Duration {
	return c.maxDelay
}

// ExponentialBackoff returns exponentially increasing backoffs by a power of 2.
func ExponentialBackoff(i uint, _ *Config) time.Duration {
	return time.Duration(1<<i) * time.Second
}

// FixedBackoff returns a fixed delay duration.
func FixedBackoff(_ uint, c *Config) time.Duration {
	return c.maxDelay
}

// JitterBackoff returns a random delay duration up to maxJitter.
func JitterBackoff(_ uint, c *Config) time.Duration {
	return time.Duration(rand.Int63n(int64(c.maxJitter)))
}
