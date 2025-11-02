package core

import (
	"github.com/lengzhao/govm/consensus"
	"github.com/lengzhao/govm/storage"
	"github.com/lengzhao/govm/types"
)

// Core 区块链核心模块接口
type Core interface {
	// Start 启动核心模块
	Start() error

	// Stop 停止核心模块
	Stop() error

	// GetConsensus 获取共识模块实例
	GetConsensus() consensus.PoAConsensus

	// GetStorage 获取存储模块实例
	GetStorage() storage.Storage

	// GetLastBlock 获取最新的区块
	GetLastBlock() *types.Block

	// GetHeight 获取当前区块高度
	GetHeight() uint64

	// GetBlockByHash 根据哈希获取区块
	GetBlockByHash(hash types.Hash) (*types.Block, error)

	// GetBlockByHeight 根据高度获取区块
	GetBlockByHeight(height uint64) (*types.Block, error)

	// AddBlock 添加新区块到区块链
	AddBlock(block *types.Block) error

	// ValidateTransaction 验证交易的有效性
	ValidateTransaction(tx *types.TransactionWithSign) error

	// ApplyTransaction 应用交易到状态
	ApplyTransaction(tx *types.TransactionWithSign) error

	// GetTransactionByHash 根据哈希获取交易
	GetTransactionByHash(hash types.Hash) (*types.TransactionWithSign, error)

	// CalculateBlockHash 计算区块哈希
	CalculateBlockHash(block *types.Block) types.Hash
}

// CoreConfig 核心模块配置
type CoreConfig struct {
	ShardID uint64               // 分片ID
	DataDir string               // 数据目录
	Genesis *types.GenesisConfig // 创世区块配置
}
