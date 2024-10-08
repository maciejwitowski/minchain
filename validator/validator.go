package validator

import (
	"fmt"
	"github.com/pkg/errors"
	"log"
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
	log.Println("Validating block", block.BlockHash().Hex())
	blockHash := block.BlockHash()
	foundBlock, err := v.db.GetBlockByHash(blockHash)
	if err != nil && !errors.Is(err, database.ErrorBlockNotFound) {
		return err
	}

	if foundBlock != nil {
		return errors.Wrap(ErrorKnownBlock, fmt.Sprintf("Block hash %s", blockHash.Hex()))
	}

	_, err = v.db.GetBlockByHash(block.Header.ParentHash)
	if errors.Is(err, database.ErrorBlockNotFound) {
		return ErrorUnknownParent
	}

	if err != nil {
		return err
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
