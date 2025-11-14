package core

import (
	"fmt"
	"sync"

	"github.com/lengzhao/binary"
	"github.com/lengzhao/govm/consensus"
	"github.com/lengzhao/govm/storage"
	"github.com/lengzhao/govm/types"
	"github.com/lengzhao/network"
)

// DefaultCore 默认核心模块实现
type DefaultCore struct {
	config      *CoreConfig
	blockchain  *Blockchain
	txProcessor TxProcessor
	consensus   consensus.PoAConsensus
	storage     storage.Storage
	network     network.NetworkInterface // 添加网络接口字段
	nodeID      int                      // 添加节点ID字段

	// 运行状态
	running bool
	mutex   sync.RWMutex
}

// NewCore 创建新的核心模块实例
func NewCore(config *CoreConfig, consensus consensus.PoAConsensus, storage storage.Storage) (Core, error) {
	// 创建区块链实例
	blockchain := NewBlockchain(storage, consensus)

	// 设置创世区块配置
	if config.Genesis != nil {
		blockchain.SetGenesisConfig(config.Genesis)
	}

	// 初始化区块链
	if err := blockchain.Init(); err != nil {
		return nil, fmt.Errorf("failed to initialize blockchain: %w", err)
	}

	// 创建交易处理器
	txProcessor := NewTxProcessor(storage)

	core := &DefaultCore{
		config:      config,
		blockchain:  blockchain,
		txProcessor: txProcessor,
		consensus:   consensus,
		storage:     storage,
		network:     nil, // 网络接口将在Start方法中设置
		nodeID:      1,   // 默认节点ID
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

	// 启动存储模块
	if err := c.storage.Start(); err != nil {
		return fmt.Errorf("failed to start storage: %w", err)
	}

	// 启动区块链
	// 区块链已经在NewCore中初始化了

	// 注册网络消息处理器（如果网络接口已设置）
	if c.network != nil {
		if err := c.registerNetworkHandlers(); err != nil {
			return fmt.Errorf("failed to register network handlers: %w", err)
		}
	}

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

	// 停止存储模块
	if err := c.storage.Stop(); err != nil {
		return fmt.Errorf("failed to stop storage: %w", err)
	}

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

// CalculateBlockHash 计算区块哈希
func (c *DefaultCore) CalculateBlockHash(block *types.Block) types.Hash {
	return c.blockchain.CalculateBlockHash(block)
}

// SetNetwork 设置网络接口并注册消息处理器
func (c *DefaultCore) SetNetwork(net network.NetworkInterface, nodeID int) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.network = net
	c.nodeID = nodeID

	// 如果模块已启动，立即注册网络消息处理器
	if c.running {
		return c.registerNetworkHandlers()
	}

	return nil
}

// registerNetworkHandlers 注册网络消息处理器
func (c *DefaultCore) registerNetworkHandlers() error {
	if c.network == nil {
		return fmt.Errorf("network interface is not set")
	}

	// 注册新区块消息处理器
	c.network.RegisterMessageHandler("new_block", func(from string, topic string, data []byte) error {
		fmt.Printf("Received new block from %s\n", from)
		var block types.Block
		if err := binary.Unmarshal(data, &block); err != nil {
			fmt.Printf("反序列化区块失败: %v\n", err)
			return err
		}
		fmt.Printf("Received block: %d\n", block.Header.BlockNumber)

		// 检查是否正在同步，如果正在同步则不处理新区块
		// 注意：这里需要通过其他方式获取syncer的状态，暂时简化处理
		// 在实际实现中，可能需要通过参数传递syncer或者通过其他方式获取同步状态

		err := c.AddBlock(&block)
		if err != nil {
			fmt.Printf("添加区块失败: %v\n", err)
			return err
		}

		// 这里应该处理接收到的新区块
		return nil
	})

	// 注册高度请求消息处理器
	c.network.RegisterRequestHandler("height_request", func(from string, topic string, data []byte) ([]byte, error) {
		fmt.Printf("Received height request from %s\n", from)

		// 获取当前节点的高度
		height := c.GetHeight()

		// 创建高度响应（使用与sync.HeightResponse相同的结构以保持兼容性）
		response := &struct {
			NodeID string // 响应节点ID
			Height uint64 // 区块链高度
			Error  string // 错误信息
		}{
			NodeID: fmt.Sprintf("node-%d", c.nodeID),
			Height: height,
			Error:  "",
		}

		// 序列化响应
		responseData, err := binary.Marshal(response)
		if err != nil {
			fmt.Printf("序列化高度响应失败: %v\n", err)
			return nil, err
		}

		return responseData, nil
	})

	return nil
}
