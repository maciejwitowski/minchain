package chain

import (
	"encoding/json"
	"fmt"
)

type Tx struct {
	From      string `json:"from"`
	Data      string `json:"data"`
	Signature []byte `json:"sig"`
}

// ToJSON serializes the Transaction to JSON
func (t *Tx) ToJSON() ([]byte, error) {
	return json.Marshal(t)
}

// FromJSON deserializes JSON data into a Transaction
func FromJSON(data []byte) (*Tx, error) {
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
