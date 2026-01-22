# 📋 Project Milestones

## Milestone 1: Project Setup & Configuration ⏱️ (Est. 30 min)

**Goal**: Initialize project structure and configuration management

- [ ] Initialize Go module: `go mod init <project_name>`
- [ ] Create project directory structure:
  ```
  ├── main.go
  ├── config.json
  ├── internal/
  │   ├── models/
  │   ├── proxy/
  │   ├── health/
  │   └── admin/
  └── README.md
  ```
- [ ] Create `config.json` with initial backend URLs
- [ ] Implement config loading logic
- [ ] Test: Verify config loads successfully

**Deliverable**: Clean project structure with configuration management

---

## Milestone 2: Core Data Models ⏱️ (Est. 45 min)

**Goal**: Define all structs, interfaces, and type definitions

- [ ] Define `Backend` struct with:
  - [ ] URL field (`*url.URL`)
  - [ ] Alive status (`bool`)
  - [ ] CurrentConns counter (`int64`)
  - [ ] RWMutex for thread safety
- [ ] Define `ServerPool` struct with:
  - [ ] Backends slice
  - [ ] Current index for round-robin (`uint64`)
- [ ] Define `ProxyConfig` struct
- [ ] Create `LoadBalancer` interface
- [ ] Add getter/setter methods to `Backend` for thread-safe access
- [ ] Test: Unit tests for Backend thread-safety

**Deliverable**: Complete, thread-safe data models

---

## Milestone 3: Load Balancing Logic ⏱️ (Est. 1-2 hours)

**Goal**: Implement core load balancing algorithms

- [ ] Implement Round-Robin algorithm:
  - [ ] `GetNextValidPeer()` using atomic operations
  - [ ] Skip dead backends
  - [ ] Handle case when no backends available
- [ ] Implement `AddBackend()` method
- [ ] Implement `SetBackendStatus()` method
- [ ] Add connection counting logic:
  - [ ] Increment on request start
  - [ ] Decrement on request end
- [ ] Test: Unit tests for load balancing logic
- [ ] Test: Concurrent request simulation

**Deliverable**: Working load balancer with round-robin support

---

## Milestone 4: Reverse Proxy Core ⏱️ (Est. 1-2 hours)

**Goal**: Build the HTTP proxy handler

- [ ] Create proxy handler using `httputil.ReverseProxy`
- [ ] Implement request forwarding logic:
  - [ ] Get next available backend
  - [ ] Update connection count
  - [ ] Forward request with context
  - [ ] Handle response
- [ ] Implement custom error handler:
  - [ ] Detect connection failures
  - [ ] Mark backends as dead on errors
  - [ ] Return 503 when no backends available
- [ ] Add context propagation for cancellation
- [ ] Add request timeout handling
- [ ] Test: Manual testing with curl
- [ ] Test: Backend failure scenarios

**Deliverable**: Functional reverse proxy

---

## Milestone 5: Health Monitoring System ⏱️ (Est. 1 hour)

**Goal**: Background health checker with periodic pings

- [ ] Create health checker goroutine
- [ ] Implement periodic ticker (configurable interval)
- [ ] Implement backend ping logic:
  - [ ] HTTP GET health check OR
  - [ ] TCP dial check
- [ ] Update backend alive status
- [ ] Add structured logging:
  - [ ] Log status changes (UP/DOWN)
  - [ ] Log health check results
- [ ] Implement graceful shutdown for health checker
- [ ] Test: Verify health checks run periodically
- [ ] Test: Backend recovery detection

**Deliverable**: Automated health monitoring

---

## Milestone 6: Admin API ⏱️ (Est. 1-1.5 hours)

**Goal**: Management interface for runtime configuration

- [ ] Create separate HTTP server for admin (port 8081)
- [ ] Implement `GET /status` endpoint:
  - [ ] Return JSON with all backends
  - [ ] Include alive status
  - [ ] Include connection counts
- [ ] Implement `POST /backends` endpoint:
  - [ ] Parse JSON body
  - [ ] Validate URL
  - [ ] Add to server pool
- [ ] Implement `DELETE /backends` endpoint:
  - [ ] Find and remove backend
  - [ ] Handle not-found case
- [ ] Add proper error handling and status codes
- [ ] Test: All endpoints with curl/Postman

**Deliverable**: Fully functional admin API

---

## Milestone 7: Integration & Testing ⏱️ (Est. 1-2 hours)

**Goal**: End-to-end testing and bug fixes

- [ ] Set up test backend servers (simple HTTP servers)
- [ ] Test complete flow:
  - [ ] Start proxy with config
  - [ ] Send requests through proxy
  - [ ] Verify load distribution
  - [ ] Simulate backend failure
  - [ ] Verify health checker detects failure
  - [ ] Verify traffic redirects to healthy backends
- [ ] Test admin API integration:
  - [ ] Add backend dynamically
  - [ ] Verify it receives traffic
  - [ ] Remove backend
  - [ ] Check status endpoint
- [ ] Test concurrent load (use tool like `ab` or `hey`)
- [ ] Test graceful shutdown
- [ ] Fix any race conditions (run with `-race` flag)

**Deliverable**: Fully tested, production-ready proxy

---

## Milestone 8: Documentation & Polish ⏱️ (Est. 30-45 min)

**Goal**: Clean code and comprehensive documentation

- [ ] Add code comments for complex logic
- [ ] Write comprehensive README with:
  - [ ] Architecture overview
  - [ ] Setup instructions
  - [ ] Usage examples
  - [ ] API documentation
- [ ] Add example `config.json`
- [ ] Create sample backend server code
- [ ] Add logging throughout application
- [ ] Code cleanup and formatting (`gofmt`)
- [ ] Final review of thread safety

**Deliverable**: Well-documented, maintainable code

---

## 🎯 Optional Enhancements (After Core Completion)

### Enhancement 1: Least-Connections Algorithm

- [ ] Implement least-connections strategy
- [ ] Add strategy selection in config
- [ ] Update load balancer to support multiple strategies

### Enhancement 2: Sticky Sessions

- [ ] Implement session tracking (by IP or cookie)
- [ ] Store client-to-backend mapping
- [ ] Add session timeout logic

### Enhancement 3: Weighted Load Balancing

- [ ] Add weight field to Backend
- [ ] Implement weighted round-robin
- [ ] Update admin API to set weights

### Enhancement 4: HTTPS/TLS Support

- [ ] Add TLS configuration
- [ ] Generate or load certificates
- [ ] Update proxy to serve HTTPS

---

## 📊 Grading Checklist

- [ ] **Proxy Functionality (30 pts)**: Correctly forwards traffic and returns responses
- [ ] **Concurrency & Safety (25 pts)**: No race conditions, proper mutex usage
- [ ] **Background Job (15 pts)**: Health checker works correctly
- [ ] **Load Balancing (10 pts)**: Round-robin implemented properly
- [ ] **Context & Timeouts (10 pts)**: Proper cancellation handling
- [ ] **Code Organization (10 pts)**: Clean structure and documentation

---

## 🚀 Quick Reference Commands

```bash
# Run with race detector
go run -race main.go --config=config.json

# Run tests with race detector
go test ./... -race -v

# Format code
gofmt -w .

# Build binary
go build -o proxy main.go
```
