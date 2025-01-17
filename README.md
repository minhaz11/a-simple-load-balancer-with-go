# Go Load Balancer

A simple HTTP load balancer implementation in Go that distributes incoming HTTP requests across multiple backend servers using a round-robin algorithm.

## Features

- Round-robin load balancing
- Automatic health checks for backend servers
- Connection pooling for improved performance
- Configurable through YAML configuration
- Health check monitoring with automatic server exclusion
- Concurrent request handling with mutex-protected shared resources

## Requirements

- Go 1.18 or higher
- `gopkg.in/yaml.v3` package

## Installation

Clone the repository:

```bash
git clone https://github.com/minhaz11/a-simple-load-balancer-with-go.git
cd <your-directory>
```
### Configuration
The load balancer can be configured through `config.yml`:

```yml 
port: 8080
servers:
  - http://127.0.0.1:8000 //set your server here
  - http://127.0.0.1:8001 //set your server here
  - Add more servers here if needed
```

### Build and run the project:
```bash
go run .
or 
go build -o <name>
```

## Usage
You can test the load balancer using any of the following server implementations for example I'm adding template for 
Go, Python and Node.js (`make two or more server and test the load balancer`) :

#### Go Server
```GO
package main

import (
    "fmt"
    "net/http"
)

func main() {
    port := "8000" // Use 8001 for second server

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Response from server port: %s\n", port)
    })

    http.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
    })

    fmt.Printf("Server listening on port %s\n", port)
    http.ListenAndServe(":"+port, nil)
}
```
### Python Server (using Flask)
```js
from flask import Flask
app = Flask(__name__)

PORT = 8000  # Use 8001 for second server

@app.route('/')
def hello():
    return f"Response from server port: {PORT}\n"

@app.route('/healthcheck')
def healthcheck():
    return "OK"

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=PORT)
```

### Node.js Server (using Express)

```javascript
const express = require('express');
const app = express();
const port = 8000; // Use 8001 for second server

app.get('/', (req, res) => {
    res.send(`Response from server port: ${port}\n`);
});

app.get('/healthcheck', (req, res) => {
    res.sendStatus(200);
});

app.listen(port, () => {
    console.log(`Server listening on port ${port}`);
});
```
### Test in Docker Setup
You can also run test servers using Docker. Here's an example using the Node.js implementation:

```dockerfile
FROM node:16-alpine

WORKDIR /app
COPY server.js package.json ./
RUN npm install
EXPOSE 8000
CMD ["node", "server.js"]
```

Build and run multiple instances:

```bash
docker build -t node-server .
docker run -d -p 8000:8000 --name node-server-1 node-server
docker run -d -p 8001:8000 --name node-server-2 node-server
```

### Testing the Load Balancer

1. Start two or more backend servers using any of the methods above
2. Ensure the server URLs are correctly configured in config.yml
3. Start the load balancer using make run
4. Test with curl:

```bash
# Send multiple requests to see load balancing in action
for i in {1..10}; do curl http://localhost:8080/; done
```

## Architecture
The project consists of several key components:
#### Connection Pool (`connection_pool.go`)

* Manages HTTP client connections
* Implements connection reuse for better performance
* Configurable maximum connections and timeouts

#### Load Balancer (`load_balancer.go`)

* Implements round-robin algorithm
* Performs health checks every 10 seconds
* Automatically excludes unhealthy servers
* Forwards requests to healthy backend servers

#### Configuration (`config_parser.go`)

* Parses YAML configuration
* Validates server configurations
* Provides typed configuration access

### Health Checking
The load balancer performs health checks on all configured backend servers every 10 seconds. A server is considered healthy if it responds to a GET request at `/healthcheck` with a 200 OK status code.

### Performance Optimization

* Connection pooling to reduce the overhead of creating new connections
* Mutex-protected concurrent operations
* Efficient round-robin server selection
* Configurable connection timeouts and pool sizes