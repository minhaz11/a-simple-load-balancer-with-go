package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

type Server struct {
	URL     string
	healthy bool
}

type LoadBalancer struct {
	Servers []*Server
	idx     int
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

func (lb * LoadBalancer) RunHealthCheck(){
	ticker := time.NewTicker(10 * time.Second)

	go func() {
		for range ticker.C {
			lb.ServerHealthCheck()
		}
	}()
}

func (lb *LoadBalancer) NextServer() *Server {
	server := lb.Servers[lb.idx]

	lb.idx = (lb.idx + 1) % len(lb.Servers)

	return server
}

func (lb *LoadBalancer) ForwardRequest(server *Server, uri string) (*http.Response, error) {
	fmt.Printf("Forwarding request to : %s\n", server.URL)

	parsedUrl, err := url.Parse(server.URL)

	if err != nil {
		log.Fatal("Error parsing server url : ", err.Error())
		return nil, err
	}

	fullUrl := parsedUrl.ResolveReference(&url.URL{Path: uri})

	response, err := http.Get(fullUrl.String())

	if err != nil {
		log.Fatal("Error sending request to server : ", fullUrl)
		return nil, err
	}

	return response, nil
}

func (lb *LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	nextServer := lb.NextServer()
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
