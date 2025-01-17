package main

import (
	"net/http"
	"time"
)

func main() {

	urls := []string{"http://127.0.0.1:8000", "http://127.0.0.1:8001"}

	opts := NewOpts().MaxConnections(100).SetTimeout(15 * time.Second)

	lb := NewLoadBalancer(urls, opts)

	lb.RunHealthCheck()

	err := http.ListenAndServe(":8080", lb)

	if err != nil {
		panic(err)
	}

}
