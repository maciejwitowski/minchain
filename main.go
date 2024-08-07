package main

import (
	"context"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"log"
	"minchain/core"
	"minchain/core/types"
	"minchain/genesis"
	"minchain/lib"
	"minchain/p2p"
	"minchain/services"
	"time"
)

var Dependencies = lib.InitApplicationDependencies(lib.InitConfig())

func main() {
	//logging.SetAllLoggers(logging.LevelDebug)

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// TODO Better dependency injection
	initializeGenesisState(Dependencies)

	node, err := p2p.InitNode(ctx, Dependencies.Config)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(node.String())

	launchTransactionsProcessing(ctx, node)

	blkSubscription, err := node.SubscribeToBlocks()
	if err != nil {
		log.Println("Error subscribing to blocks:", err)
		return
	}
	onSubscribedToBlocks(ctx, blkSubscription)

	go lib.Monitor(ctx, Dependencies.Mempool, 1*time.Second)

	isBlockProducer := Dependencies.Config.IsBlockProducer
	log.Println("IsBlockProducer=", isBlockProducer)
	if isBlockProducer {
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

func launchTransactionsProcessing(ctx context.Context, node *p2p.Node) {
	txSubscription, err := node.SubscribeToTransactions()
	if err != nil {
		log.Println("Error subscribing to transactions:", err)
		return
	}

	userInput := lib.UserInput(ctx)
	processTransactions := services.NewProcessTransactionsService(
		Dependencies.Mempool,
		Dependencies.Wallet,
		userInput,
		node.TxTopic,
	)
	processTransactions.Start(ctx, txSubscription)
}
func onSubscribedToBlocks(ctx context.Context, sub *pubsub.Subscription) {
	blocksProcessor := make(chan types.Block, 1)
	go consumeBlocksFromMempool(ctx, sub, blocksProcessor)
	go processBlocks(ctx, blocksProcessor)
}

func consumeBlocksFromMempool(ctx context.Context, sub *pubsub.Subscription, blocksProcessor chan types.Block) {
	for {
		msg, err := sub.Next(ctx)
		if err != nil {
			log.Println("Subscription error:", err)
			return
		}
		blkJson, err := types.BlockFromJson(msg.Data)
		if err != nil {
			log.Println("Error deserializing block:", err)
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
			log.Println("processBlocks cancelled")
			return
		}
	}
}
