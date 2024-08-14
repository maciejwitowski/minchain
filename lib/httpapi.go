package lib

import (
	"context"
	"io"
	"log"
	"net/http"
)

type HttpApi struct {
	server   *http.Server
	mux      *http.ServeMux
	messages chan string
}

func NewHttpApi(addr string) *HttpApi {
	log.Println("HTTP will listen on", addr)
	mux := http.NewServeMux()
	api := &HttpApi{
		mux:      mux,
		messages: make(chan string, 100),
		server:   &http.Server{Addr: addr, Handler: mux},
	}

	mux.HandleFunc("/", api.handleRequest)
	return api
}

func (api *HttpApi) handleRequest(w http.ResponseWriter, r *http.Request) {
	log.Println("received HTTP request")

	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	api.messages <- string(body)
}

func (api *HttpApi) Start() error {
	return api.server.ListenAndServe()
}

func (api *HttpApi) Stop(ctx context.Context) error {
	return api.server.Shutdown(ctx)
}

func (api *HttpApi) InputChannel(ctx context.Context) <-chan string {
	return api.messages
}
