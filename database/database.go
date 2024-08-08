package database

import (
	"github.com/ethereum/go-ethereum/common"
	"minchain/core/types"
)

type Database interface {
	PutBlock(block *types.Block) error
	GetBlockByHash(hash common.Hash) (*types.Block, error)
	Close() error
}

type MemoryDatabase struct {
	blocks map[common.Hash]*types.Block
}

func NewMemoryDatabase() Database {
	return &MemoryDatabase{blocks: make(map[common.Hash]*types.Block)}
}

func (db *MemoryDatabase) PutBlock(block *types.Block) error {
	db.blocks[block.BlockHash()] = block
	return nil
}

func (db *MemoryDatabase) GetBlockByHash(hash common.Hash) (*types.Block, error) {
	block := db.blocks[hash]
	return block, nil
}

func (db *MemoryDatabase) Close() error {
	return nil // no op
}
