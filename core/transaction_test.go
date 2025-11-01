package core

import (
	"testing"

	"github.com/lengzhao/govm/types"
	"github.com/stretchr/testify/assert"
)

func TestNewTxProcessor(t *testing.T) {
	// 创建交易处理器
	txProcessor := NewTxProcessor()

	// 验证实例创建成功
	assert.NotNil(t, txProcessor)
}

func TestTxProcessor_ValidateTransaction(t *testing.T) {
	// 创建交易处理器
	txProcessor := NewTxProcessor()

	// 测试nil交易
	err := txProcessor.ValidateTransaction(nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "transaction is nil")

	// 创建一个有效的交易
	tx := &types.TransactionWithSign{
		Transaction: types.Transaction{
			ShardID: 1,
			From:    types.Address{1},
			To:      types.Address{2},
			Amount:  100,
			Nonce:   1,
		},
		Signature: []byte("signature"),
	}

	// 验证交易
	err = txProcessor.ValidateTransaction(tx)
	assert.NoError(t, err)
}

func TestTxProcessor_ValidateTransaction_InvalidSignature(t *testing.T) {
	// 创建交易处理器
	txProcessor := NewTxProcessor()

	// 创建一个没有签名的交易
	tx := &types.TransactionWithSign{
		Transaction: types.Transaction{
			ShardID: 1,
			From:    types.Address{1},
			To:      types.Address{2},
			Amount:  100,
			Nonce:   1,
		},
		Signature: []byte{}, // 空签名
	}

	// 验证交易
	err := txProcessor.ValidateTransaction(tx)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "transaction signature is missing")
}

func TestTxProcessor_ValidateTransaction_InvalidAmount(t *testing.T) {
	// 创建交易处理器
	txProcessor := NewTxProcessor()

	// 创建一个金额为0的交易
	tx := &types.TransactionWithSign{
		Transaction: types.Transaction{
			ShardID: 1,
			From:    types.Address{1},
			To:      types.Address{2},
			Amount:  0, // 无效金额
			Nonce:   1,
		},
		Signature: []byte("signature"),
	}

	// 验证交易
	err := txProcessor.ValidateTransaction(tx)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "transaction amount must be positive")
}

func TestTxProcessor_ValidateTransaction_SameAddresses(t *testing.T) {
	// 创建交易处理器
	txProcessor := NewTxProcessor()

	// 创建发送方和接收方地址相同的交易
	addr := types.Address{1}
	tx := &types.TransactionWithSign{
		Transaction: types.Transaction{
			ShardID: 1,
			From:    addr,
			To:      addr, // 相同地址
			Amount:  100,
			Nonce:   1,
		},
		Signature: []byte("signature"),
	}

	// 验证交易
	err := txProcessor.ValidateTransaction(tx)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "sender and receiver addresses cannot be the same")
}

func TestTxProcessor_ValidateTransaction_InvalidShardID(t *testing.T) {
	// 创建交易处理器
	txProcessor := NewTxProcessor()

	// 创建分片ID为0的交易
	tx := &types.TransactionWithSign{
		Transaction: types.Transaction{
			ShardID: 0, // 无效分片ID
			From:    types.Address{1},
			To:      types.Address{2},
			Amount:  100,
			Nonce:   1,
		},
		Signature: []byte("signature"),
	}

	// 验证交易
	err := txProcessor.ValidateTransaction(tx)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "shard ID must be positive")
}

func TestTxProcessor_ApplyTransaction(t *testing.T) {
	// 创建交易处理器
	txProcessor := NewTxProcessor()

	// 创建一个有效的交易
	tx := &types.TransactionWithSign{
		Transaction: types.Transaction{
			ShardID: 1,
			From:    types.Address{1},
			To:      types.Address{2},
			Amount:  100,
			Nonce:   1,
		},
		Signature: []byte("signature"),
	}

	// 应用交易
	err := txProcessor.ApplyTransaction(tx)
	assert.NoError(t, err)
}

func TestTxProcessor_GetTransactionByHash(t *testing.T) {
	// 创建交易处理器
	txProcessor := NewTxProcessor()

	// 创建一个哈希
	hash := types.Hash{1, 2, 3}

	// 根据哈希获取交易
	tx, err := txProcessor.GetTransactionByHash(hash)
	assert.NoError(t, err)
	assert.NotNil(t, tx)
}

func TestTxProcessor_calculateTransactionHash(t *testing.T) {
	// 创建交易处理器
	txProcessor := NewTxProcessor().(*DefaultTxProcessor)

	// 创建一个交易
	tx := &types.TransactionWithSign{
		Transaction: types.Transaction{
			ShardID: 1,
			From:    types.Address{1},
			To:      types.Address{2},
			Amount:  100,
			Nonce:   1,
		},
		Signature: []byte("signature"),
	}

	// 计算交易哈希
	hash := txProcessor.calculateTransactionHash(tx)
	assert.NotNil(t, hash)
}
