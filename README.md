# Concurrent Load-Balancing Reverse Proxy

A production-grade reverse proxy server built in Go featuring intelligent load balancing, automated health monitoring, and dynamic backend management.

## 🎓 Project Context

This is my final project for an **Introduction to Go** course, demonstrating advanced concepts including:
- Concurrent programming with goroutines
- Thread-safe state management with mutexes and atomic operations
- Network programming with `net/http`
- Context propagation and graceful shutdowns
- RESTful API design

## ✨ Features

- **🔄 Load Balancing**: Round-robin distribution across healthy backends
- **💚 Health Monitoring**: Automatic background health checks with configurable intervals
- **🔧 Dynamic Configuration**: Add/remove backends at runtime via Admin API
- **🧵 Thread-Safe**: Concurrent request handling with proper synchronization
- **⏱️ Timeout Handling**: Request cancellation and backend timeout management
- **📊 Monitoring**: Real-time statistics on backend health and connection counts

## 🏗️ Architecture

```
┌─────────────┐
│   Clients   │
└──────┬──────┘
       │
       ▼
┌─────────────────────────────┐
│    Reverse Proxy Server     │
│  ┌───────────────────────┐  │
│  │  Load Balancer        │  │
│  │  (Round-Robin)        │  │
│  └───────────────────────┘  │
│  ┌───────────────────────┐  │
│  │  Health Checker       │  │
│  │  (Background Job)     │  │
│  └───────────────────────┘  │
│  ┌───────────────────────┐  │
│  │  Admin API            │  │
│  │  (Port 8081)          │  │
│  └───────────────────────┘  │
└─────────────┬───────────────┘
              │
    ┌─────────┴─────────┐
    ▼         ▼         ▼
┌────────┐ ┌────────┐ ┌────────┐
│Backend1│ │Backend2│ │Backend3│
└────────┘ └────────┘ └────────┘
```

## 🚀 Quick Start

### Prerequisites

- Go 1.21 or higher
- Basic understanding of HTTP and networking concepts

### Installation

```bash
# Clone the repository
git clone https://github.com/KhadijaLekbiri/GoBalancer.git
cd reverse-proxy

# Initialize Go modules
go mod init reverse-proxy

# Download dependencies
go mod tidy
```

### Configuration

Create a `config.json` file in the project root:

```json
{
  "port": 8080,
  "strategy": "round-robin",
  "health_check_frequency": "30s",
  "backends": [
    "http://localhost:8082",
    "http://localhost:8083",
    "http://localhost:8084"
  ]
}
```

### Running the Proxy

```bash
# Run with default config
go run main.go --config=config.json

# Run with race detector (recommended during development)
go run -race main.go --config=config.json

# Build and run
go build -o proxy main.go
./proxy --config=config.json
```

### Setting Up Test Backends

Create simple test servers to proxy to:

```bash
# Terminal 1
go run examples/backend.go -port 8082

# Terminal 2
go run examples/backend.go -port 8083

# Terminal 3
go run examples/backend.go -port 8084
```

## 📖 Usage

### Making Requests Through the Proxy

```bash
# Send a request through the proxy
curl http://localhost:8080/api/users

# The proxy will forward to one of the healthy backends
```

### Admin API Endpoints

The Admin API runs on port `8081` by default.

#### Check System Status

```bash
curl http://localhost:8081/status
```

**Response:**
```json
{
  "total_backends": 3,
  "active_backends": 2,
  "backends": [
    {
      "url": "http://localhost:8082",
      "alive": true,
      "current_connections": 5
    },
    {
      "url": "http://localhost:8083",
      "alive": false,
      "current_connections": 0
    },
    {
      "url": "http://localhost:8084",
      "alive": true,
      "current_connections": 3
    }
  ]
}
```

#### Add a Backend

```bash
curl -X POST http://localhost:8081/backends \
  -H "Content-Type: application/json" \
  -d '{"url": "http://localhost:8085"}'
```

#### Remove a Backend

```bash
curl -X DELETE http://localhost:8081/backends \
  -H "Content-Type: application/json" \
  -d '{"url": "http://localhost:8082"}'
```

## 🏛️ Project Structure

```
reverse-proxy/
├── main.go                 # Entry point
├── config.json            # Configuration file
├── TODO.md               # Development milestones
├── internal/
│   ├── models/           # Data structures and interfaces
│   │   ├── backend.go
│   │   ├── pool.go
│   │   └── config.go
│   ├── proxy/            # Proxy handler logic
│   │   └── handler.go
│   ├── health/           # Health checking system
│   │   └── checker.go
│   └── admin/            # Admin API handlers
│       └── api.go
├── examples/
│   └── backend.go        # Sample backend server
└── tests/
    └── integration_test.go
```

## 🧪 Testing

### Run Unit Tests

```bash
go test ./... -v
```

### Run with Race Detector

```bash
go test ./... -race -v
```

### Integration Testing

```bash
# Start the proxy and backends, then run
go test ./tests -integration -v
```

### Load Testing

```bash
# Using Apache Bench
ab -n 10000 -c 100 http://localhost:8080/

# Using hey
hey -n 10000 -c 100 http://localhost:8080/
```

## 🔧 Configuration Options

| Field | Type | Description | Default |
|-------|------|-------------|---------|
| `port` | int | Proxy server port | 8080 |
| `strategy` | string | Load balancing strategy (`round-robin`, `least-conn`) | round-robin |
| `health_check_frequency` | duration | Health check interval | 30s |
| `backends` | []string | Initial backend URLs | [] |
| `admin_port` | int | Admin API port | 8081 |
| `request_timeout` | duration | Backend request timeout | 10s |

## 🎯 Project Status

✅ **Completed** - All core features implemented and tested

See [TODO.md](TODO.md) for detailed development milestones and progress tracking.

## 📚 Learning Outcomes

Through this project, I gained hands-on experience with:

- **Concurrency Patterns**: Goroutines, channels, and synchronization primitives
- **Thread Safety**: Proper use of `sync.Mutex`, `sync.RWMutex`, and `sync/atomic`
- **HTTP Programming**: Building robust HTTP servers and clients
- **Context Management**: Propagating context for cancellation and timeouts
- **System Design**: Architecting a distributed system component
- **Testing**: Unit tests, integration tests, and race condition detection

## 🚀 Future Enhancements

- [ ] **Least-Connections Algorithm**: More intelligent load distribution
- [ ] **Sticky Sessions**: Client affinity based on IP or cookies
- [ ] **Weighted Load Balancing**: Assign capacity-based weights to backends
- [ ] **HTTPS/TLS Support**: Secure proxy connections
- [ ] **Metrics & Observability**: Prometheus metrics, structured logging
- [ ] **Circuit Breaker**: Prevent cascading failures
- [ ] **Rate Limiting**: Per-client request throttling

## 📄 License

This project is part of an academic assignment and is available for educational purposes.



⭐ If you found this project helpful, please consider giving it a star!
