package core

import (
	"minchain/database"
	"strings"
)

var NoHeadMessage = "[No head]"

func PrintBlockHashes(database database.Database) (string, error) {
	blockHash, err := database.GetHead()
	if err != nil {
		return NoHeadMessage, nil
	}

	hashes := make([]string, 0)

	for {
		hashes = append(hashes, blockHash.Hex())
		if blockHash == GenesisBlock.BlockHash() {
			break
		}
		currentBlock, err := database.GetBlockByHash(blockHash)
		blockHash = currentBlock.Header.ParentHash
		if err != nil {
			return "", err
		}
	}

	return strings.Join(hashes, " -> "), nil
}
