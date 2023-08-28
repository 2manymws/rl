package rl

import (
	"errors"
	"fmt"
	"math"
	"net/http"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

type Limiter interface {
	// Name returns the name of the limiter
	Name() string
	// KeyAndRateLimit returns the key and rate limit for the request
	// If the rate limit is negative, the limiter is skipped
	KeyAndRateLimit(r *http.Request) (key string, reqLimit int, windowLen time.Duration, err error)
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

type Config struct {
	Limiters []Limiter
	// Skipper is the function to skip middleware
	Skipper Skipper
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

type Skipper func(r *http.Request) bool
type rateLimiter struct {
	limiters []Limiter
	skipper  Skipper
}

func newRateLimiter(limiters []Limiter, skipper Skipper) *rateLimiter {
	return &rateLimiter{
		limiters: limiters,
		skipper:  skipper,
	}
}

func (rl *rateLimiter) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if rl.skipper != nil && rl.skipper(r) {
			next.ServeHTTP(w, r)
			return
		}
		now := time.Now().UTC()
		var lastLH *limitHandler
		eg := new(errgroup.Group)
		for _, limiter := range rl.limiters {
			key, reqLimit, windowLen, err := limiter.KeyAndRateLimit(r)
			if reqLimit < 0 {
				// If the request limit is negative, skip this limiter
				continue
			}
			if err != nil {
				http.Error(w, err.Error(), http.StatusPreconditionRequired)
				return
			}
			lh := &limitHandler{
				key:       key,
				reqLimit:  reqLimit,
				windowLen: windowLen,
				limiter:   limiter,
			}
			lastLH = lh
			eg.Go(func() error {
				lh.mu.Lock()
				defer lh.mu.Unlock()
				currentWindow := now.Truncate(lh.windowLen)

				lh.rateLimitRemaining = 0
				lh.rateLimitReset = int(currentWindow.Add(lh.windowLen).Unix())

				_, rate, err := lh.status(now, currentWindow)
				if err != nil {
					return newLimitError(http.StatusPreconditionRequired, err, lh)
				}
				nrate := int(math.Round(rate))
				if lh.reqLimit > nrate {
					lh.rateLimitRemaining = lh.reqLimit - nrate
				}

				if nrate >= lh.reqLimit {
					return newLimitError(http.StatusTooManyRequests, ErrRateLimitExceeded, lh)
				}

				if err := lh.limiter.Increment(lh.key, currentWindow); err != nil {
					return newLimitError(http.StatusInternalServerError, err, lh)
				}
				return nil
			})
		}

		// Wait for all limiters to finish
		if err := eg.Wait(); err != nil {
			// Handle first error
			if e, ok := err.(*LimitError); ok {
				http.Error(w, e.Error(), e.statusCode)
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
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if lastLH != nil {
			// Set X-RateLimit-* headers using the last limiter
			if lastLH.limiter.ShouldSetXRateLimitHeaders(nil) {
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
	rl := newRateLimiter(limiters, nil)
	return rl.Handler
}

func NewWithConfig(c *Config) func(next http.Handler) http.Handler {
	rl := newRateLimiter(c.Limiters, c.Skipper)
	return rl.Handler
}
