package main

import (
	"encoding/json"
	"fmt"
	// "io"
	"log"
	// "net/http"
	"reverse-proxy/services/models"
	"reverse-proxy/services/proxy"
	"os"
	"flag"
	"sync"
)

type State struct {
	total_backends int
	active_backends int
	backends [] models.Backend
}

// var Proxy = struct {
// 					Port  int
// 					Strategy string
// 					Health_check_frequency string
// 					Backends []string
// 					Admin_port int
// 					Request_timeout string
// 						}{}

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

	wg.Add(1)

	go func (){
			proxy.StartProxy(initial_proxy)
			defer wg.Done()
		}()

	wg.Wait()

	fmt.Println("hani")
	
	
	// resp, err := http.Get(fmt.Sprintf("http://localhost:%d/status",proxy.Admin_port))
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// defer resp.Body.Close()

	// bytes, err := io.ReadAll(resp.Body)

	// state := State{}
	// err = json.Unmarshal(bytes, &state)
		
	// fmt.Println(state)
}