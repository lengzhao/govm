# Sync 同步模块

## 概述

Sync模块负责实现新节点的启动和区块同步功能。当新节点加入网络时，它需要从其他节点同步已有的区块数据，以确保整个网络的数据一致性。

## 功能特性

1. **节点发现**: 自动发现网络中的其他节点
2. **区块同步**: 从其他节点同步历史区块数据
3. **状态管理**: 管理同步过程的状态和进度
4. **错误处理**: 处理同步过程中可能出现的网络异常和数据不一致问题
5. **断点续传**: 支持从中断的位置继续同步

## 核心组件

### Syncer 接口

Syncer接口定义了同步器的核心功能：

```go
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
```

### SyncState 结构

SyncState结构用于跟踪同步过程的状态：

```go
type SyncState struct {
    StartHeight   uint64     // 同步起始高度
    CurrentHeight uint64     // 当前同步高度
    TargetHeight  uint64     // 目标同步高度
    Status        SyncStatus // 同步状态
    LastUpdate    time.Time  // 最后更新时间
    Error         string     // 错误信息
}
```

### 消息类型

#### SyncRequest
用于请求区块数据：
```go
type SyncRequest struct {
    StartHeight uint64 // 起始区块高度
    EndHeight   uint64 // 结束区块高度
}
```

#### SyncResponse
用于响应区块数据请求：
```go
type SyncResponse struct {
    Blocks []*types.Block // 区块数据
    Error  string         // 错误信息
}
```

## 使用方法

### 初始化同步器

```go
// 创建同步器实例
syncer := sync.NewSyncer(core, network, storage)

// 启动同步器
if err := syncer.StartSync(); err != nil {
    log.Printf("Failed to start syncer: %v", err)
}
```

### 检查同步状态

```go
// 获取同步状态
state := syncer.GetSyncState()
fmt.Printf("Sync status: %s, Progress: %d/%d\n", 
    state.Status, state.CurrentHeight, state.TargetHeight)

// 检查是否正在同步
if syncer.IsSyncing() {
    fmt.Println("Node is currently syncing")
}
```

### 停止同步器

```go
// 停止同步器
if err := syncer.StopSync(); err != nil {
    log.Printf("Failed to stop syncer: %v", err)
}
```

## 同步流程

1. **节点启动**: 新节点启动时初始化同步器
2. **节点发现**: 通过网络模块发现其他节点
3. **高度协商**: 获取网络中最新的区块高度
4. **区块下载**: 请求并下载缺失的区块数据
5. **区块验证**: 验证下载的区块是否有效
6. **区块存储**: 将验证通过的区块存储到本地
7. **状态更新**: 更新同步状态和进度

## 错误处理

同步模块实现了完善的错误处理机制：

- 网络异常时自动重连
- 数据不一致时重新同步
- 对频繁发送无效数据的节点进行临时屏蔽
- 实现请求频率限制机制防止拒绝服务攻击

## 性能优化

- 支持从多个节点并行下载区块数据
- 实现区块数据的并行验证和存储
- 支持可配置的同步带宽限制
- 对频繁访问的区块数据进行缓存

## 测试

同步模块包含完整的单元测试，可以通过以下命令运行：

```bash
go test ./sync -v
```

测试覆盖了以下方面：
- 同步状态管理
- 消息序列化和反序列化
- 同步器的启动和停止
- 同步状态检查