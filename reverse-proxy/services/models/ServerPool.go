package models

import (
	"fmt"
	"net/url"
	"sync/atomic"
)

type ServerPool struct {
	Backends []*Backend `json:"backends"`
	Current uint64 `json:"current"` // Used for Round-Robin
}


func (pool *ServerPool) GetNextValidPeer() *Backend {

	backends := pool.Backends

	n := len(backends)
	
	if n == 0 {
		return nil
	}

	for i := 0 ;i < n; i++ {
		next := atomic.AddUint64(&pool.Current,1)
		index := int(next) % n

		if backends[index].IsAlive() {
			fmt.Println("slm: ",pool.Current) 
			return backends[index]
		}
	}
	
	return  nil
}

func (pool *ServerPool) AddBackend(backend *Backend) {
	pool.Backends = append(pool.Backends, backend)
}


func (pool *ServerPool) SetBackendStatus(uri *url.URL, alive bool) {
	for _ ,server := range pool.Backends {
		if server.URL.String() == uri.String() {
			server.SetAlive(alive)
		}
	}
}