package admin

import (
	"encoding/json"
	"fmt"
	"io"
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
	
	api.Pool.Mux.RLock()
	defer api.Pool.Mux.RUnlock()
	
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
		http.Error(w, err.Error(), http.StatusBadRequest)
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
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.Unmarshal([]byte(body), &new_url)

	url, err := url.ParseRequestURI(new_url.Url)
	if err != nil {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	for _, b := range api.Pool.Backends {
		if b.URL.String() == url.String() {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(map[string]string{
                "error": "Backend already exists",
            })
			return
		}
	}

	newBackend := models.Backend{
		URL: url,
		Alive: true,
		CurrentConns: 0,
	}

	api.Pool.AddBackend(&newBackend)

	w.WriteHeader(http.StatusCreated)

}

func  (api *API) HandleDelete(w http.ResponseWriter, r *http.Request){

	new_url := struct {
		Url string  `json:"url"`}{
			Url: "",
		}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.Unmarshal([]byte(body), &new_url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	parsedURL, err := url.ParseRequestURI(new_url.Url)
	if err != nil {
		http.Error(w, "Invalid Url format",http.StatusBadRequest)
		return
	}
	api.Pool.Mux.Lock()
	defer api.Pool.Mux.Unlock()

	for i, backend := range api.Pool.Backends {
		if parsedURL.String() == backend.URL.String() {
			api.Pool.Backends = append(api.Pool.Backends[:i],api.Pool.Backends[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			fmt.Println("Backend successfully deleted!")
			return
		}
	}
	http.Error(w, "Backend Not Found", http.StatusNotFound)
}


func (api *API) Handler(w http.ResponseWriter, r *http.Request){
	
	switch r.Method {
		case http.MethodGet:
			api.HandleGet(w,r)
		case http.MethodPost:
			api.HandlePost(w,r)
		case http.MethodDelete:
			api.HandleDelete(w,r)
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