package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Request received")
		_, err := fmt.Fprint(w, "Hello from the server!")
		if err != nil {
			log.Fatalf("write error: %v", err)
		}
	})

	fmt.Println("Server is stating on port 9555...")
	if err := http.ListenAndServe("0.0.0.0:9555", nil); err != nil {
		log.Fatalf("Server failed to start:  %v", err)
	}
}
