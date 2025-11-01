#!/bin/bash

# 网络功能测试脚本
# 该脚本用于测试网络模块的基本功能

echo "=== Network Module Test ==="

# 确保在项目根目录
cd "$(dirname "$0")/.."

# 创建测试目录
mkdir -p test_output

# 运行网络测试
echo "Running network test..."
go test ./network -v

if [ $? -eq 0 ]; then
    echo "Network test passed!"
else
    echo "Network test failed!"
    exit 1
fi

echo "=== Network test completed ==="