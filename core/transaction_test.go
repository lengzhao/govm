package core

import (
	"testing"

	"github.com/lengzhao/govm/crypto"
	"github.com/lengzhao/govm/storage"
	"github.com/lengzhao/govm/types"
	"github.com/stretchr/testify/assert"
)

func TestTxProcessor_ValidateTransaction(t *testing.T) {
	// 创建存储
	store := storage.NewMemoryStorage("")
	err := store.Start()
	assert.NoError(t, err)
	defer store.Stop()

	// 创建交易处理器
	txProcessor := NewTxProcessor(store)

	// 生成测试密钥对
	cryptoInstance := crypto.NewCrypto()
	_, pubKey, err := cryptoInstance.GenerateKeyPair(crypto.Ed25519)
	assert.NoError(t, err)

	// 创建有效的交易
	validTx := &types.TransactionWithSign{
		Transaction: types.Transaction{
			ShardID:   1,
			From:      types.Address{1},
			To:        types.Address{2},
			Amount:    100,
			Nonce:     1,
			Data:      []byte("test data"),
			GasPrice:  1,
			GasLimit:  21000,
			GasFeeCap: 1000,
			PublicKey: pubKey.Bytes(),
		},
		Signature: []byte("test signature"),
	}

	// 测试有效交易
	err = txProcessor.ValidateTransaction(validTx)
	assert.NoError(t, err)

	// 测试无效交易（金额为0）
	invalidTx1 := &types.TransactionWithSign{
		Transaction: types.Transaction{
			ShardID:   1,
			From:      types.Address{1},
			To:        types.Address{2},
			Amount:    0,
			Nonce:     1,
			Data:      []byte("test data"),
			GasPrice:  1,
			GasLimit:  21000,
			GasFeeCap: 1000,
			PublicKey: pubKey.Bytes(),
		},
		Signature: []byte("test signature"),
	}

	err = txProcessor.ValidateTransaction(invalidTx1)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "transaction amount must be positive")

	// 测试无效交易（发送方和接收方相同）
	invalidTx2 := &types.TransactionWithSign{
		Transaction: types.Transaction{
			ShardID:   1,
			From:      types.Address{1},
			To:        types.Address{1},
			Amount:    100,
			Nonce:     1,
			Data:      []byte("test data"),
			GasPrice:  1,
			GasLimit:  21000,
			GasFeeCap: 1000,
			PublicKey: pubKey.Bytes(),
		},
		Signature: []byte("test signature"),
	}

	err = txProcessor.ValidateTransaction(invalidTx2)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "sender and receiver addresses cannot be the same")

	// 测试无效交易（缺少公钥）
	invalidTx3 := &types.TransactionWithSign{
		Transaction: types.Transaction{
			ShardID:   1,
			From:      types.Address{1},
			To:        types.Address{2},
			Amount:    100,
			Nonce:     1,
			Data:      []byte("test data"),
			GasPrice:  1,
			GasLimit:  21000,
			GasFeeCap: 1000,
			PublicKey: []byte{},
		},
		Signature: []byte("test signature"),
	}

	err = txProcessor.ValidateTransaction(invalidTx3)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "public key is required")
}

func TestTxProcessor_ApplyTransaction(t *testing.T) {
	// 创建存储
	store := storage.NewMemoryStorage("")
	err := store.Start()
	assert.NoError(t, err)
	defer store.Stop()

	// 创建交易处理器
	txProcessor := NewTxProcessor(store)

	// 生成测试密钥对
	cryptoInstance := crypto.NewCrypto()
	_, pubKey, err := cryptoInstance.GenerateKeyPair(crypto.Ed25519)
	assert.NoError(t, err)

	// 创建发送方和接收方地址
	fromAddr := types.Address{1}
	toAddr := types.Address{2}

	// 初始化发送方账户余额
	err = txProcessor.(*DefaultTxProcessor).setBalance(fromAddr, 1000)
	assert.NoError(t, err)

	// 创建交易
	tx := &types.TransactionWithSign{
		Transaction: types.Transaction{
			ShardID:   1,
			From:      fromAddr,
			To:        toAddr,
			Amount:    100,
			Nonce:     1,
			Data:      []byte("test data"),
			GasPrice:  1,
			GasLimit:  21000,
			GasFeeCap: 1000,
			PublicKey: pubKey.Bytes(),
		},
		Signature: []byte("test signature"),
	}

	// 应用交易
	err = txProcessor.ApplyTransaction(tx)
	assert.NoError(t, err)

	// 验证发送方余额
	fromBalance, err := txProcessor.GetBalance(fromAddr)
	assert.NoError(t, err)
	// 余额应该是 1000 - 100 - (21000 * 1) = 790
	assert.Equal(t, uint64(790), fromBalance)

	// 验证接收方余额
	toBalance, err := txProcessor.GetBalance(toAddr)
	assert.NoError(t, err)
	assert.Equal(t, uint64(100), toBalance)

	// 验证发送方nonce
	fromNonce, err := txProcessor.GetNonce(fromAddr)
	assert.NoError(t, err)
	assert.Equal(t, uint64(1), fromNonce)
}

func TestTxProcessor_GetBalance(t *testing.T) {
	// 创建存储
	store := storage.NewMemoryStorage("")
	err := store.Start()
	assert.NoError(t, err)
	defer store.Stop()

	// 创建交易处理器
	txProcessor := NewTxProcessor(store)

	// 测试不存在的账户
	addr := types.Address{1}
	balance, err := txProcessor.GetBalance(addr)
	assert.NoError(t, err)
	assert.Equal(t, uint64(0), balance)

	// 设置账户余额
	err = txProcessor.(*DefaultTxProcessor).setBalance(addr, 1000)
	assert.NoError(t, err)

	// 测试存在的账户
	balance, err = txProcessor.GetBalance(addr)
	assert.NoError(t, err)
	assert.Equal(t, uint64(1000), balance)
}

func TestTxProcessor_GetNonce(t *testing.T) {
	// 创建存储
	store := storage.NewMemoryStorage("")
	err := store.Start()
	assert.NoError(t, err)
	defer store.Stop()

	// 创建交易处理器
	txProcessor := NewTxProcessor(store)

	// 测试不存在的账户
	addr := types.Address{1}
	nonce, err := txProcessor.GetNonce(addr)
	assert.NoError(t, err)
	assert.Equal(t, uint64(0), nonce)

	// 设置账户nonce
	err = txProcessor.(*DefaultTxProcessor).setNonce(addr, 5)
	assert.NoError(t, err)

	// 测试存在的账户
	nonce, err = txProcessor.GetNonce(addr)
	assert.NoError(t, err)
	assert.Equal(t, uint64(5), nonce)
}

func TestTxProcessor_CalculateGasFee(t *testing.T) {
	// 创建存储
	store := storage.NewMemoryStorage("")
	err := store.Start()
	assert.NoError(t, err)
	defer store.Stop()

	// 创建交易处理器
	txProcessor := NewTxProcessor(store)

	// 创建交易
	tx := &types.Transaction{
		ShardID:   1,
		From:      types.Address{1},
		To:        types.Address{2},
		Amount:    100,
		Nonce:     1,
		Data:      []byte("test data"), // 9 bytes
		GasPrice:  10,
		GasLimit:  30000,
		GasFeeCap: 100000,
	}

	// 计算Gas费用
	gasFee := txProcessor.(*DefaultTxProcessor).calculateGasFee(tx)
	// 基础Gas: 21000, 数据Gas: 9, 总Gas: 21009
	// Gas费用: 21009 * 10 = 210090
	// 但不应超过Gas费用上限，所以应该是100000
	assert.Equal(t, uint64(100000), gasFee)

	// 测试Gas限制的情况
	tx2 := &types.Transaction{
		ShardID:   1,
		From:      types.Address{1},
		To:        types.Address{2},
		Amount:    100,
		Nonce:     1,
		Data:      []byte("test data"), // 9 bytes
		GasPrice:  10,
		GasLimit:  21005, // 低于总Gas使用量
		GasFeeCap: 100000,
	}

	gasFee2 := txProcessor.(*DefaultTxProcessor).calculateGasFee(tx2)
	// 基础Gas: 21000, 数据Gas: 9, 总Gas: 21009
	// 但受限于GasLimit: 21005
	// Gas费用: 21005 * 10 = 210050
	// 但不应超过Gas费用上限，所以应该是100000
	assert.Equal(t, uint64(100000), gasFee2)
}
