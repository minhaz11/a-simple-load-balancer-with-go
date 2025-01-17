package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type Server struct {
	URL     string
	healthy bool
}

type LoadBalancer struct {
	Servers []*Server
	mu      sync.Mutex
	cp      *ConnectionPool
	idx     int
}

func NewLoadBalancer(urls []string, opts *Opts) *LoadBalancer {
	servers := make([]*Server, len(urls))

	for i, url := range urls {
		servers[i] = &Server{
			URL:     url,
			healthy: true,
		}
	}

	return &LoadBalancer{
		Servers: servers,
		cp:      NewConnectionPool(opts),
	}

}

func (lb *LoadBalancer) hasUnhealthyServer() bool {
	for _, server := range lb.Servers {
		if !server.healthy {
			return true
		}
	}

	return false
}

func (lb *LoadBalancer) ServerHealthCheck() {
	for _, server := range lb.Servers {
		res, err := http.Get(server.URL + "/heathcheck")

		if err != nil || res.StatusCode != http.StatusOK {
			server.healthy = false
			fmt.Printf("Server [%s] is down\n", server.URL)
		} else {
			server.healthy = true
			fmt.Printf("Server [%s] is up and running\n", server.URL)
		}
	}
}

func (lb *LoadBalancer) RunHealthCheck() {
	ticker := time.NewTicker(10 * time.Second)

	go func() {
		for range ticker.C {
			lb.ServerHealthCheck()
		}
	}()
}

func (lb *LoadBalancer) NextServer() (*Server, error) {

	lb.mu.Lock()

	defer lb.mu.Unlock()

	if lb.hasUnhealthyServer() {
		for idx := 0; idx < len(lb.Servers); idx++ {
			if lb.Servers[idx].healthy {
				break
			}

			lb.idx = idx
		}
	}

	if lb.idx == len(lb.Servers) {
		lb.idx = 0

		return nil, errors.New("no healthy server found")
	}

	server := lb.Servers[lb.idx]

	lb.idx = (lb.idx + 1) % len(lb.Servers)

	return server, nil
}

func (lb *LoadBalancer) ForwardRequest(server *Server, uri string) (*http.Response, error) {
	fmt.Printf("Forwarding request to : %s\n", server.URL)

	client := lb.cp.Get(server.URL)
	defer lb.cp.Put(server.URL, client)

	parsedUrl, err := url.Parse(server.URL)

	if err != nil {
		log.Fatal("Error parsing server url : ", err.Error())
		return nil, err
	}

	fullUrl := parsedUrl.ResolveReference(&url.URL{Path: uri})

	response, err := client.Get(fullUrl.String())

	if err != nil {
		log.Fatal("Error sending request to server : ", fullUrl)
		return nil, err
	}

	return response, nil
}

func (lb *LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	nextServer, err := lb.NextServer()
	if err != nil {
		panic(err)
	}

	res, err := lb.ForwardRequest(nextServer, r.RequestURI)

	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if err != nil {
		panic(err)
	}

	_, err = w.Write(body)

	if err != nil {
		panic(err)
	}

}
