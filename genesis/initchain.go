package genesis

import (
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"minchain/core"
	"minchain/core/types"
	"minchain/database"
)

func InitializeGenesisState(db database.Database, store core.Chainstore) error {
	genesisBlock := types.Block{
		Header: types.BlockHeader{
			ParentHash:      common.Hash{},
			TransactionHash: common.Hash{},
		},
		Transactions: make([]types.Tx, 0),
	}

	blockHash := genesisBlock.BlockHash()
	if blockHash != core.GenesisBlockHash {
		errText := log.Sprintf("Incorrect genesis hash %s", blockHash)
		return errors.New(errText)
	}

	db.PutBlock(&genesisBlock)
	store.SetHead(&genesisBlock)

	return nil
}
