package rl

import "errors"

var _ error = &LimitError{}
var ErrRateLimitExceeded error = errors.New("rate limit exceeded")

type LimitError struct {
	statusCode int
	err        error
	lh         *limitHandler
}

func newLimitError(statusCode int, err error, lh *limitHandler) *LimitError {
	return &LimitError{
		statusCode: statusCode,
		err:        err,
		lh:         lh,
	}
}

func (e *LimitError) Error() string {
	return e.err.Error()
}
