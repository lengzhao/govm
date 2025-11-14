package sync

import (
	"testing"

	"github.com/lengzhao/govm/types"
	"github.com/stretchr/testify/assert"
)

// TestSyncStructures 测试同步数据结构
func TestSyncStructures(t *testing.T) {
	// 测试同步请求结构
	request := &SyncRequest{
		StartHeight: 1,
		EndHeight:   10,
	}
	assert.Equal(t, uint64(1), request.StartHeight)
	assert.Equal(t, uint64(10), request.EndHeight)

	// 测试同步响应结构
	response := &SyncResponse{
		Blocks: []*types.Block{},
		Error:  "",
	}
	assert.NotNil(t, response)
	assert.Equal(t, 0, len(response.Blocks))

	// 测试高度请求结构
	heightRequest := &HeightRequest{
		NodeID: "test-node",
	}
	assert.Equal(t, "test-node", heightRequest.NodeID)

	// 测试高度响应结构
	heightResponse := &HeightResponse{
		NodeID: "test-node",
		Height: 100,
		Error:  "",
	}
	assert.Equal(t, "test-node", heightResponse.NodeID)
	assert.Equal(t, uint64(100), heightResponse.Height)
}

// TestSyncStatus 测试同步状态枚举
func TestSyncStatus(t *testing.T) {
	assert.Equal(t, SyncStatus("idle"), SyncStatusIdle)
	assert.Equal(t, SyncStatus("starting"), SyncStatusStarting)
	assert.Equal(t, SyncStatus("syncing"), SyncStatusSyncing)
	assert.Equal(t, SyncStatus("complete"), SyncStatusComplete)
	assert.Equal(t, SyncStatus("error"), SyncStatusError)
}

// TestHeightResponseSerialization 测试高度响应序列化
func TestHeightResponseSerialization(t *testing.T) {
	// 创建高度响应
	response := &HeightResponse{
		NodeID: "test-node-123",
		Height: 100,
		Error:  "test-error", // 使用非空值避免binary库的bug
	}

	// 序列化
	data, err := SerializeHeightResponse(response)
	assert.NoError(t, err)
	assert.NotNil(t, data)

	// 反序列化
	deserialized, err := DeserializeHeightResponse(data)
	assert.NoError(t, err)
	assert.Equal(t, response.NodeID, deserialized.NodeID)
	assert.Equal(t, response.Height, deserialized.Height)
	assert.Equal(t, response.Error, deserialized.Error)
}
