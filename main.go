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

	subChan, errChan := node.SubscribeToTransactions(ctx)
	var sub *pubsub.Subscription
	select {
	case sub = <-subChan:
		fmt.Println("Subscribed to transactions.")
		onSubscribedToTransactions(ctx, node, sub, mpool, wallet)
	case err := <-errChan:
		fmt.Println("Transactions subscription error:", err)
	case <-ctx.Done():
		fmt.Println("Transactions subscription timed out")
	}

	subChan, errChan = node.SubscribeToBlocks(ctx)
	select {
	case sub = <-subChan:
		fmt.Println("Subscribed to blocks.")
		onSubscribedToBlocks(ctx, node, sub, mpool, wallet)
	case err := <-errChan:
		fmt.Println("Block subscription error:", err)
	case <-ctx.Done():
		fmt.Println("Block subscription timed out")
	}

	go lib.Monitor(ctx, mpool, 1*time.Second)

	log.Println("IsBlockProducer=", config.IsBlockProducer)
	if config.IsBlockProducer {
		go chain.BlockProducer(mpool)
	}

	select {}
}

// Extract to separate service
func onSubscribedToTransactions(ctx context.Context, node *p2p.Node, sub *pubsub.Subscription, mpool *chain.Mempool, wallet *chain.Wallet) {
	messageProcessor := make(chan chain.Tx, 1)
	go consumeTransactionsFromMempool(ctx, sub, messageProcessor)
	go processMessages(ctx, messageProcessor, mpool)
	if node.MpoolTopic != nil {
		messages := make(chan string)
		go readUserInput(messages)
		go publishToMpool(ctx, node.MpoolTopic, wallet, messages)
	}
}

func onSubscribedToBlocks(ctx context.Context, node *p2p.Node, sub *pubsub.Subscription, mpool *chain.Mempool, wallet *chain.Wallet) {
	blocksProcessor := make(chan chain.Block, 1)
	go consumeBlocksFromMempool(ctx, sub, blocksProcessor)
	go processBlocks(ctx, blocksProcessor, mpool)
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

func consumeTransactionsFromMempool(ctx context.Context, sub *pubsub.Subscription, messageProcessor chan<- chain.Tx) {
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

func consumeBlocksFromMempool(ctx context.Context, sub *pubsub.Subscription, blocksProcessor chan chain.Block) {
	for {
		msg, err := sub.Next(ctx)
		if err != nil {
			fmt.Println("Subscription error:", err)
			return
		}
		blkJson, err := chain.BlockFromJson(msg.Data)
		if err != nil {
			fmt.Println("Error deserializing block:", err)
			return
		}
		blocksProcessor <- *blkJson
	}
}

func processBlocks(ctx context.Context, processor chan chain.Block, mpool *chain.Mempool) {
	for {
		select {
		case blk := <-processor:
			// Add Tx to mpool
			log.Println("received block: ", blk.PrettyPrint())
		case <-ctx.Done():
			fmt.Println("processBlocks cancelled")
			return
		}
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
