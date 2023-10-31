package testutil

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/httprate"
	"github.com/k1LoW/rl"
)

type Limiter struct {
	counts      map[string]map[time.Time]int
	reqLimit    int
	keyFunc     httprate.KeyFunc
	statusCode  int
	ignoreAfter func(*http.Request) bool
}

func KeyByHost(r *http.Request) (string, error) {
	return r.Host, nil
}

func NewLimiter(reqLimit int, keyFunc httprate.KeyFunc, statusCode int) *Limiter {
	return &Limiter{
		counts:      map[string]map[time.Time]int{},
		reqLimit:    reqLimit,
		keyFunc:     keyFunc,
		statusCode:  statusCode,
		ignoreAfter: func(r *http.Request) bool { return false },
	}
}

func NewSkipper(hostname string) *Limiter {
	return &Limiter{
		counts:      map[string]map[time.Time]int{},
		reqLimit:    -1,
		keyFunc:     KeyByHost,
		statusCode:  0,
		ignoreAfter: func(r *http.Request) bool { return r.Host == hostname },
	}
}

func (l *Limiter) Name() string {
	return "testutil.Limiter"
}

func (l *Limiter) Rule(r *http.Request) (*rl.Rule, error) {
	key, err := l.keyFunc(r)
	if err != nil {
		return nil, err
	}
	return &rl.Rule{
		Key:         key,
		ReqLimit:    l.reqLimit,
		WindowLen:   time.Second,
		IgnoreAfter: l.ignoreAfter(r),
	}, nil
}

func (l *Limiter) ShouldSetXRateLimitHeaders(le *rl.LimitError) bool {
	return true
}

func (l *Limiter) OnRequestLimit(le *rl.LimitError) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if l.statusCode != 0 {
			w.WriteHeader(l.statusCode)
		} else {
			w.WriteHeader(http.StatusTooManyRequests)
		}
		msg := fmt.Sprintf("Too many requests. name: %s, ratelimit: %d req/%s, ratelimit-ramaining: %d, ratelimit-reset: %d", le.Limiter.Name(), le.RequestLimit, le.WindowLen, le.RateLimitRemaining, le.RateLimitReset)
		_, _ = w.Write([]byte(msg)) //nostyle:handlerrors
	}
}

func (l *Limiter) Get(key string, window time.Time) (count int, err error) { //nostyle:getters
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
