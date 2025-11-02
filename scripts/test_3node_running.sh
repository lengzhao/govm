#!/bin/bash

# 测试3节点网络是否正常运行的脚本

echo "=== Testing 3-Node Network Status ==="

# 检查节点进程是否在运行
echo "Checking if node processes are running..."
NODE1_PID=$(pgrep -f "govm.*--node-id=1")
NODE2_PID=$(pgrep -f "govm.*--node-id=2")
NODE3_PID=$(pgrep -f "govm.*--node-id=3")

if [ ! -z "$NODE1_PID" ]; then
    echo "Node 1 is running (PID: $NODE1_PID)"
else
    echo "Node 1 is not running"
fi

if [ ! -z "$NODE2_PID" ]; then
    echo "Node 2 is running (PID: $NODE2_PID)"
else
    echo "Node 2 is not running"
fi

if [ ! -z "$NODE3_PID" ]; then
    echo "Node 3 is running (PID: $NODE3_PID)"
else
    echo "Node 3 is not running"
fi

# 检查数据目录是否存在
echo "Checking data directories..."
if [ -d "node1/data" ]; then
    echo "Node 1 data directory exists"
else
    echo "Node 1 data directory not found"
fi

if [ -d "node2/data" ]; then
    echo "Node 2 data directory exists"
else
    echo "Node 2 data directory not found"
fi

if [ -d "node3/data" ]; then
    echo "Node 3 data directory exists"
else
    echo "Node 3 data directory not found"
fi

echo "=== Test completed ==="