package main

import (

	"encoding/json"

	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/helloWorld", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))

}

func handler(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(response(helloWorldService{}))
}


type helloWorldResponse struct {
	V   string `json:"v"`
	Err string `json:"err,omitempty"`
}

func response(svc HelloWorldService) any {
	
		
		v, err := svc.HelloWorld()
		if err != nil {
			return helloWorldResponse{v, err.Error()}
		}

		return helloWorldResponse{v, ""}
	
}
