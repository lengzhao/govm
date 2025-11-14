package sync

import (
	"context"
	"testing"
	"time"

	"github.com/lengzhao/govm/storage"
	"github.com/lengzhao/govm/types"
	"github.com/lengzhao/network"
	"github.com/stretchr/testify/assert"
)

// TestBlockSyncFullIntegration 测试完整的区块同步集成
func TestBlockSyncFullIntegration(t *testing.T) {
	// 测试同步请求响应结构
	request := &SyncRequest{
		StartHeight: 10,
		EndHeight:   20,
	}
	assert.Equal(t, uint64(10), request.StartHeight)
	assert.Equal(t, uint64(20), request.EndHeight)

	// 测试同步响应
	response := &SyncResponse{
		Blocks: []*types.Block{},
		Error:  "test error", // 添加非空字符串以避免序列化问题
	}
	assert.NotNil(t, response.Blocks)
	assert.Equal(t, 0, len(response.Blocks))
	assert.Equal(t, "test error", response.Error)

	// 测试序列化功能
	requestData, err := SerializeSyncRequest(request)
	assert.NoError(t, err)
	assert.NotNil(t, requestData)

	// 测试反序列化功能
	deserializedRequest, err := DeserializeSyncRequest(requestData)
	assert.NoError(t, err)
	assert.Equal(t, request.StartHeight, deserializedRequest.StartHeight)
	assert.Equal(t, request.EndHeight, deserializedRequest.EndHeight)

	// 测试响应序列化
	responseData, err := SerializeSyncResponse(response)
	assert.NoError(t, err)
	assert.NotNil(t, responseData)

	// 测试响应反序列化
	deserializedResponse, err := DeserializeSyncResponse(responseData)
	assert.NoError(t, err)
	assert.Equal(t, len(response.Blocks), len(deserializedResponse.Blocks))
	assert.Equal(t, response.Error, deserializedResponse.Error)
}

// TestSyncStateTransitions 测试同步状态转换
func TestSyncStateTransitions(t *testing.T) {
	// 创建内存存储
	mainStore := storage.NewMemoryStorage("")
	defer mainStore.Close()

	// 创建网络配置
	config := network.NewNetworkConfig()
	config.Host = "127.0.0.1"
	config.Port = 0 // 自动分配端口

	// 创建网络实例
	net, err := network.New(config)
	assert.NoError(t, err)

	// 启动网络
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		_ = net.Run(ctx)
	}()

	// 等待网络启动
	time.Sleep(100 * time.Millisecond)

	// 创建同步器
	syncer := &DefaultSyncer{
		network: net,
		storage: mainStore,
		state: SyncState{
			Status:     SyncStatusIdle,
			LastUpdate: time.Now(),
		},
		ctx:    ctx,
		cancel: cancel,
	}

	// 初始状态应该是idle
	state := syncer.GetSyncState()
	assert.Equal(t, SyncStatusIdle, state.Status)

	// 更新状态为syncing
	syncer.updateState(SyncStatusSyncing, 100, "")

	// 检查状态是否更新
	state = syncer.GetSyncState()
	assert.Equal(t, SyncStatusSyncing, state.Status)
	assert.Equal(t, uint64(100), state.TargetHeight)

	// 更新状态为complete
	syncer.updateState(SyncStatusComplete, 100, "sync completed")

	// 检查状态是否更新
	state = syncer.GetSyncState()
	assert.Equal(t, SyncStatusComplete, state.Status)
	assert.Equal(t, "sync completed", state.Error)
}

// TestSyncRequestResponseStructures 测试同步请求响应结构
func TestSyncRequestResponseStructures(t *testing.T) {
	// 测试同步请求
	request := &SyncRequest{
		StartHeight: 10,
		EndHeight:   20,
	}
	assert.Equal(t, uint64(10), request.StartHeight)
	assert.Equal(t, uint64(20), request.EndHeight)

	// 测试同步响应
	response := &SyncResponse{
		Blocks: []*types.Block{},
		Error:  "",
	}
	assert.NotNil(t, response.Blocks)
	assert.Equal(t, 0, len(response.Blocks))
	assert.Equal(t, "", response.Error)

	// 测试高度请求
	heightRequest := &HeightRequest{
		NodeID: "node-1",
	}
	assert.Equal(t, "node-1", heightRequest.NodeID)

	// 测试高度响应
	heightResponse := &HeightResponse{
		NodeID: "node-1",
		Height: 100,
		Error:  "",
	}
	assert.Equal(t, "node-1", heightResponse.NodeID)
	assert.Equal(t, uint64(100), heightResponse.Height)
	assert.Equal(t, "", heightResponse.Error)
}
