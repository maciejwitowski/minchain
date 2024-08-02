package main

import (
	"bufio"
	"context"
	"fmt"
	logg "github.com/ipfs/go-log/v2"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
	"log"
	"os"
)

func main() {
	logg.SetAllLoggers(logg.LevelDebug)
	// Create a new libp2p host
	h, err := libp2p.New()
	if err != nil {
		log.Fatal(err)
	}
	defer h.Close()

	// Parse the multiaddr of the first server
	s1 := "/ip4/127.0.0.1/tcp/58510/p2p/12D3KooWDmEA4MP4ZENNG561pYCeh7Lky7GFJefmtfTgWoBEdj1n" //os.Args[1]
	fmt.Println("S1: ", s1)
	maddr, err := multiaddr.NewMultiaddr(s1)

	if err != nil {
		log.Fatal(err)
	}

	// Extract the peer ID from the multiaddr
	info, err := peer.AddrInfoFromP2pAddr(maddr)
	if err != nil {
		log.Fatal(err)
	}

	// Connect to the first server
	err = h.Connect(context.Background(), *info)
	if err != nil {
		log.Fatal(err)
	}

	// Open a stream to the first server
	s, err := h.NewStream(context.Background(), info.ID, "/chat/1.0.0")
	if err != nil {
		log.Fatal(err)
	}

	// Create a buffer stream for non-blocking read and write
	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

	go readData(rw)
	go writeData(rw)

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
	stdReader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		sendData, err := stdReader.ReadString('\n')
		if err != nil {
			log.Println("Error reading from stdin")
			return
		}

		_, err = rw.WriteString(fmt.Sprintf("%s\n", sendData))
		if err != nil {
			log.Println("Error writing to buffer")
			return
		}
		err = rw.Flush()
		if err != nil {
			log.Println("Error flushing buffer")
			return
		}
	}
}
