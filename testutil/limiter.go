package testutil

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/httprate"
	"github.com/k1LoW/rl"
)

type Limiter struct {
	counts     map[string]map[time.Time]int
	reqLimit   int
	keyFunc    httprate.KeyFunc
	statusCode int
}

func KeyByHost(r *http.Request) (string, error) {
	return r.Host, nil
}

func NewLimiter(reqLimit int, keyFunc httprate.KeyFunc, statusCode int) *Limiter {
	return &Limiter{
		counts:     map[string]map[time.Time]int{},
		reqLimit:   reqLimit,
		keyFunc:    keyFunc,
		statusCode: statusCode,
	}
}

func (l *Limiter) Name() string {
	return "testutil.Limiter"
}

func (l *Limiter) KeyAndRateLimit(r *http.Request) (string, int, time.Duration, error) {
	key, err := l.keyFunc(r)
	if err != nil {
		return "", 0, 0, err
	}
	return key, l.reqLimit, time.Second, nil
}

func (l *Limiter) ShouldSetXRateLimitHeaders(err error) bool {
	return true
}

func (l *Limiter) OnRequestLimit(err error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
		le, ok := err.(*rl.LimitError)
		if !ok {
			w.Write([]byte("Too many requests"))
			return
		}
		msg := fmt.Sprintf("Too many requests. name: %s, ratelimit: %d req/%s, ratelimit-ramaining: %d, ratelimit-reset: %d", le.LimierName(), le.RequestLimit(), le.WindowLen(), le.RateLimitRemaining(), le.RateLimitReset())
		w.Write([]byte(msg))
		if l.statusCode != 0 {
			w.WriteHeader(l.statusCode)
		}
	}
}

func (l *Limiter) Get(key string, window time.Time) (count int, err error) {
	if _, ok := l.counts[key]; !ok {
		l.counts[key] = map[time.Time]int{}
		return 0, nil
	}
	if _, ok := l.counts[key][window]; !ok {
		l.counts[key][window] = 0
		return 0, nil
	}
	return l.counts[key][window], nil
}

func (l *Limiter) Increment(key string, currentWindow time.Time) error {
	if _, ok := l.counts[key]; !ok {
		l.counts[key] = map[time.Time]int{}
	}
	if _, ok := l.counts[key][currentWindow]; !ok {
		l.counts[key][currentWindow] = 0
	}
	l.counts[key][currentWindow]++
	return nil
}
