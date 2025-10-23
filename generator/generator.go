package generator

import (
	"github.com/lengzhao/govm/consensus"
	"github.com/lengzhao/govm/storage"
	"github.com/lengzhao/govm/types"
)

// BlockGenerator 区块生成器接口
type BlockGenerator interface {
	// GenerateBlock 从交易池中选择交易并生成新区块
	GenerateBlock() (*types.Block, error)

	// SelectTransactions 从交易池中选择交易
	SelectTransactions() ([]*types.Transaction, error)

	// BuildBlockHeader 构建区块头
	BuildBlockHeader(transactions []*types.Transaction) (*types.BlockHeader, error)

	// AssembleBlock 组装完整区块
	AssembleBlock(header *types.BlockHeader, transactions []*types.Transaction) (*types.Block, error)

	// BroadcastBlock 广播区块到网络
	BroadcastBlock(block *types.Block) error
}

// NewBlockGenerator 创建新的区块生成器实例
func NewBlockGenerator(cons consensus.Consensus, store storage.Storage) BlockGenerator {
	// 实现将在具体结构体中提供
	return nil
}
