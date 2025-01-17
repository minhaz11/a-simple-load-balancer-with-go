package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

type Opts struct {
	maxConnections int
	timeout        time.Duration
}

func NewOpts() *Opts {
	return &Opts{
		maxConnections: 10,
		timeout: 10 * time.Second,
	}
}

func (opts *Opts) MaxConnections(maxConns int) *Opts {
	opts.maxConnections = maxConns

	return opts
}

func (opts *Opts) SetTimeout(timeout time.Duration) *Opts {
	opts.timeout = timeout

	return opts
}


//connection pools

type ConnectionPool struct {
	*Opts
	clients map[string][]*http.Client
	mu sync.Mutex
}


func NewConnectionPool(opts *Opts) *ConnectionPool  {
	return &ConnectionPool{
		Opts: opts,
		clients: make(map[string][]*http.Client),
	}
}

func (cp *ConnectionPool) Get(serverAddr string) *http.Client {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	clients, ok := cp.clients[serverAddr]

	if ok && len(clients) > 0 {
		client := clients[len(clients)-1]
		clients = clients[:len(clients)-1]

		cp.clients[serverAddr] = clients

		return client
	}

	return &http.Client{
		Timeout: cp.timeout,
	}
}

func (cp *ConnectionPool) Put(serverAddr string, client *http.Client) error {
	cp.mu.Lock()
	defer cp.mu.Unlock()
	
	if len(cp.clients[serverAddr]) > cp.maxConnections {
		return fmt.Errorf("connection pool limit exceeded for serer '%s'", serverAddr)
	}

	cp.clients[serverAddr] = append(cp.clients[serverAddr], client)

	return nil
}
