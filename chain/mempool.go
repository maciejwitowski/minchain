package chain

import (
	"fmt"
	crypto "github.com/ethereum/go-ethereum/crypto"
	"strings"
	"sync"
)

type Mempool struct {
	lock  sync.Mutex
	mpool []Tx
}

func InitMempool() *Mempool {
	return &Mempool{
		lock:  sync.Mutex{},
		mpool: make([]Tx, 0),
	}
}

func (m *Mempool) HandleTransaction(tx Tx) {
	m.lock.Lock()
	defer m.lock.Unlock()

	if IsValid(tx) {
		m.mpool = append(m.mpool, tx)
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

	fmt.Printf("Current transactions (%d)\n", len(m.mpool))
	for _, tx := range m.mpool {
		fmt.Println(tx.PrettyPrint())
	}
	if len(m.mpool) > 0 {
		fmt.Println("--------")
	}
}

func (m *Mempool) Size() int {
	return len(m.mpool)
}
