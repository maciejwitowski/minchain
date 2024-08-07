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
}

func InitApplicationDependencies() *App {
	db := database.NewMemoryDatabase()

	return &App{
		Mempool:        core.InitMempool(),
		Database:       db,
		Chainstore:     core.NewChainstore(db),
		BlockValidator: validator.NewBlockValidator(db),
	}
}
