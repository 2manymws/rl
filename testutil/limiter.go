package testutil

import (
	"fmt"
	"net/http"
	"time"

	"github.com/2manymws/rl"
	"github.com/go-chi/httprate"
)

type Limiter struct {
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
		reqLimit:    reqLimit,
		keyFunc:     keyFunc,
		statusCode:  statusCode,
		ignoreAfter: func(r *http.Request) bool { return false },
	}
}

func NewSkipper(hostname string) *Limiter {
	return &Limiter{
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

func (l *Limiter) ShouldSetXRateLimitHeaders(le *rl.Context) bool {
	return true
}

func (l *Limiter) OnRequestLimit(le *rl.Context) http.HandlerFunc {
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
