package wait

import "time"

type Config struct {
	timeout           time.Duration
	perAttemptTimeout time.Duration
	retryMaxDelay     time.Duration
}

func getConfig(options []Option) Config {
	cfg := &Config{}
	for _, opt := range options {
		opt(cfg)
	}
	return *cfg
}

type Option func(*Config)

func Timeout(timeout time.Duration) Option {
	return func(c *Config) {
		c.timeout = timeout
	}
}

func PerAttemptTimeout(timeout time.Duration) Option {
	return func(c *Config) {
		c.perAttemptTimeout = timeout
	}
}

func RetryMaxDelay(maxDelay time.Duration) Option {
	return func(c *Config) {
		c.retryMaxDelay = maxDelay
	}
}
