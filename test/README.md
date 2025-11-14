# API测试说明

这个目录包含了用于测试govm区块链API功能的测试用例。

## 测试文件说明

1. [api_test.go](file:///Volumes/ssd/myproject/govm/test/api_test.go) - API功能集成测试
2. [node_test.go](file:///Volumes/ssd/myproject/govm/test/node_test.go) - 完整节点启动和API验证测试
3. [sync_test.go](file:///Volumes/ssd/myproject/govm/test/sync_test.go) - 区块同步功能测试

## 运行测试

### 运行所有测试

```bash
go test ./test/...
```

### 运行特定测试

```bash
# 运行API集成测试
go test ./test/ -run TestAPIIntegration

# 运行完整节点测试
go test ./test/ -run TestFullNode

# 运行区块同步测试
go test ./test/ -run TestBlockSyncWithAPI
```

## 测试内容

### TestAPIIntegration
- 测试API服务的基本功能
- 验证各个端点是否正常工作
- 包括节点信息、节点列表、区块查询、账户余额等API端点

### TestFullNode
- 启动一个完整的区块链节点
- 验证API服务是否能与核心模块、交易池、存储和网络模块正确集成
- 确保API服务能正常启动并响应请求

### TestBlockSyncWithAPI
- 验证新节点的区块同步功能
- 启动两个节点并连接它们
- 使用API方式确认区块同步的正确性
- 确保目标节点能从源节点同步区块数据

## API端点测试

测试覆盖了以下API端点：

- `GET /node/info` - 获取节点信息
- `GET /node/peers` - 获取连接的节点列表
- `GET /block/hash/{hash}` - 根据哈希获取区块
- `GET /block/number/{number}` - 根据区块号获取区块
- `GET /transaction/{hash}` - 获取交易信息
- `POST /transaction/send` - 发送交易
- `GET /account/balance/{address}` - 获取账户余额

## 测试环境

测试使用以下组件：

- 内存存储（MemoryStorage）
- PoA共识机制
- 真实网络接口
- 核心模块
- 交易池
- 同步器

所有测试都是独立运行的，不会影响实际的区块链数据。