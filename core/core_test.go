package core

import (
	"testing"

	"github.com/lengzhao/govm/consensus"
	"github.com/lengzhao/govm/storage"
	"github.com/lengzhao/govm/types"
	"github.com/stretchr/testify/assert"
)

// mockConsensus 是一个模拟的共识实现
type mockConsensus struct {
	consensus.PoAConsensus
}

// ValidateBlock 实现共识接口的ValidateBlock方法
func (m *mockConsensus) ValidateBlock(block *types.Block) error {
	// 简单实现，总是返回nil
	return nil
}

// GetValidator 实现共识接口的GetValidator方法
func (m *mockConsensus) GetValidator() interface{} {
	return nil
}

// GetValidators 实现共识接口的GetValidators方法
func (m *mockConsensus) GetValidators() []interface{} {
	return nil
}

// GetRound 实现PoA共识接口的GetRound方法
func (m *mockConsensus) GetRound() uint64 {
	return 0
}

// GetTurn 实现PoA共识接口的GetTurn方法
func (m *mockConsensus) GetTurn() uint64 {
	return 0
}

// IsValidator 实现PoA共识接口的IsValidator方法
func (m *mockConsensus) IsValidator(addr types.Address) bool {
	return true
}

// UpdateValidators 实现PoA共识接口的UpdateValidators方法
func (m *mockConsensus) UpdateValidators(validators []types.Address) error {
	return nil
}

func TestNewCore(t *testing.T) {
	// 准备测试数据
	config := &CoreConfig{
		ShardID: types.DefaultShardID,
		DataDir: "",
	}

	mockStore := storage.NewMemoryStorage("")
	err := mockStore.Start()
	assert.NoError(t, err)
	defer mockStore.Stop()

	mockCons := &mockConsensus{}

	// 创建核心模块实例
	core, err := NewCore(config, mockCons, mockStore)
	assert.NoError(t, err)
	assert.NotNil(t, core)
}

func TestCore_StartStop(t *testing.T) {
	// 准备测试数据
	config := &CoreConfig{
		ShardID: types.DefaultShardID,
		DataDir: "",
	}

	mockStore := storage.NewMemoryStorage("")
	err := mockStore.Start()
	assert.NoError(t, err)
	defer mockStore.Stop()

	mockCons := &mockConsensus{}

	// 创建核心模块实例
	core, err := NewCore(config, mockCons, mockStore)
	assert.NoError(t, err)
	assert.NotNil(t, core)

	// 启动核心模块
	err = core.Start()
	assert.NoError(t, err)

	// 再次启动应该返回错误
	err = core.Start()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "core module is already running")

	// 停止核心模块
	err = core.Stop()
	assert.NoError(t, err)

	// 再次停止应该返回错误
	err = core.Stop()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "core module is not running")
}

func TestCore_GetConsensus(t *testing.T) {
	// 准备测试数据
	config := &CoreConfig{
		ShardID: types.DefaultShardID,
		DataDir: "",
	}

	mockStore := storage.NewMemoryStorage("")
	err := mockStore.Start()
	assert.NoError(t, err)
	defer mockStore.Stop()

	mockCons := &mockConsensus{}

	// 创建核心模块实例
	core, err := NewCore(config, mockCons, mockStore)
	assert.NoError(t, err)
	assert.NotNil(t, core)

	// 获取共识模块
	cons := core.GetConsensus()
	assert.Equal(t, mockCons, cons)
}

func TestCore_GetStorage(t *testing.T) {
	// 准备测试数据
	config := &CoreConfig{
		ShardID: types.DefaultShardID,
		DataDir: "",
	}

	mockStore := storage.NewMemoryStorage("")
	err := mockStore.Start()
	assert.NoError(t, err)
	defer mockStore.Stop()

	mockCons := &mockConsensus{}

	// 创建核心模块实例
	core, err := NewCore(config, mockCons, mockStore)
	assert.NoError(t, err)
	assert.NotNil(t, core)

	// 获取存储模块
	store := core.GetStorage()
	assert.Equal(t, mockStore, store)
}

func TestCore_BlockchainMethods(t *testing.T) {
	// 准备测试数据
	config := &CoreConfig{
		ShardID: types.DefaultShardID,
		DataDir: "",
	}

	mockStore := storage.NewMemoryStorage("")
	err := mockStore.Start()
	assert.NoError(t, err)
	defer mockStore.Stop()

	mockCons := &mockConsensus{}

	// 创建核心模块实例
	core, err := NewCore(config, mockCons, mockStore)
	assert.NoError(t, err)
	assert.NotNil(t, core)

	// 启动核心模块
	err = core.Start()
	assert.NoError(t, err)

	// 测试区块链方法
	height := core.GetHeight()
	assert.Equal(t, uint64(0), height)

	lastBlock := core.GetLastBlock()
	assert.NotNil(t, lastBlock)
	assert.Equal(t, uint64(0), lastBlock.Header.BlockNumber)
}

func TestCore_TxProcessorMethods(t *testing.T) {
	// 准备测试数据
	config := &CoreConfig{
		ShardID: types.DefaultShardID,
		DataDir: "",
	}

	mockStore := storage.NewMemoryStorage("")
	err := mockStore.Start()
	assert.NoError(t, err)
	defer mockStore.Stop()

	mockCons := &mockConsensus{}

	// 创建核心模块实例
	core, err := NewCore(config, mockCons, mockStore)
	assert.NoError(t, err)
	assert.NotNil(t, core)

	// 启动核心模块
	err = core.Start()
	assert.NoError(t, err)

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
	err = core.ValidateTransaction(tx)
	assert.NoError(t, err)

	// 应用交易
	err = core.ApplyTransaction(tx)
	assert.NoError(t, err)

	// 根据哈希获取交易
	hash := types.Hash{1, 2, 3}
	retrievedTx, err := core.GetTransactionByHash(hash)
	assert.NoError(t, err)
	assert.NotNil(t, retrievedTx)
}
