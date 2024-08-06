package p2p

import (
	"context"
	"fmt"
	"github.com/libp2p/go-libp2p"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
	"minchain/lib"
	"strings"
)

var addressTemplate = "/ip4/0.0.0.0/tcp/%d"
var transactionsTopic = "transactions"
var blocksTopic = "blocks"

const DiscoveryServiceTag = "p2p-service"

type Node struct {
	TxTopic     *pubsub.Topic
	BlocksTopic *pubsub.Topic

	p2pHost   host.Host
	gossipSub *pubsub.PubSub
}

func InitNode(ctx context.Context, config *lib.Config) (*Node, error) {
	options := libp2p.ListenAddrStrings(fmt.Sprintf(addressTemplate, config.ListeningPort))
	p2pHost, err := libp2p.New(options)
	if err != nil {
		return nil, err
	}

	gossipSub, err := pubsub.NewGossipSub(ctx, p2pHost)
	if err != nil {
		return nil, err
	}

	node := &Node{
		p2pHost:   p2pHost,
		gossipSub: gossipSub,
	}

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

func (n *Node) SubscribeToTransactions() (*pubsub.Subscription, error) {
	subscription, topic, err := n.subscribeToTopic(transactionsTopic)
	if err != nil {
		return nil, err
	}
	n.TxTopic = topic
	return subscription, nil
}

func (n *Node) SubscribeToBlocks() (*pubsub.Subscription, error) {
	subscription, topic, err := n.subscribeToTopic(blocksTopic)
	if err != nil {
		return nil, err
	}
	n.BlocksTopic = topic
	return subscription, nil
}

func (n *Node) subscribeToTopic(topic string) (*pubsub.Subscription, *pubsub.Topic, error) {
	joinedTopic, err := n.gossipSub.Join(topic)
	if err != nil {
		return nil, nil, err
	}

	subscription, err := joinedTopic.Subscribe()
	if err != nil {
		return nil, nil, err
	}

	return subscription, joinedTopic, nil
}

func (n *Node) Hostname() string {
	return n.p2pHost.ID().String()
}
