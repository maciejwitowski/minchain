package chain

import (
	"fmt"
	crypto "github.com/ethereum/go-ethereum/crypto"
	"strings"
	"sync"
)

type Mempool struct {
	lock       sync.Mutex
	pendingTxs []Tx
}

func InitMempool() *Mempool {
	return &Mempool{
		lock:       sync.Mutex{},
		pendingTxs: make([]Tx, 0),
	}
}

func (m *Mempool) HandleTransaction(tx Tx) {
	m.lock.Lock()
	defer m.lock.Unlock()

	if IsValid(tx) {
		m.pendingTxs = append(m.pendingTxs, tx)
	}
}

func IsValid(tx Tx) bool {
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

	fmt.Printf("Current transactions (%d)\n", len(m.pendingTxs))
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

func (m *Mempool) GetPendingTransactions() []Tx {
	m.lock.Lock()
	defer m.lock.Unlock()

	txCopy := make([]Tx, len(m.pendingTxs))
	// TODO Remove from pending
	copy(txCopy, m.pendingTxs)
	return txCopy
}
