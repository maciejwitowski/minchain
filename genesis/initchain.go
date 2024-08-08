package genesis

import (
	"errors"
	"fmt"
	"github.com/dgraph-io/badger/v4"
	"github.com/ethereum/go-ethereum/common"
	"log"
	"minchain/core"
	"minchain/core/types"
	"minchain/database"
)

func InitializeGenesisState(db database.Database, store core.Chainstore) error {
	_, err := db.GetBlockByHash(core.GenesisBlockHash)
	if err != nil {
		if errors.Is(err, badger.ErrKeyNotFound) {
			log.Println("Genesis already initialised. Skipping.")
			return nil
		} else {
			return err
		}
	}

	log.Println("Initializing genesis")

	genesisBlock := types.Block{
		Header: types.BlockHeader{
			ParentHash:      common.Hash{},
			TransactionHash: common.Hash{},
		},
		Transactions: make([]types.Tx, 0),
	}

	blockHash := genesisBlock.BlockHash()
	if blockHash != core.GenesisBlockHash {
		errText := fmt.Sprintf("Incorrect genesis hash %s", blockHash)
		return errors.New(errText)
	}

	err = db.PutBlock(&genesisBlock)
	if err != nil {
		return err
	}
	store.SetHead(&genesisBlock)

	return nil
}
