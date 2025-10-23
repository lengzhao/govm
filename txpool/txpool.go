package txpool

import (
	"github.com/lengzhao/govm/types"
)

// TxPool 交易池接口
type TxPool interface {
	// Start 启动交易池服务
	Start() error

	// Stop 停止交易池服务
	Stop() error

	// AddTransaction 添加交易到交易池
	AddTransaction(tx *types.Transaction) error

	// GetTransaction 获取交易
	GetTransaction(hash types.Hash) (*types.Transaction, error)

	// RemoveTransaction 从交易池移除交易
	RemoveTransaction(hash types.Hash) error

	// GetTransactions 获取所有交易
	GetTransactions() ([]*types.Transaction, error)

	// GetPendingTransactions 获取待处理交易
	GetPendingTransactions() ([]*types.Transaction, error)

	// GetTransactionCount 获取交易数量
	GetTransactionCount() int

	// HasTransaction 检查交易是否存在
	HasTransaction(hash types.Hash) bool

	// ValidateTransaction 验证交易有效性
	ValidateTransaction(tx *types.Transaction) error
}
