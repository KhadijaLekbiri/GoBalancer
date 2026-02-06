package models

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type Backend struct {
	URL          *url.URL `json:"url"`
	Alive        bool     `json:"alive"`
	CurrentConns int64    `json:"current_connections"`
	mux          sync.RWMutex
	server 		*http.Server
}

func (b *Backend) IsAlive() bool {
	b.mux.RLock()
	defer b.mux.RUnlock()
	return b.Alive
}

func (b *Backend) SetAlive(alive bool)  {
	b.mux.Lock()
	b.Alive = alive
	b.mux.Unlock()

	if alive {
		b.StartBackend()
	} else {
		b.StopBackend()
	}
}

func (b *Backend) AddConnection()  {
	b.mux.Lock()
	defer b.mux.Unlock()
	b.CurrentConns++
}

//  This Must function get ride of the err and only return the url we want

func Must(rawURL string) *url.URL {
	url, err := url.Parse(rawURL)
	if err != nil {
		log.Fatal("hhh", rawURL, err)
	}
	return url
}

// var backends = []*Backend{{
// 		URL:          Must("http://localhost:8082"),
// 		Alive:        true,
// 		CurrentConns: 0,
// 	},
// 	{
// 		URL:          Must("http://localhost:8083"),
// 		Alive:        true,
// 		CurrentConns: 0,
// 	},
// 	{
// 		URL:          Must("http://localhost:8084"),
// 		Alive:        false,
// 		CurrentConns: 0,
// 	},
// 	{
// 		URL:          Must("http://localhost:8085"),
// 		Alive:        true,
// 		CurrentConns: 0,
// 	},
// }

var backends = []*Backend{}

func Backends() []*Backend {
	return backends
}

func (b *Backend) StartBackend() {
	b.mux.Lock()
	defer b.mux.Unlock() 

	if b.server != nil {
		return
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/",
	func(w http.ResponseWriter, r *http.Request) {
			jsonResp , err := json.Marshal(fmt.Sprintf("Welcome to Active Backend %s lololololoy", b.URL.Port()))
			if err != nil {
				log.Fatal("hnaaaaaa",err)
			}
		w.Write(jsonResp)
	})

	srv :=  &http.Server{
		Addr: ":"+ b.URL.Port(),
		Handler: mux,
	}

	b.server = srv

	go func() {
		if err :=  srv.ListenAndServe(); err != nil  && err != http.ErrServerClosed {
			log.Println("Backend error:", err)
		}
	}()
}

func (b *Backend) StopBackend() {
	b.mux.Lock()
	defer b.mux.Unlock() 

	if b.server == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	
	b.server.Shutdown(ctx)
	b.server = nil
}


// func Activate_backends(backends []*Backend){
// 	var wg sync.WaitGroup
	
// 	for i, backend := range backends {
// 		if backend.Alive {
// 			wg.Add(1)
// 			StartBackend(i,backend)
// 		}
// 	}
	
// 	wg.Wait()
// }