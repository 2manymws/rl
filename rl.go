package rl

import (
	"context"
	"errors"
	"fmt"
	"math"
	"net/http"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

// Rule is a rate limit rule
type Rule struct {
	// Key for the rate limit
	Key string
	// ReqLimit is the request limit for the window
	// If ReqLimit is negative, the limiter is skipped
	ReqLimit int
	// WindowLen is the length of the window
	WindowLen time.Duration
	// IgnoreAfter is true if skip all limiters after this limiter
	IgnoreAfter bool
}

type Limiter interface {
	// Name returns the name of the limiter
	Name() string
	// Rule returns the key and rate limit rule for the request
	Rule(r *http.Request) (rule *Rule, err error)
	// ShouldSetXRateLimitHeaders returns true if the X-RateLimit-* headers should be set
	ShouldSetXRateLimitHeaders(err error) bool
	// OnRequestLimit returns the handler to be called when the rate limit is exceeded
	OnRequestLimit(err error) http.HandlerFunc

	// Get returns the current count for the key and window
	Get(key string, window time.Time) (count int, err error)
	// Increment increments the count for the key and window
	Increment(key string, currentWindow time.Time) error
}

type limitHandler struct {
	key                string
	reqLimit           int
	windowLen          time.Duration
	limiter            Limiter
	rateLimitRemaining int
	rateLimitReset     int
	mu                 sync.Mutex
}

func (lh *limitHandler) status(now, currentWindow time.Time) (bool, float64, error) {
	previousWindow := currentWindow.Add(-lh.windowLen)

	currCount, err := lh.limiter.Get(lh.key, currentWindow)
	if err != nil {
		return false, 0, err
	}
	prevCount, err := lh.limiter.Get(lh.key, previousWindow)
	if err != nil {
		return false, 0, err
	}

	diff := now.Sub(currentWindow)
	rate := float64(prevCount)*(float64(lh.windowLen)-float64(diff))/float64(lh.windowLen) + float64(currCount)
	if rate > float64(lh.reqLimit) {
		return false, rate, nil
	}
	return true, rate, nil
}

type limitMw struct {
	limiters []Limiter
}

func newLimitMw(limiters []Limiter) *limitMw {
	return &limitMw{
		limiters: limiters,
	}
}

func (rl *limitMw) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		now := time.Now().UTC()
		var lastLH *limitHandler
		eg, ctx := errgroup.WithContext(context.Background())
		for _, limiter := range rl.limiters {
			rule, err := limiter.Rule(r)
			if rule.ReqLimit < 0 {
				// If the request limit is negative, skip this limiter
				if rule.IgnoreAfter {
					// Skip all limiters after this limiter.
					break
				}
				continue
			}
			if err != nil {
				http.Error(w, err.Error(), http.StatusPreconditionRequired)
				return
			}
			lh := &limitHandler{
				key:       rule.Key,
				reqLimit:  rule.ReqLimit,
				windowLen: rule.WindowLen,
				limiter:   limiter,
			}
			lastLH = lh
			eg.Go(func() error {
				lh.mu.Lock()
				defer lh.mu.Unlock()
				currentWindow := now.Truncate(lh.windowLen)
				lh.rateLimitRemaining = 0
				lh.rateLimitReset = int(currentWindow.Add(lh.windowLen).Unix())
				select {
				// Check if the request limit already exceeded before calling lh.status()
				case <-ctx.Done():
					// Increment must be called even if the request limit is already exceeded
					if err := lh.limiter.Increment(lh.key, currentWindow); err != nil {
						return newLimitError(http.StatusInternalServerError, err, lh)
					}
					return nil
				default:
				}
				_, rate, err := lh.status(now, currentWindow)
				if err != nil {
					return newLimitError(http.StatusPreconditionRequired, err, lh)
				}
				nrate := int(math.Round(rate))
				if nrate >= lh.reqLimit {
					return newLimitError(http.StatusTooManyRequests, ErrRateLimitExceeded, lh)
				}

				lh.rateLimitRemaining = lh.reqLimit - nrate
				if err := lh.limiter.Increment(lh.key, currentWindow); err != nil {
					return newLimitError(http.StatusInternalServerError, err, lh)
				}
				return nil
			})
			if rule.IgnoreAfter {
				// Skip all limiters after this limiter.
				break
			}
		}

		// Wait for all limiters to finish
		if err := eg.Wait(); err != nil {
			// Handle first error
			if e, ok := err.(*LimitError); ok {
				if e.lh.limiter.ShouldSetXRateLimitHeaders(err) {
					w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", e.lh.reqLimit))
					w.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", e.lh.rateLimitRemaining))
					w.Header().Set("X-RateLimit-Reset", fmt.Sprintf("%d", e.lh.rateLimitReset))
				}
				if errors.Is(e.err, ErrRateLimitExceeded) {
					// Rate limit exceeded
					if e.lh.limiter.ShouldSetXRateLimitHeaders(err) {
						w.Header().Set("Retry-After", fmt.Sprintf("%d", int(e.lh.windowLen.Seconds()))) // RFC 6585
					}
					e.lh.limiter.OnRequestLimit(e)(w, r)
					return
				}
				http.Error(w, e.Error(), e.statusCode)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if lastLH != nil {
			// Set X-RateLimit-* headers using the last limiter
			if lastLH.limiter.ShouldSetXRateLimitHeaders(nil) && lastLH.reqLimit >= 0 {
				w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", lastLH.reqLimit))
				w.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", lastLH.rateLimitRemaining))
				w.Header().Set("X-RateLimit-Reset", fmt.Sprintf("%d", lastLH.rateLimitReset))
			}
		}

		next.ServeHTTP(w, r)
	})
}

// New returns a new rate limiter middleware.
// The order of the limitters should be arranged in **reverse** order of Limitter with strict rate limit to return appropriate X-RateLimit-* headers to the client.
func New(limiters ...Limiter) func(next http.Handler) http.Handler {
	rl := newLimitMw(limiters)
	return rl.Handler
}
