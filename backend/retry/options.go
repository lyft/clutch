package retry

import (
	"crypto/rand"
	"math/big"
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
		maxDelay:     100 * time.Millisecond,
		maxJitter:    100 * time.Millisecond,
		backoff:      DefaultBackoff,
		retryCounter: scope.Counter("retry_attempts"),
	}
}

// Option is an option for retry.
type Option func(*Config)

// Retries sets the total number of attempts.
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

// Jitter sets the max jitter.
// Defaults to 100ms if provided value <= 0.
func Jitter(jitter time.Duration) Option {
	return func(c *Config) {
		c.maxJitter = jitter
		if c.maxJitter <= 0 {
			c.maxJitter = 100 * time.Millisecond
		}
	}
}

// DefaultBackoff always returns a fixed delay duration.
func DefaultBackoff(_ uint, c *Config) time.Duration {
	return c.maxDelay
}

// ExponentialBackoff returns exponentially increasing backoffs by a power of 2.
func ExponentialBackoff(i uint, _ *Config) time.Duration {
	return 1 << i * time.Second
}

// JitterBackoff returns a random delay duration up to maxJitter.
func JitterBackoff(_ uint, c *Config) time.Duration {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(c.maxJitter)))
	if err != nil {
		return c.maxJitter
	}
	return time.Duration(n.Int64())
}
