package test

import (
	"context"
	"fmt"
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

// TestFullNode 启动一个完整节点并验证API功能
func TestFullNode(t *testing.T) {
	// 创建存储
	store := storage.NewMemoryStorage("")
	err := store.Start()
	assert.NoError(t, err)
	defer store.Stop()

	// 创建网络配置
	netConfig := network.NewNetworkConfig()
	netConfig.Host = "127.0.0.1"
	netConfig.Port = 0 // 自动分配端口

	// 创建网络实例
	net, err := network.New(netConfig)
	assert.NoError(t, err)

	// 启动网络
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		_ = net.Run(ctx)
	}()

	// 等待网络启动
	time.Sleep(100 * time.Millisecond)

	// 创建验证节点列表
	validators := make([]types.Address, 3)
	for i := 0; i < 3; i++ {
		var addr types.Address
		copy(addr[:], []byte{byte(i + 1)})
		validators[i] = addr
	}

	// 创建共识配置
	config := &consensus.PoAConfig{
		Validators:    validators,
		BlockInterval: 2000,
		RoundLength:   3,
	}

	// 创建共识实例
	cons := consensus.NewPoAConsensus(config, store)

	// 创建核心配置
	coreConfig := &core.CoreConfig{
		ShardID: types.DefaultShardID,
		DataDir: "./test_data",
		Genesis: &types.GenesisConfig{
			Timestamp: uint64(time.Now().Unix()),
		},
	}

	// 创建核心实例
	coreModule, err := core.NewCore(coreConfig, cons, store)
	assert.NoError(t, err)

	// 设置网络接口
	err = coreModule.SetNetwork(net, 1)
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
	apiServer := api.NewAPI(coreModule, txPool, store, net)

	// 启动API服务
	err = apiServer.Start()
	assert.NoError(t, err)
	defer apiServer.Stop()

	// 等待API服务启动
	time.Sleep(100 * time.Millisecond)

	// 验证API服务是否正常运行
	resp, err := http.Get("http://localhost:8080/node/info")
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	fmt.Printf("Full node test passed! API is running on port 8080\n")
}
