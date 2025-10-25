package storage

import (
	"testing"
)

func TestMemoryStorage(t *testing.T) {
	// 创建存储实例
	storage := NewMemoryStorage("test")

	// 启动存储服务
	err := storage.Start()
	if err != nil {
		t.Fatalf("Failed to start storage: %v", err)
	}

	// 测试Put和Get
	testKey := []byte("test_key")
	testValue := []byte("test_value")

	err = storage.Put(testKey, testValue)
	if err != nil {
		t.Errorf("Failed to put key-value pair: %v", err)
	}

	value, err := storage.Get(testKey)
	if err != nil {
		t.Errorf("Failed to get value: %v", err)
	}

	if string(value) != string(testValue) {
		t.Errorf("Expected value %s, got %s", string(testValue), string(value))
	}

	// 测试Has
	exists, err := storage.Has(testKey)
	if err != nil {
		t.Errorf("Failed to check if key exists: %v", err)
	}

	if !exists {
		t.Error("Expected key to exist")
	}

	// 测试Delete
	err = storage.Delete(testKey)
	if err != nil {
		t.Errorf("Failed to delete key: %v", err)
	}

	// 确认键已被删除
	exists, err = storage.Has(testKey)
	if err != nil {
		t.Errorf("Failed to check if key exists after deletion: %v", err)
	}

	if exists {
		t.Error("Expected key to be deleted")
	}

	// 测试BatchPut
	pairs := map[string][]byte{
		"batch_key1": []byte("batch_value1"),
		"batch_key2": []byte("batch_value2"),
		"batch_key3": []byte("batch_value3"),
	}

	err = storage.BatchPut(pairs)
	if err != nil {
		t.Errorf("Failed to batch put key-value pairs: %v", err)
	}

	// 验证批量插入的数据
	for key, expectedValue := range pairs {
		value, err := storage.Get([]byte(key))
		if err != nil {
			t.Errorf("Failed to get batch inserted value for key %s: %v", key, err)
		}

		if string(value) != string(expectedValue) {
			t.Errorf("Expected batch value %s for key %s, got %s", string(expectedValue), key, string(value))
		}
	}

	// 测试命名空间功能
	blockStorage, err := storage.NewStorage("block")
	if err != nil {
		t.Errorf("Failed to create block storage: %v", err)
	}

	txStorage, err := storage.NewStorage("tx")
	if err != nil {
		t.Errorf("Failed to create tx storage: %v", err)
	}

	// 在不同命名空间中存储相同键
	commonKey := []byte("common_key")
	blockValue := []byte("block_value")
	txValue := []byte("tx_value")

	err = blockStorage.Put(commonKey, blockValue)
	if err != nil {
		t.Errorf("Failed to put value in block storage: %v", err)
	}

	err = txStorage.Put(commonKey, txValue)
	if err != nil {
		t.Errorf("Failed to put value in tx storage: %v", err)
	}

	// 验证命名空间隔离
	retrievedBlockValue, err := blockStorage.Get(commonKey)
	if err != nil {
		t.Errorf("Failed to get value from block storage: %v", err)
	}

	if string(retrievedBlockValue) != string(blockValue) {
		t.Errorf("Expected block value %s, got %s", string(blockValue), string(retrievedBlockValue))
	}

	retrievedTxValue, err := txStorage.Get(commonKey)
	if err != nil {
		t.Errorf("Failed to get value from tx storage: %v", err)
	}

	if string(retrievedTxValue) != string(txValue) {
		t.Errorf("Expected tx value %s, got %s", string(txValue), string(retrievedTxValue))
	}

	// 停止存储服务
	err = storage.Stop()
	if err != nil {
		t.Errorf("Failed to stop storage: %v", err)
	}
}

func TestMemoryStorageErrors(t *testing.T) {
	// 创建存储实例
	storage := NewMemoryStorage("test")

	// 测试在未启动存储时的操作
	testKey := []byte("test_key")
	testValue := []byte("test_value")

	// 这些操作应该返回错误，因为存储尚未启动
	_, err := storage.Get(testKey)
	if err == nil {
		t.Error("Expected error when getting from unstarted storage")
	}

	// 启动存储服务
	err = storage.Start()
	if err != nil {
		t.Fatalf("Failed to start storage: %v", err)
	}

	// 存储一些数据
	err = storage.Put(testKey, testValue)
	if err != nil {
		t.Fatalf("Failed to put key-value pair: %v", err)
	}

	// 停止存储服务
	err = storage.Stop()
	if err != nil {
		t.Fatalf("Failed to stop storage: %v", err)
	}

	// 测试在关闭存储后的操作
	_, err = storage.Get(testKey)
	if err == nil {
		t.Error("Expected error when getting from closed storage")
	}
}

func TestMemoryStorageDataIsolation(t *testing.T) {
	// 创建两个完全独立的存储实例
	storage1 := NewMemoryStorage("store1")
	storage2 := NewMemoryStorage("store2")

	err := storage1.Start()
	if err != nil {
		t.Fatalf("Failed to start storage1: %v", err)
	}

	err = storage2.Start()
	if err != nil {
		t.Fatalf("Failed to start storage2: %v", err)
	}

	// 在两个存储中存储相同的键但不同的值
	key := []byte("shared_key")
	value1 := []byte("value_in_store1")
	value2 := []byte("value_in_store2")

	err = storage1.Put(key, value1)
	if err != nil {
		t.Fatalf("Failed to put value in storage1: %v", err)
	}

	err = storage2.Put(key, value2)
	if err != nil {
		t.Fatalf("Failed to put value in storage2: %v", err)
	}

	// 验证数据隔离
	result1, err := storage1.Get(key)
	if err != nil {
		t.Errorf("Failed to get value from storage1: %v", err)
	}

	if string(result1) != string(value1) {
		t.Errorf("Expected %s from storage1, got %s", string(value1), string(result1))
	}

	result2, err := storage2.Get(key)
	if err != nil {
		t.Errorf("Failed to get value from storage2: %v", err)
	}

	if string(result2) != string(value2) {
		t.Errorf("Expected %s from storage2, got %s", string(value2), string(result2))
	}
}
