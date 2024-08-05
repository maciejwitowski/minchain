package main

import (
	"bufio"
	"context"
	"fmt"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"log"
	"minchain/chain"
	"minchain/common"
	"minchain/p2p"
	"os"
	"time"
)

func main() {
	config := common.InitConfig()
	ctx := context.Background()
	mpool := chain.InitMempool()

	node, err := p2p.InitNode(ctx, config)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(node.String())

	subChan, errChan := node.Subscribe(ctx, config.PubSubTopic)
	var sub *pubsub.Subscription
	select {
	case sub = <-subChan:
		fmt.Println("Subscribed.")
		onSubscribed(ctx, node, sub, mpool)
	case err := <-errChan:
		fmt.Println("Subscription error:", err)
	case <-ctx.Done():
		fmt.Println("Subscription timed out")
	}

	ticker := time.NewTicker(5 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				mpool.DumpTx()
			}
		}
	}()

	select {}
}

func onSubscribed(ctx context.Context, node *p2p.Node, sub *pubsub.Subscription, mpool *chain.Mempool) {
	messageProcessor := make(chan chain.Tx, 1)
	go readMessages(ctx, sub, messageProcessor)
	go processMessages(ctx, messageProcessor, mpool)
	if node.Topic != nil {
		go func() {
			NewPublisher(node.Topic, node.Hostname()).StartPublishing(ctx, mpool)
		}()
	}
}

func processMessages(ctx context.Context, processor chan chain.Tx, mpool *chain.Mempool) {
	for {
		select {
		case tx := <-processor:
			fmt.Println("received ", tx)
			// Adding inbound tx to mpool
			mpool.HandleTransaction(tx)
		case <-ctx.Done():
			fmt.Println("processMessages cancelled")
			return
		}
	}
}

func readMessages(ctx context.Context, sub *pubsub.Subscription, messageProcessor chan<- chain.Tx) {
	for {
		msg, err := sub.Next(ctx)
		if err != nil {
			fmt.Println("Subscription error:", err)
			return
		}
		txJson, err := chain.FromJSON(msg.Data)
		if err != nil {
			fmt.Println("Error deserializing tx:", err)
			return
		}
		messageProcessor <- *txJson
		fmt.Printf("%s: %s", msg.ReceivedFrom.ShortString(), string(msg.Data))
	}
}

type ReaderPublisher struct {
	reader   *bufio.Reader
	topic    *pubsub.Topic
	hostname string
}

func NewPublisher(topic *pubsub.Topic, hostname string) *ReaderPublisher {
	return &ReaderPublisher{
		reader:   bufio.NewReader(os.Stdin),
		topic:    topic,
		hostname: hostname,
	}
}

func (rp *ReaderPublisher) StartPublishing(ctx context.Context, mpool *chain.Mempool) {
	for {
		fmt.Print("> ")
		message, err := rp.reader.ReadString('\n')
		if err != nil {
			return
		}

		tx := chain.Tx{
			From: rp.hostname,
			Data: message,
		}

		txJson, err := tx.ToJSON()
		if err != nil {
			fmt.Println("Serialization error :", err)
			return
		}

		if err := rp.topic.Publish(ctx, txJson); err != nil {
			fmt.Println("Publish error:", err)
		}
	}
}
