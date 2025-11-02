package core

import (
	"fmt"

	lzbinary "github.com/lengzhao/binary"
	"github.com/lengzhao/govm/crypto"
	"github.com/lengzhao/govm/storage"
	"github.com/lengzhao/govm/types"
)

// TxProcessor 交易处理器接口
type TxProcessor interface {
	// ValidateTransaction 验证交易的有效性
	ValidateTransaction(tx *types.TransactionWithSign) error

	// ApplyTransaction 应用交易到状态
	ApplyTransaction(tx *types.TransactionWithSign) error

	// GetTransactionByHash 根据哈希获取交易
	GetTransactionByHash(hash types.Hash) (*types.TransactionWithSign, error)

	// GetBalance 获取账户余额
	GetBalance(addr types.Address) (uint64, error)

	// GetNonce 获取账户nonce值
	GetNonce(addr types.Address) (uint64, error)
}

// DefaultTxProcessor 默认交易处理器实现
type DefaultTxProcessor struct {
	crypto  crypto.Crypto
	storage storage.Storage
}

// NewTxProcessor 创建新的交易处理器
func NewTxProcessor(store storage.Storage) TxProcessor {
	return &DefaultTxProcessor{
		crypto:  crypto.NewCrypto(),
		storage: store,
	}
}

// ValidateTransaction 验证交易的有效性
func (tp *DefaultTxProcessor) ValidateTransaction(tx *types.TransactionWithSign) error {
	// 检查交易基本结构
	if tx == nil {
		return fmt.Errorf("transaction is nil")
	}

	// 检查签名是否存在
	if len(tx.Signature) == 0 {
		return fmt.Errorf("transaction signature is missing")
	}

	// 验证交易签名
	if err := tp.verifyTransactionSignature(&tx.Transaction, tx.Signature); err != nil {
		return fmt.Errorf("transaction signature verification failed: %w", err)
	}

	// 验证交易字段
	if err := tp.validateTransactionFields(&tx.Transaction); err != nil {
		return fmt.Errorf("transaction field validation failed: %w", err)
	}

	return nil
}

// verifyTransactionSignature 验证交易签名
func (tp *DefaultTxProcessor) verifyTransactionSignature(tx *types.Transaction, signature []byte) error {
	// 序列化交易用于签名验证
	data, err := lzbinary.Marshal(tx)
	if err != nil {
		return fmt.Errorf("failed to marshal transaction: %w", err)
	}

	// 注意：这是一个简化的实现，在实际应用中，我们需要从交易的发送方地址推导出公钥
	// 或者交易中直接包含公钥信息。这里我们创建一个新的密钥对仅用于演示
	_, pubKey, err := tp.crypto.GenerateKeyPair(crypto.Ed25519)
	if err != nil {
		return fmt.Errorf("failed to generate key pair: %w", err)
	}

	// 验证签名
	if !tp.crypto.Verify(data, signature, pubKey) {
		return fmt.Errorf("invalid transaction signature")
	}

	return nil
}

// validateTransactionFields 验证交易字段
func (tp *DefaultTxProcessor) validateTransactionFields(tx *types.Transaction) error {
	// 检查金额是否有效
	if tx.Amount <= 0 {
		return fmt.Errorf("transaction amount must be positive")
	}

	// 检查发送方和接收方地址是否相同
	if tx.From == tx.To {
		return fmt.Errorf("sender and receiver addresses cannot be the same")
	}

	// 检查分片ID是否有效
	if tx.ShardID == 0 {
		return fmt.Errorf("shard ID must be positive")
	}

	// 可以添加更多验证规则

	return nil
}

// ApplyTransaction 应用交易到状态
func (tp *DefaultTxProcessor) ApplyTransaction(tx *types.TransactionWithSign) error {
	// 首先验证交易
	if err := tp.ValidateTransaction(tx); err != nil {
		return fmt.Errorf("transaction validation failed: %w", err)
	}

	// 应用交易到状态数据库
	// 1. 从发送方账户扣除金额
	fromBalance, err := tp.GetBalance(tx.From)
	if err != nil {
		fromBalance = 0 // 如果账户不存在，余额为0
	}

	if fromBalance < tx.Amount {
		return fmt.Errorf("insufficient balance: %d < %d", fromBalance, tx.Amount)
	}

	if err := tp.setBalance(tx.From, fromBalance-tx.Amount); err != nil {
		return fmt.Errorf("failed to update sender balance: %w", err)
	}

	// 2. 向接收方账户增加金额
	toBalance, err := tp.GetBalance(tx.To)
	if err != nil {
		toBalance = 0 // 如果账户不存在，余额为0
	}

	if err := tp.setBalance(tx.To, toBalance+tx.Amount); err != nil {
		return fmt.Errorf("failed to update receiver balance: %w", err)
	}

	// 3. 更新账户nonce值
	if err := tp.updateNonce(tx.From, tx.Nonce); err != nil {
		return fmt.Errorf("failed to update nonce: %w", err)
	}

	// 4. 记录交易日志
	if err := tp.storeTransaction(tx); err != nil {
		return fmt.Errorf("failed to store transaction: %w", err)
	}

	return nil
}

// GetTransactionByHash 根据哈希获取交易
func (tp *DefaultTxProcessor) GetTransactionByHash(hash types.Hash) (*types.TransactionWithSign, error) {
	// 从存储中获取交易
	txStorage, err := tp.storage.NewStorage(types.SNTx)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction storage: %w", err)
	}

	data, err := txStorage.Get(hash[:])
	if err != nil {
		return nil, fmt.Errorf("transaction not found: %w", err)
	}

	var tx types.TransactionWithSign
	if err := lzbinary.Unmarshal(data, &tx); err != nil {
		return nil, fmt.Errorf("failed to unmarshal transaction: %w", err)
	}

	return &tx, nil
}

// GetBalance 获取账户余额
func (tp *DefaultTxProcessor) GetBalance(addr types.Address) (uint64, error) {
	// 从存储中获取账户余额
	statusStorage, err := tp.storage.NewStorage(types.SNStatus)
	if err != nil {
		return 0, fmt.Errorf("failed to create status storage: %w", err)
	}

	key := append([]byte("balance_"), addr[:]...)
	data, err := statusStorage.Get(key)
	if err != nil {
		return 0, fmt.Errorf("balance not found: %w", err)
	}

	if len(data) != 8 {
		return 0, fmt.Errorf("invalid balance data")
	}

	balance := uint64(0)
	for i, b := range data {
		balance |= uint64(b) << (8 * i)
	}

	return balance, nil
}

// GetNonce 获取账户nonce值
func (tp *DefaultTxProcessor) GetNonce(addr types.Address) (uint64, error) {
	// 从存储中获取账户nonce值
	statusStorage, err := tp.storage.NewStorage(types.SNStatus)
	if err != nil {
		return 0, fmt.Errorf("failed to create status storage: %w", err)
	}

	key := append([]byte("nonce_"), addr[:]...)
	data, err := statusStorage.Get(key)
	if err != nil {
		return 0, nil // 如果nonce不存在，返回0
	}

	if len(data) != 8 {
		return 0, fmt.Errorf("invalid nonce data")
	}

	nonce := uint64(0)
	for i, b := range data {
		nonce |= uint64(b) << (8 * i)
	}

	return nonce, nil
}

// setBalance 设置账户余额
func (tp *DefaultTxProcessor) setBalance(addr types.Address, balance uint64) error {
	// 存储账户余额
	statusStorage, err := tp.storage.NewStorage(types.SNStatus)
	if err != nil {
		return fmt.Errorf("failed to create status storage: %w", err)
	}

	key := append([]byte("balance_"), addr[:]...)
	data := make([]byte, 8)
	for i := 0; i < 8; i++ {
		data[i] = byte(balance >> (8 * i))
	}

	return statusStorage.Put(key, data)
}

// setNonce 设置账户nonce值
func (tp *DefaultTxProcessor) setNonce(addr types.Address, nonce uint64) error {
	// 存储账户nonce值
	statusStorage, err := tp.storage.NewStorage(types.SNStatus)
	if err != nil {
		return fmt.Errorf("failed to create status storage: %w", err)
	}

	key := append([]byte("nonce_"), addr[:]...)
	data := make([]byte, 8)
	for i := 0; i < 8; i++ {
		data[i] = byte(nonce >> (8 * i))
	}

	return statusStorage.Put(key, data)
}

// storeTransaction 存储交易
func (tp *DefaultTxProcessor) storeTransaction(tx *types.TransactionWithSign) error {
	// 存储交易
	txStorage, err := tp.storage.NewStorage(types.SNTx)
	if err != nil {
		return fmt.Errorf("failed to create transaction storage: %w", err)
	}

	data, err := lzbinary.Marshal(tx)
	if err != nil {
		return fmt.Errorf("failed to marshal transaction: %w", err)
	}

	txHash := tp.calculateTransactionHash(tx)
	return txStorage.Put(txHash[:], data)
}

// updateNonce 更新账户nonce值
func (tp *DefaultTxProcessor) updateNonce(addr types.Address, nonce uint64) error {
	// 获取当前nonce值
	currentNonce, err := tp.GetNonce(addr)
	if err != nil {
		currentNonce = 0 // 如果账户不存在，nonce为0
	}

	// 检查nonce是否有效
	if nonce <= currentNonce {
		return fmt.Errorf("invalid nonce: expected %d, got %d", currentNonce+1, nonce)
	}

	// 更新nonce值
	return tp.setNonce(addr, nonce)
}

// calculateTransactionHash 计算交易哈希
func (tp *DefaultTxProcessor) calculateTransactionHash(tx *types.TransactionWithSign) types.Hash {
	// 序列化交易（排除签名字段）
	txCopy := *tx
	txCopy.Signature = nil

	data, err := lzbinary.Marshal(&txCopy.Transaction)
	if err != nil {
		return types.Hash{} // 返回空哈希
	}

	return tp.crypto.Hash(data)
}
