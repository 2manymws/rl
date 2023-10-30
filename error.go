package rl

import (
	"errors"
	"time"
)

type LimitError interface {
	error
	LimierName() string
	RequestLimit() int
	WindowLen() time.Duration
	RateLimitReset() int
	RateLimitRemaining() int
}

var _ LimitError = &limitError{}

var ErrRateLimitExceeded error = errors.New("rate limit exceeded")

// limitError is the error returned by the middleware.
type limitError struct {
	statusCode int
	err        error
	lh         *limitHandler
}

func newLimitError(statusCode int, err error, lh *limitHandler) *limitError {
	return &limitError{
		statusCode: statusCode,
		err:        err,
		lh:         lh,
	}
}

// Error returns the error message.
func (e *limitError) Error() string {
	return e.err.Error()
}

// LimierName returns the name of the limiter that caused the error.
func (e *limitError) LimierName() string {
	return e.lh.limiter.Name()
}

// RequestLimit returns the request limit of the limiter that caused the error.
func (e *limitError) RequestLimit() int {
	return e.lh.reqLimit
}

// WindowLen returns the window length of the limiter that caused the error.
func (e *limitError) WindowLen() time.Duration {
	return e.lh.windowLen
}

// RateLimitReset returns the number of seconds until the rate limit resets.
func (e *limitError) RateLimitReset() int {
	return e.lh.rateLimitReset
}

// RateLimitRemaining returns the number of requests remaining in the current window.
func (e *limitError) RateLimitRemaining() int {
	return e.lh.rateLimitRemaining
}
