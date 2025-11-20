package core

import (
	"testing"
	"time"

	"github.com/lengzhao/govm/consensus"
	"github.com/lengzhao/govm/storage"
	"github.com/lengzhao/govm/types"
	"github.com/stretchr/testify/assert"
)

func TestBlockchain_Init(t *testing.T) {
	// 创建存储
	store := storage.NewMemoryStorage("")
	err := store.Start()
	assert.NoError(t, err)
	defer store.Stop()

	// 创建共识模块
	config := &consensus.PoAConfig{
		Validators:    []types.Address{{1}, {2}, {3}},
		BlockInterval: 2000,
		RoundLength:   3,
	}
	cons := consensus.NewPoAConsensus(config, store)

	// 创建区块链
	blockchain := NewBlockchain(store, cons)

	// 初始化区块链
	err = blockchain.Init()
	assert.NoError(t, err)

	// 验证创世区块已创建
	lastBlock := blockchain.GetLastBlock()
	assert.NotNil(t, lastBlock)
	assert.Equal(t, uint64(0), lastBlock.Header.BlockNumber)
	assert.Equal(t, types.DefaultShardID, lastBlock.Header.ShardID)
}

func TestBlockchain_AddBlock(t *testing.T) {
	// 创建存储
	store := storage.NewMemoryStorage("")
	err := store.Start()
	assert.NoError(t, err)
	defer store.Stop()

	// 创建共识模块
	// 注意：这里我们将验证者设置为与区块中的验证者匹配
	config := &consensus.PoAConfig{
		Validators:    []types.Address{{1}}, // 只有一个验证者，与区块中的验证者匹配
		BlockInterval: 2000,
		RoundLength:   3,
	}
	cons := consensus.NewPoAConsensus(config, store)

	// 创建区块链
	blockchain := NewBlockchain(store, cons)

	// 初始化区块链
	err = blockchain.Init()
	assert.NoError(t, err)

	// 创建空区块（没有签名，避免签名验证问题）
	newBlock := &types.Block{
		Header: types.BlockHeaderWithSign{
			BlockHeader: types.BlockHeader{
				ShardID:       types.DefaultShardID,
				BlockNumber:   1,
				Timestamp:     uint64(time.Now().UnixMilli()),
				Validator:     types.Address{1}, // 与配置中的验证者匹配
				PrevHash:      blockchain.lastBlockHash,
				MerkleRoot:    types.Hash{},
				StateRootHash: types.Hash{},
				OtherShards:   [3]types.Hash{},
			},
			Signature: nil, // 空签名，创建空区块
		},
		Transactions: []types.Hash{},
	}

	// 添加区块
	err = blockchain.AddBlock(newBlock)
	assert.NoError(t, err)

	// 验证区块已添加
	lastBlock := blockchain.GetLastBlock()
	assert.NotNil(t, lastBlock)
	assert.Equal(t, uint64(1), lastBlock.Header.BlockNumber)
	assert.Equal(t, newBlock.Header.BlockNumber, lastBlock.Header.BlockNumber)
}

func TestBlockchain_GetBlockByHash(t *testing.T) {
	// 创建存储
	store := storage.NewMemoryStorage("")
	err := store.Start()
	assert.NoError(t, err)
	defer store.Stop()

	// 创建共识模块
	config := &consensus.PoAConfig{
		Validators:    []types.Address{{1}, {2}, {3}},
		BlockInterval: 2000,
		RoundLength:   3,
	}
	cons := consensus.NewPoAConsensus(config, store)

	// 创建区块链
	blockchain := NewBlockchain(store, cons)

	// 初始化区块链
	err = blockchain.Init()
	assert.NoError(t, err)

	// 获取创世区块哈希
	genesisHash := blockchain.lastBlockHash

	// 通过哈希获取区块
	block, err := blockchain.GetBlockByHash(genesisHash)
	assert.NoError(t, err)
	assert.NotNil(t, block)
	assert.Equal(t, uint64(0), block.Header.BlockNumber)
}

func TestBlockchain_GetBlockByHeight(t *testing.T) {
	// 创建存储
	store := storage.NewMemoryStorage("")
	err := store.Start()
	assert.NoError(t, err)
	defer store.Stop()

	// 创建共识模块
	config := &consensus.PoAConfig{
		Validators:    []types.Address{{1}, {2}, {3}},
		BlockInterval: 2000,
		RoundLength:   3,
	}
	cons := consensus.NewPoAConsensus(config, store)

	// 创建区块链
	blockchain := NewBlockchain(store, cons)

	// 初始化区块链
	err = blockchain.Init()
	assert.NoError(t, err)

	// 通过高度获取区块
	block, err := blockchain.GetBlockByHeight(0)
	assert.NoError(t, err)
	assert.NotNil(t, block)
	assert.Equal(t, uint64(0), block.Header.BlockNumber)
}

func TestBlockchain_CalculateBlockHash(t *testing.T) {
	// 创建存储
	store := storage.NewMemoryStorage("")
	err := store.Start()
	assert.NoError(t, err)
	defer store.Stop()

	// 创建共识模块
	config := &consensus.PoAConfig{
		Validators:    []types.Address{{1}}, // 只有一个验证者，与区块中的验证者匹配
		BlockInterval: 2000,
		RoundLength:   3,
	}
	cons := consensus.NewPoAConsensus(config, store)

	// 创建区块链
	blockchain := NewBlockchain(store, cons)

	// 初始化区块链
	err = blockchain.Init()
	assert.NoError(t, err)

	// 创建测试区块（空区块，没有签名）
	testBlock := &types.Block{
		Header: types.BlockHeaderWithSign{
			BlockHeader: types.BlockHeader{
				ShardID:       types.DefaultShardID,
				BlockNumber:   1,
				Timestamp:     uint64(time.Now().UnixMilli()),
				Validator:     types.Address{1}, // 与配置中的验证者匹配
				PrevHash:      blockchain.lastBlockHash,
				MerkleRoot:    types.Hash{},
				StateRootHash: types.Hash{},
				OtherShards:   [3]types.Hash{},
			},
			Signature: nil, // 空签名
		},
		Transactions: []types.Hash{},
	}

	// 计算区块哈希
	hash := blockchain.CalculateBlockHash(testBlock)
	assert.NotEqual(t, types.Hash{}, hash)

	// 验证哈希一致性
	hash2 := blockchain.CalculateBlockHash(testBlock)
	assert.Equal(t, hash, hash2)
}

func TestBlockchain_ValidateBlock(t *testing.T) {
	// 创建存储
	store := storage.NewMemoryStorage("")
	err := store.Start()
	assert.NoError(t, err)
	defer store.Stop()

	// 创建共识模块
	config := &consensus.PoAConfig{
		Validators:    []types.Address{{1}, {2}, {3}},
		BlockInterval: 2000,
		RoundLength:   3,
	}
	cons := consensus.NewPoAConsensus(config, store)

	// 创建区块链
	blockchain := NewBlockchain(store, cons)

	// 初始化区块链
	err = blockchain.Init()
	assert.NoError(t, err)

	// 创建有效的区块
	validBlock := &types.Block{
		Header: types.BlockHeaderWithSign{
			BlockHeader: types.BlockHeader{
				ShardID:       types.DefaultShardID,
				BlockNumber:   1,
				Timestamp:     uint64(time.Now().UnixMilli()),
				Validator:     types.Address{1},
				PrevHash:      blockchain.lastBlockHash,
				MerkleRoot:    types.Hash{},
				StateRootHash: types.Hash{},
				OtherShards:   [3]types.Hash{},
			},
			Signature: []byte{}, // 空区块签名
		},
		Transactions: []types.Hash{},
	}

	// 验证有效区块
	err = blockchain.validateBlock(validBlock)
	assert.NoError(t, err)

	// 创建无效区块（错误的分片ID）
	invalidBlock := &types.Block{
		Header: types.BlockHeaderWithSign{
			BlockHeader: types.BlockHeader{
				ShardID:       999, // 错误的分片ID
				BlockNumber:   1,
				Timestamp:     uint64(time.Now().UnixMilli()),
				Validator:     types.Address{1},
				PrevHash:      blockchain.lastBlockHash,
				MerkleRoot:    types.Hash{},
				StateRootHash: types.Hash{},
				OtherShards:   [3]types.Hash{},
			},
			Signature: []byte{}, // 空区块签名
		},
		Transactions: []types.Hash{},
	}

	// 验证无效区块
	err = blockchain.validateBlock(invalidBlock)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "block shard ID mismatch")
}
