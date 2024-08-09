package database

import (
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"minchain/core/types"
)

var ErrorHeadBlockNotSet = errors.New("head block not set")
var ErrorBlockNotFound = errors.New("head block not set")

type Database interface {
	SetHead(blockHash common.Hash) error
	GetHead() (common.Hash, error)
	PutBlock(block *types.Block) error
	GetBlockByHash(hash common.Hash) (*types.Block, error)
	Close() error
}

type MemoryDatabase struct {
	blocks    map[common.Hash]*types.Block
	headBlock common.Hash
}

func NewMemoryDatabase() Database {
	return &MemoryDatabase{blocks: make(map[common.Hash]*types.Block)}
}

func (db *MemoryDatabase) PutBlock(block *types.Block) error {
	db.blocks[block.BlockHash()] = block
	return nil
}

func (db *MemoryDatabase) SetHead(blockHash common.Hash) error {
	db.headBlock = blockHash
	return nil
}

func (db *MemoryDatabase) GetHead() (common.Hash, error) {
	var zeroHash common.Hash
	if db.headBlock == zeroHash {
		return zeroHash, ErrorHeadBlockNotSet
	}

	return db.headBlock, nil
}

func (db *MemoryDatabase) GetBlockByHash(hash common.Hash) (*types.Block, error) {
	block, exists := db.blocks[hash]
	if !exists {
		return nil, ErrorBlockNotFound
	}
	return block, nil
}

func (db *MemoryDatabase) Close() error {
	return nil // no op
}
