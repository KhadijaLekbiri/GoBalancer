package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"reverse-proxy/services/models"
	"reverse-proxy/services/proxy"
	"reverse-proxy/services/admin"
	"sync"
	"time"
)




func Healthcheck(server_pool *models.ServerPool, timeout time.Duration){
	// here we could assume that we go through the backends and ping them

	for i, backend := range server_pool.Backends {
		var wg sync.WaitGroup

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

		wg.Wait()
	}
}

func main() {
	var wg sync.WaitGroup

	result := flag.String("config","","a path to config")
	
	flag.Parse()

	if *result == "" {
		fmt.Print("No Config file given")
		return
	}
	fmt.Println(*result)

	data, _ := os.ReadFile(*result)

	initial_proxy := models.ProxyConfig{}
	err := json.Unmarshal(data, &initial_proxy)

	initial_proxy.HealthCheckFreq = 10*time.Second
	initial_proxy.Timeout = 2*time.Second

	if err != nil {
		log.Fatal(err)
	}				

	fmt.Println(initial_proxy.Strategy)

	server_pool := models.ServerPool {
			Backends: models.Backends(),
			Current: ^uint64(0) , // max unit64 + 1 gives 0 thanks to the overflow
		}
	wg.Add(1)

	go func (){
			proxy.StartProxy(initial_proxy, &server_pool)
			defer wg.Done()
		}()

	wg.Add(1)

	go func (){
		defer wg.Done()
		ticker := time.NewTicker(initial_proxy.HealthCheckFreq)

		for range ticker.C {
			Healthcheck(&server_pool,initial_proxy.Timeout)
		}
	}()
	

	wg.Add(1)

	go func (){
		admin.CheckStatus(&server_pool)
		defer wg.Done()
	}()

	wg.Add(1)

	go func (){
		defer wg.Done()
		ticker2 := time.NewTicker(20*time.Second)

		for range ticker2.C {
			backend := &models.Backend{
				URL:          models.Must(fmt.Sprintf("http://localhost:%d",8087)),
				Alive:        true,
				CurrentConns: 0,
			}
			server_pool.AddBackend(backend)
		}
	}()
	wg.Wait()



	fmt.Println("hani")
	

}