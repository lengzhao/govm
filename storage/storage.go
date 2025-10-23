package storage

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
}
