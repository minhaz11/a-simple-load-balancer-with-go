package main

import "net/http"

func main() {

	servers := []*Server{
		{
			URL: "http://127.0.0.1:8000",
		},

		{
			URL: "http://127.0.0.1:8001",
		},
	}

	loadBalancer := &LoadBalancer{
		Servers: servers,
	}

	loadBalancer.RunHealthCheck()

	err := http.ListenAndServe(":8080", loadBalancer)

	if err != nil {
		panic(err)
	}

}
