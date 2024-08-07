package core

import (
	"fmt"
	crypto "github.com/ethereum/go-ethereum/crypto"
	"minchain/core/types"
	"strings"
	"sync"
)

type Mempool struct {
	lock       sync.Mutex
	pendingTxs []types.Tx
}

func InitMempool() *Mempool {
	return &Mempool{
		lock:       sync.Mutex{},
		pendingTxs: make([]types.Tx, 0),
	}
}

func (m *Mempool) HandleTransaction(tx types.Tx) {
	m.lock.Lock()
	defer m.lock.Unlock()

	if IsValid(tx) {
		m.pendingTxs = append(m.pendingTxs, tx)
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

	fmt.Printf("Pending transactions (%d)\n", len(m.pendingTxs))
	for _, tx := range m.pendingTxs {
		fmt.Println(tx.PrettyPrint())
	}
	if len(m.pendingTxs) > 0 {
		fmt.Println("--------")
	}
}

func (m *Mempool) Size() int {
	m.lock.Lock()
	defer m.lock.Unlock()

	return len(m.pendingTxs)
}

func (m *Mempool) BuildBlockFromTransactions(builder func(txs []types.Tx) (*types.Block, error)) (*types.Block, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	if len(m.pendingTxs) == 0 {
		return nil, nil
	}

	txCopy := make([]types.Tx, len(m.pendingTxs))
	copy(txCopy, m.pendingTxs)

	block, err := builder(m.pendingTxs)
	if err != nil {
		return nil, err
	}

	// Block created successfully. We assume all pending have been handled and can be cleared
	m.pendingTxs = m.pendingTxs[:0]

	return block, err
}
