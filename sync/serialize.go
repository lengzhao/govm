package sync

import (
	"encoding/json"
	"fmt"

	"github.com/lengzhao/govm/types"
)

// SerializeSyncRequest 序列化同步请求
func SerializeSyncRequest(req *SyncRequest) ([]byte, error) {
	if req == nil {
		return nil, fmt.Errorf("sync request is nil")
	}

	return json.Marshal(req)
}

// DeserializeSyncRequest 反序列化同步请求
func DeserializeSyncRequest(data []byte) (*SyncRequest, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("empty data")
	}

	var req SyncRequest
	if err := json.Unmarshal(data, &req); err != nil {
		return nil, fmt.Errorf("failed to unmarshal sync request: %w", err)
	}

	return &req, nil
}

// SerializeSyncResponse 序列化同步响应
func SerializeSyncResponse(resp *SyncResponse) ([]byte, error) {
	if resp == nil {
		return nil, fmt.Errorf("sync response is nil")
	}

	return json.Marshal(resp)
}

// DeserializeSyncResponse 反序列化同步响应
func DeserializeSyncResponse(data []byte) (*SyncResponse, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("empty data")
	}

	var resp SyncResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal sync response: %w", err)
	}

	return &resp, nil
}

// SerializeBlock 序列化区块
func SerializeBlock(block *types.Block) ([]byte, error) {
	if block == nil {
		return nil, fmt.Errorf("block is nil")
	}

	return json.Marshal(block)
}

// DeserializeBlock 反序列化区块
func DeserializeBlock(data []byte) (*types.Block, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("empty data")
	}

	var block types.Block
	if err := json.Unmarshal(data, &block); err != nil {
		return nil, fmt.Errorf("failed to unmarshal block: %w", err)
	}

	return &block, nil
}

// SerializeBlocks 序列化区块列表
func SerializeBlocks(blocks []*types.Block) ([][]byte, error) {
	if blocks == nil {
		return nil, fmt.Errorf("blocks is nil")
	}

	result := make([][]byte, len(blocks))
	for i, block := range blocks {
		data, err := SerializeBlock(block)
		if err != nil {
			return nil, fmt.Errorf("failed to serialize block %d: %w", i, err)
		}
		result[i] = data
	}

	return result, nil
}

// DeserializeBlocks 反序列化区块列表
func DeserializeBlocks(data [][]byte) ([]*types.Block, error) {
	if data == nil {
		return nil, fmt.Errorf("data is nil")
	}

	result := make([]*types.Block, len(data))
	for i, blockData := range data {
		block, err := DeserializeBlock(blockData)
		if err != nil {
			return nil, fmt.Errorf("failed to deserialize block %d: %w", i, err)
		}
		result[i] = block
	}

	return result, nil
}

// SerializeHeightRequest 序列化高度请求
func SerializeHeightRequest(req *HeightRequest) ([]byte, error) {
	if req == nil {
		return nil, fmt.Errorf("height request is nil")
	}

	return json.Marshal(req)
}

// DeserializeHeightRequest 反序列化高度请求
func DeserializeHeightRequest(data []byte) (*HeightRequest, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("empty data")
	}

	var req HeightRequest
	if err := json.Unmarshal(data, &req); err != nil {
		return nil, fmt.Errorf("failed to unmarshal height request: %w", err)
	}

	return &req, nil
}

// SerializeHeightResponse 序列化高度响应
func SerializeHeightResponse(resp *HeightResponse) ([]byte, error) {
	if resp == nil {
		return nil, fmt.Errorf("height response is nil")
	}

	return json.Marshal(resp)
}

// DeserializeHeightResponse 反序列化高度响应
func DeserializeHeightResponse(data []byte) (*HeightResponse, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("empty data")
	}

	var resp HeightResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal height response: %w", err)
	}

	return &resp, nil
}
