package core

import (
	"fmt"
	"sync"

	"github.com/lengzhao/govm/consensus"
	"github.com/lengzhao/govm/storage"
	"github.com/lengzhao/govm/types"
)

// DefaultCore 默认核心模块实现
type DefaultCore struct {
	config      *CoreConfig
	blockchain  *Blockchain
	txProcessor TxProcessor
	consensus   consensus.PoAConsensus
	storage     storage.Storage

	// 运行状态
	running bool
	mutex   sync.RWMutex
}

// NewCore 创建新的核心模块实例
func NewCore(config *CoreConfig, consensus consensus.PoAConsensus, storage storage.Storage) (Core, error) {
	// 创建区块链实例
	blockchain := NewBlockchain(storage, consensus)

	// 初始化区块链
	if err := blockchain.Init(); err != nil {
		return nil, fmt.Errorf("failed to initialize blockchain: %w", err)
	}

	// 创建交易处理器
	txProcessor := NewTxProcessor()

	core := &DefaultCore{
		config:      config,
		blockchain:  blockchain,
		txProcessor: txProcessor,
		consensus:   consensus,
		storage:     storage,
		running:     false,
	}

	return core, nil
}

// Start 启动核心模块
func (c *DefaultCore) Start() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.running {
		return fmt.Errorf("core module is already running")
	}

	// 启动区块链
	// 区块链已经在NewCore中初始化了

	// 启动其他组件
	// TODO: 启动其他需要的组件

	c.running = true
	return nil
}

// Stop 停止核心模块
func (c *DefaultCore) Stop() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if !c.running {
		return fmt.Errorf("core module is not running")
	}

	// 停止各组件
	// TODO: 停止各组件

	c.running = false
	return nil
}

// GetConsensus 获取共识模块实例
func (c *DefaultCore) GetConsensus() consensus.PoAConsensus {
	return c.consensus
}

// GetStorage 获取存储模块实例
func (c *DefaultCore) GetStorage() storage.Storage {
	return c.storage
}

// GetLastBlock 获取最新的区块
func (c *DefaultCore) GetLastBlock() *types.Block {
	return c.blockchain.GetLastBlock()
}

// GetHeight 获取当前区块高度
func (c *DefaultCore) GetHeight() uint64 {
	return c.blockchain.GetHeight()
}

// GetBlockByHash 根据哈希获取区块
func (c *DefaultCore) GetBlockByHash(hash types.Hash) (*types.Block, error) {
	return c.blockchain.GetBlockByHash(hash)
}

// GetBlockByHeight 根据高度获取区块
func (c *DefaultCore) GetBlockByHeight(height uint64) (*types.Block, error) {
	return c.blockchain.GetBlockByHeight(height)
}

// AddBlock 添加新区块到区块链
func (c *DefaultCore) AddBlock(block *types.Block) error {
	return c.blockchain.AddBlock(block)
}

// ValidateTransaction 验证交易的有效性
func (c *DefaultCore) ValidateTransaction(tx *types.TransactionWithSign) error {
	return c.txProcessor.ValidateTransaction(tx)
}

// ApplyTransaction 应用交易到状态
func (c *DefaultCore) ApplyTransaction(tx *types.TransactionWithSign) error {
	return c.txProcessor.ApplyTransaction(tx)
}

// GetTransactionByHash 根据哈希获取交易
func (c *DefaultCore) GetTransactionByHash(hash types.Hash) (*types.TransactionWithSign, error) {
	return c.txProcessor.GetTransactionByHash(hash)
}
