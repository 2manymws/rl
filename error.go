package rl

import (
	"errors"
	"time"
)

var _ error = &LimitError{}
var ErrRateLimitExceeded error = errors.New("rate limit exceeded")

// LimitError is the error returned by the middleware.
type LimitError struct {
	statusCode int
	err        error
	lh         *limitHandler
}

func NewLimitError(statusCode int, err error, lh *limitHandler) *LimitError {
	return &LimitError{
		statusCode: statusCode,
		err:        err,
		lh:         lh,
	}
}

// Error returns the error message.
func (e *LimitError) Error() string {
	return e.err.Error()
}

// LimierName returns the name of the limiter that caused the error.
func (e *LimitError) LimierName() string {
	return e.lh.limiter.Name()
}

// RequestLimit returns the request limit of the limiter that caused the error.
func (e *LimitError) RequestLimit() int {
	return e.lh.reqLimit
}

// WindowLen returns the window length of the limiter that caused the error.
func (e *LimitError) WindowLen() time.Duration {
	return e.lh.windowLen
}

// RateLimitReset returns the number of seconds until the rate limit resets.
func (e *LimitError) RateLimitReset() int {
	return e.lh.rateLimitReset
}

// RateLimitRemaining returns the number of requests remaining in the current window.
func (e *LimitError) RateLimitRemaining() int {
	return e.lh.rateLimitRemaining
}
