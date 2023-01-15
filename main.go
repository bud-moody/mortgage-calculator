package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
)

func main() {
	helloWorldService := helloWorldService{}

	httpHandler := httptransport.NewServer(
		makeHelloWorldEndpoint(helloWorldService),
		decodeRequest,
		encodeResponse,)

	http.Handle("/helloWorld", httpHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))

}

type helloWorldRequest struct {
	S string `json:"s"`
}

type helloWorldResponse struct {
	V   string `json:"v"`
	Err string `json:"err,omitempty"`
}


func decodeRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request helloWorldRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

func makeHelloWorldEndpoint(svc HelloWorldService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		
		v, err := svc.HelloWorld()
		if err != nil {
			return helloWorldResponse{v, err.Error()}, nil
		}
		return helloWorldResponse{v, ""}, nil
	}
}
