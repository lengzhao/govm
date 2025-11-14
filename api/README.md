# API 模块

API模块为govm区块链平台提供RESTful API接口，允许外部应用与区块链进行交互。

## 功能特性

- 区块查询（按哈希或高度）
- 交易查询和提交
- 账户余额查询
- 节点信息查询
- 网络节点列表查询

## API端点

### 区块查询

- `GET /block/hash/{hash}` - 根据区块哈希获取区块信息
- `GET /block/number/{number}` - 根据区块高度获取区块信息

### 交易查询和提交

- `GET /transaction/{hash}` - 根据交易哈希获取交易信息
- `POST /transaction/send` - 提交新交易

### 账户查询

- `GET /account/balance/{address}` - 获取账户余额

### 节点信息

- `GET /node/info` - 获取当前节点信息
- `GET /node/peers` - 获取连接的节点列表

## 使用示例

启动节点后，API服务将在端口8080上运行。

### 查询区块信息

```bash
# 根据区块高度查询区块
curl http://localhost:8080/block/number/1

# 根据区块哈希查询区块
curl http://localhost:8080/block/hash/abc123...
```

### 提交交易

```bash
# 提交交易
curl -X POST http://localhost:8080/transaction/send \
  -H "Content-Type: application/json" \
  -d '{
    "ShardID": 1,
    "From": "address1",
    "To": "address2",
    "Amount": 100,
    "Nonce": 1,
    "Data": ""
  }'
```

### 查询账户余额

```bash
# 查询账户余额
curl http://localhost:8080/account/balance/address1
```

## 架构设计

API模块采用以下接口设计：

- `API` - 核心API接口
- `WalletAPI` - 钱包相关API接口
- `AdminAPI` - 管理相关API接口

## 依赖关系

API模块依赖以下核心模块：

- `core` - 区块链核心功能
- `txpool` - 交易池
- `storage` - 存储
- `network` - 网络通信