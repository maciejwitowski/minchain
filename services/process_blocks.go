package services

import (
	"context"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"log"
	"minchain/core"
	"minchain/core/types"
	"minchain/database"
	"minchain/validator"
)

type ProcessBlocks struct {
	blockValidator validator.Validator
	database       database.Database
	chainstore     core.Chainstore
	mempool        core.Mempool
}

func NewProcessBlocksService(blockValidator validator.Validator, database database.Database, chainstore core.Chainstore, mempool core.Mempool) *ProcessBlocks {
	return &ProcessBlocks{
		blockValidator: blockValidator,
		database:       database,
		chainstore:     chainstore,
		mempool:        mempool,
	}
}

func (p *ProcessBlocks) Start(ctx context.Context, sub *pubsub.Subscription) {
	for {
		msg, err := sub.Next(ctx)
		if err != nil {
			log.Println("Subscription error:", err)
			return
		}
		block, err := types.BlockFromJson(msg.Data)
		if err != nil {
			log.Println("Error deserializing block:", err)
			return
		}

		log.Println("received block: ", block.BlockHash())
		err = p.blockValidator.Validate(block)
		if err != nil {
			log.Println("validator error ", err)
			continue
		}

		p.database.PutBlock(block)
		p.chainstore.SetHead(block)
		p.mempool.PruneTransactions(block.Transactions)
	}
}
