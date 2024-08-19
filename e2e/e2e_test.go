package test

import (
	"context"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	"minchain/app"
	"minchain/core"
	"minchain/core/types"
	"minchain/database"
	"minchain/lib"
	"minchain/validator"
	"sync"
	"testing"
	"time"
)

func TestE2E(t *testing.T) {
	var ctx = context.Background()
	var db = database.NewMemoryDatabase()
	var mempool = core.NewMempool()
	var pk, _ = crypto.HexToECDSA("ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80")

	var testConfig = lib.Config{
		PrivateKey:      pk,
		IsBlockProducer: true,
		BlockTime:       1 * time.Millisecond,
	}

	var publisher = TestPublisher{
		publishedBlocks:       make([]*types.Block, 0),
		publishedTransactions: make([]*types.Tx, 0),
	}

	var consumer = TestConsumer{
		make(chan *types.Block),
		make(chan *types.Tx),
	}

	var input = TestTransactionsInput{input: make(chan string)}

	var testApp = app.NewApp(
		mempool,
		db,
		validator.NewBlockValidator(db),
		core.NewWallet(testConfig.PrivateKey),
		testConfig,
		&publisher,
		&consumer,
		[]lib.TransactionsInput{&input},
	)

	testApp.Start(ctx)

	// Simulate new transaction from a user
	input.NewUserInput("hello world")
	waitForPropagation()

	require.Equal(t, 1, len(publisher.publishedTransactions))
	require.Equal(t, "hello world", publisher.publishedTransactions[0].Data)

	// Simulate the transaction has been received from p2p
	publishedTx := publisher.publishedTransactions[0]
	consumer.TxChannel <- publishedTx

	waitForPropagation()

	require.Equal(t, 1, len(publisher.publishedBlocks))
	require.Equal(t, "hello world", publisher.publishedBlocks[0].Transactions[0].Data)

	// Simulate the block has been received from p2p
	publishedBlock := publisher.publishedBlocks[0]
	consumer.BlocksChannel <- publishedBlock

	waitForPropagation()

	blockStoredInDb, _ := db.GetBlockByHash(publishedBlock.BlockHash())
	headBlock, _ := db.GetHead()
	require.Equal(t, publishedBlock.BlockHash(), blockStoredInDb.BlockHash())
	require.Equal(t, publishedBlock.BlockHash(), headBlock)
	require.Equal(t, publishedBlock.Header.Height, int64(1))
	require.Equal(t, 0, len(mempool.ListPendingTransactions()))
}

func waitForPropagation() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(10 * time.Millisecond)
	}()
	wg.Wait()
}

type TestPublisher struct {
	publishedBlocks       []*types.Block
	publishedTransactions []*types.Tx
}

func (p *TestPublisher) PublishBlock(ctx context.Context, block *types.Block) error {
	p.publishedBlocks = append(p.publishedBlocks, block)
	return nil
}

func (p *TestPublisher) PublishTransaction(ctx context.Context, transaction *types.Tx) error {
	p.publishedTransactions = append(p.publishedTransactions, transaction)
	return nil
}

type TestConsumer struct {
	BlocksChannel chan *types.Block
	TxChannel     chan *types.Tx
}

func (c *TestConsumer) ConsumeTransaction(ctx context.Context) (*types.Tx, error) {
	return <-c.TxChannel, nil
}

func (c *TestConsumer) ConsumeBlock(ctx context.Context) (*types.Block, error) {
	return <-c.BlocksChannel, nil
}

type TestTransactionsInput struct {
	input chan string
}

func (ui *TestTransactionsInput) InputChannel(ctx context.Context) <-chan string {
	return ui.input
}

func (ui *TestTransactionsInput) NewUserInput(message string) {
	ui.input <- message
}
