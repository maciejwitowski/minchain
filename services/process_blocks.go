package services

import (
	"context"
	"log"
	"minchain/core"
	"minchain/database"
	"minchain/p2p"
	"minchain/validator"
)

type ProcessBlocks struct {
	blockValidator validator.Validator
	database       database.Database
	mempool        core.Mempool
	consumer       p2p.Consumer
}

func NewProcessBlocksService(blockValidator validator.Validator, database database.Database, mempool core.Mempool, consumer p2p.Consumer) *ProcessBlocks {
	return &ProcessBlocks{
		blockValidator: blockValidator,
		database:       database,
		mempool:        mempool,
		consumer:       consumer,
	}
}

func (p *ProcessBlocks) Start(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Println("context cancelled, stopping processing blocks")
				return
			default:
				block, err := p.consumer.ConsumeBlock(ctx)
				if err != nil {
					return
				}
				err = p.blockValidator.Validate(block)
				if err != nil {
					log.Println("validator error ", err)
					continue
				}

				log.Println("Valid block becomes new head", block.BlockHash().Hex())
				err = p.database.PutBlock(block)
				if err != nil {
					log.Println("Validator.PutBlock ", err)
				}

				err = p.database.SetHead(block.BlockHash())
				if err != nil {
					log.Println("Validator.SetHead ", err)
					return
				}
				p.mempool.PruneTransactions(block.Transactions)

				blockchainHashes, err := core.PrintBlockHashes(p.database)
				log.Println("Block added. New blockchain:", blockchainHashes)
			}
		}
	}()
}
