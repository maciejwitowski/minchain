package chain

import (
	"fmt"
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
