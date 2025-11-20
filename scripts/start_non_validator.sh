#!/bin/bash

# 启动非验证节点脚本

# 创建数据目录和节点私钥目录
mkdir -p ./data/node102
mkdir -p ./node102

# 启动非验证节点 (使用相同的验证节点配置，但不同的节点ID)
echo "启动非验证节点 (Node ID: 102, Port: 8102)"
go run ./main.go -node-id=102 -port=8102 -data-dir=./data/node102 -config=./config/validators.json -genesis=./config/genesis.json

echo "非验证节点已退出"
