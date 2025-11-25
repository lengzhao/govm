package txpool

import (
	"fmt"
	"sync"

	lzbinary "github.com/lengzhao/binary"
	"github.com/lengzhao/govm/core"
	"github.com/lengzhao/govm/crypto"
	"github.com/lengzhao/govm/storage"
	"github.com/lengzhao/govm/types"
)

// DefaultTxPool 默认交易池实现
type DefaultTxPool struct {
	core    core.Core
	storage storage.Storage

	// 内存中的交易池
	pendingTxs map[types.Hash]*types.TransactionWithSign
	mutex      sync.RWMutex

	// 运行状态
	running bool
}

// NewTxPool 创建新的交易池实例
func NewTxPool(core core.Core, store storage.Storage) TxPool {
	return &DefaultTxPool{
		core:       core,
		storage:    store,
		pendingTxs: make(map[types.Hash]*types.TransactionWithSign),
		running:    false,
	}
}

// Start 启动交易池
func (tp *DefaultTxPool) Start() error {
	tp.mutex.Lock()
	defer tp.mutex.Unlock()

	if tp.running {
		return fmt.Errorf("txpool is already running")
	}

	tp.running = true
	return nil
}

// Stop 停止交易池
func (tp *DefaultTxPool) Stop() error {
	tp.mutex.Lock()
	defer tp.mutex.Unlock()

	if !tp.running {
		return fmt.Errorf("txpool is not running")
	}

	tp.running = false
	return nil
}

// AddTransaction 添加交易到交易池
func (tp *DefaultTxPool) AddTransaction(tx *types.TransactionWithSign) error {
	tp.mutex.Lock()
	defer tp.mutex.Unlock()

	if !tp.running {
		return fmt.Errorf("txpool is not running")
	}

	// 验证交易
	if err := tp.ValidateTransaction(tx); err != nil {
		return fmt.Errorf("transaction validation failed: %w", err)
	}

	// 计算交易哈希
	txHash := tp.calculateTransactionHash(tx)

	// 检查交易是否已存在
	if _, exists := tp.pendingTxs[txHash]; exists {
		return fmt.Errorf("transaction already exists in pool")
	}

	// 添加到内存交易池
	tp.pendingTxs[txHash] = tx

	// 持久化到存储
	if err := tp.storeTransaction(tx, txHash); err != nil {
		// 如果存储失败，从内存中移除
		delete(tp.pendingTxs, txHash)
		return fmt.Errorf("failed to store transaction: %w", err)
	}

	return nil
}

// GetTransaction 获取交易
func (tp *DefaultTxPool) GetTransaction(hash types.Hash) (*types.TransactionWithSign, error) {
	tp.mutex.RLock()
	defer tp.mutex.RUnlock()

	if !tp.running {
		return nil, fmt.Errorf("txpool is not running")
	}

	// 首先从内存中查找
	if tx, exists := tp.pendingTxs[hash]; exists {
		return tx, nil
	}

	// 如果内存中没有，从存储中查找
	return tp.loadTransactionFromStorage(hash)
}

// RemoveTransaction 从交易池中移除交易
func (tp *DefaultTxPool) RemoveTransaction(hash types.Hash) error {
	tp.mutex.Lock()
	defer tp.mutex.Unlock()

	if !tp.running {
		return fmt.Errorf("txpool is not running")
	}

	// 从内存中移除
	delete(tp.pendingTxs, hash)

	// 从存储中移除
	return tp.removeTransactionFromStorage(hash)
}

// SelectTransactions 选择一批交易用于区块生成
func (tp *DefaultTxPool) SelectTransactions(maxCount int) ([]*types.TransactionWithSign, error) {
	tp.mutex.RLock()
	defer tp.mutex.RUnlock()

	if !tp.running {
		return nil, fmt.Errorf("txpool is not running")
	}

	// 选择最多maxCount个交易
	selected := make([]*types.TransactionWithSign, 0, maxCount)
	count := 0

	for _, tx := range tp.pendingTxs {
		if count >= maxCount {
			break
		}
		selected = append(selected, tx)
		count++
	}

	return selected, nil
}

// ValidateTransaction 验证交易
func (tp *DefaultTxPool) ValidateTransaction(tx *types.TransactionWithSign) error {
	if !tp.running {
		return fmt.Errorf("txpool is not running")
	}

	// 使用核心模块的交易验证功能
	return tp.core.ValidateTransaction(tx)
}

// GetTransactionCount 获取交易池中的交易数量
func (tp *DefaultTxPool) GetTransactionCount() int {
	tp.mutex.RLock()
	defer tp.mutex.RUnlock()

	return len(tp.pendingTxs)
}

// calculateTransactionHash 计算交易哈希
func (tp *DefaultTxPool) calculateTransactionHash(tx *types.TransactionWithSign) types.Hash {
	// 序列化交易（排除签名字段）
	txCopy := *tx
	txCopy.Signature = nil

	data, err := lzbinary.Marshal(&txCopy.Transaction)
	if err != nil {
		return types.Hash{} // 返回空哈希
	}

	return crypto.Hash(data)
}

// storeTransaction 持久化交易到存储
func (tp *DefaultTxPool) storeTransaction(tx *types.TransactionWithSign, hash types.Hash) error {
	// 创建交易存储实例
	txStore, err := tp.storage.NewStorage(types.SNTx)
	if err != nil {
		return fmt.Errorf("failed to create transaction storage: %w", err)
	}

	// 序列化交易
	data, err := lzbinary.Marshal(tx)
	if err != nil {
		return fmt.Errorf("failed to marshal transaction: %w", err)
	}

	// 存储交易
	return txStore.Put(hash[:], data)
}

// loadTransactionFromStorage 从存储中加载交易
func (tp *DefaultTxPool) loadTransactionFromStorage(hash types.Hash) (*types.TransactionWithSign, error) {
	// 创建交易存储实例
	txStore, err := tp.storage.NewStorage(types.SNTx)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction storage: %w", err)
	}

	// 从存储中获取交易
	data, err := txStore.Get(hash[:])
	if err != nil {
		return nil, fmt.Errorf("transaction not found: %w", err)
	}

	// 反序列化交易
	var tx types.TransactionWithSign
	if err := lzbinary.Unmarshal(data, &tx); err != nil {
		return nil, fmt.Errorf("failed to unmarshal transaction: %w", err)
	}

	return &tx, nil
}

// removeTransactionFromStorage 从存储中移除交易
func (tp *DefaultTxPool) removeTransactionFromStorage(hash types.Hash) error {
	// 创建交易存储实例
	txStore, err := tp.storage.NewStorage(types.SNTx)
	if err != nil {
		return fmt.Errorf("failed to create transaction storage: %w", err)
	}

	// 从存储中删除交易
	return txStore.Delete(hash[:])
}
