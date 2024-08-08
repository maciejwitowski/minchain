package database

import (
	"github.com/dgraph-io/badger/v4"
	"github.com/ethereum/go-ethereum/common"
	"minchain/core/types"
)

var badgerFilePath = "/tmp/badger"

type DiskDatabase struct {
	inner *badger.DB
}

func NewDiskDatabase() (Database, error) {
	open, err := badger.Open(badger.DefaultOptions(badgerFilePath))
	if err != nil {
		return nil, err
	}
	return &DiskDatabase{inner: open}, nil
}

func (db *DiskDatabase) PutBlock(block *types.Block) error {
	blockJson, err := block.ToJson()
	if err != nil {
		return err
	}

	return db.inner.Update(func(txn *badger.Txn) error {
		err := txn.Set(block.BlockHash().Bytes(), blockJson)
		if err != nil {
			return err
		}
		return nil
	})
}

func (db *DiskDatabase) GetBlockByHash(hash common.Hash) (*types.Block, error) {
	var bytes []byte

	err := db.inner.View(func(txn *badger.Txn) error {
		item, err := txn.Get(hash.Bytes())
		if err != nil {
			return err
		}
		err = item.Value(func(val []byte) error {
			bytes = val
			return nil
		})
		return err
	})

	if err != nil {
		return nil, err
	}

	return types.BlockFromJson(bytes)
}

func (db *DiskDatabase) Close() error {
	return db.inner.Close()
}
