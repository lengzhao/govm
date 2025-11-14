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

// TestBlockSyncIntegration 测试区块同步集成
func TestBlockSyncIntegration(t *testing.T) {
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

	// 验证网络已启动
	assert.NotNil(t, net)

	// 创建同步器（使用nil作为core，因为我们只测试结构创建）
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

	// 验证同步器创建成功
	assert.NotNil(t, syncer)

	// 测试同步状态获取
	state := syncer.GetSyncState()
	assert.Equal(t, SyncStatusIdle, state.Status)

	// 测试IsSyncing方法
	assert.False(t, syncer.IsSyncing())
}

// TestSyncRequestSerializationNew 测试同步请求序列化
func TestSyncRequestSerializationNew(t *testing.T) {
	// 创建同步请求
	request := &SyncRequest{
		StartHeight: 100,
		EndHeight:   200,
	}

	// 序列化
	data, err := SerializeSyncRequest(request)
	assert.NoError(t, err)
	assert.NotNil(t, data)

	// 反序列化
	deserialized, err := DeserializeSyncRequest(data)
	assert.NoError(t, err)
	assert.Equal(t, request.StartHeight, deserialized.StartHeight)
	assert.Equal(t, request.EndHeight, deserialized.EndHeight)
}

// TestSyncResponseSerializationNew 测试同步响应序列化
func TestSyncResponseSerializationNew(t *testing.T) {
	// 创建同步响应
	response := &SyncResponse{
		Blocks: []*types.Block{},
		Error:  "test error message",
	}

	// 序列化
	data, err := SerializeSyncResponse(response)
	assert.NoError(t, err)
	assert.NotNil(t, data)

	// 反序列化
	deserialized, err := DeserializeSyncResponse(data)
	assert.NoError(t, err)
	assert.Equal(t, len(response.Blocks), len(deserialized.Blocks))
	assert.Equal(t, response.Error, deserialized.Error)
}

// TestHeightRequestSerializationNew 测试高度请求序列化
func TestHeightRequestSerializationNew(t *testing.T) {
	// 创建高度请求
	request := &HeightRequest{
		NodeID: "test-node-123",
	}

	// 序列化
	data, err := SerializeHeightRequest(request)
	assert.NoError(t, err)
	assert.NotNil(t, data)

	// 反序列化
	deserialized, err := DeserializeHeightRequest(data)
	assert.NoError(t, err)
	assert.Equal(t, request.NodeID, deserialized.NodeID)
}
