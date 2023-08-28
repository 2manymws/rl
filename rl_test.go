package rl_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/httprate"
	"github.com/k1LoW/rl"
	"github.com/k1LoW/rl/testutil"
)

var _ rl.Limiter = (*testutil.Limiter)(nil)

func TestRL(t *testing.T) {
	const noLimitReq = 100
	tests := []struct {
		name           string
		limiter        rl.Limiter
		hosts          []string
		wantReqCount   int
		wantStatusCode int
	}{
		{"key by ip", testutil.NewLimiter(10, httprate.KeyByIP, 0), []string{"a.example.com", "b.example.com"}, 10, http.StatusTooManyRequests},
		{"key by host", testutil.NewLimiter(10, testutil.KeyByHost, 0), []string{"a.example.com", "b.example.com"}, 20, http.StatusTooManyRequests},
		{"no limit", testutil.NewLimiter(-1, httprate.KeyByIP, 0), []string{"a.example.com", "b.example.com"}, noLimitReq, http.StatusTooManyRequests},
		{"set other statusCode", testutil.NewLimiter(10, httprate.KeyByIP, http.StatusOK), []string{"a.example.com", "b.example.com"}, 10, http.StatusOK},
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
					b, err := io.ReadAll(res.Body)
					if err != nil {
						t.Fatal(err)
					}
					defer res.Body.Close()
					if strings.Contains(string(b), "Too many requests") {
						if res.StatusCode != tt.wantStatusCode {
							t.Errorf("got %v want %v", res.StatusCode, tt.wantStatusCode)
						}
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

func BenchmarkRL(b *testing.B) {
	r := http.NewServeMux()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world"))
	})
	m := rl.New(
		testutil.NewLimiter(10, httprate.KeyByIP, 0),
		testutil.NewLimiter(10, testutil.KeyByHost, 0),
	)
	ts := httptest.NewServer(m(r))
	b.Cleanup(func() {
		ts.Close()
	})

	for i := 0; i < b.N; i++ {
		req, err := http.NewRequest("GET", ts.URL, nil)
		if err != nil {
			b.Fatal(err)
		}
		req.Host = "a.example.com"
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			b.Fatal(err)
		}
		res.Body.Close()
	}
}
