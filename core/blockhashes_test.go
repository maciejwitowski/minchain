package core

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"minchain/core/types"
	"minchain/database"
	"testing"
)

func TestPrintBlockHashes(t *testing.T) {
	db := database.NewMemoryDatabase()
	hashes, _ := PrintBlockHashes(db)
	require.Equal(t, NoHeadMessage, hashes)

	_ = db.PutBlock(&GenesisBlock)
	_ = db.SetHead(GenesisBlock.BlockHash())
	hashes, _ = PrintBlockHashes(db)
	require.Equal(t, GenesisBlock.BlockHash().Hex(), hashes)

	nextBlock := types.Block{
		Header: types.BlockHeader{
			ParentHash:      GenesisBlock.BlockHash(),
			TransactionHash: common.Hash{},
		},
		Transactions: make([]types.Tx, 0),
	}

	_ = db.PutBlock(&nextBlock)
	_ = db.SetHead(nextBlock.BlockHash())

	expected := nextBlock.BlockHash().Hex() + " -> " + GenesisBlock.BlockHash().Hex()
	hashes, _ = PrintBlockHashes(db)
	require.Equal(t, expected, hashes)
}
