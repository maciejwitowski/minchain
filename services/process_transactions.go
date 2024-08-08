package services

import (
	"context"
	"log"
	"minchain/core"
	"minchain/p2p"
)

type ProcessTransactions struct {
	mempool   core.Mempool
	wallet    *core.Wallet
	publisher p2p.Publisher
	consumer  p2p.Consumer
	userInput <-chan string
}

func NewProcessTransactionsService(mempool core.Mempool, wallet *core.Wallet, publisher p2p.Publisher, consumer p2p.Consumer, userInput <-chan string) *ProcessTransactions {
	return &ProcessTransactions{
		wallet:    wallet,
		mempool:   mempool,
		publisher: publisher,
		consumer:  consumer,
		userInput: userInput,
	}
}

func (p *ProcessTransactions) Start(ctx context.Context) {
	go p.publishInputToMempool(ctx)
	go p.processMempoolTransactions(ctx)
}

func (p *ProcessTransactions) publishInputToMempool(ctx context.Context) {
	for message := range p.userInput {
		log.Println("user input: ", message)
		tx, err := p.wallet.SignedTransaction(message)
		if err != nil {
			log.Println("Error building transaction:", err)
			return
		}

		log.Println("Publishing: ", tx.PrettyPrint())
		if err := p.publisher.PublishTransaction(ctx, tx); err != nil {
			log.Println("Publish error:", err)
		}
	}
}

func (p *ProcessTransactions) processMempoolTransactions(ctx context.Context) {
	for {
		tx, err := p.consumer.ConsumeTransaction(ctx)
		if err != nil {
			log.Println("Error deserializing tx:", err)
			return
		}
		p.mempool.ValidateAndStorePending(tx)
	}
}
