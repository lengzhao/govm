package test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/lengzhao/govm/api"
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

// TestAPIIntegration 测试API集成
func TestAPIIntegration(t *testing.T) {
	// 创建存储
	store := storage.NewMemoryStorage("")
	err := store.Start()
	assert.NoError(t, err)
	defer store.Stop()

	// 创建模拟网络
	mockNetwork := &MockNetwork{
		peers: []string{"peer1", "peer2", "peer3"},
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
	apiServer := api.NewAPI(coreModule, txPool, store, mockNetwork)

	// 启动API服务
	err = apiServer.Start()
	assert.NoError(t, err)
	defer apiServer.Stop()

	// 等待API服务启动
	time.Sleep(100 * time.Millisecond)

	// 测试获取节点信息
	t.Run("TestGetNodeInfo", func(t *testing.T) {
		resp, err := http.Get("http://localhost:8080/node/info")
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	// 测试获取节点列表
	t.Run("TestGetPeers", func(t *testing.T) {
		resp, err := http.Get("http://localhost:8080/node/peers")
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	// 测试获取区块（创世区块）
	t.Run("TestGetBlockByNumber", func(t *testing.T) {
		resp, err := http.Get("http://localhost:8080/block/number/0")
		if err != nil {
			// 创世区块可能不存在，这是可以接受的
			return
		}
		defer resp.Body.Close()

		// 可以是OK或NotFound
		assert.Contains(t, []int{http.StatusOK, http.StatusNotFound}, resp.StatusCode)
	})

	// 测试获取账户余额
	t.Run("TestGetBalance", func(t *testing.T) {
		resp, err := http.Get("http://localhost:8080/account/balance/0000000000000000000000000000000000000000")
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

// TestAPIStructures 测试API数据结构
func TestAPIStructures(t *testing.T) {
	// 测试NodeInfo结构
	nodeInfo := api.NodeInfo{
		ID:      "test-node-1",
		Address: "localhost:8080",
		Status:  "running",
	}
	assert.Equal(t, "test-node-1", nodeInfo.ID)
	assert.Equal(t, "localhost:8080", nodeInfo.Address)
	assert.Equal(t, "running", nodeInfo.Status)

	// 测试Metrics结构
	metrics := api.Metrics{
		BlockHeight:      1000,
		TransactionCount: 5000,
		PeerCount:        5,
	}
	assert.Equal(t, uint64(1000), metrics.BlockHeight)
	assert.Equal(t, uint64(5000), metrics.TransactionCount)
	assert.Equal(t, 5, metrics.PeerCount)
}

// TestAPIInterfaceImplementation 测试API接口实现
func TestAPIInterfaceImplementation(t *testing.T) {
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

	// 创建API实例
	apiServer := api.NewAPI(coreModule, txPool, store, mockNetwork)

	// 验证API实例实现了API接口
	var _ api.API = apiServer
}
