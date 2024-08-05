package p2p

import (
	"bufio"
	"fmt"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/network"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	node, err := libp2p.New(
		libp2p.ListenAddrStrings("/ip4/127.0.0.1/tcp/0"),
	)
	if err != nil {
		panic(err)
	}

	node.SetStreamHandler("/chat/1.0.0", handleStream)

	// Print the p2pHost's addresses
	for _, addr := range node.Addrs() {
		fmt.Printf("Listening on %s/p2p/%s\n", addr, node.ID())
	}
	// wait for a SIGINT or SIGTERM signal
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	fmt.Println("Received signal, shutting down...")

	if err := node.Close(); err != nil {
		panic(err)
	}
}

func handleStream(s network.Stream) {
	log.Println("Got a new stream!")

	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
	go readData(rw)
	go writeData(rw)
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
