package storage

import (
	"sync"

	"github.com/syndtr/goleveldb/leveldb"
)

// MemoryStorage 内存存储实现
type MemoryStorage struct {
	data      map[string][]byte
	namespace string
	closed    bool
	mutex     sync.RWMutex
}

// NewMemoryStorage 创建新的内存存储实例
func NewMemoryStorage(namespace string) *MemoryStorage {
	return &MemoryStorage{
		data:      make(map[string][]byte),
		namespace: namespace,
		closed:    true,
	}
}

// Start 启动存储服务
func (s *MemoryStorage) Start() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if !s.closed {
		return nil
	}

	s.closed = false
	return nil
}

// Stop 停止存储服务
func (s *MemoryStorage) Stop() error {
	return s.Close()
}

// Close 关闭当前存储实例
func (s *MemoryStorage) Close() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.closed {
		return nil
	}

	s.closed = true
	// 不清除数据，以便后续可能的重启
	return nil
}

// getDB 获取底层LevelDB实例（仅供内部使用）
// 对于内存存储，这个方法返回nil
func (s *MemoryStorage) getDB() *leveldb.DB {
	return nil
}

// GetNamespace 获取当前存储实例的命名空间
func (s *MemoryStorage) GetNamespace() string {
	return s.namespace
}

// addNamespacePrefix 为键添加命名空间前缀
func (s *MemoryStorage) addNamespacePrefix(key []byte) []byte {
	if s.namespace == "" {
		return key
	}

	// 创建带命名空间前缀的键
	prefixedKey := make([]byte, len(s.namespace)+1+len(key))
	copy(prefixedKey, s.namespace)
	prefixedKey[len(s.namespace)] = '_'
	copy(prefixedKey[len(s.namespace)+1:], key)

	return prefixedKey
}

// Put 存储键值对
func (s *MemoryStorage) Put(key []byte, value []byte) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.closed {
		return &StorageError{Code: ErrClosed, Message: "storage is closed"}
	}

	prefixedKey := s.addNamespacePrefix(key)
	// 创建值的副本以避免外部修改影响存储的数据
	valueCopy := make([]byte, len(value))
	copy(valueCopy, value)
	s.data[string(prefixedKey)] = valueCopy
	return nil
}

// Get 根据键获取值
func (s *MemoryStorage) Get(key []byte) ([]byte, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if s.closed {
		return nil, &StorageError{Code: ErrClosed, Message: "storage is closed"}
	}

	prefixedKey := s.addNamespacePrefix(key)
	value, exists := s.data[string(prefixedKey)]
	if !exists {
		return nil, &StorageError{Code: ErrNotFound, Message: "key not found"}
	}

	// 返回值的副本以避免外部修改影响存储的数据
	valueCopy := make([]byte, len(value))
	copy(valueCopy, value)
	return valueCopy, nil
}

// Has 检查键是否存在
func (s *MemoryStorage) Has(key []byte) (bool, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if s.closed {
		return false, &StorageError{Code: ErrClosed, Message: "storage is closed"}
	}

	prefixedKey := s.addNamespacePrefix(key)
	_, exists := s.data[string(prefixedKey)]
	return exists, nil
}

// Delete 删除键值对
func (s *MemoryStorage) Delete(key []byte) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.closed {
		return &StorageError{Code: ErrClosed, Message: "storage is closed"}
	}

	prefixedKey := s.addNamespacePrefix(key)
	delete(s.data, string(prefixedKey))
	return nil
}

// BatchPut 批量存储键值对
func (s *MemoryStorage) BatchPut(pairs map[string][]byte) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.closed {
		return &StorageError{Code: ErrClosed, Message: "storage is closed"}
	}

	for key, value := range pairs {
		prefixedKey := s.addNamespacePrefix([]byte(key))
		// 创建值的副本以避免外部修改影响存储的数据
		valueCopy := make([]byte, len(value))
		copy(valueCopy, value)
		s.data[string(prefixedKey)] = valueCopy
	}

	return nil
}

// NewStorage 创建指定命名空间的新存储实例
func (s *MemoryStorage) NewStorage(namespace string) (Storage, error) {
	// 检查命名空间是否冲突
	if namespace == s.namespace {
		return nil, &StorageError{Code: ErrNamespaceConflict, Message: "namespace conflict"}
	}

	// 创建新的存储实例，共享相同的数据映射但使用不同的命名空间
	newStorage := &MemoryStorage{
		data:      s.data, // 共享数据映射
		namespace: namespace,
		closed:    s.closed,
	}

	return newStorage, nil
}
