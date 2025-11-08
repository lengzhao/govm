package generator

import (
	"testing"
	"time"

	"github.com/lengzhao/govm/consensus"
	"github.com/lengzhao/govm/core"
	"github.com/lengzhao/govm/storage"
	"github.com/lengzhao/govm/txpool"
	"github.com/lengzhao/govm/types"
	"github.com/stretchr/testify/assert"
)

// TestBlockGeneratorWithTxPoolIntegration 测试区块生成器与交易池的集成
func TestBlockGeneratorWithTxPoolIntegration(t *testing.T) {
	// 创建存储实例
	store := storage.NewMemoryStorage("")
	err := store.Start()
	assert.NoError(t, err)
	defer store.Stop()

	// 创建验证者地址
	validatorAddr := types.Address{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}

	// 创建共识配置
	config := &consensus.PoAConfig{
		Validators:    []types.Address{validatorAddr},
		BlockInterval: 2000,
		RoundLength:   1,
	}

	// 创建共识实例
	cons := consensus.NewPoAConsensus(config, store)

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
	blockGenerator := NewBlockGenerator(cons, store, txPool)

	// 验证区块生成器创建成功
	assert.NotNil(t, blockGenerator)

	// 测试从空交易池选择交易
	transactions, err := blockGenerator.SelectTransactions()
	assert.NoError(t, err)
	assert.NotNil(t, transactions)
	assert.Equal(t, 0, len(transactions), "应该从空交易池中选择0个交易")
}

// TestBlockGenerationWithoutTransactions 测试不使用交易生成区块
func TestBlockGenerationWithoutTransactions(t *testing.T) {
	// 创建存储实例
	store := storage.NewMemoryStorage("")
	err := store.Start()
	assert.NoError(t, err)
	defer store.Stop()

	// 创建验证者地址
	validatorAddr := types.Address{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}

	// 创建共识配置
	config := &consensus.PoAConfig{
		Validators:    []types.Address{validatorAddr},
		BlockInterval: 2000,
		RoundLength:   1,
	}

	// 创建共识实例
	cons := consensus.NewPoAConsensus(config, store)

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
	blockGenerator := NewBlockGenerator(cons, store, txPool)

	// 生成新区块（没有交易）
	block, err := blockGenerator.GenerateBlock(nil)
	assert.NoError(t, err)
	assert.NotNil(t, block)

	// 验证区块不包含交易
	assert.Equal(t, 0, len(block.Transactions), "生成的区块不应该包含交易哈希")

	// 验证区块头信息
	assert.Equal(t, uint64(1), block.Header.BlockNumber)
	assert.Equal(t, types.DefaultShardID, block.Header.ShardID)
	assert.Equal(t, validatorAddr, block.Header.Validator)
}
