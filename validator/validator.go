package validator

import (
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

	foundParent := v.db.GetBlockByHash(block.Header.ParentHash)
	if foundParent == nil {
		return ErrorUnknownParent
	}

	hash, err := types.CombinedHash(block.Transactions)
	if err != nil {
		return err
	}

	if hash != block.Header.TransactionHash {
		return IncorrectTxHash
	}

	return nil
}
