package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"

	"github.com/go-chi/httprate"
	"github.com/k1LoW/rl"
	"github.com/k1LoW/rl/testutil"
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
	ts := httptest.NewServer(m(r))
	defer ts.Close()

	log.Printf("start %d requests", n)
	for i := 0; i < n; i++ {
		req, err := http.NewRequest("GET", ts.URL, nil)
		if err != nil {
			log.Fatal(err)
		}
		req.Host = "a.example.com"
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		res.Body.Close()
	}
	log.Println("done")
}
