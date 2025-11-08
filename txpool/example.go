package txpool

import (
	"fmt"

	"github.com/lengzhao/govm/core"
	"github.com/lengzhao/govm/storage"
	"github.com/lengzhao/govm/types"
)

// ExampleUsage 交易池使用示例
func ExampleUsage(core core.Core, store storage.Storage) {
	// 创建交易池实例
	txPool := NewTxPool(core, store)

	// 启动交易池
	if err := txPool.Start(); err != nil {
		fmt.Printf("Failed to start txpool: %v\n", err)
		return
	}
	defer txPool.Stop()

	// 创建示例交易
	tx := &types.TransactionWithSign{
		Transaction: types.Transaction{
			ShardID: types.DefaultShardID,
			From:    types.Address{1, 2, 3}, // 示例地址
			To:      types.Address{4, 5, 6}, // 示例地址
			Amount:  100,
			Nonce:   1,
			Data:    []byte("example transaction"),
		},
		Signature: []byte("example signature"),
	}

	// 添加交易到交易池
	if err := txPool.AddTransaction(tx); err != nil {
		fmt.Printf("Failed to add transaction: %v\n", err)
		return
	}

	fmt.Printf("Transaction added to pool, pool size: %d\n", txPool.GetTransactionCount())

	// 选择交易用于区块生成
	selectedTxs, err := txPool.SelectTransactions(10)
	if err != nil {
		fmt.Printf("Failed to select transactions: %v\n", err)
		return
	}

	fmt.Printf("Selected %d transactions for block generation\n", len(selectedTxs))

	// 获取交易
	txHash := calculateExampleTxHash(tx)
	retrievedTx, err := txPool.GetTransaction(txHash)
	if err != nil {
		fmt.Printf("Failed to get transaction: %v\n", err)
		return
	}

	fmt.Printf("Retrieved transaction: %+v\n", retrievedTx)

	// 从交易池中移除交易
	if err := txPool.RemoveTransaction(txHash); err != nil {
		fmt.Printf("Failed to remove transaction: %v\n", err)
		return
	}

	fmt.Printf("Transaction removed from pool, pool size: %d\n", txPool.GetTransactionCount())
}

// calculateExampleTxHash 计算示例交易哈希
func calculateExampleTxHash(tx *types.TransactionWithSign) types.Hash {
	// 简化实现，实际应该使用crypto模块计算哈希
	var hash types.Hash
	copy(hash[:], []byte("example_hash"))
	return hash
}
