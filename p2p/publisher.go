package p2p

import (
	"context"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"minchain/core/types"
)

type Publisher interface {
	PublishBlock(ctx context.Context, block *types.Block) error
	PublishTransaction(ctx context.Context, transaction *types.Tx) error
}

type P2pPublisher struct {
	txTopic     *pubsub.Topic
	blocksTopic *pubsub.Topic
}

func NewP2pPublisher(txTopic *pubsub.Topic, blocksTopic *pubsub.Topic) Publisher {
	return &P2pPublisher{
		txTopic:     txTopic,
		blocksTopic: blocksTopic,
	}
}

func (p *P2pPublisher) PublishBlock(ctx context.Context, block *types.Block) error {
	json, err := block.ToJson()
	if err != nil {
		return err
	}

	return p.blocksTopic.Publish(ctx, json)
}

func (p *P2pPublisher) PublishTransaction(ctx context.Context, transaction *types.Tx) error {
	json, err := transaction.ToJson()
	if err != nil {
		return err
	}

	return p.blocksTopic.Publish(ctx, json)
}
