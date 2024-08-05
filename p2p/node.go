package p2p

import (
	"context"
	"fmt"
	"github.com/libp2p/go-libp2p"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
	"minchain/common"
	"strings"
)

var addressTemplate = "/ip4/0.0.0.0/tcp/%d"

const DiscoveryServiceTag = "p2p-service"

type Node struct {
	Topic   *pubsub.Topic
	p2pHost host.Host
}

func InitNode(ctx context.Context, config common.Config) (*Node, error) {
	options := libp2p.ListenAddrStrings(fmt.Sprintf(addressTemplate, config.ListeningPort))
	p2pHost, err := libp2p.New(options)
	if err != nil {
		return nil, err
	}

	node := &Node{p2pHost: p2pHost}

	s := mdns.NewMdnsService(node.p2pHost, DiscoveryServiceTag, &discoveryNotifee{ctx: ctx, h: p2pHost})

	if err := s.Start(); err != nil {
		panic(err)
	}

	return node, nil
}

type discoveryNotifee struct {
	ctx context.Context
	h   host.Host
}

func (n *discoveryNotifee) HandlePeerFound(pi peer.AddrInfo) {
	fmt.Printf("Discovered new peer %s\n", pi.ID.String())
	err := n.h.Connect(n.ctx, pi)
	if err != nil {
		fmt.Printf("Error connecting to peer %s: %s\n", pi.ID.String(), err)
		return
	}
}

func (n *Node) String() string {
	sb := strings.Builder{}
	sb.WriteString("Host ID:" + n.p2pHost.ID().String())
	sb.WriteString("\nHost Addresses:")
	for _, addr := range n.p2pHost.Addrs() {
		sb.WriteString(fmt.Sprintf("  %s/p2p/%s\n", addr, n.p2pHost.ID()))
	}
	return sb.String()
}

func (n *Node) Subscribe(ctx context.Context, topic string) (<-chan *pubsub.Subscription, <-chan error) {
	fmt.Println("Subscribing...")
	subChan := make(chan *pubsub.Subscription, 1)
	errChan := make(chan error, 1)

	var joinedTopic *pubsub.Topic

	go func() {
		defer close(subChan)
		defer close(errChan)

		pubSub, err := pubsub.NewGossipSub(ctx, n.p2pHost)
		if err != nil {
			errChan <- err
			return
		}

		joinedTopic, err = pubSub.Join(topic)
		if err != nil {
			errChan <- err
			return
		}
		fmt.Println("Joined topic ", topic)

		subscription, err := joinedTopic.Subscribe()
		if err != nil {
			errChan <- err
			return
		}

		n.Topic = joinedTopic

		select {
		case subChan <- subscription:
		case <-ctx.Done():
			errChan <- ctx.Err()
		}
	}()

	return subChan, errChan
}

func (n *Node) Hostname() string {
	return n.p2pHost.ID().String()
}
