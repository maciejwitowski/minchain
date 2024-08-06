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
	digest := crypto.Keccak256([]byte(tx.Data))
	publicKey, err := crypto.Ecrecover(digest, tx.Signature)
	if err != nil {
		fmt.Println("publicKey error")
		return false
	}

	if !crypto.VerifySignature(publicKey, digest, tx.Signature) {
		return false
	}

	return len(strings.TrimSpace(tx.Data)) > 0
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
