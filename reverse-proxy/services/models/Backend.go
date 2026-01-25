package models

import (
	"log"
	"net/url"
	"sync"
)

type Backend struct {
	URL          *url.URL `json:"url"`
	Alive        bool     `json:"alive"`
	CurrentConns int64    `json:"current_connections"`
	mux          sync.RWMutex
}

func (b *Backend) IsAlive() bool {
	b.mux.RLock()
	defer b.mux.RUnlock()
	return b.Alive
}

func (b *Backend) SetAlive(alive bool)  {
	b.mux.Lock()
	defer b.mux.Unlock()
	b.Alive = alive
}
//  This Must function get ride of the err and only return the url we want

func Must(rawURL string) *url.URL {
	url, err := url.Parse(rawURL)
	if err != nil {
		log.Fatal("hhh", rawURL, err)
	}
	return url
}

	var backends = []*Backend{{
		URL:          Must("http://localhost:8082"),
		Alive:        true,
		CurrentConns: 0,
	},
	{
		URL:          Must("http://localhost:8083"),
		Alive:        true,
		CurrentConns: 0,
	},
	{
		URL:          Must("http://localhost:8084"),
		Alive:        false,
		CurrentConns: 0,
	},
	{
		URL:          Must("http://localhost:8085"),
		Alive:        true,
		CurrentConns: 0,
	},
}

func Backends() []*Backend {
	return backends
}
