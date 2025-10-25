package storage

import (
	"path/filepath"
	"sync"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

// LevelDBStorage LevelDB存储实现
type LevelDBStorage struct {
	db        *leveldb.DB
	namespace string
	path      string
	closed    bool
	mutex     sync.RWMutex
}

// NewLevelDBStorage 创建新的LevelDB存储实例
func NewLevelDBStorage(path string, namespace string) (*LevelDBStorage, error) {
	storage := &LevelDBStorage{
		namespace: namespace,
		path:      filepath.Join(path, namespace),
		closed:    true,
	}

	return storage, nil
}

// Start 启动存储服务
func (s *LevelDBStorage) Start() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if !s.closed {
		return nil
	}

	// 创建数据库选项
	options := &opt.Options{
		// 缓存大小: 8MB
		BlockCacheCapacity: 8 * opt.MiB,
		// 写缓冲区大小: 4MB
		WriteBuffer: 4 * opt.MiB,
		// 打开文件数量限制: 1000
		OpenFilesCacheCapacity: 1000,
	}

	// 打开数据库
	db, err := leveldb.OpenFile(s.path, options)
	if err != nil {
		return err
	}

	s.db = db
	s.closed = false
	return nil
}

// Stop 停止存储服务
func (s *LevelDBStorage) Stop() error {
	return s.Close()
}

// Close 关闭当前存储实例
func (s *LevelDBStorage) Close() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.closed || s.db == nil {
		return nil
	}

	err := s.db.Close()
	if err != nil {
		return err
	}

	s.closed = true
	s.db = nil
	return nil
}

// getDB 获取底层LevelDB实例（仅供内部使用）
func (s *LevelDBStorage) getDB() *leveldb.DB {
	return s.db
}

// GetNamespace 获取当前存储实例的命名空间
func (s *LevelDBStorage) GetNamespace() string {
	return s.namespace
}

// addNamespacePrefix 为键添加命名空间前缀
func (s *LevelDBStorage) addNamespacePrefix(key []byte) []byte {
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
func (s *LevelDBStorage) Put(key []byte, value []byte) error {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if s.closed {
		return &StorageError{Code: ErrClosed, Message: "storage is closed"}
	}

	prefixedKey := s.addNamespacePrefix(key)
	return s.db.Put(prefixedKey, value, nil)
}

// Get 根据键获取值
func (s *LevelDBStorage) Get(key []byte) ([]byte, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if s.closed {
		return nil, &StorageError{Code: ErrClosed, Message: "storage is closed"}
	}

	prefixedKey := s.addNamespacePrefix(key)
	return s.db.Get(prefixedKey, nil)
}

// Has 检查键是否存在
func (s *LevelDBStorage) Has(key []byte) (bool, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if s.closed {
		return false, &StorageError{Code: ErrClosed, Message: "storage is closed"}
	}

	prefixedKey := s.addNamespacePrefix(key)
	return s.db.Has(prefixedKey, nil)
}

// Delete 删除键值对
func (s *LevelDBStorage) Delete(key []byte) error {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if s.closed {
		return &StorageError{Code: ErrClosed, Message: "storage is closed"}
	}

	prefixedKey := s.addNamespacePrefix(key)
	return s.db.Delete(prefixedKey, nil)
}

// BatchPut 批量存储键值对
func (s *LevelDBStorage) BatchPut(pairs map[string][]byte) error {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if s.closed {
		return &StorageError{Code: ErrClosed, Message: "storage is closed"}
	}

	batch := new(leveldb.Batch)
	for key, value := range pairs {
		prefixedKey := s.addNamespacePrefix([]byte(key))
		batch.Put(prefixedKey, value)
	}

	return s.db.Write(batch, nil)
}

// NewStorage 创建指定命名空间的新存储实例
func (s *LevelDBStorage) NewStorage(namespace string) (Storage, error) {
	// 检查命名空间是否冲突
	if namespace == s.namespace {
		return nil, &StorageError{Code: ErrNamespaceConflict, Message: "namespace conflict"}
	}

	// 创建新的存储实例，共享相同的数据库路径但使用不同的命名空间
	newStorage := &LevelDBStorage{
		db:        s.db,
		namespace: namespace,
		path:      filepath.Dir(s.path), // 使用相同的父目录
		closed:    s.closed,
	}

	return newStorage, nil
}
