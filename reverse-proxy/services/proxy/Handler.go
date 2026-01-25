package proxy

import (
	// "encoding/json"
	"fmt"
	// "io"
	"net/http"
	"net/http/httputil"
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

type Handler struct{
	proxy *httputil.ReverseProxy
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request){  
	h.proxy.ServeHTTP(w,r)
}



func StartProxy(proxy models.ProxyConfig) {

	fmt.Println("llllollllllll")


	server_pool := models.ServerPool {
			Backends: models.Backends(),
			Current: ^uint64(0) , // max unit64 + 1 gives 0 thanks to the overflow
		}

	
	director := func (req *http.Request){
		selected_backend := server_pool.GetNextValidPeer()

		req.URL.Scheme = selected_backend.URL.Scheme
		req.URL.Host = selected_backend.URL.Host
		
		fmt.Println("Url of the selected active backend: ",selected_backend.URL.String()) 

	}

	reverse_proxy := &httputil.ReverseProxy{Director: director}
	handler := Handler{proxy: reverse_proxy}


	http.Handle("/api", handler)
	http.ListenAndServe(fmt.Sprintf(":%d",proxy.Port),nil)	

// go func (){
// 	http.HandleFunc("/status",Handler)
// 	http.ListenAndServe(fmt.Sprintf(":%d",proxy.Admin_port),nil)
// }()


	// http.ListenAndServe(fmt.Sprintf(":%d",proxy.Port),nil)
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