package models

import (
	"net/url"
	"sync"
	"sync/atomic"
)

type ServerPool struct {
	Backends []*Backend `json:"backends"`
	Current uint64 `json:"current"` // Used for Round-Robin
	Mux sync.RWMutex
}


func (pool *ServerPool) GetNextValidPeer() *Backend {
	pool.Mux.RLock()
	defer pool.Mux.RUnlock()
	
	backends := pool.Backends

	n := len(backends)
	
	if n == 0 {
		return nil
	}

	for i := 0 ;i < n; i++ {
		next := atomic.AddUint64(&pool.Current,1)
		index := int(next) % n

		if backends[index].IsAlive() {
			return backends[index]
		}
	}
	
	return  nil
}

func (pool *ServerPool) AddBackend(backend *Backend) {
	pool.Mux.Lock()
	defer pool.Mux.Unlock()
	pool.Backends = append(pool.Backends, backend)
	backend.SetAlive(backend.Alive)
}


func (pool *ServerPool) SetBackendStatus(uri *url.URL, alive bool) {

	pool.Mux.Lock()
	defer pool.Mux.Unlock()

	backends := pool.Backends

	for _ ,server := range backends {
		if server.URL.String() == uri.String() {
			server.SetAlive(alive)
		}
	}
	
}

func (pool *ServerPool) Activate_backends(){
	pool.Mux.RLock()

	backends := make([]*Backend, len(pool.Backends))
	copy(backends,pool.Backends)

	pool.Mux.RUnlock()

	for _, backend := range backends {
		if backend.Alive {
			backend.StartBackend()
		}
	}
}