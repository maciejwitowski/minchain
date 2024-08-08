package main

import (
	"context"
	"log"
	"minchain/core"
	"minchain/database"
	"minchain/genesis"
	"minchain/lib"
	"minchain/monitor"
	"minchain/p2p"
	"minchain/services"
	"minchain/validator"
	"time"
)

var Dependencies = InitApplicationDependencies(lib.InitConfig())

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

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
		go core.NewBlockProducer(Dependencies.Mempool, node.Publisher, Dependencies.Chainstore, Dependencies.Config).BuildAndPublishBlock(ctx)
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
	processTransactions := services.NewProcessTransactionsService(
		Dependencies.Mempool,
		Dependencies.Wallet,
		node.Publisher,
		node.Consumer,
		lib.NewUserInput(),
	)
	processTransactions.Start(ctx)
}

func launchBlocksProcessing(ctx context.Context, node *p2p.Node) {
	blocksProcessing := services.NewProcessBlocksService(
		Dependencies.BlockValidator,
		Dependencies.Database,
		Dependencies.Chainstore,
		Dependencies.Mempool,
		node.Consumer,
	)
	blocksProcessing.Start(ctx)
}

// App type keeps applevel dependencies (singletons)
type App struct {
	Mempool        core.Mempool
	Database       database.Database
	Chainstore     core.Chainstore
	BlockValidator validator.Validator
	Wallet         *core.Wallet
	Config         lib.Config
}

func InitApplicationDependencies(config lib.Config) *App {
	db := database.NewMemoryDatabase()

	return &App{
		Mempool:        core.InitMempool(),
		Database:       db,
		Chainstore:     core.NewChainstore(db),
		BlockValidator: validator.NewBlockValidator(db),
		Wallet:         core.NewWallet(config.PrivateKey),
		Config:         config,
	}
}
