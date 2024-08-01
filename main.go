package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "9555"
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Request received")
		_, err := fmt.Fprint(w, "Hello from the server!")
		if err != nil {
			log.Fatalf("write error: %v", err)
		}
	})

	server := &http.Server{Addr: "0.0.0.0:" + port}

	go func() {
		log.Printf("Server is starting on port %s...", port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	if err := server.Close(); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
