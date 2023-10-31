package rl

import (
	"errors"
	"time"
)

var _ error = &LimitError{}
var ErrRateLimitExceeded error = errors.New("rate limit exceeded")

// LimitError is the error returned by the middleware.
type LimitError struct {
	StatusCode         int
	Err                error
	LimiterName        string
	RequestLimit       int
	WindowLen          time.Duration
	RateLimitRemaining int
	RateLimitReset     int
	lh                 *limitHandler
}

func newLimitError(statusCode int, err error, lh *limitHandler) *LimitError {
	return &LimitError{
		StatusCode:         statusCode,
		Err:                err,
		LimiterName:        lh.limiter.Name(),
		RequestLimit:       lh.reqLimit,
		WindowLen:          lh.windowLen,
		RateLimitRemaining: lh.rateLimitRemaining,
		RateLimitReset:     lh.rateLimitReset,
		lh:                 lh,
	}
}

// Error returns the error message.
func (e *LimitError) Error() string {
	return e.Err.Error()
}
