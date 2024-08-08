package p2p

import (
	"context"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"log"
	"minchain/core/types"
)

type Consumer interface {
	ConsumeTransaction(ctx context.Context) (*types.Tx, error)
	ConsumeBlock(ctx context.Context) (*types.Block, error)
}

type P2pConsumer struct {
	txSubscription     *pubsub.Subscription
	blocksSubscription *pubsub.Subscription
}

func NewP2pConsumer(txSubscription *pubsub.Subscription, blocksSubscription *pubsub.Subscription) Consumer {
	return &P2pConsumer{
		txSubscription:     txSubscription,
		blocksSubscription: blocksSubscription,
	}
}

func (c *P2pConsumer) ConsumeTransaction(ctx context.Context) (*types.Tx, error) {
	msg, err := c.txSubscription.Next(ctx)
	if err != nil {
		return nil, err
	}

	tx, err := types.TransactionFromJSON(msg.Data)
	if err != nil {
		log.Println("Error deserializing transaction:", err)
		return nil, err
	}
	return tx, nil
}

func (c *P2pConsumer) ConsumeBlock(ctx context.Context) (*types.Block, error) {
	msg, err := c.blocksSubscription.Next(ctx)
	if err != nil {
		return nil, err
	}

	tx, err := types.BlockFromJson(msg.Data)
	if err != nil {
		log.Println("Error deserializing block:", err)
		return nil, err
	}
	return tx, nil
}
