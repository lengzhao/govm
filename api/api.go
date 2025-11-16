package api

import (
	"net/http"

	"github.com/lengzhao/govm/types"
)

// API API接口
type API interface {
	// Start 启动API服务
	Start() error

	// Stop 停止API服务
	Stop() error

	// RegisterEndpoints 注册API端点
	RegisterEndpoints(mux *http.ServeMux)

	// GetBlockByHash 根据哈希获取区块
	GetBlockByHash(hash types.Hash) (*types.Block, error)

	// GetBlockByNumber 根据区块号获取区块
	GetBlockByNumber(number uint64) (*types.Block, error)

	// GetTransaction 获取交易
	GetTransaction(hash types.Hash) (*types.Transaction, error)

	// SendTransaction 发送交易
	SendTransaction(tx *types.Transaction) (types.Hash, error)

	// GetBalance 获取余额
	GetBalance(address types.Address) (uint64, error)

	// GetAccount 获取账户信息
	GetAccount(address types.Address) (*types.Account, error)

	// CreateAccount 创建账户
	CreateAccount() (types.Address, error)

	// GetNodeInfo 获取节点信息
	GetNodeInfo() (NodeInfo, error)

	// GetPeers 获取节点列表
	GetPeers() ([]NodeInfo, error)

	// SetPort 设置API服务端口
	SetPort(port string)
}

// WalletAPI 钱包API接口
type WalletAPI interface {
	// CreateAccount 创建账户
	CreateAccount() (types.Address, error)

	// ImportAccount 导入账户
	ImportAccount(privateKey []byte) (types.Address, error)

	// ExportAccount 导出账户
	ExportAccount(address types.Address) ([]byte, error)

	// SignTransaction 签名交易
	SignTransaction(tx *types.Transaction, address types.Address) (*types.Transaction, error)

	// GetAccounts 获取账户列表
	GetAccounts() ([]types.Address, error)
}

// AdminAPI 管理API接口
type AdminAPI interface {
	// AddPeer 添加节点
	AddPeer(peerAddr string) error

	// RemovePeer 移除节点
	RemovePeer(peerAddr string) error

	// StartMining 开始挖矿
	StartMining() error

	// StopMining 停止挖矿
	StopMining() error

	// GetMetrics 获取节点指标
	GetMetrics() (Metrics, error)
}

// NodeInfo 节点信息结构
type NodeInfo struct {
	ID      string
	Address string
	Status  string
}

// Metrics 节点指标结构
type Metrics struct {
	BlockHeight      uint64
	TransactionCount uint64
	PeerCount        int
}
