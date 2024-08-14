package services

import (
	"context"
	"log"
	"minchain/core"
	"minchain/lib"
	"minchain/p2p"
)

type ProcessTransactions struct {
	mempool   core.Mempool
	wallet    *core.Wallet
	publisher p2p.Publisher
	consumer  p2p.Consumer
	inputs    []lib.TransactionsInput
}

func NewProcessTransactionsService(mempool core.Mempool, wallet *core.Wallet, publisher p2p.Publisher, consumer p2p.Consumer, inputs []lib.TransactionsInput) *ProcessTransactions {
	return &ProcessTransactions{
		wallet:    wallet,
		mempool:   mempool,
		publisher: publisher,
		consumer:  consumer,
		inputs:    inputs,
	}
}

func (p *ProcessTransactions) Start(ctx context.Context) {
	for _, input := range p.inputs {
		go p.publishTransactionsToNetwork(ctx, input)
	}
	go p.consumeTransactionsFromNetwork(ctx)
}

func (p *ProcessTransactions) publishTransactionsToNetwork(ctx context.Context, input lib.TransactionsInput) {
	for message := range input.InputChannel(ctx) {
		tx, err := p.wallet.SignedTransaction(message)
		if err != nil {
			log.Println("Error building transaction:", err)
			return
		}

		if err := p.publisher.PublishTransaction(ctx, tx); err != nil {
			log.Println("Publish error:", err)
		}
	}
}

func (p *ProcessTransactions) consumeTransactionsFromNetwork(ctx context.Context) {
	for {
		tx, err := p.consumer.ConsumeTransaction(ctx)
		log.Println("Received tx from the network: ", tx.PrettyPrint())
		if err != nil {
			log.Println("Error deserializing tx:", err)
			return
		}
		p.mempool.ValidateAndStorePending(tx)
	}
}
