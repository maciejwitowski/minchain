package services

import (
	"context"
	"fmt"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
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
	messageProcessor := make(chan types.Tx, 1)
	go p.publishToMpool(ctx)
	go p.consumeTransactionsFromMempool(ctx, sub, messageProcessor)
	go p.processMessages(ctx, messageProcessor)
}

func (p *ProcessTransactions) consumeTransactionsFromMempool(ctx context.Context, sub *pubsub.Subscription, messageProcessor chan<- types.Tx) {
	for {
		msg, err := sub.Next(ctx)
		if err != nil {
			fmt.Println("Subscription error:", err)
			return
		}
		txJson, err := types.FromJSON(msg.Data)
		if err != nil {
			fmt.Println("Error deserializing tx:", err)
			return
		}
		messageProcessor <- *txJson
	}
}

func (p *ProcessTransactions) processMessages(ctx context.Context, processor chan types.Tx) {
	for {
		select {
		case tx := <-processor:
			// Add Tx to mpool
			p.mempool.ValidateAndStorePending(tx)
		case <-ctx.Done():
			fmt.Println("processMessages cancelled")
			return
		}
	}
}

func (p *ProcessTransactions) publishToMpool(ctx context.Context) {
	for message := range p.userInput {
		tx, err := p.wallet.SignedTransaction(message)
		if err != nil {
			fmt.Println("Error building transaction:", err)
			return
		}

		txJson, err := tx.ToJSON()
		if err != nil {
			fmt.Println("Serialization error :", err)
			return
		}

		if err := p.pubSubTopic.Publish(ctx, txJson); err != nil {
			fmt.Println("Publish error:", err)
		}
	}
}
