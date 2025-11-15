package test

import (
	"testing"
	"time"

	"github.com/lengzhao/govm/consensus"
	"github.com/lengzhao/govm/core"
	"github.com/lengzhao/govm/generator"
	"github.com/lengzhao/govm/storage"
	"github.com/lengzhao/govm/txpool"
	"github.com/lengzhao/govm/types"
	"github.com/stretchr/testify/assert"
)

// mockSyncChecker 模拟同步检查器
type mockSyncChecker struct {
	isSyncing bool
}

func (m *mockSyncChecker) IsSyncing() bool {
	return m.isSyncing
}

// TestThreeNodeRoundRobinIntegration 测试3节点轮流出块集成
func TestThreeNodeRoundRobinIntegration(t *testing.T) {
	// 创建存储实例
	store := storage.NewMemoryStorage("")
	err := store.Start()
	assert.NoError(t, err)
	defer store.Stop()

	// 创建3个验证者地址
	validator1 := types.Address{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}
	validator2 := types.Address{2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2}
	validator3 := types.Address{3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 3}

	// 创建共识配置（3个验证节点）
	config := &consensus.PoAConfig{
		Validators:    []types.Address{validator1, validator2, validator3},
		BlockInterval: 100, // 100毫秒区块间隔（测试用）
		RoundLength:   3,   // 3个区块为一轮
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

	// 创建交易池
	txPool := txpool.NewTxPool(coreModule, store)
	err = txPool.Start()
	assert.NoError(t, err)
	defer txPool.Stop()

	// 创建区块生成器
	blockGenerator := generator.NewBlockGenerator(cons, store, txPool)

	// 验证区块生成器创建成功
	assert.NotNil(t, blockGenerator)

	// 验证初始状态
	lastBlock := coreModule.GetLastBlock()
	assert.NotNil(t, lastBlock)
	assert.Equal(t, uint64(0), lastBlock.Header.BlockNumber) // 创世区块高度为0

	// 模拟3个节点轮流生成区块
	// 节点1生成高度1的区块
	t.Log("节点1生成高度1的区块")
	block1, err := blockGenerator.GenerateBlock(lastBlock)
	assert.NoError(t, err)
	assert.NotNil(t, block1)
	assert.Equal(t, uint64(1), block1.Header.BlockNumber)

	// 验证应该由验证者2出块（索引1）
	expectedValidator1 := validator2
	assert.Equal(t, expectedValidator1, block1.Header.Validator)

	// 添加区块到区块链
	err = coreModule.AddBlock(block1)
	assert.NoError(t, err)

	// 验证区块链高度更新
	assert.Equal(t, uint64(1), coreModule.GetHeight())

	// 节点2生成高度2的区块
	t.Log("节点2生成高度2的区块")
	block2, err := blockGenerator.GenerateBlock(block1)
	assert.NoError(t, err)
	assert.NotNil(t, block2)
	assert.Equal(t, uint64(2), block2.Header.BlockNumber)

	// 验证应该由验证者3出块（索引2）
	expectedValidator2 := validator3
	assert.Equal(t, expectedValidator2, block2.Header.Validator)

	// 添加区块到区块链
	err = coreModule.AddBlock(block2)
	assert.NoError(t, err)

	// 验证区块链高度更新
	assert.Equal(t, uint64(2), coreModule.GetHeight())

	// 节点3生成高度3的区块
	t.Log("节点3生成高度3的区块")
	block3, err := blockGenerator.GenerateBlock(block2)
	assert.NoError(t, err)
	assert.NotNil(t, block3)
	assert.Equal(t, uint64(3), block3.Header.BlockNumber)

	// 验证应该由验证者1出块（索引0）
	expectedValidator3 := validator1
	assert.Equal(t, expectedValidator3, block3.Header.Validator)

	// 添加区块到区块链
	err = coreModule.AddBlock(block3)
	assert.NoError(t, err)

	// 验证区块链高度更新
	assert.Equal(t, uint64(3), coreModule.GetHeight())

	// 节点1生成高度4的区块（新一轮开始）
	t.Log("节点1生成高度4的区块")
	block4, err := blockGenerator.GenerateBlock(block3)
	assert.NoError(t, err)
	assert.NotNil(t, block4)
	assert.Equal(t, uint64(4), block4.Header.BlockNumber)

	// 验证应该由验证者2出块（索引1）
	expectedValidator4 := validator2
	assert.Equal(t, expectedValidator4, block4.Header.Validator)

	// 添加区块到区块链
	err = coreModule.AddBlock(block4)
	assert.NoError(t, err)

	// 验证区块链高度更新
	assert.Equal(t, uint64(4), coreModule.GetHeight())

	// 验证所有区块都已正确添加
	for height := uint64(1); height <= 4; height++ {
		block, err := coreModule.GetBlockByHeight(height)
		assert.NoError(t, err)
		assert.NotNil(t, block)
		assert.Equal(t, height, block.Header.BlockNumber)
	}

	// 验证出块顺序符合轮流出块规则
	block1FromChain, _ := coreModule.GetBlockByHeight(1)
	assert.Equal(t, validator2, block1FromChain.Header.Validator)

	block2FromChain, _ := coreModule.GetBlockByHeight(2)
	assert.Equal(t, validator3, block2FromChain.Header.Validator)

	block3FromChain, _ := coreModule.GetBlockByHeight(3)
	assert.Equal(t, validator1, block3FromChain.Header.Validator)

	block4FromChain, _ := coreModule.GetBlockByHeight(4)
	assert.Equal(t, validator2, block4FromChain.Header.Validator)
}

// TestIsMyTurnFunction 测试IsMyTurn函数在不同场景下的行为
func TestIsMyTurnFunction(t *testing.T) {
	// 创建存储实例
	store := storage.NewMemoryStorage("")
	err := store.Start()
	assert.NoError(t, err)
	defer store.Stop()

	// 创建3个验证者地址
	validator1 := types.Address{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}
	validator2 := types.Address{2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2}
	validator3 := types.Address{3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 3}

	// 创建共识配置（3个验证节点）
	config := &consensus.PoAConfig{
		Validators:    []types.Address{validator1, validator2, validator3},
		BlockInterval: 2000, // 2秒区块间隔
		RoundLength:   3,    // 3个区块为一轮
	}

	// 创建共识实例
	cons := consensus.NewPoAConsensus(config, store)
	assert.NotNil(t, cons)

	// 测试各种高度下各节点的出块权限
	testCases := []struct {
		height      uint64
		nodeAddress types.Address
		expected    bool
		description string
	}{
		{1, validator1, false, "高度1，节点1不应该出块"},
		{1, validator2, true, "高度1，节点2应该出块"},
		{1, validator3, false, "高度1，节点3不应该出块"},
		{2, validator1, false, "高度2，节点1不应该出块"},
		{2, validator2, false, "高度2，节点2不应该出块"},
		{2, validator3, true, "高度2，节点3应该出块"},
		{3, validator1, true, "高度3，节点1应该出块"},
		{3, validator2, false, "高度3，节点2不应该出块"},
		{3, validator3, false, "高度3，节点3不应该出块"},
		{4, validator1, false, "高度4，节点1不应该出块"},
		{4, validator2, true, "高度4，节点2应该出块"},
		{4, validator3, false, "高度4，节点3不应该出块"},
		{100, validator1, false, "高度100，节点1不应该出块"},
		{100, validator2, true, "高度100，节点2应该出块"},
		{100, validator3, false, "高度100，节点3不应该出块"},
		{101, validator1, false, "高度101，节点1不应该出块"},
		{101, validator2, false, "高度101，节点2不应该出块"},
		{101, validator3, true, "高度101，节点3应该出块"},
		{102, validator1, true, "高度102，节点1应该出块"},
		{102, validator2, false, "高度102，节点2不应该出块"},
		{102, validator3, false, "高度102，节点3不应该出块"},
	}

	for _, tc := range testCases {
		result := cons.IsMyTurn(tc.height, tc.nodeAddress)
		assert.Equal(t, tc.expected, result, tc.description)
	}
}
