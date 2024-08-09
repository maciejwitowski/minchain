package app

import (
	"context"
	"log"
	"minchain/core"
	"minchain/database"
	"minchain/genesis"
	"minchain/lib"
	"minchain/p2p"
	"minchain/services"
	"minchain/validator"
)

type App struct {
	mempool           core.Mempool
	database          database.Database
	blockValidator    validator.Validator
	wallet            *core.Wallet
	config            lib.Config
	publisher         p2p.Publisher
	consumer          p2p.Consumer
	transactionsInput lib.TransactionsInput
}

func NewApp(
	mempool core.Mempool,
	database database.Database,
	blockValidator validator.Validator,
	wallet *core.Wallet,
	config lib.Config,
	publisher p2p.Publisher,
	consumer p2p.Consumer,
	transactionsInput lib.TransactionsInput,
) *App {
	return &App{
		mempool:           mempool,
		database:          database,
		blockValidator:    blockValidator,
		wallet:            wallet,
		config:            config,
		publisher:         publisher,
		consumer:          consumer,
		transactionsInput: transactionsInput,
	}
}

func (app *App) Start(ctx context.Context) {
	app.initializeGenesisState()
	app.launchTransactionsProcessing(ctx)
	app.launchBlocksProcessing(ctx)

	if app.config.IsBlockProducer {
		go core.NewBlockProducer(app.mempool, app.database, app.publisher, app.config).BuildAndPublishBlock(ctx)
	}
}

func (app *App) initializeGenesisState() {
	err := genesis.InitializeGenesisState(app.database)
	if err != nil {
		log.Fatal(err)
	}
}

func (app *App) launchTransactionsProcessing(ctx context.Context) {
	processTransactions := services.NewProcessTransactionsService(
		app.mempool,
		app.wallet,
		app.publisher,
		app.consumer,
		app.transactionsInput,
	)
	processTransactions.Start(ctx)
}

func (app *App) launchBlocksProcessing(ctx context.Context) {
	blocksProcessing := services.NewProcessBlocksService(
		app.blockValidator,
		app.database,
		app.mempool,
		app.consumer,
	)
	blocksProcessing.Start(ctx)
}
