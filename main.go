package main

import (
	"context"
	"log"
	"minchain/app"
	"minchain/core"
	"minchain/database"
	"minchain/lib"
	"minchain/monitor"
	"minchain/p2p"
	"minchain/validator"
	"time"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var db database.Database
	db, err := database.NewDiskDatabase()
	if err != nil {
		log.Fatal(err)
	}

	defer func(db database.Database) {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(db)

	config := lib.InitConfig()
	node, err := p2p.InitNode(ctx, config)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Initalized node: ", node.String())

	mempool := core.NewMempool()
	go monitor.Monitor(ctx, mempool, 1*time.Second)

	application := app.NewApp(
		mempool,
		db,
		core.NewChainhead(db),
		validator.NewBlockValidator(db),
		core.NewWallet(config.PrivateKey),
		config,
		node.Publisher,
		node.Consumer,
		lib.NewUserInput(),
	)
	application.Start(ctx)

	select {}
}
