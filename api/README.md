# API 模块

API模块为govm区块链平台提供RESTful API接口，允许外部应用与区块链进行交互。

## 功能特性

- 区块查询（按哈希或高度）
- 交易查询和提交
- 账户余额查询
- 节点信息查询
- 网络节点列表查询
- 钱包功能（账户创建、导入、导出、交易签名）
- 节点管理功能（添加/移除节点、挖矿控制、指标查询）

## API端点

### 区块查询

- `GET /block/hash/{hash}` - 根据区块哈希获取区块信息
- `GET /block/number/{number}` - 根据区块高度获取区块信息

### 交易查询和提交

- `GET /transaction/{hash}` - 根据交易哈希获取交易信息
- `POST /transaction/send` - 提交新交易

### 账户查询

- `GET /account/balance/{address}` - 获取账户余额
- `POST /account/create` - 创建新账户

### 钱包功能

- `GET /wallet/accounts` - 获取账户列表
- `POST /wallet/account/import` - 导入账户
- `GET /wallet/account/export/{address}` - 导出账户
- `POST /wallet/transaction/sign` - 签名交易

### 节点信息

- `GET /node/info` - 获取当前节点信息
- `GET /node/peers` - 获取连接的节点列表

### 节点管理

- `POST /admin/peer/add` - 添加节点
- `POST /admin/peer/remove` - 移除节点
- `POST /admin/mining/start` - 开始挖矿
- `POST /admin/mining/stop` - 停止挖矿
- `GET /admin/metrics` - 获取节点指标

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

### 钱包操作

```bash
# 创建新账户
curl -X POST http://localhost:8080/account/create

# 获取账户列表
curl http://localhost:8080/wallet/accounts

# 导入账户
curl -X POST http://localhost:8080/wallet/account/import \
  -H "Content-Type: application/json" \
  -d '{"private_key": "private_key_hex_string"}'

# 导出账户
curl http://localhost:8080/wallet/account/export/account_address

# 签名交易
curl -X POST http://localhost:8080/wallet/transaction/sign \
  -H "Content-Type: application/json" \
  -d '{
    "transaction": {
      "ShardID": 1,
      "From": "from_address",
      "To": "to_address",
      "Amount": 100,
      "Nonce": 1,
      "Data": ""
    },
    "address": "account_address"
  }'
```

### 节点管理

```bash
# 添加节点
curl -X POST http://localhost:8080/admin/peer/add \
  -H "Content-Type: application/json" \
  -d '{"peer_addr": "peer_address"}'

# 获取节点指标
curl http://localhost:8080/admin/metrics
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