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

// TestBlockHeightConsistency 测试区块高度一致性
func TestBlockHeightConsistency(t *testing.T) {
	// 创建存储实例
	store := storage.NewMemoryStorage("")
	err := store.Start()
	assert.NoError(t, err)
	defer store.Stop()

	// 创建验证者地址
	validator := types.Address{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}

	// 创建共识配置
	config := &consensus.PoAConfig{
		Validators:    []types.Address{validator},
		BlockInterval: 100, // 100毫秒区块间隔（测试用）
		RoundLength:   1,   // 1个区块为一轮
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
	initialHeight := coreModule.GetHeight()
	assert.Equal(t, uint64(0), initialHeight) // 创世区块高度为0

	// 获取创世区块
	genesisBlock, err := coreModule.GetBlockByHeight(0)
	assert.NoError(t, err)
	assert.NotNil(t, genesisBlock)
	assert.Equal(t, uint64(0), genesisBlock.Header.BlockNumber)

	// 验证创世区块哈希一致性
	genesisHash1 := coreModule.CalculateBlockHash(genesisBlock)
	genesisHash2 := coreModule.CalculateBlockHash(genesisBlock)
	assert.Equal(t, genesisHash1, genesisHash2, "同一区块的哈希应该一致")

	// 验证获取最新区块
	lastBlock := coreModule.GetLastBlock()
	assert.Equal(t, genesisBlock.Header.BlockNumber, lastBlock.Header.BlockNumber, "最新区块高度应该等于创世区块高度")
	assert.Equal(t, genesisBlock.Header.Validator, lastBlock.Header.Validator, "最新区块验证者应该等于创世区块验证者")
}

// TestNewBlockGeneration 测试新区块生成
func TestNewBlockGeneration(t *testing.T) {
	// 创建存储实例
	store := storage.NewMemoryStorage("")
	err := store.Start()
	assert.NoError(t, err)
	defer store.Stop()

	// 创建验证者地址
	validator := types.Address{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}

	// 创建共识配置
	config := &consensus.PoAConfig{
		Validators:    []types.Address{validator},
		BlockInterval: 100, // 100毫秒区块间隔（测试用）
		RoundLength:   1,   // 1个区块为一轮
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

	// 验证初始高度
	initialHeight := coreModule.GetHeight()
	assert.Equal(t, uint64(0), initialHeight)

	// 创建新区块
	genesisBlock := coreModule.GetLastBlock()
	newBlock := &types.Block{
		Header: types.BlockHeaderWithSign{
			BlockHeader: types.BlockHeader{
				ShardID:       types.DefaultShardID,
				BlockNumber:   1,
				Timestamp:     uint64(time.Now().UnixMilli()),
				Validator:     validator,
				PrevHash:      coreModule.CalculateBlockHash(genesisBlock),
				MerkleRoot:    types.Hash{},
				StateRootHash: types.Hash{},
				OtherShards:   [3]types.Hash{},
			},
			Signature: []byte{}, // 空签名表示空区块
		},
		Transactions: []types.Hash{},
	}

	// 验证区块结构
	assert.Equal(t, uint64(1), newBlock.Header.BlockNumber)
	assert.Equal(t, validator, newBlock.Header.Validator)

	// 添加新区块到区块链
	err = coreModule.AddBlock(newBlock)
	assert.NoError(t, err)

	// 验证高度更新
	updatedHeight := coreModule.GetHeight()
	assert.Equal(t, uint64(1), updatedHeight, "添加新区块后高度应该更新")

	// 验证最新区块
	latestBlock := coreModule.GetLastBlock()
	assert.Equal(t, newBlock.Header.BlockNumber, latestBlock.Header.BlockNumber, "最新区块高度应该匹配")
	assert.Equal(t, newBlock.Header.Validator, latestBlock.Header.Validator, "最新区块验证者应该匹配")

	// 通过高度获取区块
	blockByHeight, err := coreModule.GetBlockByHeight(1)
	assert.NoError(t, err)
	assert.NotNil(t, blockByHeight)
	assert.Equal(t, uint64(1), blockByHeight.Header.BlockNumber)

	// 验证区块哈希一致性
	blockHash1 := coreModule.CalculateBlockHash(blockByHeight)
	blockHash2 := coreModule.CalculateBlockHash(blockByHeight)
	assert.Equal(t, blockHash1, blockHash2, "同一区块的哈希应该一致")

	// 创建第二个新区块
	secondBlock := &types.Block{
		Header: types.BlockHeaderWithSign{
			BlockHeader: types.BlockHeader{
				ShardID:       types.DefaultShardID,
				BlockNumber:   2,
				Timestamp:     uint64(time.Now().UnixMilli()),
				Validator:     validator,
				PrevHash:      coreModule.CalculateBlockHash(newBlock),
				MerkleRoot:    types.Hash{},
				StateRootHash: types.Hash{},
				OtherShards:   [3]types.Hash{},
			},
			Signature: []byte{}, // 空签名表示空区块
		},
		Transactions: []types.Hash{},
	}

	// 添加第二个区块
	err = coreModule.AddBlock(secondBlock)
	assert.NoError(t, err)

	// 验证高度再次更新
	finalHeight := coreModule.GetHeight()
	assert.Equal(t, uint64(2), finalHeight, "添加第二个区块后高度应该更新到2")

	// 验证区块序列一致性
	for height := uint64(0); height <= 2; height++ {
		block, err := coreModule.GetBlockByHeight(height)
		assert.NoError(t, err, "应该能够通过高度%d获取区块", height)
		assert.NotNil(t, block, "高度%d的区块不应该为nil", height)
		assert.Equal(t, height, block.Header.BlockNumber, "区块高度应该匹配")

		// 验证区块哈希一致性
		blockHash := coreModule.CalculateBlockHash(block)
		assert.NotEqual(t, types.Hash{}, blockHash, "区块哈希不应该为空")
	}
}

// TestMultipleBlocksConsistency 测试多个区块的一致性
func TestMultipleBlocksConsistency(t *testing.T) {
	// 创建存储实例
	store := storage.NewMemoryStorage("")
	err := store.Start()
	assert.NoError(t, err)
	defer store.Stop()

	// 创建验证者地址
	validator := types.Address{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}

	// 创建共识配置
	config := &consensus.PoAConfig{
		Validators:    []types.Address{validator},
		BlockInterval: 100, // 100毫秒区块间隔（测试用）
		RoundLength:   1,   // 1个区块为一轮
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

	// 验证初始高度
	assert.Equal(t, uint64(0), coreModule.GetHeight())

	// 生成多个区块
	prevBlock := coreModule.GetLastBlock()
	for i := uint64(1); i <= 5; i++ {
		newBlock := &types.Block{
			Header: types.BlockHeaderWithSign{
				BlockHeader: types.BlockHeader{
					ShardID:       types.DefaultShardID,
					BlockNumber:   i,
					Timestamp:     uint64(time.Now().UnixMilli()),
					Validator:     validator,
					PrevHash:      coreModule.CalculateBlockHash(prevBlock),
					MerkleRoot:    types.Hash{},
					StateRootHash: types.Hash{},
					OtherShards:   [3]types.Hash{},
				},
				Signature: []byte{}, // 空签名表示空区块
			},
			Transactions: []types.Hash{},
		}

		// 添加区块
		err = coreModule.AddBlock(newBlock)
		assert.NoError(t, err, "添加高度%d的区块时应该没有错误", i)

		// 更新前一个区块
		prevBlock = newBlock
	}

	// 验证最终高度
	assert.Equal(t, uint64(5), coreModule.GetHeight(), "最终高度应该为5")

	// 验证所有区块的一致性
	for height := uint64(0); height <= 5; height++ {
		// 通过高度获取区块
		blockByHeight, err := coreModule.GetBlockByHeight(height)
		assert.NoError(t, err, "应该能够通过高度%d获取区块", height)
		assert.NotNil(t, blockByHeight, "高度%d的区块不应该为nil", height)
		assert.Equal(t, height, blockByHeight.Header.BlockNumber, "区块高度应该匹配")

		// 验证通过哈希获取区块的一致性
		blockHash := coreModule.CalculateBlockHash(blockByHeight)
		blockByHash, err := coreModule.GetBlockByHash(blockHash)
		assert.NoError(t, err, "应该能够通过哈希获取高度%d的区块", height)
		assert.NotNil(t, blockByHash, "通过哈希获取的高度%d的区块不应该为nil", height)
		assert.Equal(t, blockByHeight.Header.BlockNumber, blockByHash.Header.BlockNumber, "通过高度和哈希获取的区块高度应该一致")
		assert.Equal(t, blockByHeight.Header.Validator, blockByHash.Header.Validator, "通过高度和哈希获取的区块验证者应该一致")

		// 验证区块哈希一致性
		hash1 := coreModule.CalculateBlockHash(blockByHeight)
		hash2 := coreModule.CalculateBlockHash(blockByHash)
		assert.Equal(t, hash1, hash2, "同一区块通过不同方式获取后计算的哈希应该一致")
	}

	// 验证最新区块
	lastBlock := coreModule.GetLastBlock()
	assert.Equal(t, uint64(5), lastBlock.Header.BlockNumber, "最新区块高度应该为5")
}
