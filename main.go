package main

import (
	"context"
	"log"
	"minchain/core"
	"minchain/genesis"
	"minchain/lib"
	"minchain/monitor"
	"minchain/p2p"
	"minchain/services"
	"time"
)

var Dependencies = InitApplicationDependencies(lib.InitConfig())

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
	launchBlocksProcessing(ctx, node)

	go monitor.Monitor(ctx, Dependencies.Mempool, 1*time.Second)

	if Dependencies.Config.IsBlockProducer {
		go core.NewBlockProducer(Dependencies.Mempool, node.BlocksTopic, Dependencies.Chainstore, Dependencies.Config).BuildAndPublishBlock(ctx)
	}

	select {}
}

func initializeGenesisState(app *App) {
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

func launchBlocksProcessing(ctx context.Context, node *p2p.Node) {
	blockSubscription, err := node.SubscribeToBlocks()
	if err != nil {
		log.Println("Error subscribing to blocks:", err)
		return
	}

	blocksProcessing := services.NewProcessBlocksService(
		Dependencies.BlockValidator,
		Dependencies.Database,
		Dependencies.Chainstore,
		Dependencies.Mempool,
	)
	blocksProcessing.Start(ctx, blockSubscription)
}
