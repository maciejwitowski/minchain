package lib

import (
	"minchain/core"
	"minchain/database"
	"minchain/validator"
)

// App type keeps app-level dependencies (singletons)
type App struct {
	Mempool        core.Mempool
	Database       database.Database
	Chainstore     core.Chainstore
	BlockValidator validator.Validator
	Wallet         *core.Wallet
	Config         Config
}

func InitApplicationDependencies(config Config) *App {
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
