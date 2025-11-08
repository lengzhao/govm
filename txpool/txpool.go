package txpool

import (
	"github.com/lengzhao/govm/types"
)

// TxPool 交易池接口
type TxPool interface {
	// AddTransaction 添加交易到交易池
	AddTransaction(tx *types.TransactionWithSign) error

	// GetTransaction 获取交易
	GetTransaction(hash types.Hash) (*types.TransactionWithSign, error)

	// RemoveTransaction 从交易池中移除交易
	RemoveTransaction(hash types.Hash) error

	// SelectTransactions 选择一批交易用于区块生成
	SelectTransactions(maxCount int) ([]*types.TransactionWithSign, error)

	// ValidateTransaction 验证交易
	ValidateTransaction(tx *types.TransactionWithSign) error

	// GetTransactionCount 获取交易池中的交易数量
	GetTransactionCount() int

	// Start 启动交易池
	Start() error

	// Stop 停止交易池
	Stop() error
}
