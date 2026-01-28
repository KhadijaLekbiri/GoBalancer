package admin

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"reverse-proxy/services/models"
)


type API struct {
	Pool *models.ServerPool
}

func NewAPI(pool *models.ServerPool) *API {
	return &API{Pool: pool}
}

type StateResponse struct {
	Total_backends int
	Active_backends int
	Backends []BackendResponse
}

type BackendResponse struct {
	URL          string `json:"url"`
	Alive        bool     `json:"alive"`
	CurrentConnections int64    `json:"current_connections"`
}

func (api *API) HandleGet(w http.ResponseWriter, _ *http.Request){
	fmt.Printf("helooooooo")

	backends := api.Pool.Backends

	state := StateResponse {
		Total_backends: len(backends),
		Backends: make([]BackendResponse, 0, len(backends)),
	}

	for _, backend := range backends {
		alive := backend.IsAlive()
		if alive {
			state.Active_backends++
		}

		state.Backends = append(state.Backends, BackendResponse {
			URL: backend.URL.String(),
			Alive: alive,
			CurrentConnections: backend.CurrentConns,
		})
	}

	jsonResp , err := json.MarshalIndent(state,"","	")
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type","application/json")
	w.Write(jsonResp)
}

func (api *API) HandlePost(w http.ResponseWriter, r *http.Request) {

	new_url := struct {
		Url string  `json:"url"`}{
			Url: "",
		}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal([]byte(body), &new_url)

	url, _ := url.ParseRequestURI(new_url.Url)

	newBackend := models.Backend{
		URL: url,
		Alive: false,
		CurrentConns: 0,
	}

	api.Pool.AddBackend(&newBackend)

	w.WriteHeader(http.StatusCreated)

}

func (api *API) Handler(w http.ResponseWriter, r *http.Request){
	
	switch r.Method {
		case http.MethodGet:
			api.HandleGet(w,r)
		case http.MethodPost:
			api.HandlePost(w,r)
		default:
			fmt.Print("Oops")
	}
}

func CheckStatus(pool *models.ServerPool) {
	my_api :=  NewAPI(pool)
		
	http.HandleFunc("/backends",my_api.Handler)
	http.HandleFunc("/status",my_api.Handler)

	http.ListenAndServe(":8081",nil)
}	