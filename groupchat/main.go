package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"github.com/libp2p/go-libp2p"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
	"os"
)

const (
	DiscoveryServiceTag = "group-chat"
)

func main() {
	ctx := context.Background()
	topicF := flag.String("t", "", "pubsub topic")
	listenF := flag.Int("l", 0, "wait for incoming connections")
	flag.Parse()

	h, err := libp2p.New(
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", *listenF)),
	)
	if err != nil {
		panic(err)
	}
	defer func(h host.Host) {
		err := h.Close()
		if err != nil {
			panic(err)
		}
	}(h)

	fmt.Println("Host ID:", h.ID())
	fmt.Println("Host Addresses:")

	for _, addr := range h.Addrs() {
		fmt.Printf("  %s/p2p/%s\n", addr, h.ID())
	}

	// Create a new PubSub service using the GossipSub router
	ps, err := pubsub.NewGossipSub(context.Background(), h)
	if err != nil {
		panic(err)
	}

	// Join the pubsub topic
	topic, err := ps.Join(*topicF)
	if err != nil {
		panic(err)
	}

	// Create a new subscription to the topic
	sub, err := topic.Subscribe()
	if err != nil {
		panic(err)
	}

	// Setup local mDNS discovery
	setupDiscovery(ctx, h)

	go handleMessages(ctx, sub)

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		message, err := reader.ReadString('\n')
		if err != nil {
			return
		}
		if err := topic.Publish(ctx, []byte(message)); err != nil {
			fmt.Println("Publish error:", err)
		}
	}
}

// setupDiscovery creates an mDNS discovery service and attaches it to the libp2p Host.
func setupDiscovery(ctx context.Context, h host.Host) {
	// Setup local mDNS discovery
	s := mdns.NewMdnsService(h, DiscoveryServiceTag, &discoveryNotifee{ctx: ctx, h: h})
	if err := s.Start(); err != nil {
		panic(err)
	}
}

// discoveryNotifee gets notified when we find a new peer via mDNS discovery
type discoveryNotifee struct {
	ctx context.Context
	h   host.Host
}

func (n *discoveryNotifee) HandlePeerFound(pi peer.AddrInfo) {
	fmt.Printf("Discovered new peer %s\n", pi.ID.String())
	err := n.h.Connect(n.ctx, pi)
	if err != nil {
		fmt.Printf("Error connecting to peer %s: %s\n", pi.ID.String(), err)
	}
}

func handleMessages(ctx context.Context, sub *pubsub.Subscription) {
	for {
		msg, err := sub.Next(ctx)
		if err != nil {
			fmt.Println("Subscription error:", err)
			return
		}
		fmt.Printf("%s: %s", msg.ReceivedFrom.ShortString(), string(msg.Data))
	}
}
