package main

import (
	"bufio"
	"context"
	"fmt"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"log"
	"minchain/core"
	"minchain/core/types"
	"minchain/genesis"
	"minchain/lib"
	"minchain/p2p"
	"os"
	"time"
)

var Dependencies = lib.InitApplicationDependencies()

func main() {
	//logging.SetAllLoggers(logging.LevelDebug)

	config, err := lib.InitConfig()
	if err != nil {
		log.Fatal("Init config error:", err)
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// TODO Better dependency injection
	initializeGenesisState(Dependencies)

	node, err := p2p.InitNode(ctx, config)
	wallet := core.NewWallet(config.PrivateKey)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(node.String())

	txSubscription, err := node.SubscribeToTransactions()
	if err != nil {
		fmt.Println("Error subscribing to transactions:", err)
		return
	}
	onSubscribedToTransactions(ctx, node, txSubscription, wallet)

	blkSubscription, err := node.SubscribeToBlocks()
	if err != nil {
		fmt.Println("Error subscribing to blocks:", err)
		return
	}
	onSubscribedToBlocks(ctx, blkSubscription)

	go lib.Monitor(ctx, Dependencies.Mempool, 1*time.Second)

	log.Println("IsBlockProducer=", config.IsBlockProducer)
	if config.IsBlockProducer {
		go core.NewBlockProducer(Dependencies.Mempool, node.BlocksTopic, Dependencies.Chainstore).BuildAndPublishBlock(ctx)
	}

	select {}
}

func initializeGenesisState(app *lib.App) {
	err := genesis.InitializeGenesisState(app.Database, app.Chainstore)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Initialised genesis")
}

// Extract to separate service
func onSubscribedToTransactions(ctx context.Context, node *p2p.Node, sub *pubsub.Subscription, wallet *core.Wallet) {
	messageProcessor := make(chan types.Tx, 1)
	go consumeTransactionsFromMempool(ctx, sub, messageProcessor)
	go processMessages(ctx, messageProcessor)
	messages := make(chan string)
	go readUserInput(messages)
	go publishToMpool(ctx, node.TxTopic, wallet, messages)
}

func onSubscribedToBlocks(ctx context.Context, sub *pubsub.Subscription) {
	blocksProcessor := make(chan types.Block, 1)
	go consumeBlocksFromMempool(ctx, sub, blocksProcessor)
	go processBlocks(ctx, blocksProcessor)
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

func processMessages(ctx context.Context, processor chan types.Tx) {
	for {
		select {
		case tx := <-processor:
			// Add Tx to mpool
			Dependencies.Mempool.ValidateAndStorePending(tx)
		case <-ctx.Done():
			fmt.Println("processMessages cancelled")
			return
		}
	}
}

func consumeTransactionsFromMempool(ctx context.Context, sub *pubsub.Subscription, messageProcessor chan<- types.Tx) {
	for {
		msg, err := sub.Next(ctx)
		if err != nil {
			fmt.Println("Subscription error:", err)
			return
		}
		txJson, err := types.FromJSON(msg.Data)
		if err != nil {
			fmt.Println("Error deserializing tx:", err)
			return
		}
		messageProcessor <- *txJson
	}
}

func consumeBlocksFromMempool(ctx context.Context, sub *pubsub.Subscription, blocksProcessor chan types.Block) {
	for {
		msg, err := sub.Next(ctx)
		if err != nil {
			fmt.Println("Subscription error:", err)
			return
		}
		blkJson, err := types.BlockFromJson(msg.Data)
		if err != nil {
			fmt.Println("Error deserializing block:", err)
			return
		}
		blocksProcessor <- *blkJson
	}
}

func processBlocks(ctx context.Context, processor chan types.Block) {
	for {
		select {
		case blk := <-processor:
			log.Println("received block: ", blk.BlockHash())
			err := Dependencies.BlockValidator.Validate(&blk)
			if err != nil {
				log.Println("validator error ", err)
				continue
			}
			Dependencies.Database.PutBlock(&blk)
			Dependencies.Chainstore.SetHead(&blk)
			Dependencies.Mempool.PruneTransactions(blk.Transactions)
		case <-ctx.Done():
			fmt.Println("processBlocks cancelled")
			return
		}
	}
}

func publishToMpool(ctx context.Context, topic *pubsub.Topic, wallet *core.Wallet, userInput <-chan string) {
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
