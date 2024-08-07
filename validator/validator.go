package validator

import (
	"minchain/core"
	"minchain/core/types"
	"minchain/database"
)

type Validator interface {
	Validate(block *types.Block) error
}

type BlockValidator struct {
	db database.Database
}

func NewBlockValidator(db database.Database) *BlockValidator {
	return &BlockValidator{
		db: db,
	}
}

func (v *BlockValidator) Validate(block *types.Block) error {
	blockHash := block.BlockHash()
	foundBlock := v.db.GetBlockByHash(blockHash)
	if foundBlock != nil {
		return ErrorKnownBlock
	}

	// TODO requires genesis setup (adding genesis block to db so the parent case doesn't have to be handled separately
	foundParent := v.db.GetBlockByHash(block.Header.ParentHash)
	if foundParent == nil {
		return ErrorUnknownParent
	}

	hash, err := core.Hash(block.Transactions)
	if err != nil {
		return err
	}

	if hash != block.Header.TransactionHash {
		return IncorrectTxHash
	}

	return nil
}
