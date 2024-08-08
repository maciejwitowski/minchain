package core

import (
	"bytes"
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"log"
	"minchain/core/types"
	"minchain/lib"
	"time"
)

// BlockProducer reads mempool and then produces and publishes a block
type BlockProducer struct {
	mempool    Mempool
	topic      *pubsub.Topic
	chainstore Chainstore
	config     lib.Config
}

func NewBlockProducer(mempool Mempool, topic *pubsub.Topic, chainstore Chainstore, config lib.Config) *BlockProducer {
	return &BlockProducer{
		mempool:    mempool,
		topic:      topic,
		chainstore: chainstore,
		config:     config,
	}
}

// TODO Split block production and publishing
func (bp *BlockProducer) BuildAndPublishBlock(ctx context.Context) {
	blocktimeTicker := time.NewTicker(bp.config.BlockTime * time.Second)
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
			if err != nil {
				log.Println("error building the block")
				continue
			}

			if block != nil {
				log.Println("Produced block: ", block.PrettyPrint())

				blkJson, err := block.ToJson()
				if err != nil {
					log.Println("Serialization error :", err)
					continue
				}

				if err := bp.topic.Publish(ctx, blkJson); err != nil {
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
	txHash, err := Hash(txs)
	if err != nil {
		log.Println("Block production failed. Skipping") // TODO error handling
		return nil, err
	}

	parent := bp.chainstore.GetHead()
	block := types.Block{
		Header: types.BlockHeader{
			ParentHash:      parent.BlockHash(),
			TransactionHash: txHash,
		},
		Transactions: txs,
	}
	return &block, nil
}

func Hash(txs []types.Tx) (common.Hash, error) {
	buffer := bytes.Buffer{}
	for _, tx := range txs {
		hashBytes, err := tx.HashBytes()
		if err != nil {
			return common.Hash{}, err
		}
		buffer.Write(hashBytes)
	}

	combinedHash := crypto.Keccak256(buffer.Bytes())
	return common.BytesToHash(combinedHash), nil
}
