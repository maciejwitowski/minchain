package main

import (
	"bufio"
	"context"
	"fmt"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"log"
	"minchain/common"
	"minchain/p2p"
	"os"
)

func main() {
	config := common.InitConfig()
	ctx := context.Background()
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
		onSubscribed(ctx, node, sub)
	case err := <-errChan:
		fmt.Println("Subscription error:", err)
	case <-ctx.Done():
		fmt.Println("Subscription timed out")
	}

	select {}
}

func onSubscribed(ctx context.Context, node *p2p.Node, sub *pubsub.Subscription) {
	go handleMessages(ctx, sub)
	if node.Topic != nil {
		go func() {
			NewReaderPublisher(node.Topic).Start(ctx)
		}()
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

type ReaderPublisher struct {
	reader *bufio.Reader
	topic  *pubsub.Topic
}

func NewReaderPublisher(topic *pubsub.Topic) *ReaderPublisher {
	return &ReaderPublisher{
		reader: bufio.NewReader(os.Stdin),
		topic:  topic,
	}
}

func (rp *ReaderPublisher) Start(ctx context.Context) {
	for {
		fmt.Print("> ")
		message, err := rp.reader.ReadString('\n')
		if err != nil {
			return
		}

		if err := rp.topic.Publish(ctx, []byte(message)); err != nil {
			fmt.Println("Publish error:", err)
		}
	}
}
