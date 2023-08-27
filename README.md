# rl [![Go Reference](https://pkg.go.dev/badge/github.com/k1LoW/rl.svg)](https://pkg.go.dev/github.com/k1LoW/rl) ![Coverage](https://raw.githubusercontent.com/k1LoW/octocovs/main/badges/k1LoW/rl/coverage.svg) ![Code to Test Ratio](https://raw.githubusercontent.com/k1LoW/octocovs/main/badges/k1LoW/rl/ratio.svg) ![Test Execution Time](https://raw.githubusercontent.com/k1LoW/octocovs/main/badges/k1LoW/rl/time.svg)

`rl` is a rate limit middleware for multiple limit rules.

## Usage

Prepare an instance that satisfies [`rl.Limiter`](https://pkg.go.dev/github.com/k1LoW/rl#Limiter) interface.

Then, generate the middleware ( `func(next http.Handler) http.Handler` ) with [`rl.New`](https://pkg.go.dev/github.com/k1LoW/rl#New)

```go
package main

import (
    "log"
    "net/http"

    "github.com/k1LoW/rl"
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

## Reference

- [go-chi/httprate](https://github.com/go-chi/httprate)
    - **Most of rl's rate limit implementations refer to httprate**
