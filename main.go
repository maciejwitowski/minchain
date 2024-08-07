package main

import (
	"bufio"
	"context"
	"fmt"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"log"
	"minchain/core"
	"minchain/core/types"
	"minchain/database"
	"minchain/genesis"
	"minchain/lib"
	"minchain/p2p"
	"minchain/validator"
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

	// TODO Better dependency injection
	mpool := core.InitMempool()
	db := database.NewMemoryDatabase()
	chainstore := core.NewChainstore(db)

	err = genesis.InitializeGenesisState(db, chainstore)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Initialised genesis")

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
	onSubscribedToTransactions(ctx, node, txSubscription, mpool, wallet)

	blkSubscription, err := node.SubscribeToBlocks()
	if err != nil {
		fmt.Println("Error subscribing to blocks:", err)
		return
	}
	onSubscribedToBlocks(ctx, blkSubscription, mpool, db, chainstore)

	go lib.Monitor(ctx, mpool, 1*time.Second)

	log.Println("IsBlockProducer=", config.IsBlockProducer)
	if config.IsBlockProducer {
		go core.NewBlockProducer(mpool, node.BlocksTopic, chainstore).BuildAndPublishBlock(ctx)
	}

	select {}
}

// Extract to separate service
func onSubscribedToTransactions(ctx context.Context, node *p2p.Node, sub *pubsub.Subscription, mpool *core.Mempool, wallet *core.Wallet) {
	messageProcessor := make(chan types.Tx, 1)
	go consumeTransactionsFromMempool(ctx, sub, messageProcessor)
	go processMessages(ctx, messageProcessor, mpool)
	messages := make(chan string)
	go readUserInput(messages)
	go publishToMpool(ctx, node.TxTopic, wallet, messages)
}

func onSubscribedToBlocks(ctx context.Context, sub *pubsub.Subscription, mpool *core.Mempool, db database.Database, store *core.Chainstore) {
	blocksProcessor := make(chan types.Block, 1)
	go consumeBlocksFromMempool(ctx, sub, blocksProcessor)
	go processBlocks(ctx, blocksProcessor, mpool, db, store)
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

func processMessages(ctx context.Context, processor chan types.Tx, mpool *core.Mempool) {
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

func processBlocks(ctx context.Context, processor chan types.Block, mpool *core.Mempool, db database.Database, store *core.Chainstore) {
	for {
		select {
		case blk := <-processor:
			log.Println("received block: ", blk.BlockHash())
			blockValidator := validator.NewBlockValidator(db)
			err := blockValidator.Validate(&blk)
			if err != nil {
				log.Println("validator error ", err)
				continue
			}
			db.PutBlock(&blk)
			store.SetHead(&blk)
			mpool.PruneTransactions(blk.Transactions)
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
