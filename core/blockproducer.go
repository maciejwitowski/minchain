package core

import (
	"context"
	"log"
	"minchain/core/types"
	"minchain/database"
	"minchain/lib"
	"minchain/p2p"
	"time"
)

// BlockProducer reads mempool and then produces and publishes a block
type BlockProducer struct {
	mempool      Mempool
	database     database.Database
	config       lib.Config
	p2pPublisher p2p.Publisher
}

func NewBlockProducer(mempool Mempool, database database.Database, p2pPublisher p2p.Publisher, config lib.Config) *BlockProducer {
	return &BlockProducer{
		mempool:      mempool,
		database:     database,
		p2pPublisher: p2pPublisher,
		config:       config,
	}
}

// TODO Split block production and publishing
func (bp *BlockProducer) BuildAndPublishBlock(ctx context.Context) {
	blocktimeTicker := time.NewTicker(bp.config.BlockTime)
	defer blocktimeTicker.Stop()

	for {
		select {
		case <-blocktimeTicker.C:
			// TODO more advanced selection logic
			transactions := bp.mempool.ListPendingTransactions()
			if len(transactions) == 0 {
				continue
			}

			block, err := bp.buildBlock(transactions)
			log.Println("Building block. Block hash:", block.BlockHash())
			if err != nil {
				log.Println("error building the block")
				continue
			}

			if block != nil {
				log.Println("Produced block: ", block.PrettyPrint())
				if err := bp.p2pPublisher.PublishBlock(ctx, block); err != nil {
					log.Println("Publish error:", err)
				}

				bp.mempool.PruneTransactions(block.Transactions)
			}
		}
	}
}

type BlockBuilder interface {
	buildBlock([]types.Tx) (*types.Block, error)
}

func (bp *BlockProducer) buildBlock(txs []types.Tx) (*types.Block, error) {
	txHash, err := types.CombinedHash(txs)
	if err != nil {
		log.Println("Block production failed. Skipping") // TODO error handling
		return nil, err
	}

	parent, err := bp.database.GetHead()
	if err != nil {
		log.Fatal("No parent in database due to incorrect node initialization. Should never happen")
	}

	parentBlock, err := bp.database.GetBlockByHash(parent)
	if err != nil {
		log.Fatal("error getting parentBlock", err)
	}

	block := types.Block{
		Header: types.BlockHeader{
			ParentHash:      parent,
			TransactionHash: txHash,
			Height:          parentBlock.Header.Height + 1,
		},
		Transactions: txs,
	}
	return &block, nil
}
