package core

import (
	"github.com/ethereum/go-ethereum/common"
	"minchain/core/types"
)

var ChainstoreInstance = NewChainstore()

type Chainstore struct {
	blocks       []*types.Block
	blocksLookup map[common.Hash]types.Block
}

func NewChainstore() *Chainstore {
	return &Chainstore{
		blocksLookup: make(map[common.Hash]types.Block),
	}
}

func (c *Chainstore) AddBlock(block *types.Block) {

}

func (c *Chainstore) GetHead() *types.Block {
	if len(c.blocks) > 0 {
		return c.blocks[len(c.blocks)-1]
	} else {
		// empty "genesis" block
		return &types.Block{
			Header: types.BlockHeader{
				ParentHash:      common.Hash{},
				TransactionHash: common.Hash{},
			},
			Txs: make([]types.Tx, 0),
		}
	}
}
