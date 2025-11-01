# Core 模块

Core模块是govm区块链平台的核心组件，负责区块链的核心逻辑处理，包括区块管理、交易处理、状态维护等功能。

## 模块概述

Core模块是区块链系统的核心，它整合了区块链的基本功能，包括：

1. **区块管理** - 负责区块的创建、验证、存储和检索
2. **交易处理** - 负责交易的验证、执行和状态更新
3. **状态维护** - 维护区块链的全局状态
4. **模块协调** - 协调其他模块（共识、存储、网络等）的工作

## 核心组件

### Blockchain (区块链)
负责区块链的基本操作：
- 创世区块创建
- 区块验证和存储
- 区块查询（按哈希或高度）
- 区块链状态维护

### TxProcessor (交易处理器)
负责交易相关操作：
- 交易验证（签名、字段检查等）
- 交易执行和状态更新
- 交易查询

## 接口设计

### Core 接口
Core接口定义了核心模块的完整功能：

```go
type Core interface {
    // 区块链操作
    GetLastBlock() *types.Block
    GetHeight() uint64
    GetBlockByHash(hash types.Hash) (*types.Block, error)
    GetBlockByHeight(height uint64) (*types.Block, error)
    AddBlock(block *types.Block) error
    
    // 交易处理
    ValidateTransaction(tx *types.TransactionWithSign) error
    ApplyTransaction(tx *types.TransactionWithSign) error
    GetTransactionByHash(hash types.Hash) (*types.TransactionWithSign, error)
    
    // 模块管理
    Start() error
    Stop() error
    GetConsensus() consensus.PoAConsensus
    GetStorage() storage.Storage
}
```

## 使用示例

```go
// 创建核心模块配置
config := &core.CoreConfig{
    ShardID: types.DefaultShardID,
    DataDir: "./data",
}

// 创建核心模块实例
coreModule, err := core.NewCore(config, consensus, storage)
if err != nil {
    log.Fatal(err)
}

// 启动核心模块
err = coreModule.Start()
if err != nil {
    log.Fatal(err)
}

// 使用核心模块进行操作
lastBlock := coreModule.GetLastBlock()
height := coreModule.GetHeight()

// 停止核心模块
coreModule.Stop()
```

## 依赖关系

Core模块依赖以下模块：
- **consensus** - 共识模块，用于区块验证
- **storage** - 存储模块，用于数据持久化
- **crypto** - 加密模块，用于签名验证
- **types** - 类型定义模块

## 设计原则

1. **模块化设计** - 各功能模块职责清晰，便于维护和扩展
2. **接口抽象** - 通过接口定义模块功能，便于测试和替换实现
3. **状态安全** - 使用读写锁保护共享状态，确保并发安全
4. **错误处理** - 完善的错误处理机制，提供详细的错误信息

## 扩展性考虑

Core模块在设计时考虑了以下扩展性：

1. **插件化共识** - 支持不同的共识算法
2. **可配置存储** - 支持不同的存储后端
3. **分片支持** - 数据结构支持分片特性
4. **交易类型扩展** - 易于添加新的交易类型