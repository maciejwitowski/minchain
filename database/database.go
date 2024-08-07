package database

import (
	"github.com/ethereum/go-ethereum/common"
	"minchain/core/types"
)

type Database interface {
	PutBlock(block *types.Block)
	GetBlockByHash(hash common.Hash) *types.Block
}

type MemoryDatabase struct {
	blocks map[common.Hash]*types.Block
}

func NewMemoryDatabase() Database {
	return &MemoryDatabase{blocks: make(map[common.Hash]*types.Block)}
}

func (db *MemoryDatabase) PutBlock(block *types.Block) {
	db.blocks[block.BlockHash()] = block
}

func (db *MemoryDatabase) GetBlockByHash(hash common.Hash) *types.Block {
	block := db.blocks[hash]
	return block
}
