package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/lengzhao/govm/consensus"
	"github.com/lengzhao/govm/core"
	"github.com/lengzhao/govm/storage"
	"github.com/lengzhao/govm/txpool"
	"github.com/lengzhao/govm/types"
	"github.com/lengzhao/network"
	"github.com/stretchr/testify/assert"
)

// MockNetwork 模拟网络接口
type MockNetwork struct {
	peers []string
}

func (m *MockNetwork) Run(ctx context.Context) error {
	return nil
}

func (m *MockNetwork) GetLocalPeerID() string {
	return "test-node"
}

func (m *MockNetwork) GetLocalAddresses() []string {
	return []string{"localhost:8081"}
}

func (m *MockNetwork) GetPeers() []string {
	return m.peers
}

func (m *MockNetwork) SendRequest(peerID, msgType string, data []byte) ([]byte, error) {
	return []byte{}, nil
}

func (m *MockNetwork) BroadcastMessage(topic string, data []byte) error {
	return nil
}

func (m *MockNetwork) RegisterMessageHandler(msgType string, handler network.MessageHandler) {
	// 空实现
}

func (m *MockNetwork) RegisterRequestHandler(msgType string, handler network.RequestHandler) {
	// 空实现
}

func (m *MockNetwork) RegisterMessageFilter(topic string, filter network.MessageFilter) {
	// 空实现
}

func (m *MockNetwork) ConnectToPeer(addr string) error {
	return nil
}

// TestWalletAPI 测试钱包API功能
func TestWalletAPI(t *testing.T) {
	// 创建存储
	store := storage.NewMemoryStorage("")
	err := store.Start()
	assert.NoError(t, err)
	defer store.Stop()

	// 创建模拟网络
	mockNetwork := &MockNetwork{
		peers: []string{"peer1", "peer2"},
	}

	// 创建共识模块
	config := &consensus.PoAConfig{
		Validators:    []types.Address{{1}, {2}, {3}},
		BlockInterval: 2000,
		RoundLength:   3,
	}
	cons := consensus.NewPoAConsensus(config, store)

	// 创建核心模块
	coreConfig := &core.CoreConfig{
		ShardID: types.DefaultShardID,
		DataDir: "./test_data",
		Genesis: &types.GenesisConfig{
			Timestamp: uint64(time.Now().Unix()),
		},
	}
	coreModule, err := core.NewCore(coreConfig, cons, store)
	assert.NoError(t, err)

	err = coreModule.Start()
	assert.NoError(t, err)
	defer coreModule.Stop()

	// 创建交易池
	txPool := txpool.NewTxPool(coreModule, store)
	err = txPool.Start()
	assert.NoError(t, err)
	defer txPool.Stop()

	// 创建API实例
	apiServer := NewAPI(coreModule, txPool, store, mockNetwork)

	// 启动API服务
	err = apiServer.Start()
	assert.NoError(t, err)
	defer apiServer.Stop()

	// 等待API服务启动
	time.Sleep(100 * time.Millisecond)

	// 测试创建账户
	t.Run("TestCreateAccount", func(t *testing.T) {
		resp, err := http.Post("http://localhost:8080/account/create", "application/json", nil)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	// 测试获取账户列表
	t.Run("TestGetAccounts", func(t *testing.T) {
		resp, err := http.Get("http://localhost:8080/wallet/accounts")
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	// 测试获取节点指标
	t.Run("TestGetMetrics", func(t *testing.T) {
		resp, err := http.Get("http://localhost:8080/admin/metrics")
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var metrics Metrics
		err = json.NewDecoder(resp.Body).Decode(&metrics)
		assert.NoError(t, err)
		assert.Equal(t, uint64(0), metrics.BlockHeight) // 创世区块高度为0
		assert.Equal(t, uint64(0), metrics.TransactionCount)
		assert.Equal(t, 2, metrics.PeerCount)
	})
}

// TestWalletAPIInterface 测试WalletAPI接口实现
func TestWalletAPIInterface(t *testing.T) {
	// 创建临时钱包目录
	tmpWalletDir := "./test_wallet_" + fmt.Sprintf("%d", time.Now().UnixNano())
	os.MkdirAll(tmpWalletDir, 0755)
	defer os.RemoveAll(tmpWalletDir)

	// 创建存储
	store := storage.NewMemoryStorage("")
	err := store.Start()
	assert.NoError(t, err)
	defer store.Stop()

	// 创建模拟网络
	mockNetwork := &MockNetwork{}

	// 创建共识模块
	config := &consensus.PoAConfig{
		Validators:    []types.Address{{1}, {2}, {3}},
		BlockInterval: 2000,
		RoundLength:   3,
	}
	cons := consensus.NewPoAConsensus(config, store)

	// 创建核心模块
	coreConfig := &core.CoreConfig{
		ShardID: types.DefaultShardID,
		DataDir: "./test_data",
		Genesis: &types.GenesisConfig{
			Timestamp: uint64(time.Now().Unix()),
		},
	}
	coreModule, err := core.NewCore(coreConfig, cons, store)
	assert.NoError(t, err)

	err = coreModule.Start()
	assert.NoError(t, err)
	defer coreModule.Stop()

	// 创建交易池
	txPool := txpool.NewTxPool(coreModule, store)
	err = txPool.Start()
	assert.NoError(t, err)
	defer txPool.Stop()

	// 创建API实例（使用临时钱包目录）
	apiServer := &DefaultAPI{
		core:      coreModule,
		txPool:    txPool,
		storage:   store,
		network:   mockNetwork,
		port:      ":8080",
		accounts:  make(map[types.Address][]byte),
		walletDir: tmpWalletDir,
	}

	// 验证API实例实现了WalletAPI接口
	var _ WalletAPI = apiServer.WalletAPI()

	// 验证API实例实现了AdminAPI接口
	var _ AdminAPI = apiServer.AdminAPI()

	// 测试创建账户
	address, err := apiServer.CreateAccount()
	assert.NoError(t, err)
	assert.NotEqual(t, types.Address{}, address)

	// 测试获取账户列表
	accounts, err := apiServer.WalletAPI().GetAccounts()
	assert.NoError(t, err)
	assert.Contains(t, accounts, address)

	// 测试获取指标
	metrics, err := apiServer.AdminAPI().GetMetrics()
	assert.NoError(t, err)
	assert.Equal(t, uint64(0), metrics.BlockHeight)
	assert.Equal(t, uint64(0), metrics.TransactionCount)
	assert.Equal(t, 0, metrics.PeerCount)
}
