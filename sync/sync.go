package sync

import (
	"time"

	"github.com/lengzhao/govm/types"
)

// SyncStatus 同步状态枚举
type SyncStatus string

const (
	SyncStatusIdle     SyncStatus = "idle"     // 空闲状态
	SyncStatusStarting SyncStatus = "starting" // 启动中
	SyncStatusSyncing  SyncStatus = "syncing"  // 同步中
	SyncStatusComplete SyncStatus = "complete" // 同步完成
	SyncStatusError    SyncStatus = "error"    // 同步错误
)

// SyncState 同步状态结构
type SyncState struct {
	StartHeight   uint64     // 同步起始高度
	CurrentHeight uint64     // 当前同步高度
	TargetHeight  uint64     // 目标同步高度
	Status        SyncStatus // 同步状态
	LastUpdate    time.Time  // 最后更新时间
	Error         string     // 错误信息
}

// Syncer 同步器接口
type Syncer interface {
	// StartSync 启动同步过程
	StartSync() error

	// StopSync 停止同步过程
	StopSync() error

	// GetSyncState 获取同步状态
	GetSyncState() *SyncState

	// IsSyncing 检查是否正在同步
	IsSyncing() bool
}

// SyncRequest 同步请求消息
type SyncRequest struct {
	StartHeight uint64 // 起始区块高度
	EndHeight   uint64 // 结束区块高度
}

// SyncResponse 同步响应消息
type SyncResponse struct {
	Blocks []*types.Block // 区块数据
	Error  string         // 错误信息
}

// HeightRequest 高度请求消息
type HeightRequest struct {
	NodeID string // 请求节点ID
}

// HeightResponse 高度响应消息
type HeightResponse struct {
	NodeID string // 响应节点ID
	Height uint64 // 区块链高度
	Error  string // 错误信息
}

// SyncMessageHandler 同步消息处理器接口
type SyncMessageHandler interface {
	// HandleSyncRequest 处理同步请求
	HandleSyncRequest(from string, request *SyncRequest) error

	// HandleSyncResponse 处理同步响应
	HandleSyncResponse(from string, response *SyncResponse) error
}
