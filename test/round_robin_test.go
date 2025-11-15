package test

import (
	"testing"
	"time"

	"github.com/lengzhao/govm/consensus"
	"github.com/lengzhao/govm/core"
	"github.com/lengzhao/govm/storage"
	"github.com/lengzhao/govm/types"
	"github.com/stretchr/testify/assert"
)

// TestRoundRobinBlockGeneration 测试3节点轮流出块功能
func TestRoundRobinBlockGeneration(t *testing.T) {
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

	// 验证初始状态
	lastBlock := coreModule.GetLastBlock()
	assert.NotNil(t, lastBlock)
	assert.Equal(t, uint64(0), lastBlock.Header.BlockNumber) // 创世区块高度为0

	// 验证共识机制的轮流出块逻辑
	// 测试高度1应该由验证者2出块（索引1）
	validatorAtHeight1 := cons.GetCurrentValidator(1)
	assert.Equal(t, validator2, validatorAtHeight1)

	// 测试高度2应该由验证者3出块（索引2）
	validatorAtHeight2 := cons.GetCurrentValidator(2)
	assert.Equal(t, validator3, validatorAtHeight2)

	// 测试高度3应该由验证者1出块（索引0）
	validatorAtHeight3 := cons.GetCurrentValidator(3)
	assert.Equal(t, validator1, validatorAtHeight3)

	// 测试高度4应该由验证者2出块（索引1，新一轮开始）
	validatorAtHeight4 := cons.GetCurrentValidator(4)
	assert.Equal(t, validator2, validatorAtHeight4)

	// 验证IsMyTurn函数
	assert.True(t, cons.IsMyTurn(1, validator2))
	assert.True(t, cons.IsMyTurn(2, validator3))
	assert.True(t, cons.IsMyTurn(3, validator1))
	assert.True(t, cons.IsMyTurn(4, validator2))
	assert.False(t, cons.IsMyTurn(1, validator1))
	assert.False(t, cons.IsMyTurn(2, validator1))
	assert.False(t, cons.IsMyTurn(3, validator2))
}

// TestRoundRobinWithMultipleBlocks 测试生成多个区块并验证轮流出块
func TestRoundRobinWithMultipleBlocks(t *testing.T) {
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

	// 验证创世区块
	lastBlock := coreModule.GetLastBlock()
	assert.NotNil(t, lastBlock)
	assert.Equal(t, uint64(0), lastBlock.Header.BlockNumber)

	// 验证前10个区块的出块顺序
	expectedValidators := []types.Address{
		validator2, // 高度1（索引1）
		validator3, // 高度2（索引2）
		validator1, // 高度3（索引0）
		validator2, // 高度4（索引1）
		validator3, // 高度5（索引2）
		validator1, // 高度6（索引0）
		validator2, // 高度7（索引1）
		validator3, // 高度8（索引2）
		validator1, // 高度9（索引0）
		validator2, // 高度10（索引1）
	}

	// 验证每个高度应该由哪个验证者出块
	for height := uint64(1); height <= 10; height++ {
		expectedValidator := expectedValidators[height-1]
		actualValidator := cons.GetCurrentValidator(height)
		assert.Equal(t, expectedValidator, actualValidator, "高度%d应该由验证者%x出块", height, expectedValidator)
	}
}
