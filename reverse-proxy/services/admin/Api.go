package admin

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reverse-proxy/services/models"

	// "flag"
)


type API struct {
	Pool *models.ServerPool
}

func NewAPI(pool *models.ServerPool) *API {
	return &API{Pool: pool}
}



func (api *API) HandleGet(w http.ResponseWriter, _ *http.Request){
	fmt.Printf("helooooooo")

	backends := api.Pool.Backends

	jsonResp , err := json.MarshalIndent(backends,"","	")
	if err != nil {
		log.Fatal(err)
	}
	w.Write(jsonResp)
}

func HandlePost(w http.ResponseWriter, r *http.Request){

}

func (api *API) Handler(w http.ResponseWriter, r *http.Request){
	
	switch r.Method {
		case http.MethodGet:
			api.HandleGet(w,r)
		case http.MethodPost:
				HandlePost(w,r)
		default:
			fmt.Print("Oops")
	}
}

func CheckStatus(pool *models.ServerPool) {
	my_api :=  NewAPI(pool)

	http.HandleFunc("/status",my_api.Handler)
	// http.HandleFunc("/backends",Handler)
	http.ListenAndServe(":8081",nil)
}


// SOME OLD STUFF 
// func(w http.ResponseWriter, r *http.Request) {

				
// 					resp, err := http.Get(selected_backend.URL.String())
// 					if err != nil {
// 						log.Fatal(err)
// 					}

// 					defer resp.Body.Close()

// 					bytes, err := io.ReadAll(resp.Body)

// 					result := []string{}

// 					_  = json.Unmarshal(bytes,&result)
// 					if err != nil {
// 						log.Fatal(err)}

// 				w.Write(bytes)
// 			}