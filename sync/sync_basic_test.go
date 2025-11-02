package sync

import (
	"testing"

	"github.com/lengzhao/govm/types"
	"github.com/stretchr/testify/assert"
)

func TestSyncState(t *testing.T) {
	// 测试同步状态枚举值
	assert.Equal(t, SyncStatus("idle"), SyncStatusIdle)
	assert.Equal(t, SyncStatus("starting"), SyncStatusStarting)
	assert.Equal(t, SyncStatus("syncing"), SyncStatusSyncing)
	assert.Equal(t, SyncStatus("complete"), SyncStatusComplete)
	assert.Equal(t, SyncStatus("error"), SyncStatusError)
}

func TestSyncRequestSerialization(t *testing.T) {
	// 测试同步请求序列化
	request := &SyncRequest{
		StartHeight: 10,
		EndHeight:   20,
	}

	data, err := SerializeSyncRequest(request)
	assert.NoError(t, err)
	assert.NotNil(t, data)

	// 测试反序列化
	deserialized, err := DeserializeSyncRequest(data)
	assert.NoError(t, err)
	assert.Equal(t, request.StartHeight, deserialized.StartHeight)
	assert.Equal(t, request.EndHeight, deserialized.EndHeight)
}

func TestSyncResponseSerialization(t *testing.T) {
	// 测试同步响应序列化
	response := &SyncResponse{
		Blocks: []*types.Block{},
		Error:  "test error",
	}

	data, err := SerializeSyncResponse(response)
	assert.NoError(t, err)
	assert.NotNil(t, data)

	// 测试反序列化
	deserialized, err := DeserializeSyncResponse(data)
	assert.NoError(t, err)
	assert.Equal(t, len(response.Blocks), len(deserialized.Blocks))
	assert.Equal(t, response.Error, deserialized.Error)
}