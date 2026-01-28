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

type State struct {
	Total_backends int
	Active_backends int
	Backends [] models.Backend
}

// var Proxy = struct {
// 					Port  int
// 					Strategy string
// 					Health_check_frequency string
// 					Backends []string
// 					Admin_port int
// 					Request_timeout string
// 						}{}

func Healthcheck(server_pool *models.ServerPool){
	// here we could assume that we go through the backends and ping them

	for i, backend := range server_pool.Backends {
		var wg sync.WaitGroup

		wg.Add(1)

		go func (b *models.Backend) {
			defer wg.Done()

			client := &http.Client{Timeout: 2*time.Second}
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
		ticker := time.NewTicker(10*time.Second)

		for range ticker.C {
			Healthcheck(&server_pool)
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