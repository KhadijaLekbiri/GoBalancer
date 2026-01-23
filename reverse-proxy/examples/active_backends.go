package main

import(
	"sync"
"net/http"
"encoding/json"
"log")


func main(){
	var wg sync.WaitGroup

	wg.Add(1)
	go func ()  {
		defer wg.Done()

		http.HandleFunc("/",
		func(w http.ResponseWriter, r *http.Request) {
				jsonResp , err := json.Marshal([]string{"Welcome to Active Backend 2"})
				if err != nil {
					log.Fatal(err)
			}
			w.Write(jsonResp)
		})
		http.ListenAndServe(":8082",nil)
	}()
	wg.Wait()


}