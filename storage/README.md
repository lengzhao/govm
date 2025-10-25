# Storage Module

存储模块是govm区块链平台的核心组件之一，负责数据的持久化存储。该模块提供了两种实现方式：
1. 基于LevelDB的持久化存储
2. 纯内存存储（适用于测试或特殊场景）

模块提供高性能的键值对存储能力，并支持创建不同维度的存储实例，用于存储区块、交易、合约状态等不同类型的数据。

## 功能特性

1. 支持LevelDB持久化存储和内存存储两种模式
2. 支持创建不同命名空间的存储实例
3. 提供统一的存储接口供其他模块使用
4. 支持批量操作以提高性能
5. 线程安全的操作
6. 命名空间隔离确保不同类型数据的独立性

## 安装

```bash
go get github.com/lengzhao/govm/storage
```

## 快速开始

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/lengzhao/govm/storage"
)

func main() {
    // 创建LevelDB存储实例
    store, err := storage.NewLevelDBStorage("./data", "main")
    if err != nil {
        log.Fatal(err)
    }
    
    // 启动存储服务
    err = store.Start()
    if err != nil {
        log.Fatal(err)
    }
    defer store.Stop()
    
    // 存储数据
    err = store.Put([]byte("key"), []byte("value"))
    if err != nil {
        log.Fatal(err)
    }
    
    // 获取数据
    value, err := store.Get([]byte("key"))
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Retrieved value: %s\n", string(value))
}
```

### 内存存储示例

```go
// 创建内存存储实例
memStore := storage.NewMemoryStorage("main")

// 启动存储服务
err := memStore.Start()
if err != nil {
    log.Fatal(err)
}
defer memStore.Stop()

// 存储数据
err = memStore.Put([]byte("key"), []byte("value"))
if err != nil {
    log.Fatal(err)
}

// 获取数据
value, err := memStore.Get([]byte("key"))
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Retrieved value: %s\n", string(value))
```

## 接口说明

### Storage 接口

```go
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
}
```

## 命名空间支持

存储模块支持创建具有不同命名空间的存储实例，以实现不同类型数据的隔离存储：

```go
// 创建区块存储实例
blockStore, err := mainStore.NewStorage("block")
if err != nil {
    log.Fatal(err)
}

// 创建交易存储实例
txStore, err := mainStore.NewStorage("tx")
if err != nil {
    log.Fatal(err)
}

// 在不同命名空间中可以使用相同的键而不会冲突
err = blockStore.Put([]byte("data"), []byte("block_data"))
err = txStore.Put([]byte("data"), []byte("tx_data"))
```

## 批量操作

为了提高性能，存储模块支持批量操作：

```go
pairs := map[string][]byte{
    "key1": []byte("value1"),
    "key2": []byte("value2"),
    "key3": []byte("value3"),
}

err := store.BatchPut(pairs)
if err != nil {
    log.Fatal(err)
}
```

## 错误处理

存储模块定义了以下错误类型：

- `ErrNotStarted`: 存储服务未启动
- `ErrNotFound`: 键值对不存在
- `ErrClosed`: 存储实例已关闭
- `ErrNamespaceConflict`: 命名空间冲突
- `ErrBatchPartialFailure`: 批量操作部分失败

## 配置选项

LevelDB的配置选项：
- 缓存大小: 8MB
- 写缓冲区大小: 4MB
- 打开文件数量限制: 1000

## 性能优化

1. 使用批量操作减少磁盘I/O
2. 利用LevelDB内置缓存机制
3. 命名空间前缀优化
4. 读写锁确保并发安全

## 测试

运行存储模块的单元测试：

```bash
go test ./storage -v
```