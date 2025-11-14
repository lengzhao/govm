#!/bin/bash

# 启动非验证节点脚本

# 创建数据目录和节点私钥目录
mkdir -p ./data/node2
mkdir -p ./node2

# 启动非验证节点 (使用相同的验证节点配置，但不同的节点ID)
echo "启动非验证节点 (Node ID: 2, Port: 8002)"
go run ../main.go -node-id=2 -port=8002 -data-dir=./data/node2 -config=../config/validators.json -genesis=../config/genesis.json > node2.log 2>&1 &

echo "非验证节点已启动，请等待几秒钟初始化..."
sleep 5

echo "检查节点日志:"
tail -n 20 node2.log