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

// TestTenthBlockConsistency 测试第10个区块在多节点间的一致性
func TestTenthBlockConsistency(t *testing.T) {
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

	// 生成区块直到第10个区块
	prevHash := coreModule1.CalculateBlockHash(coreModule1.GetLastBlock())
	
	// 创建10个区块
	for i := uint64(1); i <= 10; i++ {
		// 根据轮流出块规则确定当前应该出块的验证者
		var currentValidator types.Address
		switch i % 3 {
		case 1:
			currentValidator = validator2 // 索引1
		case 2:
			currentValidator = validator3 // 索引2
		case 0:
			currentValidator = validator1 // 索引0
		}
		
		// 创建新区块
		newBlock := &types.Block{
			Header: types.BlockHeaderWithSign{
				BlockHeader: types.BlockHeader{
					ShardID:       types.DefaultShardID,
					BlockNumber:   i,
					Timestamp:     uint64(time.Now().UnixMilli()),
					Validator:     currentValidator,
					PrevHash:      prevHash,
					MerkleRoot:    types.Hash{},
					StateRootHash: types.Hash{},
					OtherShards:   [3]types.Hash{},
				},
				Signature: []byte{}, // 空签名表示空区块
			},
			Transactions: []types.Hash{},
		}

		// 将区块添加到所有节点
		err = coreModule1.AddBlock(newBlock)
		err2 = coreModule2.AddBlock(newBlock)
		err3 = coreModule3.AddBlock(newBlock)
		
		assert.NoError(t, err, "节点1添加高度%d的区块时应该没有错误", i)
		assert.NoError(t, err2, "节点2添加高度%d的区块时应该没有错误", i)
		assert.NoError(t, err3, "节点3添加高度%d的区块时应该没有错误", i)

		// 更新前一个区块哈希
		prevHash = coreModule1.CalculateBlockHash(newBlock)
	}

	// 验证所有节点高度更新一致
	height1 = coreModule1.GetHeight()
	height2 = coreModule2.GetHeight()
	height3 = coreModule3.GetHeight()
	
	assert.Equal(t, uint64(10), height1, "节点1高度应该更新为10")
	assert.Equal(t, uint64(10), height2, "节点2高度应该更新为10")
	assert.Equal(t, uint64(10), height3, "节点3高度应该更新为10")
	assert.Equal(t, height1, height2, "节点1和节点2高度应该一致")
	assert.Equal(t, height2, height3, "节点2和节点3高度应该一致")

	// 特别验证第10个区块的一致性
	tenthBlockNode1, err := coreModule1.GetBlockByHeight(10)
	tenthBlockNode2, err2 := coreModule2.GetBlockByHeight(10)
	tenthBlockNode3, err3 := coreModule3.GetBlockByHeight(10)
	
	assert.NoError(t, err, "节点1应该能够获取第10个区块")
	assert.NoError(t, err2, "节点2应该能够获取第10个区块")
	assert.NoError(t, err3, "节点3应该能够获取第10个区块")
	
	assert.NotNil(t, tenthBlockNode1, "节点1第10个区块不应该为nil")
	assert.NotNil(t, tenthBlockNode2, "节点2第10个区块不应该为nil")
	assert.NotNil(t, tenthBlockNode3, "节点3第10个区块不应该为nil")

	// 验证第10个区块内容一致性
	assert.Equal(t, uint64(10), tenthBlockNode1.Header.BlockNumber, "节点1第10个区块高度应该为10")
	assert.Equal(t, uint64(10), tenthBlockNode2.Header.BlockNumber, "节点2第10个区块高度应该为10")
	assert.Equal(t, uint64(10), tenthBlockNode3.Header.BlockNumber, "节点3第10个区块高度应该为10")

	// 根据轮流出块规则，高度10应该由验证者2出块（10 % 3 = 1，对应索引1，即验证者2）
	assert.Equal(t, validator2, tenthBlockNode1.Header.Validator, "节点1第10个区块验证者应该为验证者2")
	assert.Equal(t, validator2, tenthBlockNode2.Header.Validator, "节点2第10个区块验证者应该为验证者2")
	assert.Equal(t, validator2, tenthBlockNode3.Header.Validator, "节点3第10个区块验证者应该为验证者2")

	// 验证第10个区块哈希一致性
	tenthBlockHash1 := coreModule1.CalculateBlockHash(tenthBlockNode1)
	tenthBlockHash2 := coreModule2.CalculateBlockHash(tenthBlockNode2)
	tenthBlockHash3 := coreModule3.CalculateBlockHash(tenthBlockNode3)
	
	assert.Equal(t, tenthBlockHash1, tenthBlockHash2, "节点1和节点2第10个区块哈希应该一致")
	assert.Equal(t, tenthBlockHash2, tenthBlockHash3, "节点2和节点3第10个区块哈希应该一致")

	// 验证通过哈希获取区块的一致性
	tenthBlockByHashNode1, err := coreModule1.GetBlockByHash(tenthBlockHash1)
	tenthBlockByHashNode2, err2 := coreModule2.GetBlockByHash(tenthBlockHash2)
	tenthBlockByHashNode3, err3 := coreModule3.GetBlockByHash(tenthBlockHash3)
	
	assert.NoError(t, err, "节点1应该能够通过哈希获取第10个区块")
	assert.NoError(t, err2, "节点2应该能够通过哈希获取第10个区块")
	assert.NoError(t, err3, "节点3应该能够通过哈希获取第10个区块")
	
	assert.NotNil(t, tenthBlockByHashNode1, "节点1通过哈希获取的第10个区块不应该为nil")
	assert.NotNil(t, tenthBlockByHashNode2, "节点2通过哈希获取的第10个区块不应该为nil")
	assert.NotNil(t, tenthBlockByHashNode3, "节点3通过哈希获取的第10个区块不应该为nil")

	// 验证通过不同方式获取的区块一致性
	assert.Equal(t, tenthBlockNode1.Header.BlockNumber, tenthBlockByHashNode1.Header.BlockNumber, "节点1通过高度和哈希获取的第10个区块高度应该一致")
	assert.Equal(t, tenthBlockNode2.Header.BlockNumber, tenthBlockByHashNode2.Header.BlockNumber, "节点2通过高度和哈希获取的第10个区块高度应该一致")
	assert.Equal(t, tenthBlockNode3.Header.BlockNumber, tenthBlockByHashNode3.Header.BlockNumber, "节点3通过高度和哈希获取的第10个区块高度应该一致")

	// 验证所有区块序列的一致性
	for nodeIndex, coreModule := range []core.Core{coreModule1, coreModule2, coreModule3} {
		for height := uint64(0); height <= 10; height++ {
			block, err := coreModule.GetBlockByHeight(height)
			assert.NoError(t, err, "节点%d应该能够通过高度%d获取区块", nodeIndex+1, height)
			assert.NotNil(t, block, "节点%d高度%d的区块不应该为nil", nodeIndex+1, height)
			assert.Equal(t, height, block.Header.BlockNumber, "节点%d区块高度应该匹配", nodeIndex+1)

			// 验证区块哈希一致性
			blockHash := coreModule.CalculateBlockHash(block)
			assert.NotEqual(t, types.Hash{}, blockHash, "节点%d高度%d的区块哈希不应该为空", nodeIndex+1, height)
			
			// 验证通过哈希获取的区块一致性
			blockByHash, err := coreModule.GetBlockByHash(blockHash)
			assert.NoError(t, err, "节点%d应该能够通过哈希获取高度%d的区块", nodeIndex+1, height)
			assert.NotNil(t, blockByHash, "节点%d通过哈希获取的高度%d的区块不应该为nil", nodeIndex+1, height)
			assert.Equal(t, block.Header.BlockNumber, blockByHash.Header.BlockNumber, "节点%d通过高度和哈希获取的区块高度应该一致", nodeIndex+1)
		}
	}
}

// TestTenthBlockContentConsistency 测试第10个区块内容的一致性
func TestTenthBlockContentConsistency(t *testing.T) {
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

	// 生成区块直到第10个区块
	prevHash := genesisHash
	
	// 创建10个区块
	var blocks []*types.Block
	for i := uint64(1); i <= 10; i++ {
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

	// 验证区块链高度
	assert.Equal(t, uint64(10), coreModule.GetHeight(), "最终高度应该为10")

	// 特别验证第10个区块的内容
	tenthBlock, err := coreModule.GetBlockByHeight(10)
	assert.NoError(t, err, "应该能够获取第10个区块")
	assert.NotNil(t, tenthBlock, "第10个区块不应该为nil")

	// 验证第10个区块的基本信息
	assert.Equal(t, uint64(10), tenthBlock.Header.BlockNumber, "第10个区块高度应该为10")
	assert.Equal(t, validator, tenthBlock.Header.Validator, "第10个区块验证者应该正确")
	assert.Equal(t, types.DefaultShardID, tenthBlock.Header.ShardID, "第10个区块分片ID应该正确")

	// 验证第10个区块的时间戳
	currentTime := uint64(time.Now().UnixMilli())
	assert.True(t, tenthBlock.Header.Timestamp <= currentTime, "第10个区块时间戳应该不大于当前时间")
	assert.True(t, tenthBlock.Header.Timestamp > 0, "第10个区块时间戳应该大于0")

	// 验证第10个区块的前向链接
	ninthBlock, err := coreModule.GetBlockByHeight(9)
	assert.NoError(t, err, "应该能够获取第9个区块")
	ninthBlockHash := coreModule.CalculateBlockHash(ninthBlock)
	assert.Equal(t, ninthBlockHash, tenthBlock.Header.PrevHash, "第10个区块的前向哈希应该指向前一个区块")

	// 验证第10个区块的哈希一致性
	tenthBlockHash1 := coreModule.CalculateBlockHash(tenthBlock)
	tenthBlockHash2 := coreModule.CalculateBlockHash(tenthBlock)
	assert.Equal(t, tenthBlockHash1, tenthBlockHash2, "同一区块多次计算的哈希应该一致")

	// 验证通过哈希获取的区块一致性
	tenthBlockByHash, err := coreModule.GetBlockByHash(tenthBlockHash1)
	assert.NoError(t, err, "应该能够通过哈希获取第10个区块")
	assert.NotNil(t, tenthBlockByHash, "通过哈希获取的第10个区块不应该为nil")

	// 验证内容一致性
	assert.Equal(t, tenthBlock.Header.BlockNumber, tenthBlockByHash.Header.BlockNumber, "通过高度和哈希获取的第10个区块高度应该一致")
	assert.Equal(t, tenthBlock.Header.Validator, tenthBlockByHash.Header.Validator, "通过高度和哈希获取的第10个区块验证者应该一致")
	assert.Equal(t, tenthBlock.Header.PrevHash, tenthBlockByHash.Header.PrevHash, "通过高度和哈希获取的第10个区块前向哈希应该一致")
	assert.Equal(t, tenthBlock.Header.Timestamp, tenthBlockByHash.Header.Timestamp, "通过高度和哈希获取的第10个区块时间戳应该一致")

	// 验证最新区块就是第10个区块
	lastBlock := coreModule.GetLastBlock()
	lastBlockHash := coreModule.CalculateBlockHash(lastBlock)
	assert.Equal(t, tenthBlockHash1, lastBlockHash, "最新区块哈希应该与第10个区块哈希一致")
	assert.Equal(t, uint64(10), lastBlock.Header.BlockNumber, "最新区块高度应该为10")
}