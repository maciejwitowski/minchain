package core

import (
	"minchain/core/types"
	"minchain/database"
	"sync"
)

type Chainstore interface {
	SetHead(block *types.Block)
	GetHead() *types.Block
}

type MemoryChainstore struct {
	lock sync.Mutex

	db   database.Database
	head *types.Block
}

func NewChainstore(db database.Database) Chainstore {
	return &MemoryChainstore{
		db: db,
	}
}

func (c *MemoryChainstore) SetHead(block *types.Block) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.head = block
}

func (c *MemoryChainstore) GetHead() *types.Block {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.head == nil {
		// Assumes genesis block has been initialised and exists in DB
		genesisBlock, err := c.db.GetBlockByHash(GenesisBlockHash)
		if err != nil {
			return nil
		}
		c.head = genesisBlock
	}

	return c.head
}
