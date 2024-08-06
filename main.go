package main

import (
	"bufio"
	"context"
	"fmt"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"log"
	"minchain/chain"
	"minchain/lib"
	"minchain/p2p"
	"os"
	"time"
)

func main() {
	//logging.SetAllLoggers(logging.LevelDebug)

	config, err := lib.InitConfig()
	if err != nil {
		log.Fatal("Init config error:", err)
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mpool := chain.InitMempool()

	node, err := p2p.InitNode(ctx, config)
	wallet := chain.NewWallet(config.PrivateKey)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(node.String())

	subChan, errChan := node.Subscribe(ctx, config.PubSubTopic)
	var sub *pubsub.Subscription
	select {
	case sub = <-subChan:
		fmt.Println("Subscribed.")
		onSubscribed(ctx, node, sub, mpool, wallet)
	case err := <-errChan:
		fmt.Println("Subscription error:", err)
	case <-ctx.Done():
		fmt.Println("Subscription timed out")
	}

	go lib.Monitor(ctx, mpool, 1*time.Second)
	select {}
}

// Extract to separate service
func onSubscribed(ctx context.Context, node *p2p.Node, sub *pubsub.Subscription, mpool *chain.Mempool, wallet *chain.Wallet) {
	messageProcessor := make(chan chain.Tx, 1)
	go consumeFromMpool(ctx, sub, messageProcessor)
	go processMessages(ctx, messageProcessor, mpool)
	if node.Topic != nil {
		messages := make(chan string)
		go readUserInput(messages)
		go publishToMpool(ctx, node.Topic, wallet, messages)
	}
}

func readUserInput(messages chan<- string) {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading the message:", err)
		}
		messages <- message
	}
}

func processMessages(ctx context.Context, processor chan chain.Tx, mpool *chain.Mempool) {
	for {
		select {
		case tx := <-processor:
			// Add Tx to mpool
			mpool.HandleTransaction(tx)
		case <-ctx.Done():
			fmt.Println("processMessages cancelled")
			return
		}
	}
}

func consumeFromMpool(ctx context.Context, sub *pubsub.Subscription, messageProcessor chan<- chain.Tx) {
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
	}
}

func publishToMpool(ctx context.Context, topic *pubsub.Topic, wallet *chain.Wallet, userInput <-chan string) {
	for message := range userInput {
		tx, err := wallet.SignedTransaction(message)
		if err != nil {
			fmt.Println("Error building transaction:", err)
			return
		}

		txJson, err := tx.ToJSON()
		if err != nil {
			fmt.Println("Serialization error :", err)
			return
		}

		if err := topic.Publish(ctx, txJson); err != nil {
			fmt.Println("Publish error:", err)
		}
	}
}
