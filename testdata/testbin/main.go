package main

import (
	"log"
	"os"
	"strconv"
)

func main() {
	n, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < n; i++ {
		var s int
		for i := 0; i < 10000; i++ {
			s += i
		}
	}
}
