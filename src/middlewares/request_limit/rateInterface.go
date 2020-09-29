package request_limit

import "time"

type RateLimiter interface {
	Acquire() (*Token, error)
	Release(*Token)
	IsLimitExceeded() bool
}

type Config struct {
	Limit int
	Throttle time.Duration
	TokenResetAfter time.Duration
}

