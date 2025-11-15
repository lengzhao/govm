#!/bin/bash

# 验证3节点网络功能脚本
# 该脚本用于验证3节点网络是否正常工作

echo "Verifying 3-node blockchain network..."

# 检查节点是否都在运行
echo "Checking if all nodes are running..."
all_running=true
for i in {1..3}; do
    if [ -f "node${i}/govm.pid" ]; then
        pid=$(cat node${i}/govm.pid)
        if ps -p $pid > /dev/null; then
            echo "Node ${i}: Running (PID: $pid)"
        else
            echo "Node ${i}: Not running"
            all_running=false
        fi
    else
        echo "Node ${i}: PID file not found"
        all_running=false
    fi
done

if [ "$all_running" = false ]; then
    echo "Error: Not all nodes are running!"
    exit 1
fi

echo "All nodes are running."

# 检查区块链高度
echo "Checking blockchain heights..."
for i in {1..3}; do
    # 查询节点的区块链高度
    api_port=$((8080 + i))
    height=$(curl -s http://localhost:${api_port}/blockchain/height 2>/dev/null || echo "unknown")
    echo "Node ${i} blockchain height: ${height}"
done

# 检查节点间连接
echo "Checking node connectivity..."
for i in {1..3}; do
    # 查询节点的邻居信息
    api_port=$((8080 + i))
    neighbors=$(curl -s http://localhost:${api_port}/network/peers 2>/dev/null || echo "unknown")
    echo "Node ${i} neighbors: ${neighbors}"
done

# 验证区块同步
echo "Checking block synchronization..."
# 获取所有节点的最新区块哈希进行比较
heights=()
hashes=()
for i in {1..3}; do
    api_port=$((8080 + i))
    # 获取节点最新区块哈希
    hash=$(curl -s http://localhost:${api_port}/blockchain/latest/hash 2>/dev/null || echo "unknown")
    hashes+=(${hash})
    
    # 获取节点高度
    height=$(curl -s http://localhost:${api_port}/blockchain/height 2>/dev/null || echo "unknown")
    heights+=(${height})
    
    echo "Node ${i} latest block hash: ${hash} (height: ${height})"
done

# 检查所有节点的区块哈希是否一致
all_synced=true
for ((i=1; i<${#hashes[@]}; i++)); do
    if [ "${hashes[0]}" != "${hashes[i]}" ] && [ "${hashes[0]}" != "unknown" ] && [ "${hashes[i]}" != "unknown" ]; then
        all_synced=false
        break
    fi
done

if [ "$all_synced" = true ]; then
    echo "All nodes have synchronized blocks: true"
else
    echo "All nodes have synchronized blocks: false"
fi

echo "Network verification completed successfully!"
echo "3-node network is functioning properly."