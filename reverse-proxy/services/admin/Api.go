package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	// "flag"
)


func HandleGet(w http.ResponseWriter, _ *http.Request){
	fmt.Printf("helooooooo")


	fmt.Println("oplaaa")

	jsonResp , err := json.Marshal([]string{"awdii","we are surviving"})
	if err != nil {
		log.Fatal(err)
	}
	w.Write(jsonResp)
}

func HandlePost(w http.ResponseWriter, r *http.Request){

}

func Handler(w http.ResponseWriter, r *http.Request){
	switch r.Method {
		case http.MethodGet:
			HandleGet(w,r)
		case http.MethodPost:
			HandlePost(w,r)
		default:
			fmt.Print("Oops")
	}
}

func main(){
	http.HandleFunc("/status",Handler)
	// http.HandleFunc("/backends",Handler)
	http.ListenAndServe(":8081",nil)
}