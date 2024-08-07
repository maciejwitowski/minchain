package types

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
)

type Block struct {
	Header       BlockHeader `json:"blockHeader"`
	Transactions []Tx        `json:"transactions"`
}

type BlockHeader struct {
	ParentHash      common.Hash `json:"parentHash"`
	TransactionHash common.Hash `json:"transactionHash"`
	Number          *big.Int    `json:"number"`
}

func (block *Block) BlockHash() common.Hash {
	headerBytes, err := json.Marshal(block.Header)
	if err != nil {
		return [32]byte{}
	}

	return common.BytesToHash(crypto.Keccak256(headerBytes))
}

func (block *Block) ToJson() ([]byte, error) {
	return json.Marshal(block)
}

func BlockFromJson(data []byte) (*Block, error) {
	var blk Block
	err := json.Unmarshal(data, &blk)
	if err != nil {
		return nil, err
	}
	return &blk, nil
}

func (block *Block) PrettyPrint() string {
	jsonData, err := json.MarshalIndent(block, "", "  ")
	if err != nil {
		return fmt.Sprintf("Error pretty printing transaction: %v", err)
	}
	return string(jsonData)
}
