package services

import (
	"context"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"log"
	"minchain/core"
	"minchain/core/types"
)

type ProcessTransactions struct {
	mempool     core.Mempool
	wallet      *core.Wallet
	pubSubTopic *pubsub.Topic
	userInput   <-chan string
}

func NewProcessTransactionsService(mempool core.Mempool, wallet *core.Wallet, userInput <-chan string, pubSubTopic *pubsub.Topic) *ProcessTransactions {
	return &ProcessTransactions{
		wallet:      wallet,
		mempool:     mempool,
		userInput:   userInput,
		pubSubTopic: pubSubTopic,
	}
}

func (p *ProcessTransactions) Start(ctx context.Context, sub *pubsub.Subscription) {
	go p.publishInputToMempool(ctx)
	go p.processMempoolTransactions(ctx, sub)
}

func (p *ProcessTransactions) publishInputToMempool(ctx context.Context) {
	for message := range p.userInput {
		log.Println("user input: ", message)
		tx, err := p.wallet.SignedTransaction(message)
		if err != nil {
			log.Println("Error building transaction:", err)
			return
		}

		txJson, err := tx.ToJSON()
		if err != nil {
			log.Println("Serialization error :", err)
			return
		}

		log.Println("Publishing: ", tx.PrettyPrint())
		if err := p.pubSubTopic.Publish(ctx, txJson); err != nil {
			log.Println("Publish error:", err)
		}
	}
}

func (p *ProcessTransactions) processMempoolTransactions(ctx context.Context, sub *pubsub.Subscription) {
	for {
		msg, err := sub.Next(ctx)
		if err != nil {
			log.Println("Subscription error:", err)
			return
		}
		tx, err := types.FromJSON(msg.Data)
		if err != nil {
			log.Println("Error deserializing tx:", err)
			return
		}
		p.mempool.ValidateAndStorePending(tx)
	}
}
