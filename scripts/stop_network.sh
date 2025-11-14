#!/bin/bash

# 停止所有网络节点脚本

echo "停止所有网络节点..."

# 查找并终止所有govm进程
pkill -f "main.go" || true

# 等待进程终止
sleep 2

echo "所有网络节点已停止"