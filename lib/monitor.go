package lib

import (
	"context"
	"fmt"
	"log"
	"minchain/core"
	"time"
)

// Monitor checks every interval period whether mempool size changed and if so, it prints it.
func Monitor(ctx context.Context, mpool core.Mempool, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	var previousMpoolSize = 0
	for {
		select {
		case <-ticker.C:
			pendingTransactions := mpool.ListPendingTransactions()
			if len(pendingTransactions) != previousMpoolSize {
				fmt.Printf("Pending transactions (%d)\n", len(pendingTransactions))
				for _, tx := range mpool.ListPendingTransactions() {
					fmt.Println(tx.PrettyPrint())
				}
				if len(pendingTransactions) > 0 {
					fmt.Println("--------")
				}
				previousMpoolSize = len(pendingTransactions)
			}
		case <-ctx.Done():
			log.Println("parent context closed")
			return
		}
	}
}
