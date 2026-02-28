# Concurrent Load-Balancing Reverse Proxy

A reverse proxy server built in Go featuring intelligent load balancing, automated health monitoring, and dynamic backend management.

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
cd GoBalancer

# Initialize Go modules
go mod init reverse-proxy

# Download dependencies
go mod tidy
```

### Configuration

A `config.json` file is included in the repository with the following structure:

```json
{
  "Port": 8080,
  "Strategy": "round-robin",
  "Admin_port": 8081
}
```

Adjust the values as needed before running the proxy.

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

## 📖 Usage

### Making Requests Through the Proxy

```bash
# Send a request through the proxy
curl http://localhost:8080/api

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
GoBalancer/
├── main.go                    # Entry point
├── config.json                # Configuration file
├── README.md
└── services/
    ├── models/                # Data structures and interfaces
    │   ├── Backend.go
    │   ├── LoadBalancer.go
    │   ├── ServerPool.go
    │   └── ProxyConfig.go
    ├── proxy/                 # Proxy handler logic
    │   └── Handler.go
    ├── health/                # Health checking system
    │   └── checker.go
    └── admin/                 # Admin API handlers
        └── Api.go
```

## 🔧 Configuration Options

| Field | Type | Description | Default |
|-------|------|-------------|---------|
| `Port` | int | Proxy server port | 8080 |
| `Strategy` | string | Load balancing strategy (`round-robin`) | `round-robin` |
| `health_check_frequency` | duration | Health check interval | `30s` |
| `backends` | []string | Initial backend URLs | `[]` |
| `Admin_port` | int | Admin API port | 8081 |
| `request_timeout` | duration | Backend request timeout | `10s` |

## 🎯 Project Status

✅ **Complete** — All core features implemented and tested.

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

---

⭐ If you found this project interesting, please consider giving it a star!