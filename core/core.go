package core

import (
	"github.com/lengzhao/govm/consensus"
	"github.com/lengzhao/govm/storage"
	"github.com/lengzhao/govm/types"
)

// Core 区块链核心接口
type Core interface {
	// Blockchain interface methods
	GetLastBlock() *types.Block
	GetHeight() uint64
	GetBlockByHash(hash types.Hash) (*types.Block, error)
	GetBlockByHeight(height uint64) (*types.Block, error)
	AddBlock(block *types.Block) error

	// TxProcessor interface methods
	ValidateTransaction(tx *types.TransactionWithSign) error
	ApplyTransaction(tx *types.TransactionWithSign) error
	GetTransactionByHash(hash types.Hash) (*types.TransactionWithSign, error)

	// Core specific methods
	Start() error
	Stop() error
	GetConsensus() consensus.PoAConsensus
	GetStorage() storage.Storage
}

// CoreConfig 核心模块配置
type CoreConfig struct {
	// 区块链配置
	ShardID uint64 // 分片ID

	// 存储配置
	DataDir string // 数据存储目录

	// 其他配置项可以在这里添加
}
