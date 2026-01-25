package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reverse-proxy/services/models"
	"sync"
)


// Mux = Router that maps URL paths to handler functions
// Default mux = Global singleton shared across your entire program
// Custom mux = Independent router you create and control


func main(){
	var wg sync.WaitGroup

	backends := models.Backends()
	
	for i, backend := range backends {
		if backend.Alive {
			wg.Add(1)
			go func (index int, backend *models.Backend){
				defer wg.Done()

				mux := http.NewServeMux()
				mux.HandleFunc("/",
				func(w http.ResponseWriter, r *http.Request) {
						jsonResp , err := json.Marshal(fmt.Sprintf("Welcome to Active Backend %d lololololoy",i))
						if err != nil {
							log.Fatal("hnaaaaaa",err)
					}
					w.Write(jsonResp)
				})
				http.ListenAndServe(fmt.Sprintf(":%s",backend.URL.Port()),mux)
			}(i, backend)
		}
	}
	
	wg.Wait()


}