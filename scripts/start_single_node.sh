#!/bin/bash

# 启动单验证节点网络脚本

# 创建数据目录和节点私钥目录
mkdir -p ./data/node1
mkdir -p ./node1

# 启动验证节点
echo "启动验证节点 (Node ID: 1, Port: 8001)"
go run ../main.go -node-id=1 -port=8001 -data-dir=./data/node1 -config=../config/validators.json -genesis=../config/genesis.json > node1.log 2>&1 &

echo "验证节点已启动，请等待几秒钟初始化..."
sleep 5

echo "检查节点日志:"
tail -n 20 node1.log