package proxy

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sync"
	"io"
	"reverse-proxy/services/models"
	// "flag"
)

// var channel = make(chan []byte, 1)

// func HandleGet(w http.ResponseWriter, _ *http.Request){

// 	fmt.Printf("helooooooo")


// 	fmt.Println("oplaaa")

// 	jsonResp , err := json.Marshal([]string{"hhh"})
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	channel <- jsonResp
// 	// w.Write(jsonResp)
// 	w.Write(<- channel)
// }

// func HandlePost(w http.ResponseWriter, r *http.Request){

// }

// func Handler(w http.ResponseWriter, r *http.Request){
// 	switch r.Method {
// 		case http.MethodGet:
// 			HandleGet(w,r)
// 		case http.MethodPost:
// 			HandlePost(w,r)
// 		default:
// 			fmt.Print("Oops")
// 	}
// }

func StartProxy(proxy models.ProxyConfig) {
	var wg sync.WaitGroup

	fmt.Println("llllollllllll")

	parsedUrl, _ := url.Parse("http://localhost:8082");

	server_pool := models.ServerPool {
			Backends: []*models.Backend{{
				URL: parsedUrl,
				Alive: true,
				CurrentConns: 0,
			}},
			Current: uint64(0),
		}
	if len(server_pool.Backends) == 0 {

	}
	wg.Add(1)

	go func (){
		defer wg.Done()

		http.HandleFunc("/api",
		func(w http.ResponseWriter, r *http.Request) {
				selected_backend := server_pool.Backends[server_pool.Current]

				fmt.Println("Url of the selected active backend: ",selected_backend.URL.String()) 
				
					resp, err := http.Get(selected_backend.URL.String())
					if err != nil {
						log.Fatal(err)
					}

					defer resp.Body.Close()

					bytes, err := io.ReadAll(resp.Body)

					result := []string{}

					_  = json.Unmarshal(bytes,&result)
					if err != nil {
						log.Fatal(err)}

				w.Write(bytes)
			})
		http.ListenAndServe(fmt.Sprintf(":%d",proxy.Port),nil)	
	}()

	wg.Wait()
	// go func (){
	// 	http.HandleFunc("/status",Handler)
	// 	http.ListenAndServe(fmt.Sprintf(":%d",proxy.Admin_port),nil)
	// }()

	
	// http.ListenAndServe(fmt.Sprintf(":%d",proxy.Port),nil)
}