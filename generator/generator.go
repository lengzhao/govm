package generator

import (
	"fmt"
	"time"

	lzbinary "github.com/lengzhao/binary"
	"github.com/lengzhao/govm/consensus"
	"github.com/lengzhao/govm/core"
	"github.com/lengzhao/govm/crypto"
	"github.com/lengzhao/govm/storage"
	"github.com/lengzhao/govm/types"
)

// BlockGenerator 区块生成器接口
type BlockGenerator interface {
	// GenerateBlock 从交易池中选择交易并生成新区块
	GenerateBlock(lastBlock *types.Block) (*types.Block, error)

	// SelectTransactions 从交易池中选择交易
	SelectTransactions() ([]*types.Transaction, error)

	// BuildBlockHeader 构建区块头
	BuildBlockHeader(lastBlock *types.Block, transactions []*types.Transaction) (*types.BlockHeader, error)

	// AssembleBlock 组装完整区块
	AssembleBlock(header *types.BlockHeader, transactions []*types.Transaction) (*types.Block, error)

	// SignBlock 对区块进行签名
	SignBlock(block *types.Block) error

	// BroadcastBlock 广播区块到网络
	BroadcastBlock(block *types.Block) error
}

// DefaultBlockGenerator 默认区块生成器实现
type DefaultBlockGenerator struct {
	consensus consensus.PoAConsensus
	storage   storage.Storage
	crypto    crypto.Crypto
}

// NewBlockGenerator 创建新的区块生成器实例
func NewBlockGenerator(cons consensus.PoAConsensus, store storage.Storage) BlockGenerator {
	return &DefaultBlockGenerator{
		consensus: cons,
		storage:   store,
		crypto:    crypto.NewCrypto(),
	}
}

// GenerateBlock 从交易池中选择交易并生成新区块
func (bg *DefaultBlockGenerator) GenerateBlock(lastBlock *types.Block) (*types.Block, error) {
	// 选择交易
	transactions, err := bg.SelectTransactions()
	if err != nil {
		return nil, fmt.Errorf("failed to select transactions: %w", err)
	}

	// 构建区块头
	header, err := bg.BuildBlockHeader(lastBlock, transactions)
	if err != nil {
		return nil, fmt.Errorf("failed to build block header: %w", err)
	}

	// 组装完整区块
	block, err := bg.AssembleBlock(header, transactions)
	if err != nil {
		return nil, fmt.Errorf("failed to assemble block: %w", err)
	}

	// 对区块进行签名
	if err := bg.SignBlock(block); err != nil {
		return nil, fmt.Errorf("failed to sign block: %w", err)
	}

	return block, nil
}

// SelectTransactions 从交易池中选择交易
func (bg *DefaultBlockGenerator) SelectTransactions() ([]*types.Transaction, error) {
	// 简化实现：暂不从交易池中选择交易
	// 实际实现中应该从交易池中选择有效的交易
	return []*types.Transaction{}, nil
}

// BuildBlockHeader 构建区块头
func (bg *DefaultBlockGenerator) BuildBlockHeader(lastBlock *types.Block, transactions []*types.Transaction) (*types.BlockHeader, error) {
	// 计算Merkle根
	merkleRoot := bg.calculateMerkleRoot(transactions)

	// 获取前一个区块哈希
	var prevBlockHash types.Hash
	if lastBlock != nil {
		// 创建临时区块链实例用于计算哈希
		tempBlockchain := core.NewBlockchain(bg.storage, bg.consensus)
		prevBlockHash = tempBlockchain.CalculateBlockHash(lastBlock)
	}

	// 获取当前区块高度
	blockHeight := uint64(1)
	if lastBlock != nil {
		blockHeight = lastBlock.Header.BlockNumber + 1
	}

	// 获取当前验证者
	validator := bg.consensus.GetCurrentValidator(blockHeight)

	// 获取时间戳（毫秒）
	timestamp := uint64(time.Now().UnixMilli())

	// 如果有前一个区块，确保时间戳不小于前一个区块
	if lastBlock != nil && timestamp < lastBlock.Header.Timestamp {
		timestamp = lastBlock.Header.Timestamp + 1
	}

	// 构建区块头
	header := &types.BlockHeader{
		ShardID:       types.DefaultShardID,
		BlockNumber:   blockHeight,
		Timestamp:     timestamp,
		Validator:     validator,
		PrevHash:      prevBlockHash,
		MerkleRoot:    merkleRoot,
		StateRootHash: types.Hash{},    // 状态根哈希，简化实现中为空
		OtherShards:   [3]types.Hash{}, // 相邻分片哈希，简化实现中为空
	}

	return header, nil
}

// AssembleBlock 组装完整区块
func (bg *DefaultBlockGenerator) AssembleBlock(header *types.BlockHeader, transactions []*types.Transaction) (*types.Block, error) {
	// 构建交易哈希列表
	txHashes := make([]types.Hash, len(transactions))
	for i, tx := range transactions {
		// 序列化交易以计算哈希
		data, err := lzbinary.Marshal(tx)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal transaction: %w", err)
		}
		txHashes[i] = bg.crypto.Hash(data)
	}

	// 创建区块头（带签名）
	headerWithSign := types.BlockHeaderWithSign{
		BlockHeader: *header,
		Signature:   nil, // 签名将在后续步骤中添加
	}

	// 创建区块
	block := &types.Block{
		Header:       headerWithSign,
		Transactions: txHashes,
	}

	return block, nil
}

// SignBlock 对区块进行签名
func (bg *DefaultBlockGenerator) SignBlock(block *types.Block) error {
	// 对于简化实现，我们创建一个空区块（不进行签名）
	// 空区块不包含任何交易和随机信息，所有内容都是确定的，且无需签名
	if len(block.Transactions) == 0 {
		// 创建空区块签名（空字节切片）
		block.Header.Signature = []byte{}
		return nil
	}

	// 获取当前验证者的私钥
	// 简化实现：使用测试私钥
	privKey, _, err := bg.crypto.GenerateKeyPair(crypto.Ed25519)
	if err != nil {
		return fmt.Errorf("failed to generate key pair: %w", err)
	}

	// 序列化区块头（排除签名字段）用于签名
	blockCopy := *block
	blockCopy.Header.Signature = nil

	data, err := lzbinary.Marshal(&blockCopy.Header.BlockHeader)
	if err != nil {
		return fmt.Errorf("failed to marshal block header: %w", err)
	}

	// 对区块头进行签名
	signature, err := bg.crypto.Sign(data, privKey)
	if err != nil {
		return fmt.Errorf("failed to sign block: %w", err)
	}

	// 将签名添加到区块头
	block.Header.Signature = signature

	return nil
}

// BroadcastBlock 广播区块到网络
func (bg *DefaultBlockGenerator) BroadcastBlock(block *types.Block) error {
	// 简化实现：暂不实现广播功能
	// 实际实现中应该将区块广播到网络中的其他节点
	fmt.Println("Broadcasting block...")
	return nil
}

// calculateMerkleRoot 计算交易的Merkle根
func (bg *DefaultBlockGenerator) calculateMerkleRoot(transactions []*types.Transaction) types.Hash {
	// 简化实现：如果交易为空，返回空哈希
	if len(transactions) == 0 {
		return types.Hash{}
	}

	// 如果只有一笔交易，直接返回该交易的哈希
	if len(transactions) == 1 {
		data, err := lzbinary.Marshal(transactions[0])
		if err != nil {
			return types.Hash{} // 出错时返回空哈希
		}
		return bg.crypto.Hash(data)
	}

	// 多笔交易时计算Merkle根（简化实现）
	// 实际实现中应该构建完整的Merkle树
	hashes := make([]types.Hash, len(transactions))
	for i, tx := range transactions {
		data, err := lzbinary.Marshal(tx)
		if err != nil {
			return types.Hash{} // 出错时返回空哈希
		}
		hashes[i] = bg.crypto.Hash(data)
	}

	// 简化实现：将所有哈希连接后再次哈希
	var combined []byte
	for _, hash := range hashes {
		combined = append(combined, hash[:]...)
	}
	return bg.crypto.Hash(combined)
}
