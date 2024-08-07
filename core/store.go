package core

import (
	"minchain/core/types"
	"minchain/database"
	"sync"
)

type Chainstore struct {
	lock sync.Mutex

	db   database.Database
	head *types.Block
}

func NewChainstore(db database.Database) *Chainstore {
	return &Chainstore{
		db: db,
	}
}

func (c *Chainstore) SetHead(block *types.Block) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.head = block
}

func (c *Chainstore) GetHead() *types.Block {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.head == nil {
		// Assumes genesis block has been initialised and exists in DB
		c.db.GetBlockByHash(GenesisBlockHash)
	}

	return c.head
}
