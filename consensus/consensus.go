package consensus

// Consensus 代表共识机制接口
type Consensus interface {
	// ValidateBlock 验证区块是否符合共识规则
	ValidateBlock(block interface{}) bool

	// GetValidator 获取当前验证者
	GetValidator() interface{}
}
