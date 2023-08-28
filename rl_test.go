package rl_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/httprate"
	"github.com/k1LoW/rl"
	"github.com/k1LoW/rl/testutil"
)

var _ rl.Limiter = (*testutil.Limiter)(nil)

func TestRL(t *testing.T) {
	const noLimitReq = 100
	tests := []struct {
		name         string
		limiter      rl.Limiter
		hosts        []string
		wantReqCount int
	}{
		{"key by ip", testutil.NewLimiter(10, httprate.KeyByIP), []string{"a.example.com", "b.example.com"}, 10},
		{"key by host", testutil.NewLimiter(10, testutil.KeyByHost), []string{"a.example.com", "b.example.com"}, 20},
		{"no limit", testutil.NewLimiter(-1, httprate.KeyByIP), []string{"a.example.com", "b.example.com"}, noLimitReq},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := http.NewServeMux()
			r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("Hello, world"))
			})
			m := rl.New(tt.limiter)
			ts := httptest.NewServer(m(r))
			t.Cleanup(func() {
				ts.Close()
			})
			got := 0
		L:
			for {
				for _, host := range tt.hosts {
					req, err := http.NewRequest("GET", ts.URL, nil)
					if err != nil {
						t.Fatal(err)
					}
					req.Host = host
					res, err := http.DefaultClient.Do(req)
					if err != nil {
						t.Fatal(err)
					}
					if res.StatusCode == http.StatusTooManyRequests {
						break L
					}
					got++
					if got == noLimitReq { // circuit breaker
						break L
					}
				}
			}
			if got != tt.wantReqCount {
				t.Errorf("got %v want %v", got, tt.wantReqCount)
			}
		})
	}
}
