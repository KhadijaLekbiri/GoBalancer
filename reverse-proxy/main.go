package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"reverse-proxy/services/admin"
	"reverse-proxy/services/models"
	"reverse-proxy/services/proxy"
	"reverse-proxy/services/health"
	"strings"
	"sync"
	"time"
)




func main() {
	var wg sync.WaitGroup

	result := flag.String("config","","a path to config")
	
	flag.Parse()

	if *result == "" {
		fmt.Print("No Config file given")
		return
	}

	data, _ := os.ReadFile(*result)

	initial_proxy := models.ProxyConfig{}
	err := json.Unmarshal(data, &initial_proxy)

	initial_proxy.HealthCheckFreq = 10*time.Second
	initial_proxy.Timeout = 2*time.Second

	if err != nil {
		log.Fatal(err)
	}				

	server_pool := models.ServerPool {
			Backends: models.Backends(),
			Current: ^uint64(0) , // max unit64 + 1 gives 0 thanks to the overflow
		}
	server_pool.Activate_backends()

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
			health.Healthcheck(&server_pool,initial_proxy.Timeout)
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

		port := 8087
		for range ticker2.C {
			jsonData := fmt.Sprintf(`{"url": "http://localhost:%d"}`, port)
			
			resp, err := http.Post("http://localhost:8081/backends", "application/json", strings.NewReader(jsonData))
			if err != nil {
				log.Println("Failed to add backend: ", err)
				continue
			}
			
			log.Println("Added backend on port", port, "status:", resp.Status)
			resp.Body.Close()
		
			port++
		}
	}()
	wg.Wait()
	

}