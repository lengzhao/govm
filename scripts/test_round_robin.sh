#!/bin/bash

# 轮流出块测试脚本
# 该脚本用于验证3个节点是否能够正确地轮流出块

set -e  # 遇到错误时退出

echo "=== 3-Node Round Robin Block Generation Test ==="

# 确保在项目根目录
cd "$(dirname "$0")/.."

# 清理之前的运行数据
echo "Cleaning up previous data..."
rm -rf node1 node2 node3 config data
mkdir -p node1/data node2/data node3/data node1 node2 node3 config data

# 创建3节点验证配置文件
cat > config/validators.json << 'EOF'
{
  "validators": [
    {
      "id": 1,
      "address": "00000000000000000001",
      "public_key": ""
    },
    {
      "id": 2,
      "address": "00000000000000000002",
      "public_key": ""
    },
    {
      "id": 3,
      "address": "00000000000000000003",
      "public_key": ""
    }
  ]
}
EOF

# 编译程序
echo "Building govm..."
go build -o govm main.go

# 函数：启动节点
start_node() {
    local node_id=$1
    local port=$2
    local data_dir="node${node_id}/data"
    
    echo "Starting node ${node_id} on port ${port} with data directory ${data_dir}"
    
    # 创建节点目录
    mkdir -p node${node_id}
    
    # 启动节点（后台运行）
    ./govm --node-id=${node_id} --port=${port} --data-dir=${data_dir} --config=./config/validators.json > node${node_id}/stdout.log 2> node${node_id}/stderr.log &
    
    # 保存进程ID
    echo $! > node${node_id}/govm.pid
    
    # 等待节点启动
    sleep 3
}

# 函数：停止节点
stop_node() {
    local node_id=$1
    
    if [ -f "node${node_id}/govm.pid" ]; then
        pid=$(cat node${node_id}/govm.pid)
        if ps -p $pid > /dev/null; then
            echo "Stopping node ${node_id} (PID: $pid)"
            kill $pid
            # 等待进程结束
            sleep 2
            # 如果进程仍然存在，强制杀死
            if ps -p $pid > /dev/null; then
                echo "Force killing node ${node_id}"
                kill -9 $pid
            fi
        else
            echo "Node ${node_id} is not running"
        fi
        # 删除PID文件
        rm -f node${node_id}/govm.pid
    else
        echo "Node ${node_id} PID file not found"
    fi
}

# 启动3个节点
echo "Starting 3 nodes..."
start_node 1 8001
start_node 2 8002
start_node 3 8003

# 验证节点是否正常运行
echo "Checking node status..."
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
    # 停止所有节点
    for i in {1..3}; do
        stop_node $i
    done
    exit 1
fi

echo "All nodes are running."

# 观察网络运行一段时间，足够看到轮流出块
echo "Observing network for 30 seconds to see round-robin block generation..."
for i in {1..30}; do
    echo "Second $i..."
    sleep 1
done

# 检查日志输出
echo "Checking node logs for round-robin block generation..."
for i in {1..3}; do
    echo "--- Node ${i} stdout ---"
    if [ -f "node${i}/stdout.log" ]; then
        cat node${i}/stdout.log
    else
        echo "No stdout log found"
    fi
    
    echo "--- Node ${i} stderr ---"
    if [ -f "node${i}/stderr.log" ]; then
        tail -20 node${i}/stderr.log
    else
        echo "No stderr log found"
    fi
    echo ""
done

# 停止所有节点
echo "Stopping all nodes..."
for i in {1..3}; do
    stop_node $i
done

echo "=== Round Robin Block Generation Test Completed ==="