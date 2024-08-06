package lib

import (
	"context"
	"log"
	"minchain/chain"
	"time"
)

// Monitor checks every interval period whether mempool size changed and if so, it prints it.
func Monitor(ctx context.Context, mpool *chain.Mempool, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	var previousMpoolSize = 0
	for {
		select {
		case <-ticker.C:
			if mpool.Size() != previousMpoolSize {
				mpool.DumpTx()
				previousMpoolSize = mpool.Size()
			}
		case <-ctx.Done():
			log.Println("parent context closed")
			return
		}
	}
}
