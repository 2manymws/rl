package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"

	"github.com/2manymws/rl"
	"github.com/2manymws/rl/testutil"
	"github.com/go-chi/httprate"
)

func main() {
	n, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	r := http.NewServeMux()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("Hello, world"))
	})
	m := rl.New(
		testutil.NewLimiter(10, httprate.KeyByIP, 0),
		testutil.NewLimiter(10, testutil.KeyByHost, 0),
	)
	router := m(r)

	log.Printf("start %d requests", n)
	for i := 0; i < n; i++ {
		req, err := http.NewRequest("GET", "http://example.com", nil)
		if err != nil {
			log.Fatal(err)
		}
		res := httptest.NewRecorder()
		router.ServeHTTP(res, req)
	}
	log.Println("done")
}
