package core

import (
	"github.com/ethereum/go-ethereum/common"
	"minchain/core/types"
)

// GenesisBlock Hash: 0x565d980ee06afa9edc0d2b4ed01ea6cd755a250a8b7b8508dd63bf1cd5415efd
var GenesisBlock = types.Block{
	Header: types.BlockHeader{
		ParentHash:      common.Hash{},
		TransactionHash: common.Hash{},
	},
	Transactions: make([]types.Tx, 0),
}
