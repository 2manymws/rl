package rl

import (
	"errors"
	"net/http"
	"time"
)

var _ error = &Context{}
var ErrRateLimitExceeded error = errors.New("rate limit exceeded")

// Context is the error returned by the middleware.
type Context struct {
	StatusCode         int
	Err                error
	Limiter            Limiter
	RequestLimit       int
	WindowLen          time.Duration
	RateLimitRemaining int
	RateLimitReset     int
	Next               http.Handler
	lh                 *limitHandler
}

func newContext(statusCode int, err error, lh *limitHandler, next http.Handler) *Context {
	return &Context{
		StatusCode:         statusCode,
		Err:                err,
		Limiter:            lh.limiter,
		RequestLimit:       lh.reqLimit,
		WindowLen:          lh.windowLen,
		RateLimitRemaining: lh.rateLimitRemaining,
		RateLimitReset:     lh.rateLimitReset,
		Next:               next,
		lh:                 lh,
	}
}

// Error returns the error message.
func (e *Context) Error() string {
	return e.Err.Error()
}
