package proxy

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httputil"
	"reverse-proxy/services/models"
)

type Handler struct {
	proxy *httputil.ReverseProxy
	pool *models.ServerPool
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request){  
	backend := h.pool.GetNextValidPeer()
	if backend == nil {
		http.Error(
			w,
			fmt.Sprintf("%d Service Unavailable", http.StatusServiceUnavailable),
			http.StatusServiceUnavailable,
		)
		return
	}

	ctx := context.WithValue(r.Context(),"backend",backend)

	h.proxy.ServeHTTP(w,r.WithContext(ctx))
}



func StartProxy(proxy models.ProxyConfig, server_pool *models.ServerPool) {

	fmt.Println("Proxy starts! ")

	
	director := func (req *http.Request){
		selected_backend := req.Context().Value("backend").(*models.Backend)

		req.URL.Scheme = selected_backend.URL.Scheme
		req.URL.Host = selected_backend.URL.Host
		
		fmt.Println("Forwarding to:", selected_backend.URL.String())
	}

	reverse_proxy := &httputil.ReverseProxy{Director: director}

	handler := Handler{
		proxy: reverse_proxy,
		pool: server_pool,
	}


	http.Handle("/api", handler)
	http.ListenAndServe(fmt.Sprintf(":%d",proxy.Port),nil)	
}
