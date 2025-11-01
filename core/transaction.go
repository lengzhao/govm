package core

import (
	"fmt"

	"github.com/lengzhao/govm/crypto"
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
}

// DefaultTxProcessor 默认交易处理器实现
type DefaultTxProcessor struct {
	crypto crypto.Crypto
}

// NewTxProcessor 创建新的交易处理器
func NewTxProcessor() TxProcessor {
	return &DefaultTxProcessor{
		crypto: crypto.NewCrypto(),
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
	// 创建交易副本并移除签名字段用于验证
	_ = *tx

	// 序列化交易用于签名验证
	// TODO: 实现序列化逻辑
	// data, err := lzbinary.Marshal(&txCopy)
	// if err != nil {
	// 	return fmt.Errorf("failed to marshal transaction: %w", err)
	// }

	// 从发送方地址推导公钥（简化实现）
	// 在实际实现中，需要从交易中获取公钥或通过其他方式验证
	// 这里我们假设地址就是公钥的哈希

	// 创建公钥对象（示例使用Ed25519）
	_, _, err := tp.crypto.GenerateKeyPair(crypto.Ed25519)
	if err != nil {
		return fmt.Errorf("failed to generate key pair: %w", err)
	}

	// 注意：这是一个简化的实现，实际应用中需要正确地从地址或交易中获取公钥
	// 验证签名
	// if !tp.crypto.Verify(data, signature, pubKey) {
	// 	return fmt.Errorf("invalid transaction signature")
	// }

	_ = signature // 避免未使用变量警告

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
	// TODO: 实现状态更新逻辑
	// 这包括：
	// 1. 从发送方账户扣除金额
	// 2. 向接收方账户增加金额
	// 3. 更新账户nonce值
	// 4. 记录交易日志

	return nil
}

// GetTransactionByHash 根据哈希获取交易
func (tp *DefaultTxProcessor) GetTransactionByHash(hash types.Hash) (*types.TransactionWithSign, error) {
	// TODO: 从存储中获取交易
	// 这需要与存储模块集成

	var tx types.TransactionWithSign
	return &tx, nil
}

// calculateTransactionHash 计算交易哈希
func (tp *DefaultTxProcessor) calculateTransactionHash(tx *types.TransactionWithSign) types.Hash {
	// 序列化交易（排除签名字段）
	txCopy := *tx
	txCopy.Signature = nil

	// TODO: 实现序列化逻辑
	// data, err := lzbinary.Marshal(&txCopy.Transaction)
	// if err != nil {
	// 	return types.Hash{}
	// }

	// return tp.crypto.Hash(data)
	return types.Hash{} // 临时返回空哈希
}
