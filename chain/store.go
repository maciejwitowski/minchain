package chain

import (
	"github.com/ethereum/go-ethereum/common"
)

var ChainstoreInstance = NewChainstore()

type Chainstore struct {
	blocks       []*Block
	blocksLookup map[common.Hash]Block
}

func NewChainstore() *Chainstore {
	return &Chainstore{
		blocksLookup: make(map[common.Hash]Block),
	}
}

func (c *Chainstore) AddBlock(block *Block) {

}

func (c *Chainstore) GetHead() *Block {
	if len(c.blocks) > 0 {
		return c.blocks[len(c.blocks)-1]
	} else {
		// empty "genesis" block
		return &Block{
			Header: BlockHeader{
				ParentHash:      common.Hash{},
				TransactionHash: common.Hash{},
			},
			Txs: make([]Tx, 0),
		}
	}
}
