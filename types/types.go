package types

// Hash 哈希值类型
type Hash [32]byte

// Address 地址类型
type Address [20]byte

// Block 区块结构
type Block struct {
	Header       BlockHeaderWithSign
	Transactions []Hash
}

// BlockHeader 区块头信息
type BlockHeader struct {
	ShardID       uint64  // 分片ID
	BlockNumber   uint64  // 区块编号
	Timestamp     uint64  // 时间戳
	Validator     Address // 验证者
	PrevHash      Hash    // 前一区块哈希
	MerkleRoot    Hash    // Merkle根
	StateRootHash Hash    // 状态根哈希
	OtherShards   [3]Hash // 相邻3个分片链的对应区块的哈希
}

// BlockHeader 区块头信息
type BlockHeaderWithSign struct {
	BlockHeader
	Signature []byte // 区块签名
}

// Transaction 交易结构
type Transaction struct {
	ShardID uint64  // 分片ID
	From    Address // 发送方地址
	To      Address // 接收方地址
	Amount  uint64  // 转账金额
	Nonce   uint64  // 防重放攻击 nonce
}

// Transaction 交易结构
type TransactionWithSign struct {
	Transaction
	Signature []byte // 交易签名
}

// Constants 常量定义
const (
	// DefaultShardID 默认分片ID（第一分片）
	DefaultShardID uint64 = 1

	// BlockInterval 区块间隔时间（2秒）
	BlockInterval = 2000

	// BlockSizeLimit 区块大小限制（2M）
	BlockSizeLimit = 2 * 1024 * 1024

	// ValidatorCount 验证节点数量
	ValidatorCount = 21
)

// storage namespace
const (
	SNBlock     = "1" // key=block.hash
	SNStatus    = "2" // 全局状态
	SNTx        = "3" // key=tx.hash
	SNTxLog     = "4" // key=tx.hash
	SNValidator = "5" // key=block.height
)
