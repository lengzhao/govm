package core

import (
	"fmt"
	"log/slog"
	"sync"

	lzbinary "github.com/lengzhao/binary"
	"github.com/lengzhao/govm/consensus"
	"github.com/lengzhao/govm/crypto"
	"github.com/lengzhao/govm/storage"
	"github.com/lengzhao/govm/types"
)

// Blockchain 区块链核心结构
type Blockchain struct {
	storage   storage.Storage
	consensus consensus.PoAConsensus
	crypto    crypto.Crypto

	// 区块链状态
	lastBlock     *types.Block
	lastBlockHash types.Hash
	height        uint64

	// 创世区块配置
	genesisConfig *types.GenesisConfig

	// 读写锁保护区块链状态
	mutex sync.RWMutex
}

// NewBlockchain 创建新的区块链实例
func NewBlockchain(storage storage.Storage, consensus consensus.PoAConsensus) *Blockchain {
	return &Blockchain{
		storage:   storage,
		consensus: consensus,
		crypto:    crypto.NewCrypto(),
		height:    0,
	}
}

// Init 初始化区块链
func (bc *Blockchain) Init() error {
	bc.mutex.Lock()
	defer bc.mutex.Unlock()

	// 尝试加载创世区块
	genesisBlock, err := bc.loadGenesisBlock()
	if err != nil {
		// 如果没有创世区块，则创建一个
		genesisBlock, err = bc.createGenesisBlock()
		if err != nil {
			return fmt.Errorf("failed to create genesis block: %w", err)
		}
	}

	bc.lastBlock = genesisBlock
	bc.lastBlockHash = bc.calculateBlockHash(genesisBlock)
	bc.height = genesisBlock.Header.BlockNumber

	return nil
}

// SetGenesisConfig 设置创世区块配置
func (bc *Blockchain) SetGenesisConfig(config *types.GenesisConfig) {
	bc.genesisConfig = config
}

// GetLastBlock 获取最新的区块
func (bc *Blockchain) GetLastBlock() *types.Block {
	bc.mutex.RLock()
	defer bc.mutex.RUnlock()
	return bc.lastBlock
}

// GetHeight 获取当前区块高度
func (bc *Blockchain) GetHeight() uint64 {
	bc.mutex.RLock()
	defer bc.mutex.RUnlock()
	return bc.height
}

// GetBlockByHash 根据哈希获取区块
func (bc *Blockchain) GetBlockByHash(hash types.Hash) (*types.Block, error) {
	// 从存储中获取区块
	blockStorage, err := bc.storage.NewStorage(types.SNBlock)
	if err != nil {
		return nil, fmt.Errorf("failed to create block storage: %w", err)
	}

	key := hash[:]
	data, err := blockStorage.Get(key)
	if err != nil {
		return nil, fmt.Errorf("block not found: %w", err)
	}

	// 反序列化区块
	var block types.Block
	if err := lzbinary.Unmarshal(data, &block); err != nil {
		return nil, fmt.Errorf("failed to unmarshal block: %w", err)
	}

	return &block, nil
}

// GetBlockByHeight 根据高度获取区块
func (bc *Blockchain) GetBlockByHeight(height uint64) (*types.Block, error) {
	// 首先从存储中获取区块哈希
	statusStorage, err := bc.storage.NewStorage(types.SNStatus)
	if err != nil {
		return nil, fmt.Errorf("failed to create status storage: %w", err)
	}

	// 构造键名
	key := []byte(fmt.Sprintf("height_%d", height))
	hashData, err := statusStorage.Get(key)
	if err != nil {
		return nil, fmt.Errorf("block hash not found at height %d: %w", height, err)
	}

	// 将数据转换为哈希
	var hash types.Hash
	copy(hash[:], hashData)

	// 根据哈希获取区块
	return bc.GetBlockByHash(hash)
}

// AddBlock 添加新区块到区块链
func (bc *Blockchain) AddBlock(block *types.Block) error {
	bc.mutex.Lock()
	defer bc.mutex.Unlock()

	// 验证区块
	if err := bc.validateBlock(block); err != nil {
		return fmt.Errorf("block validation failed: %w", err)
	}
	slog.Info("adding block", "height", block.Header.BlockNumber, "Validator", block.Header.Validator)

	// 将区块存储到数据库
	if err := bc.storeBlock(block); err != nil {
		return fmt.Errorf("failed to store block: %w", err)
	}

	// 更新区块链状态
	bc.lastBlock = block
	bc.lastBlockHash = bc.calculateBlockHash(block)
	bc.height = block.Header.BlockNumber

	return nil
}

// validateBlock 验证区块
func (bc *Blockchain) validateBlock(block *types.Block) error {
	// 使用共识模块验证区块
	if err := bc.consensus.ValidateBlock(block); err != nil {
		return fmt.Errorf("consensus validation failed: %w", err)
	}

	// 验证区块高度
	if block.Header.BlockNumber > 0 && block.Header.BlockNumber != bc.height+1 {
		// 如果是添加历史区块，允许高度不连续
		if block.Header.BlockNumber <= bc.height {
			// 检查该高度的区块是否已存在
			_, err := bc.GetBlockByHeight(block.Header.BlockNumber)
			if err == nil {
				return fmt.Errorf("block at height %d already exists", block.Header.BlockNumber)
			}
		} else {
			return fmt.Errorf("block height mismatch: expected %d, got %d", bc.height+1, block.Header.BlockNumber)
		}
	}

	// 验证分片ID
	if block.Header.ShardID != types.DefaultShardID {
		return fmt.Errorf("block shard ID mismatch: expected %d, got %d", types.DefaultShardID, block.Header.ShardID)
	}

	// 验证时间戳
	if block.Header.Timestamp == 0 && block.Header.BlockNumber != 0 {
		return fmt.Errorf("block timestamp cannot be zero for non-genesis block")
	}

	// 验证Merkle根（如果区块包含交易）
	if err := bc.validateMerkleRoot(block); err != nil {
		return fmt.Errorf("merkle root validation failed: %w", err)
	}

	// 验证交易
	if err := bc.validateTransactions(block); err != nil {
		return fmt.Errorf("transaction validation failed: %w", err)
	}

	return nil
}

// validateMerkleRoot 验证Merkle根
func (bc *Blockchain) validateMerkleRoot(block *types.Block) error {
	// 如果区块不包含交易，Merkle根应该为空
	if len(block.Transactions) == 0 {
		if block.Header.MerkleRoot != [32]byte{} {
			return fmt.Errorf("merkle root should be empty for block without transactions")
		}
		return nil
	}

	// 计算交易的Merkle根
	merkleTree := NewMerkleTree(block.Transactions)
	calculatedRoot := merkleTree.GetRootHash()

	// 验证Merkle根是否匹配
	if block.Header.MerkleRoot != calculatedRoot {
		return fmt.Errorf("merkle root mismatch: expected %x, got %x", calculatedRoot, block.Header.MerkleRoot)
	}

	return nil
}

// validateTransactions 验证区块中的交易
func (bc *Blockchain) validateTransactions(block *types.Block) error {
	// 如果区块不包含交易，直接返回
	if len(block.Transactions) == 0 {
		return nil
	}

	// TODO: 实现交易验证逻辑
	// 这需要从存储中获取每个交易并验证其有效性
	return nil
}

// calculateBlockHash 计算区块哈希
func (bc *Blockchain) calculateBlockHash(block *types.Block) types.Hash {
	// 序列化区块头（排除签名字段）
	blockCopy := *block
	blockCopy.Header.Signature = nil

	data, err := lzbinary.Marshal(&blockCopy.Header.BlockHeader)
	if err != nil {
		// 如果序列化失败，返回空哈希
		return types.Hash{}
	}

	return bc.crypto.Hash(data)
}

// CalculateBlockHash 公共方法计算区块哈希
func (bc *Blockchain) CalculateBlockHash(block *types.Block) types.Hash {
	return bc.calculateBlockHash(block)
}

// loadGenesisBlock 加载创世区块
func (bc *Blockchain) loadGenesisBlock() (*types.Block, error) {
	// 尝试从存储中加载创世区块
	return bc.GetBlockByHeight(0)
}

// createGenesisBlock 创建创世区块
func (bc *Blockchain) createGenesisBlock() (*types.Block, error) {
	// 使用配置的时间戳，如果没有配置则使用默认值0
	timestamp := uint64(0)
	if bc.genesisConfig != nil {
		timestamp = bc.genesisConfig.Timestamp
	}

	// 创建创世区块
	genesis := &types.Block{
		Header: types.BlockHeaderWithSign{
			BlockHeader: types.BlockHeader{
				ShardID:       types.DefaultShardID,
				BlockNumber:   0,
				Timestamp:     timestamp,       // 使用配置的时间戳
				Validator:     types.Address{}, // 创世区块不需要验证者
				PrevHash:      types.Hash{},    // 创世区块的前一个哈希为空
				MerkleRoot:    types.Hash{},    // 创世区块的Merkle根为空
				StateRootHash: types.Hash{},    // 创世区块的状态根为空
				OtherShards:   [3]types.Hash{}, // 创世区块的相邻分片哈希为空
			},
			Signature: nil, // 创世区块不需要签名
		},
		Transactions: []types.Hash{}, // 创世区块不包含交易
	}

	// 存储创世区块
	if err := bc.storeBlock(genesis); err != nil {
		return nil, fmt.Errorf("failed to store genesis block: %w", err)
	}

	return genesis, nil
}

// storeBlock 存储区块到数据库
func (bc *Blockchain) storeBlock(block *types.Block) error {
	// 创建区块存储实例
	blockStorage, err := bc.storage.NewStorage(types.SNBlock)
	if err != nil {
		return fmt.Errorf("failed to create block storage: %w", err)
	}

	// 计算区块哈希
	blockHash := bc.calculateBlockHash(block)

	// 序列化区块
	data, err := lzbinary.Marshal(block)
	if err != nil {
		return fmt.Errorf("failed to marshal block: %w", err)
	}

	// 存储区块
	if err := blockStorage.Put(blockHash[:], data); err != nil {
		return fmt.Errorf("failed to store block: %w", err)
	}

	// 更新高度到哈希的映射
	statusStorage, err := bc.storage.NewStorage(types.SNStatus)
	if err != nil {
		return fmt.Errorf("failed to create status storage: %w", err)
	}

	// 存储高度到哈希的映射
	heightKey := []byte(fmt.Sprintf("height_%d", block.Header.BlockNumber))
	if err := statusStorage.Put(heightKey, blockHash[:]); err != nil {
		return fmt.Errorf("failed to store height mapping: %w", err)
	}

	return nil
}
