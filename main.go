package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
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

	// ---------------------LIB P2P ---------------------------------------

	h, err := libp2p.New()
	if err != nil {
		log.Fatal(err)
	}
	defer func(h host.Host) {
		err := h.Close()
		if err != nil {
			log.Fatal("closing failed")
		}
	}(h)

	maddr, err := multiaddr.NewMultiaddr("/ip4/127.0.0.1/tcp/PORT/p2p/PEER_ID\"")
	if err != nil {
		log.Fatal(err)
	}

	info, err := peer.AddrInfoFromP2pAddr(maddr)
	if err != nil {
		log.Fatal(err)
	}

	// Connect to the first server
	err = h.Connect(context.Background(), *info)
	if err != nil {
		log.Fatal(err)
	}

	//s, err := h.NewStream(context.Background(), info.ID, "/chat/1.0.0")
	//if err != nil {
	//	return
	//}

	//rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriterSize(s))
	//
	//go readData(rw)
	//go writeData(rw)

	// Wait forever
	select {}
}

func readData(rw *bufio.ReadWriter) {
	for {
		str, err := rw.ReadString('\n')
		if err != nil {
			log.Println("Error reading from buffer")
			return
		}

		if str == "" {
			return
		}
		if str != "\n" {
			// Green console colour: 	\x1b[32m
			// Reset console colour: 	\x1b[0m
			fmt.Printf("\x1b[32m%s\x1b[0m> ", str)
		}
	}
}

func writeData(rw *bufio.ReadWriter) {

}
