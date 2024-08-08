package core

import (
	"minchain/core/types"
	"minchain/database"
	"sync"
)

type Chainhead interface {
	SetHead(block *types.Block)
	GetHead() *types.Block
}

type MemoryChainhead struct {
	lock sync.Mutex

	db   database.Database
	head *types.Block
}

func NewChainhead(db database.Database) Chainhead {
	return &MemoryChainhead{
		db: db,
	}
}

func (c *MemoryChainhead) SetHead(block *types.Block) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.head = block
}

func (c *MemoryChainhead) GetHead() *types.Block {
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
