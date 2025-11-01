#!/bin/bash

# 3节点网络完整测试脚本
# 该脚本会启动3个节点，验证它们的功能，然后停止它们

set -e  # 遇到错误时退出

echo "=== 3-Node Network Test ==="

# 确保在项目根目录
cd "$(dirname "$0")/.."

# 编译程序
echo "Building govm..."
go build -o govm main.go

# 启动网络
echo "Starting 3-node network..."
./scripts/start_3node_network.sh

# 等待网络稳定
echo "Waiting for network to stabilize..."
sleep 15

# 验证网络
echo "Verifying network..."
./scripts/verify_3node_network.sh

# 运行一段时间观察出块
echo "Observing block generation for 30 seconds..."
for i in {1..30}; do
    echo "Second $i..."
    sleep 1
done

# 停止网络
echo "Stopping 3-node network..."
./scripts/stop_3node_network.sh

echo "=== Test completed successfully! ==="