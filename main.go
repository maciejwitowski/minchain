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

	var inputs []lib.TransactionsInput
	for _, i := range config.Inputs {
		switch i {
		case lib.INPUT_STDIN:
			inputs = append(inputs, lib.NewUserInput())
		case lib.INPUT_API:
			httpApi := lib.NewHttpApi("0.0.0.0:8080")
			// TODO move into app.start
			log.Println("before start")

			go func() {
				err := httpApi.Start()
				if err != nil {
					log.Fatal("error starting http: ", err)
				}
			}()

			inputs = append(inputs, httpApi)
		default:
			log.Fatalf("Unknown input type: %s", i)
		}
	}

	application := app.NewApp(
		mempool,
		db,
		validator.NewBlockValidator(db),
		core.NewWallet(config.PrivateKey),
		config,
		node.Publisher,
		node.Consumer,
		inputs,
	)
	application.Start(ctx)

	select {}
}
