package core

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	crypto "github.com/ethereum/go-ethereum/crypto"
	"minchain/core/types"
	"strings"
	"sync"
)

type Mempool struct {
	lock                sync.Mutex
	pendingTransactions map[common.Hash]types.Tx
}

func InitMempool() *Mempool {
	return &Mempool{
		lock:                sync.Mutex{},
		pendingTransactions: make(map[common.Hash]types.Tx),
	}
}

func (m *Mempool) HandleTransaction(tx types.Tx) {
	m.lock.Lock()
	defer m.lock.Unlock()

	txHash, err := tx.Hash()
	if err != nil {
		fmt.Println("error getting tx hash")
		return
	}

	if IsValid(tx) {
		m.pendingTransactions[txHash] = tx
	}
}

func IsValid(tx types.Tx) bool {
	if len(strings.TrimSpace(tx.Data)) == 0 {
		return false
	}

	if len(tx.Signature) != 65 {
		fmt.Println("Invalid signature length")
		return false
	}

	digest := crypto.Keccak256([]byte(tx.Data))
	publicKey, err := crypto.Ecrecover(digest, tx.Signature)
	if err != nil {
		fmt.Println("publicKey error")
		return false
	}

	// VerifySignature expects 64-bytes long sig, without the last recovery ID byte
	sig := tx.Signature[:64]
	if !crypto.VerifySignature(publicKey, digest, sig) {
		return false
	}

	return true
}

func (m *Mempool) DumpTx() {
	m.lock.Lock()
	defer m.lock.Unlock()

	fmt.Printf("Pending transactions (%d)\n", len(m.pendingTransactions))
	for _, tx := range m.pendingTransactions {
		fmt.Println(tx.PrettyPrint())
	}
	if len(m.pendingTransactions) > 0 {
		fmt.Println("--------")
	}
}

func (m *Mempool) Size() int {
	m.lock.Lock()
	defer m.lock.Unlock()

	return len(m.pendingTransactions)
}

func (m *Mempool) BuildBlockFromTransactions(blockProducer *BlockProducer) (*types.Block, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	if len(m.pendingTransactions) == 0 {
		return nil, nil
	}

	selectTransactions := make([]types.Tx, 0)
	for _, tx := range m.pendingTransactions {
		// TODO more advanced selection logic
		selectTransactions = append(selectTransactions, tx)
	}

	// TODO Refactor to lower coupling between mempool and block producer
	block, err := blockProducer.builder(selectTransactions)
	if err != nil {
		return nil, err
	}

	// Block created successfully. We assume all pending have been handled and can be cleared
	clear(m.pendingTransactions)

	return block, err
}

func (m *Mempool) PruneTransactions(transactions []types.Tx) {
	for _, t := range transactions {
		hash, err := t.Hash()
		if err != nil {
			continue
		}
		delete(m.pendingTransactions, hash)
	}
}
