package types

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
)

type Block struct {
	Header BlockHeader `json:"blockHeader"`
	Txs    []Tx        `json:"transactions"`
}

type BlockHeader struct {
	ParentHash      common.Hash `json:"parentHash"`
	TransactionHash common.Hash `json:"transactionHash"`
}

func (blk *Block) ToJson() ([]byte, error) {
	return json.Marshal(blk)
}

func BlockFromJson(data []byte) (*Block, error) {
	var blk Block
	err := json.Unmarshal(data, &blk)
	if err != nil {
		return nil, err
	}
	return &blk, nil
}

func (blk *Block) PrettyPrint() string {
	jsonData, err := json.MarshalIndent(blk, "", "  ")
	if err != nil {
		return fmt.Sprintf("Error pretty printing transaction: %v", err)
	}
	return string(jsonData)
}
