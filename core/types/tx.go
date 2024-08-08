package types

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type Tx struct {
	From      string `json:"from"`
	Data      string `json:"data"`
	Signature []byte `json:"sig"`
}

// ToJson serializes the Transaction to JSON
func (t *Tx) ToJson() ([]byte, error) {
	return json.Marshal(t)
}

// TransactionFromJSON deserializes JSON data into a Transaction
func TransactionFromJSON(data []byte) (*Tx, error) {
	var t Tx
	err := json.Unmarshal(data, &t)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (t *Tx) PrettyPrint() string {
	jsonData, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return fmt.Sprintf("Error pretty printing transaction: %v", err)
	}
	return string(jsonData)
}

func (t *Tx) HashBytes() ([]byte, error) {
	serialized, err := t.ToJson()
	if err != nil {
		return []byte{}, err
	}
	return crypto.Keccak256(serialized), nil
}

func (t *Tx) Hash() (common.Hash, error) {
	bytes, err := t.HashBytes()
	if err != nil {
		return common.Hash{}, err
	}
	return common.BytesToHash(bytes), nil
}
