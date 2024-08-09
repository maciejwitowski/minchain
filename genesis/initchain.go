package genesis

import (
	"errors"
	"log"
	"minchain/core"
	"minchain/database"
)

func InitializeGenesisState(db database.Database) error {
	blockchainHashes, err := core.PrintBlockHashes(db)
	log.Println("InitializeGenesisState. Current blockchain:", blockchainHashes)

	head, err := db.GetHead()
	if err == nil {
		log.Println("Head exists, no need to initialise genesis. ", head.Hex())
		return nil
	}

	if err != nil && !errors.Is(err, database.ErrorHeadBlockNotSet) {
		return err
	}

	log.Println("Initializing genesis")

	err = db.PutBlock(&core.GenesisBlock)
	if err != nil {
		return err
	}
	err = db.SetHead(core.GenesisBlock.BlockHash())
	if err != nil {
		return err
	}

	return nil
}
