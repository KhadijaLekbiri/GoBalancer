package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"reverse-proxy/services/models"
	"os"
	"flag"
)

type State struct {
	total_backends int
	active_backends int
	backends [] models.Backend
}

func main() {

	result := flag.String("config","","a path to config")
	
	flag.Parse()

	if *result == "" {
		fmt.Print("No Config file given")
		return
	}
	fmt.Println(*result)

	data, _ := os.ReadFile(*result)

	proxy := struct {
						Port  int
						Strategy string
						Health_check_frequency string
						Backends []string
						Admin_port int
						Request_timeout string
						}{}

	err := json.Unmarshal(data, &proxy)
	if err != nil {
		log.Fatal(err)
	}				

	fmt.Println(proxy)

	resp, err := http.Get(fmt.Sprintf("http://localhost:%d/status",proxy.Admin_port))
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)

	state := State{}
	err = json.Unmarshal(bytes, &state)
		
	fmt.Println(state)
}