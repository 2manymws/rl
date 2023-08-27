package testutil

import (
	"net/http"
	"time"

	"github.com/go-chi/httprate"
)

type Limiter struct {
	counts   map[string]map[time.Time]int
	reqLimit int
	keyFunc  httprate.KeyFunc
}

func KeyByHost(r *http.Request) (string, error) {
	return r.Host, nil
}

func NewLimiter(reqLimit int, keyFunc httprate.KeyFunc) *Limiter {
	return &Limiter{
		counts:   map[string]map[time.Time]int{},
		reqLimit: reqLimit,
		keyFunc:  keyFunc,
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

func (l *Limiter) OnRequestLimit() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
		w.Write([]byte("Too many requests"))
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
