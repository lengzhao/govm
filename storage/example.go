package storage

import (
	"fmt"
	"os"
)

// ExampleUsage 展示如何使用存储模块
func ExampleUsage() {
	// 创建临时目录用于示例
	tempDir := "./example_data"
	defer os.RemoveAll(tempDir)

	// 创建主存储实例
	mainStorage, err := NewLevelDBStorage(tempDir, "main")
	if err != nil {
		fmt.Printf("Failed to create main storage: %v\n", err)
		return
	}

	// 启动存储服务
	err = mainStorage.Start()
	if err != nil {
		fmt.Printf("Failed to start storage: %v\n", err)
		return
	}
	defer mainStorage.Stop()

	// 存储一些数据
	err = mainStorage.Put([]byte("key1"), []byte("value1"))
	if err != nil {
		fmt.Printf("Failed to store key-value pair: %v\n", err)
		return
	}

	// 获取数据
	value, err := mainStorage.Get([]byte("key1"))
	if err != nil {
		fmt.Printf("Failed to retrieve value: %v\n", err)
		return
	}
	fmt.Printf("Retrieved value: %s\n", string(value))

	// 检查键是否存在
	exists, err := mainStorage.Has([]byte("key1"))
	if err != nil {
		fmt.Printf("Failed to check key existence: %v\n", err)
		return
	}
	fmt.Printf("Key exists: %t\n", exists)

	// 批量存储数据
	pairs := map[string][]byte{
		"batch_key1": []byte("batch_value1"),
		"batch_key2": []byte("batch_value2"),
		"batch_key3": []byte("batch_value3"),
	}
	err = mainStorage.BatchPut(pairs)
	if err != nil {
		fmt.Printf("Failed to batch store key-value pairs: %v\n", err)
		return
	}

	// 创建不同命名空间的存储实例
	blockStorage, err := mainStorage.NewStorage("block")
	if err != nil {
		fmt.Printf("Failed to create block storage: %v\n", err)
		return
	}

	txStorage, err := mainStorage.NewStorage("tx")
	if err != nil {
		fmt.Printf("Failed to create tx storage: %v\n", err)
		return
	}

	// 在不同命名空间中存储相同键
	commonKey := []byte("common_key")
	err = blockStorage.Put(commonKey, []byte("block_data"))
	if err != nil {
		fmt.Printf("Failed to store in block storage: %v\n", err)
		return
	}

	err = txStorage.Put(commonKey, []byte("tx_data"))
	if err != nil {
		fmt.Printf("Failed to store in tx storage: %v\n", err)
		return
	}

	// 从不同命名空间检索数据
	blockValue, err := blockStorage.Get(commonKey)
	if err != nil {
		fmt.Printf("Failed to retrieve from block storage: %v\n", err)
		return
	}
	fmt.Printf("Block storage value: %s\n", string(blockValue))

	txValue, err := txStorage.Get(commonKey)
	if err != nil {
		fmt.Printf("Failed to retrieve from tx storage: %v\n", err)
		return
	}
	fmt.Printf("Tx storage value: %s\n", string(txValue))

	fmt.Println("Storage example completed successfully!")
}