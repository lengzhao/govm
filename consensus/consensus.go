package consensus

import (
	"github.com/lengzhao/govm/types"
)

// Consensus 代表共识机制接口
type Consensus interface {
	// ValidateBlock 验证区块是否符合共识规则，包括签名验证和共识规则检查
	ValidateBlock(block *types.Block) error

	// GetValidator 获取当前验证者
	GetValidator() interface{}

	// GetValidators 获取所有验证者列表
	GetValidators() []interface{}
}

// PoAConsensus PoA共识机制的具体实现接口
type PoAConsensus interface {
	Consensus

	// GetRound 获取当前轮次
	GetRound() uint64

	// GetTurn 获取当前轮值的验证者
	GetTurn() interface{}

	// RotateValidators 轮换验证者顺序
	RotateValidators() error

	// IsValidator 检查节点是否为验证者
	IsValidator(nodeID string) bool
}
