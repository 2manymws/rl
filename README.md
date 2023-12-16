# rl [![Go Reference](https://pkg.go.dev/badge/github.com/2manymws/rl.svg)](https://pkg.go.dev/github.com/2manymws/rl) [![build](https://github.com/2manymws/rl/actions/workflows/ci.yml/badge.svg)](https://github.com/2manymws/rl/actions/workflows/ci.yml) ![Coverage](https://raw.githubusercontent.com/k1LoW/octocovs/main/badges/2manymws/rl/coverage.svg) ![Code to Test Ratio](https://raw.githubusercontent.com/k1LoW/octocovs/main/badges/2manymws/rl/ratio.svg) ![Test Execution Time](https://raw.githubusercontent.com/k1LoW/octocovs/main/badges/2manymws/rl/time.svg)

`rl` is a **r**ate **l**imit middleware for multiple limit rules.

## Usage

Prepare an instance that implements [`rl.Limiter`](https://pkg.go.dev/github.com/2manymws/rl#Limiter) interface.

Then, generate the middleware ( `func(next http.Handler) http.Handler` ) with [`rl.New`](https://pkg.go.dev/github.com/2manymws/rl#New)

```go
package main

import (
    "log"
    "net/http"

    "github.com/2manymws/rl"
)

func main() {
    r := http.NewServeMux()
    r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello World"))
    })

    var l rl.Limiter = newMyLimiter()
    m := rl.New(l)

    log.Fatal(http.ListenAndServe(":8080", m(r)))
}
```

## Rate limiting approach

`rl` uses the Sliding Window Counter pattern same as [go-chi/httprate](https://github.com/go-chi/httprate).

- https://blog.cloudflare.com/counting-things-a-lot-of-different-things/
- https://www.figma.com/blog/an-alternative-approach-to-rate-limiting/

## Reference

- [go-chi/httprate](https://github.com/go-chi/httprate)
    - **Most of `rl`'s rate limit implementations refer to httprate. Thanks for the simple and clean implementation!**
