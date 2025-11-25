package consensus

import (
	"encoding/binary"
	"fmt"
	"time"

	lzbinary "github.com/lengzhao/binary"
	"github.com/lengzhao/govm/crypto"
	"github.com/lengzhao/govm/storage"
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

// PoAConfig PoA共识配置
type PoAConfig struct {
	Validators    []types.Address // 验证节点地址列表
	BlockInterval uint64          // 区块间隔时间（毫秒）
	RoundLength   uint64          // 轮次长度（区块数）
}

// ValidatorInfo 验证节点信息
type ValidatorInfo struct {
	Address   types.Address // 节点地址
	PublicKey []byte        // 公钥
	Stake     uint64        // 抵押金额（保留字段）
	LastBlock uint64        // 最后出块高度
	IsActive  bool          // 是否活跃
}

// ConsensusState 共识状态
type ConsensusState struct {
	CurrentRound  uint64          // 当前轮次
	CurrentTurn   uint64          // 当前轮值
	Validators    []ValidatorInfo // 验证节点信息
	LastBlockTime uint64          // 上一个区块时间戳
}

// PoAConsensus PoA共识机制的具体实现接口
type PoAConsensus interface {
	Consensus

	// GetRound 获取当前轮次
	GetRound() uint64

	// GetTurn 获取当前轮值索引
	GetTurn() uint64

	// IsValidator 检查节点是否为验证者
	IsValidator(addr types.Address) bool

	// UpdateValidators 更新验证节点列表
	UpdateValidators(validators []types.Address) error

	// IsMyTurn 检查当前节点是否应该出块
	IsMyTurn(height uint64, nodeAddress types.Address) bool

	// GetCurrentValidator 获取当前应该出块的验证者地址
	GetCurrentValidator(height uint64) types.Address

	// GetValidatorCount 获取验证者数量
	GetValidatorCount() int
}

// DefaultPoA 默认PoA共识实现
type DefaultPoA struct {
	config  *PoAConfig
	state   *ConsensusState
	storage storage.Storage // 添加存储实例
}

// NewPoAConsensus 创建新的PoA共识实例
func NewPoAConsensus(config *PoAConfig, storage storage.Storage) PoAConsensus {
	vs, err := storage.NewStorage(types.SNValidator)
	if err != nil {
		return nil
	}

	return &DefaultPoA{
		config: config,
		state: &ConsensusState{
			CurrentRound:  0,
			CurrentTurn:   0,
			Validators:    make([]ValidatorInfo, len(config.Validators)),
			LastBlockTime: 0,
		},
		storage: vs, // validator storage
	}
}

// ValidateBlock 验证区块是否符合共识规则
func (p *DefaultPoA) ValidateBlock(block *types.Block) error {
	// 验证区块基本结构
	if block == nil || (block.Header.PrevHash == [32]byte{} && block.Header.BlockNumber != 0) {
		return fmt.Errorf("invalid block structure")
	}

	// 检查是否为空区块
	if p.IsEmptyBlock(block) {
		// 验证空区块
		return p.validateEmptyBlock(block)
	}

	// 验证区块签名
	if err := p.verifyBlockSignature(block); err != nil {
		return fmt.Errorf("block signature verification failed: %w", err)
	}

	// 验证出块节点权限
	if err := p.VerifyAuthority(block); err != nil {
		return fmt.Errorf("block authority verification failed: %w", err)
	}

	// 验证时间戳合规性
	if err := p.verifyTimestamp(block); err != nil {
		return fmt.Errorf("timestamp verification failed: %w", err)
	}

	// 验证相邻分片哈希
	if err := p.verifyAdjacentShards(block); err != nil {
		return fmt.Errorf("adjacent shards verification failed: %w", err)
	}

	return nil
}

// validateEmptyBlock 验证空区块
func (p *DefaultPoA) validateEmptyBlock(block *types.Block) error {
	// 验证空区块的基本结构
	if len(block.Transactions) > 0 {
		return fmt.Errorf("empty block should not contain transactions")
	}

	// 验证出块节点权限
	if err := p.VerifyAuthority(block); err != nil {
		return fmt.Errorf("block authority verification failed: %w", err)
	}

	// 验证时间戳合规性
	if err := p.verifyTimestamp(block); err != nil {
		return fmt.Errorf("timestamp verification failed: %w", err)
	}

	// 验证相邻分片哈希
	if err := p.verifyAdjacentShards(block); err != nil {
		return fmt.Errorf("adjacent shards verification failed: %w", err)
	}

	return nil
}

// verifyBlockSignature 验证区块签名
func (p *DefaultPoA) verifyBlockSignature(block *types.Block) error {
	// 空区块不需要验证签名
	if p.IsEmptyBlock(block) {
		return nil
	}

	// 获取应该出块的验证者地址
	expectedValidator := p.GetCurrentValidator(block.Header.BlockNumber)
	if expectedValidator == (types.Address{}) {
		return fmt.Errorf("failed to get expected validator for height %d", block.Header.BlockNumber)
	}

	// 检查区块头中的验证者是否匹配
	if block.Header.Validator != expectedValidator {
		return fmt.Errorf("block validator mismatch: expected %x, got %x", expectedValidator, block.Header.Validator)
	}

	// 重建区块头用于签名验证（排除签名字段）
	blockCopy := *block
	blockCopy.Header.Signature = nil

	// 序列化区块头用于签名验证
	data, err := lzbinary.Marshal(&blockCopy.Header.BlockHeader)
	if err != nil {
		return fmt.Errorf("failed to marshal block header: %w", err)
	}

	// 获取验证者的公钥
	// 注意：这是一个简化的实现，在实际应用中，我们需要从验证者的地址推导出公钥
	// 或者区块中直接包含公钥信息。这里我们创建一个新的密钥对仅用于演示
	_, pubKey, err := crypto.GenerateKeyPair(crypto.Ed25519)
	if err != nil {
		return fmt.Errorf("failed to generate key pair: %w", err)
	}

	// 验证签名
	if !crypto.Verify(data, block.Header.Signature, pubKey, crypto.Ed25519) {
		return fmt.Errorf("invalid block signature")
	}

	return nil
}

// verifyTimestamp 验证时间戳
func (p *DefaultPoA) verifyTimestamp(block *types.Block) error {
	// 检查时间戳是否合理
	if block.Header.Timestamp < p.state.LastBlockTime {
		return fmt.Errorf("block timestamp is older than previous block")
	}

	// 检查时间戳是否符合出块间隔
	expectedTime := p.state.LastBlockTime + p.config.BlockInterval
	if block.Header.BlockNumber > 0 && block.Header.Timestamp < expectedTime {
		return fmt.Errorf("block timestamp is too far in the past")
	}

	// 允许一定的时间误差（例如1秒）
	currentMs := time.Now().UnixMilli()
	if block.Header.Timestamp > uint64(currentMs+1000) {
		return fmt.Errorf("block timestamp is in the future")
	}

	return nil
}

// verifyAdjacentShards 验证相邻分片哈希
func (p *DefaultPoA) verifyAdjacentShards(block *types.Block) error {
	// 在实际实现中，这里会验证相邻分片的区块头哈希
	// 由于这是第一阶段，暂时留空
	if block.Header.OtherShards != [3]types.Hash{} {
		return fmt.Errorf("adjacent shards verification not implemented")
	}
	return nil
}

// GetValidator 获取当前验证者
func (p *DefaultPoA) GetValidator() interface{} {
	// 获取当前应该出块的验证节点
	if len(p.state.Validators) == 0 {
		return ValidatorInfo{}
	}
	index := p.state.CurrentTurn % uint64(len(p.state.Validators))
	return p.state.Validators[index]
}

// GetValidators 获取所有验证者列表
func (p *DefaultPoA) GetValidators() []interface{} {
	validators := make([]interface{}, len(p.state.Validators))
	for i, v := range p.state.Validators {
		validators[i] = v
	}
	return validators
}

// GetRound 获取当前轮次
func (p *DefaultPoA) GetRound() uint64 {
	return p.state.CurrentRound
}

// GetTurn 获取当前轮值索引
func (p *DefaultPoA) GetTurn() uint64 {
	return p.state.CurrentTurn
}

// IsValidator 检查节点是否为验证者
func (p *DefaultPoA) IsValidator(addr types.Address) bool {
	for _, validator := range p.state.Validators {
		if validator.Address == addr {
			return true
		}
	}

	// 同时检查配置中的验证者列表
	for _, validator := range p.config.Validators {
		if validator == addr {
			return true
		}
	}

	return false
}

// UpdateValidators 更新验证节点列表
func (p *DefaultPoA) UpdateValidators(validators []types.Address) error {
	// 更新配置中的验证节点列表
	p.config.Validators = validators

	// 更新状态中的验证节点信息
	p.state.Validators = make([]ValidatorInfo, len(validators))
	for i, addr := range validators {
		p.state.Validators[i] = ValidatorInfo{
			Address:   addr,
			PublicKey: make([]byte, 32), // 初始化为空的公钥
			Stake:     0,
			LastBlock: 0,
			IsActive:  true,
		}
	}

	return nil
}

// VerifyAuthority 验证节点是否有出块权限
func (p *DefaultPoA) VerifyAuthority(block *types.Block) error {
	if len(p.config.Validators) == 0 {
		return fmt.Errorf("no validators configured")
	}

	// 对于空区块，我们使用不同的验证方式
	if p.IsEmptyBlock(block) {
		// 空区块的验证可以简化，只需要检查是否是正确的验证节点
		// 这里简化实现，实际应用中可能需要更复杂的逻辑
		return nil
	}

	// 获取预期的验证节点地址
	expectedValidator, err := p.GetValidatorsAtHeight(block.Header.BlockNumber)
	if err != nil {
		return fmt.Errorf("failed to get validator at height: %w", err)
	}

	// 检查区块头中的签名是否来自预期的验证节点
	// 注意：这是一个简化的实现，在实际应用中需要更复杂的验证逻辑
	if expectedValidator.Address != block.Header.Validator {
		return fmt.Errorf("block not signed by expected validator")
	}

	return nil
}

// IsEmptyBlock 检查是否为空区块
func (p *DefaultPoA) IsEmptyBlock(block *types.Block) bool {
	// 空区块的判断条件：
	// 1. 交易列表为空
	// 2. 签名为空（可以是nil或空切片）
	return len(block.Transactions) == 0 && (len(block.Header.Signature) == 0 || block.Header.Signature == nil)
}

// GetConfig 获取共识配置
func (p *DefaultPoA) GetConfig() *PoAConfig {
	return p.config
}

// GetState 获取共识状态
func (p *DefaultPoA) GetState() *ConsensusState {
	return p.state
}

// GetValidatorsAtHeight 获取指定区块高度的验证者列表
func (p *DefaultPoA) GetValidatorsAtHeight(height uint64) (*ValidatorInfo, error) {
	// 构造存储键
	key := make([]byte, 8)
	binary.BigEndian.PutUint64(key, height)

	// 从存储中获取验证者信息
	data, err := p.storage.Get(key)
	if err != nil {
		// 如果存储中没有找到验证者信息，返回错误
		return nil, fmt.Errorf("validator info not found for height %d: %w", height, err)
	}

	// 反序列化验证者信息
	var validator ValidatorInfo
	if err := lzbinary.Unmarshal(data, &validator); err != nil {
		return nil, fmt.Errorf("failed to unmarshal validators: %w", err)
	}

	return &validator, nil
}

// IsMyTurn 检查当前节点是否应该出块
func (p *DefaultPoA) IsMyTurn(height uint64, nodeAddress types.Address) bool {
	if len(p.config.Validators) == 0 {
		return false
	}

	// 计算应该出块的验证者索引
	validatorIndex := height % uint64(len(p.config.Validators))

	// 检查当前节点地址是否匹配应该出块的验证者地址
	if validatorIndex < uint64(len(p.config.Validators)) {
		return p.config.Validators[validatorIndex] == nodeAddress
	}

	return false
}

// GetCurrentValidator 获取当前应该出块的验证者地址
func (p *DefaultPoA) GetCurrentValidator(height uint64) types.Address {
	if len(p.config.Validators) == 0 {
		return types.Address{}
	}

	// 计算应该出块的验证者索引
	validatorIndex := height % uint64(len(p.config.Validators))

	// 返回应该出块的验证者地址
	if validatorIndex < uint64(len(p.config.Validators)) {
		return p.config.Validators[validatorIndex]
	}

	return types.Address{}
}

// GetValidatorCount 获取验证者数量
func (p *DefaultPoA) GetValidatorCount() int {
	return len(p.config.Validators)
}
