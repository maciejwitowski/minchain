package core

import (
	"bytes"
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"log"
	"minchain/core/types"
	"time"
)

// BlockProducer reads mempool and then produces and publishes a block
func BlockProducer(ctx context.Context, mempool *Mempool, topic *pubsub.Topic) {
	blocktime := time.NewTicker(5 * time.Second)
	defer blocktime.Stop()

	for {
		select {
		case <-blocktime.C:
			log.Println("Check if block should be produced")
			block, err := mempool.BuildBlockFromTransactions(builder)

			if err != nil {
				log.Println("error building the block")
				continue
			}

			if block != nil {
				log.Println("Produced block: ", block.PrettyPrint())

				blkJson, err := block.ToJson()
				if err != nil {
					fmt.Println("Serialization error :", err)
					return
				}

				if err := topic.Publish(ctx, blkJson); err != nil {
					fmt.Println("Publish error:", err)
				}
			}
		}
	}
}

type BlockBuilder interface {
	buildBlock([]types.Tx) (*types.Block, error)
}

func builder(txs []types.Tx) (*types.Block, error) {
	txHash, err := Hash(txs)
	if err != nil {
		fmt.Println("Block production failed. Skipping") // TODO error handling
		return nil, err
	}

	parent := ChainstoreInstance.GetHead()
	block := types.Block{
		Header: types.BlockHeader{
			ParentHash:      parent.Header.ParentHash,
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
