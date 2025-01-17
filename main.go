package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {

	cfg, err := LoadConfig()

	if err != nil {
		panic(err)
	}

	opts := NewOpts().MaxConnections(100).SetTimeout(15 * time.Second)

	lb := NewLoadBalancer(cfg.Servers, opts)

	lb.RunHealthCheck()

	fmt.Printf("Load balancer listening on port :%d\n", cfg.Port)

	err = http.ListenAndServe(":8080", lb)

	if err != nil {
		panic(err)
	}

}
