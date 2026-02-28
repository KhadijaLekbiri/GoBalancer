package health

import (
	"fmt"
	"net/http"
	"reverse-proxy/services/models"
	"time"
	"sync"
)

func Healthcheck(server_pool *models.ServerPool, timeout time.Duration) {
	// here we could assume that we go through the backends and ping them

	server_pool.Mux.RLock()

	backends := make([]*models.Backend, len(server_pool.Backends))
	copy(backends,server_pool.Backends)

	server_pool.Mux.RUnlock()
		
	var wg sync.WaitGroup

	for i, backend := range backends {

		wg.Add(1)

		go func (b *models.Backend) {
			defer wg.Done()

			client := &http.Client{Timeout: timeout}
			resp, err :=  client.Get(b.URL.String())
			
			alive := err == nil && resp.StatusCode >= 200 && resp.StatusCode < 400
			fmt.Println("Backend",i, "is: ",alive)
			b.SetAlive(alive)
			
			if err == nil {
				resp.Body.Close()
			}
			
		}(backend)
	}
	wg.Wait()
}