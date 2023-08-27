package rl_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/httprate"
	"github.com/k1LoW/rl"
	"github.com/k1LoW/rl/testutil"
)

func TestRL(t *testing.T) {
	tests := []struct {
		name         string
		keyFunc      httprate.KeyFunc
		reqLimit     int
		hosts        []string
		wantReqCount int
	}{
		{"key by ip", httprate.KeyByIP, 10, []string{"a.example.com", "b.example.com"}, 10},
		{"key by host", testutil.KeyByHost, 10, []string{"a.example.com", "b.example.com"}, 20},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := http.NewServeMux()
			r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("Hello, world"))
			})
			l := testutil.NewLimiter(tt.reqLimit, tt.keyFunc)
			m := rl.New(l)
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
				}
			}
			if got != tt.wantReqCount {
				t.Errorf("got %v want %v", got, tt.wantReqCount)
			}
		})
	}
}
