package storage

import (
	"github.com/syndtr/goleveldb/leveldb"
)

// Storage 基础存储接口
type Storage interface {
	// Start 启动存储服务
	Start() error

	// Stop 停止存储服务
	Stop() error

	// Put 存储键值对
	Put(key []byte, value []byte) error

	// Get 根据键获取值
	Get(key []byte) ([]byte, error)

	// Has 检查键是否存在
	Has(key []byte) (bool, error)

	// Delete 删除键值对
	Delete(key []byte) error

	// BatchPut 批量存储键值对
	BatchPut(pairs map[string][]byte) error

	// NewStorage 创建指定命名空间的新存储实例
	NewStorage(namespace string) (Storage, error)

	// GetNamespace 获取当前存储实例的命名空间
	GetNamespace() string

	// Close 关闭当前存储实例
	Close() error

	// GetDB 获取底层LevelDB实例（仅供内部使用）
	getDB() *leveldb.DB
}

// Storage错误码定义
const (
	ErrNotStarted          = "STOR001"
	ErrNotFound            = "STOR002"
	ErrClosed              = "STOR003"
	ErrNamespaceConflict   = "STOR004"
	ErrBatchPartialFailure = "STOR005"
)

// StorageError 存储错误类型
type StorageError struct {
	Code    string
	Message string
}

func (e *StorageError) Error() string {
	return "[" + e.Code + "] " + e.Message
}
