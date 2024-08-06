package chain

import (
	"bytes"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"time"
)

// BlockProducer reads mempool and then produces and publishes a block
func BlockProducer(mempool *Mempool) {
	blocktime := time.NewTicker(5 * time.Second)
	defer blocktime.Stop()

	for {
		select {
		case <-blocktime.C:
			txs := mempool.GetPendingTransactions()
			txHash, err := Hash(txs)
			if err != nil {
				fmt.Println("Block production failed. Skipping") // TODO error handling
				continue
			}

			parent := ChainstoreInstance.GetHead()
			block := Block{
				Header: BlockHeader{
					ParentHash:      parent.Header.ParentHash,
					TransactionHash: txHash,
				},
				Txs: txs,
			}

			fmt.Println("Produced block: ", block.PrettyPrint())
		}
	}
}

func Hash(txs []Tx) (common.Hash, error) {
	buffer := bytes.Buffer{}
	for _, tx := range txs {
		serialized, err := tx.ToJSON()
		if err != nil {
			return common.Hash{}, err
		}
		txHash := crypto.Keccak256(serialized)
		buffer.Write(txHash)
	}

	combinedHash := crypto.Keccak256(buffer.Bytes())
	return common.BytesToHash(combinedHash), nil
}
