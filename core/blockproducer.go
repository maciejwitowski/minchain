package core

import (
	"context"
	"log"
	"minchain/core/types"
	"minchain/lib"
	"minchain/p2p"
	"time"
)

// BlockProducer reads mempool and then produces and publishes a block
type BlockProducer struct {
	mempool      Mempool
	chainhead    Chainhead
	config       lib.Config
	p2pPublisher p2p.Publisher
}

func NewBlockProducer(mempool Mempool, p2pPublisher p2p.Publisher, chainhead Chainhead, config lib.Config) *BlockProducer {
	return &BlockProducer{
		mempool:      mempool,
		p2pPublisher: p2pPublisher,
		chainhead:    chainhead,
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
			log.Println("Before publishing a block. Transactions cound:", len(transactions))
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

	parent := bp.chainhead.GetHead()
	block := types.Block{
		Header: types.BlockHeader{
			ParentHash:      parent.BlockHash(),
			TransactionHash: txHash,
		},
		Transactions: txs,
	}
	return &block, nil
}
