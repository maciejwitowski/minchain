package main

import (
	"minchain/core"
	"minchain/database"
	"minchain/lib"
	"minchain/validator"
)

// App type keeps app-level dependencies (singletons)
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
