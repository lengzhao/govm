package types

// SyncChecker 同步状态检查接口
type SyncChecker interface {
	// IsSyncing 检查是否正在同步
	IsSyncing() bool
}
