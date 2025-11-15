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

// TestMultiNodeBlockHeightConsistency 测试多节点区块高度一致性
func TestMultiNodeBlockHeightConsistency(t *testing.T) {
	// 创建3个节点的存储实例
	store1 := storage.NewMemoryStorage("")
	store2 := storage.NewMemoryStorage("")
	store3 := storage.NewMemoryStorage("")

	err1 := store1.Start()
	err2 := store2.Start()
	err3 := store3.Start()

	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.NoError(t, err3)

	defer store1.Stop()
	defer store2.Stop()
	defer store3.Stop()

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

	// 创建3个节点的核心模块
	cons1 := consensus.NewPoAConsensus(config, store1)
	cons2 := consensus.NewPoAConsensus(config, store2)
	cons3 := consensus.NewPoAConsensus(config, store3)

	// 创建核心模块配置
	coreConfig := &core.CoreConfig{
		ShardID: types.DefaultShardID,
		DataDir: "",
		Genesis: &types.GenesisConfig{
			Timestamp: uint64(time.Now().Unix()),
		},
	}

	// 创建3个节点的核心模块实例
	coreModule1, err := core.NewCore(coreConfig, cons1, store1)
	coreModule2, err2 := core.NewCore(coreConfig, cons2, store2)
	coreModule3, err3 := core.NewCore(coreConfig, cons3, store3)

	assert.NoError(t, err)
	assert.NoError(t, err2)
	assert.NoError(t, err3)

	assert.NotNil(t, coreModule1)
	assert.NotNil(t, coreModule2)
	assert.NotNil(t, coreModule3)

	// 启动所有核心模块
	err = coreModule1.Start()
	err2 = coreModule2.Start()
	err3 = coreModule3.Start()

	assert.NoError(t, err)
	assert.NoError(t, err2)
	assert.NoError(t, err3)

	defer coreModule1.Stop()
	defer coreModule2.Stop()
	defer coreModule3.Stop()

	// 验证所有节点初始高度一致
	height1 := coreModule1.GetHeight()
	height2 := coreModule2.GetHeight()
	height3 := coreModule3.GetHeight()

	assert.Equal(t, uint64(0), height1, "节点1初始高度应该为0")
	assert.Equal(t, uint64(0), height2, "节点2初始高度应该为0")
	assert.Equal(t, uint64(0), height3, "节点3初始高度应该为0")
	assert.Equal(t, height1, height2, "节点1和节点2初始高度应该一致")
	assert.Equal(t, height2, height3, "节点2和节点3初始高度应该一致")

	// 获取各节点的创世区块
	genesisBlock1, err := coreModule1.GetBlockByHeight(0)
	genesisBlock2, err2 := coreModule2.GetBlockByHeight(0)
	genesisBlock3, err3 := coreModule3.GetBlockByHeight(0)

	assert.NoError(t, err)
	assert.NoError(t, err2)
	assert.NoError(t, err3)

	assert.NotNil(t, genesisBlock1)
	assert.NotNil(t, genesisBlock2)
	assert.NotNil(t, genesisBlock3)

	// 验证创世区块一致性
	genesisHash1 := coreModule1.CalculateBlockHash(genesisBlock1)
	genesisHash2 := coreModule2.CalculateBlockHash(genesisBlock2)
	genesisHash3 := coreModule3.CalculateBlockHash(genesisBlock3)

	assert.Equal(t, genesisHash1, genesisHash2, "节点1和节点2创世区块哈希应该一致")
	assert.Equal(t, genesisHash2, genesisHash3, "节点2和节点3创世区块哈希应该一致")

	// 模拟节点1生成第一个区块（高度1，应该由验证者2出块，但我们简化测试）
	block1 := &types.Block{
		Header: types.BlockHeaderWithSign{
			BlockHeader: types.BlockHeader{
				ShardID:       types.DefaultShardID,
				BlockNumber:   1,
				Timestamp:     uint64(time.Now().UnixMilli()),
				Validator:     validator2, // 根据轮流出块规则，高度1应该由验证者2出块
				PrevHash:      genesisHash1,
				MerkleRoot:    types.Hash{},
				StateRootHash: types.Hash{},
				OtherShards:   [3]types.Hash{},
			},
			Signature: []byte{}, // 空签名表示空区块
		},
		Transactions: []types.Hash{},
	}

	// 将区块添加到所有节点
	err = coreModule1.AddBlock(block1)
	err2 = coreModule2.AddBlock(block1)
	err3 = coreModule3.AddBlock(block1)

	assert.NoError(t, err)
	assert.NoError(t, err2)
	assert.NoError(t, err3)

	// 验证所有节点高度更新一致
	height1 = coreModule1.GetHeight()
	height2 = coreModule2.GetHeight()
	height3 = coreModule3.GetHeight()

	assert.Equal(t, uint64(1), height1, "节点1高度应该更新为1")
	assert.Equal(t, uint64(1), height2, "节点2高度应该更新为1")
	assert.Equal(t, uint64(1), height3, "节点3高度应该更新为1")
	assert.Equal(t, height1, height2, "节点1和节点2高度应该一致")
	assert.Equal(t, height2, height3, "节点2和节点3高度应该一致")

	// 验证各节点区块内容一致性
	block1Node1, err := coreModule1.GetBlockByHeight(1)
	block1Node2, err2 := coreModule2.GetBlockByHeight(1)
	block1Node3, err3 := coreModule3.GetBlockByHeight(1)

	assert.NoError(t, err)
	assert.NoError(t, err2)
	assert.NoError(t, err3)

	assert.NotNil(t, block1Node1)
	assert.NotNil(t, block1Node2)
	assert.NotNil(t, block1Node3)

	// 验证区块内容一致性
	assert.Equal(t, block1Node1.Header.BlockNumber, block1Node2.Header.BlockNumber, "节点1和节点2区块高度应该一致")
	assert.Equal(t, block1Node2.Header.BlockNumber, block1Node3.Header.BlockNumber, "节点2和节点3区块高度应该一致")

	assert.Equal(t, block1Node1.Header.Validator, block1Node2.Header.Validator, "节点1和节点2区块验证者应该一致")
	assert.Equal(t, block1Node2.Header.Validator, block1Node3.Header.Validator, "节点2和节点3区块验证者应该一致")

	block1Hash1 := coreModule1.CalculateBlockHash(block1Node1)
	block1Hash2 := coreModule2.CalculateBlockHash(block1Node2)
	block1Hash3 := coreModule3.CalculateBlockHash(block1Node3)

	assert.Equal(t, block1Hash1, block1Hash2, "节点1和节点2区块哈希应该一致")
	assert.Equal(t, block1Hash2, block1Hash3, "节点2和节点3区块哈希应该一致")

	// 模拟节点2生成第二个区块（高度2，应该由验证者3出块）
	block2 := &types.Block{
		Header: types.BlockHeaderWithSign{
			BlockHeader: types.BlockHeader{
				ShardID:       types.DefaultShardID,
				BlockNumber:   2,
				Timestamp:     uint64(time.Now().UnixMilli()),
				Validator:     validator3,  // 根据轮流出块规则，高度2应该由验证者3出块
				PrevHash:      block1Hash1, // 指向前一个区块的哈希
				MerkleRoot:    types.Hash{},
				StateRootHash: types.Hash{},
				OtherShards:   [3]types.Hash{},
			},
			Signature: []byte{}, // 空签名表示空区块
		},
		Transactions: []types.Hash{},
	}

	// 将区块添加到所有节点
	err = coreModule1.AddBlock(block2)
	err2 = coreModule2.AddBlock(block2)
	err3 = coreModule3.AddBlock(block2)

	assert.NoError(t, err)
	assert.NoError(t, err2)
	assert.NoError(t, err3)

	// 验证所有节点高度再次更新一致
	height1 = coreModule1.GetHeight()
	height2 = coreModule2.GetHeight()
	height3 = coreModule3.GetHeight()

	assert.Equal(t, uint64(2), height1, "节点1高度应该更新为2")
	assert.Equal(t, uint64(2), height2, "节点2高度应该更新为2")
	assert.Equal(t, uint64(2), height3, "节点3高度应该更新为2")
	assert.Equal(t, height1, height2, "节点1和节点2高度应该一致")
	assert.Equal(t, height2, height3, "节点2和节点3高度应该一致")

	// 验证区块序列一致性
	for nodeIndex, coreModule := range []core.Core{coreModule1, coreModule2, coreModule3} {
		for height := uint64(0); height <= 2; height++ {
			block, err := coreModule.GetBlockByHeight(height)
			assert.NoError(t, err, "节点%d应该能够通过高度%d获取区块", nodeIndex+1, height)
			assert.NotNil(t, block, "节点%d高度%d的区块不应该为nil", nodeIndex+1, height)
			assert.Equal(t, height, block.Header.BlockNumber, "节点%d区块高度应该匹配", nodeIndex+1)

			// 验证区块哈希一致性
			blockHash := coreModule.CalculateBlockHash(block)
			assert.NotEqual(t, types.Hash{}, blockHash, "节点%d高度%d的区块哈希不应该为空", nodeIndex+1, height)
		}
	}
}

// TestBlockChainIntegrity 测试区块链完整性
func TestBlockChainIntegrity(t *testing.T) {
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

	// 获取创世区块
	genesisBlock := coreModule.GetLastBlock()
	genesisHash := coreModule.CalculateBlockHash(genesisBlock)

	// 生成一系列区块
	prevHash := genesisHash
	var blocks []*types.Block

	// 创建5个区块
	for i := uint64(1); i <= 5; i++ {
		newBlock := &types.Block{
			Header: types.BlockHeaderWithSign{
				BlockHeader: types.BlockHeader{
					ShardID:       types.DefaultShardID,
					BlockNumber:   i,
					Timestamp:     uint64(time.Now().UnixMilli()),
					Validator:     validator,
					PrevHash:      prevHash,
					MerkleRoot:    types.Hash{},
					StateRootHash: types.Hash{},
					OtherShards:   [3]types.Hash{},
				},
				Signature: []byte{}, // 空签名表示空区块
			},
			Transactions: []types.Hash{},
		}

		// 添加区块到区块链
		err = coreModule.AddBlock(newBlock)
		assert.NoError(t, err, "添加高度%d的区块时应该没有错误", i)

		// 保存区块信息
		blocks = append(blocks, newBlock)
		prevHash = coreModule.CalculateBlockHash(newBlock)
	}

	// 验证区块链完整性
	assert.Equal(t, uint64(5), coreModule.GetHeight(), "最终高度应该为5")

	// 验证每个区块的前向哈希链接
	genesisBlockFromChain, err := coreModule.GetBlockByHeight(0)
	assert.NoError(t, err)
	assert.Equal(t, genesisHash, coreModule.CalculateBlockHash(genesisBlockFromChain), "创世区块哈希应该一致")

	prevBlockHash := genesisHash
	for i := uint64(1); i <= 5; i++ {
		block, err := coreModule.GetBlockByHeight(i)
		assert.NoError(t, err, "应该能够获取高度%d的区块", i)
		assert.NotNil(t, block, "高度%d的区块不应该为nil", i)

		// 验证区块高度
		assert.Equal(t, i, block.Header.BlockNumber, "区块高度应该匹配")

		// 验证前向链接
		assert.Equal(t, prevBlockHash, block.Header.PrevHash, "高度%d区块的前向哈希应该指向前一个区块", i)

		// 更新前一个区块哈希
		prevBlockHash = coreModule.CalculateBlockHash(block)
	}

	// 验证最新区块
	lastBlock := coreModule.GetLastBlock()
	assert.Equal(t, uint64(5), lastBlock.Header.BlockNumber, "最新区块高度应该为5")
	lastBlockHash := coreModule.CalculateBlockHash(lastBlock)
	assert.Equal(t, prevBlockHash, lastBlockHash, "最新区块哈希应该与链式计算的哈希一致")
}
