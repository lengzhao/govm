package test

import (
	"testing"
	"time"

	lzbinary "github.com/lengzhao/binary"
	"github.com/lengzhao/govm/consensus"
	"github.com/lengzhao/govm/core"
	"github.com/lengzhao/govm/crypto"
	"github.com/lengzhao/govm/generator"
	"github.com/lengzhao/govm/storage"
	"github.com/lengzhao/govm/txpool"
	"github.com/lengzhao/govm/types"
	"github.com/stretchr/testify/assert"
)

// mockTxPool 创建一个模拟的交易池，实现TxPool接口
type mockTxPool struct {
	transactions []*types.TransactionWithSign
}

// 实现TxPool接口的方法
func (m *mockTxPool) AddTransaction(tx *types.TransactionWithSign) error {
	return nil
}

func (m *mockTxPool) GetTransaction(hash types.Hash) (*types.TransactionWithSign, error) {
	return nil, nil
}

func (m *mockTxPool) RemoveTransaction(hash types.Hash) error {
	return nil
}

func (m *mockTxPool) SelectTransactions(maxCount int) ([]*types.TransactionWithSign, error) {
	// 返回预设的交易
	return m.transactions, nil
}

func (m *mockTxPool) ValidateTransaction(tx *types.TransactionWithSign) error {
	return nil
}

func (m *mockTxPool) GetTransactionCount() int {
	return len(m.transactions)
}

func (m *mockTxPool) Start() error {
	return nil
}

func (m *mockTxPool) Stop() error {
	return nil
}

// TestBlockGenerationWithTransactions 测试区块生成器处理交易的功能
func TestBlockGenerationWithTransactions(t *testing.T) {
	// 创建存储实例
	store := storage.NewMemoryStorage("")
	err := store.Start()
	assert.NoError(t, err)
	defer store.Stop()

	// 创建验证者地址
	validatorAddr := types.Address{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}

	// 创建共识配置
	config := &consensus.PoAConfig{
		Validators:    []types.Address{validatorAddr},
		BlockInterval: 100, // 100毫秒区块间隔（测试用）
		RoundLength:   1,
	}

	// 创建共识实例
	cons := consensus.NewPoAConsensus(config, store)
	assert.NotNil(t, cons)

	// 创建核心模块配置
	coreConfig := &core.CoreConfig{
		ShardID: types.DefaultShardID,
		DataDir: "",
		Genesis: &types.GenesisConfig{
			Timestamp: uint64(time.Now().Unix()),
		},
	}

	// 创建核心模块
	coreModule, err := core.NewCore(coreConfig, cons, store)
	assert.NoError(t, err)
	assert.NotNil(t, coreModule)

	// 启动核心模块
	err = coreModule.Start()
	assert.NoError(t, err)
	defer coreModule.Stop()

	// 创建区块生成器（不使用交易池）
	blockGenerator := generator.NewBlockGenerator(cons, store, nil)

	// 验证区块生成器创建成功
	assert.NotNil(t, blockGenerator)

	// 验证初始状态
	lastBlock := coreModule.GetLastBlock()
	assert.NotNil(t, lastBlock)
	assert.Equal(t, uint64(0), lastBlock.Header.BlockNumber) // 创世区块高度为0

	// 测试生成不包含交易的区块
	block1, err := blockGenerator.GenerateBlock(lastBlock)
	assert.NoError(t, err)
	assert.NotNil(t, block1)

	// 验证区块不包含交易
	assert.Equal(t, 0, len(block1.Transactions), "生成的区块不应该包含交易哈希")

	// 验证区块头信息
	assert.Equal(t, uint64(1), block1.Header.BlockNumber)
	assert.Equal(t, types.DefaultShardID, block1.Header.ShardID)
	assert.Equal(t, validatorAddr, block1.Header.Validator)

	// 添加区块到区块链
	err = coreModule.AddBlock(block1)
	assert.NoError(t, err)

	// 验证区块链高度更新
	assert.Equal(t, uint64(1), coreModule.GetHeight())
}

// TestSelectTransactionsFromPool 测试从交易池选择交易的功能（模拟）
func TestSelectTransactionsFromPool(t *testing.T) {
	// 创建一些模拟的交易
	transactions := []*types.TransactionWithSign{
		{
			Transaction: types.Transaction{
				ShardID: types.DefaultShardID,
				From:    types.Address{1, 2, 3},
				To:      types.Address{4, 5, 6},
				Amount:  100,
				Nonce:   1,
				Data:    []byte("test transaction 1"),
			},
			Signature: []byte("signature1"),
		},
		{
			Transaction: types.Transaction{
				ShardID: types.DefaultShardID,
				From:    types.Address{7, 8, 9},
				To:      types.Address{10, 11, 12},
				Amount:  200,
				Nonce:   1,
				Data:    []byte("test transaction 2"),
			},
			Signature: []byte("signature2"),
		},
	}

	mockPool := &mockTxPool{
		transactions: transactions,
	}

	// 验证mockPool实现了TxPool接口
	var _ txpool.TxPool = mockPool

	// 创建存储实例
	store := storage.NewMemoryStorage("")
	err := store.Start()
	assert.NoError(t, err)
	defer store.Stop()

	// 创建验证者地址
	validatorAddr := types.Address{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}

	// 创建共识配置
	config := &consensus.PoAConfig{
		Validators:    []types.Address{validatorAddr},
		BlockInterval: 100, // 100毫秒区块间隔（测试用）
		RoundLength:   1,
	}

	// 创建共识实例
	cons := consensus.NewPoAConsensus(config, store)
	assert.NotNil(t, cons)

	// 创建区块生成器
	blockGenerator := generator.NewBlockGenerator(cons, store, mockPool)

	// 验证区块生成器创建成功
	assert.NotNil(t, blockGenerator)

	// 测试从交易池选择交易的功能
	selectedTransactions, err := blockGenerator.SelectTransactions()
	assert.NoError(t, err)
	assert.Equal(t, 2, len(selectedTransactions), "应该从交易池中选择2个交易")
}

// TestTransactionPackagingIntoBlock 测试交易打包进区块的核心功能（不进行完整的区块验证）
func TestTransactionPackagingIntoBlock(t *testing.T) {
	// 创建一些模拟的交易
	transactions := []*types.TransactionWithSign{
		{
			Transaction: types.Transaction{
				ShardID: types.DefaultShardID,
				From:    types.Address{1, 2, 3},
				To:      types.Address{4, 5, 6},
				Amount:  100,
				Nonce:   1,
				Data:    []byte("test transaction 1"),
			},
			Signature: []byte("signature1"),
		},
		{
			Transaction: types.Transaction{
				ShardID: types.DefaultShardID,
				From:    types.Address{7, 8, 9},
				To:      types.Address{10, 11, 12},
				Amount:  200,
				Nonce:   1,
				Data:    []byte("test transaction 2"),
			},
			Signature: []byte("signature2"),
		},
	}

	mockPool := &mockTxPool{
		transactions: transactions,
	}

	// 创建存储实例
	store := storage.NewMemoryStorage("")
	err := store.Start()
	assert.NoError(t, err)
	defer store.Stop()

	// 创建验证者地址
	validatorAddr := types.Address{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}

	// 创建共识配置
	config := &consensus.PoAConfig{
		Validators:    []types.Address{validatorAddr},
		BlockInterval: 100, // 100毫秒区块间隔（测试用）
		RoundLength:   1,
	}

	// 创建共识实例
	cons := consensus.NewPoAConsensus(config, store)
	assert.NotNil(t, cons)

	// 创建核心模块配置
	coreConfig := &core.CoreConfig{
		ShardID: types.DefaultShardID,
		DataDir: "",
		Genesis: &types.GenesisConfig{
			Timestamp: uint64(time.Now().Unix()),
		},
	}

	// 创建核心模块
	coreModule, err := core.NewCore(coreConfig, cons, store)
	assert.NoError(t, err)
	assert.NotNil(t, coreModule)

	// 启动核心模块
	err = coreModule.Start()
	assert.NoError(t, err)
	defer coreModule.Stop()

	// 创建区块生成器
	blockGenerator := generator.NewBlockGenerator(cons, store, mockPool)

	// 验证区块生成器创建成功
	assert.NotNil(t, blockGenerator)

	// 获取创世区块作为前一个区块
	lastBlock := coreModule.GetLastBlock()
	assert.NotNil(t, lastBlock)
	assert.Equal(t, uint64(0), lastBlock.Header.BlockNumber)

	// 直接测试区块生成的核心功能，而不进行完整的区块验证
	// 1. 选择交易
	selectedTransactions, err := blockGenerator.SelectTransactions()
	assert.NoError(t, err)
	assert.Equal(t, 2, len(selectedTransactions), "应该从交易池中选择2个交易")

	// 2. 构建区块头
	header, err := blockGenerator.BuildBlockHeader(lastBlock, selectedTransactions)
	assert.NoError(t, err)
	assert.NotNil(t, header)

	// 验证区块头信息
	assert.Equal(t, uint64(1), header.BlockNumber)
	assert.Equal(t, types.DefaultShardID, header.ShardID)
	assert.Equal(t, validatorAddr, header.Validator)
	assert.NotEqual(t, types.Hash{}, header.MerkleRoot, "Merkle根不应该为空")

	// 3. 组装完整区块
	block, err := blockGenerator.AssembleBlock(header, selectedTransactions)
	assert.NoError(t, err)
	assert.NotNil(t, block)

	// 验证区块包含交易哈希
	assert.Equal(t, 2, len(block.Transactions), "生成的区块应该包含2个交易哈希")

	// 验证交易哈希是否正确计算
	for i, tx := range selectedTransactions {
		// 序列化交易以计算哈希
		data, err := lzbinary.Marshal(tx)
		assert.NoError(t, err)
		expectedHash := crypto.Hash(data)
		assert.Equal(t, expectedHash, block.Transactions[i], "交易哈希应该正确计算")
	}

	// 4. 对区块进行签名
	err = blockGenerator.SignBlock(block)
	assert.NoError(t, err)
	assert.NotNil(t, block.Header.Signature, "区块应该被签名")

	// 验证区块头信息
	assert.Equal(t, uint64(1), block.Header.BlockNumber)
	assert.Equal(t, types.DefaultShardID, block.Header.ShardID)
	assert.Equal(t, validatorAddr, block.Header.Validator)
	assert.NotEqual(t, types.Hash{}, block.Header.MerkleRoot, "Merkle根不应该为空")
}
