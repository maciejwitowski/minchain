package core

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"log"
	"minchain/core/types"
	"strings"
	"sync"
)

type Mempool interface {
	ValidateAndStorePending(transations *types.Tx)
	ListPendingTransactions() []types.Tx
	PruneTransactions(transactions []types.Tx)
}

type MemoryMempool struct {
	lock                sync.Mutex
	pendingTransactions map[common.Hash]*types.Tx
}

func NewMempool() Mempool {
	return &MemoryMempool{
		lock:                sync.Mutex{},
		pendingTransactions: make(map[common.Hash]*types.Tx),
	}
}

func (m *MemoryMempool) ValidateAndStorePending(tx *types.Tx) {
	m.lock.Lock()
	defer m.lock.Unlock()

	txHash, err := tx.Hash()
	if err != nil {
		log.Println("error getting tx hash")
		return
	}

	if IsValid(tx) {
		m.pendingTransactions[txHash] = tx
	}
}

func IsValid(tx *types.Tx) bool {
	if len(strings.TrimSpace(tx.Data)) == 0 {
		return false
	}

	if len(tx.Signature) != 65 {
		log.Println("Invalid signature length")
		return false
	}

	digest := crypto.Keccak256([]byte(tx.Data))
	publicKey, err := crypto.Ecrecover(digest, tx.Signature)
	if err != nil {
		log.Println("publicKey error")
		return false
	}

	// VerifySignature expects 64-bytes long sig, without the last recovery ID byte
	sig := tx.Signature[:64]
	if !crypto.VerifySignature(publicKey, digest, sig) {
		return false
	}

	return true
}

func (m *MemoryMempool) ListPendingTransactions() []types.Tx {
	m.lock.Lock()
	defer m.lock.Unlock()

	transactions := make([]types.Tx, 0)
	for _, tx := range m.pendingTransactions {
		transactions = append(transactions, *tx)
	}
	return transactions
}

func (m *MemoryMempool) PruneTransactions(transactions []types.Tx) {
	m.lock.Lock()
	defer m.lock.Unlock()

	for _, t := range transactions {
		hash, err := t.Hash()
		if err != nil {
			continue
		}
		delete(m.pendingTransactions, hash)
	}
}
