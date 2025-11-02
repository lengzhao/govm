#!/bin/bash

# 单节点网络测试脚本
# 该脚本用于验证单个节点是否能够正常启动和出块

set -e  # 遇到错误时退出

echo "=== Single Node Test ==="

# 确保在项目根目录
cd "$(dirname "$0")/.."

# 清理之前的运行数据
echo "Cleaning up previous data..."
rm -rf node1 config data
mkdir -p node1/data node1 config data

# 创建创世区块配置文件，设置时间为当前时间+10秒
GENESIS_TIME=$(($(date +%s) + 10))
cat > config/genesis.json << EOF
{
  "genesis": {
    "timestamp": $GENESIS_TIME
  }
}
EOF

# 创建单节点验证配置文件
cat > config/validators.json << 'EOF'
{
  "validators": [
    {
      "id": 1,
      "address": "00000000000000000001",
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
    ./govm --node-id=${node_id} --port=${port} --data-dir=${data_dir} --config=./config/validators.json --genesis=./config/genesis.json > node${node_id}/stdout.log 2> node${node_id}/stderr.log &
    
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

# 启动单个节点
echo "Starting single node..."
start_node 1 8001

# 验证节点是否正常运行
echo "Checking node status..."
node_running=false
if [ -f "node1/govm.pid" ]; then
    pid=$(cat node1/govm.pid)
    if ps -p $pid > /dev/null; then
        echo "Node 1: Running (PID: $pid)"
        node_running=true
    else
        echo "Node 1: Not running"
    fi
else
    echo "Node 1: PID file not found"
fi

if [ "$node_running" = false ]; then
    echo "Error: Node is not running!"
    stop_node 1
    exit 1
fi

echo "Node is running."

# 观察网络运行一段时间，足够看到出块
echo "Observing network for 30 seconds to see block generation..."
for i in {1..30}; do
    echo "Second $i..."
    sleep 1
done

# 检查日志输出
echo "Checking node logs for block generation..."
echo "--- Node 1 stdout ---"
if [ -f "node1/stdout.log" ]; then
    cat node1/stdout.log
else
    echo "No stdout log found"
fi

echo "--- Node 1 stderr ---"
if [ -f "node1/stderr.log" ]; then
    tail -20 node1/stderr.log
else
    echo "No stderr log found"
fi
echo ""

# 检查是否有区块生成
echo "Checking for block generation..."
if grep -q "成功生成并添加区块" node1/stdout.log; then
    echo "SUCCESS: Blocks are being generated!"
elif grep -q "Node is syncing" node1/stdout.log && grep -q "sync completed" node1/stderr.log; then
    echo "INFO: Node was syncing but now completed sync"
    # 再等待一段时间看是否开始出块
    echo "Waiting additional 10 seconds to see if block generation starts..."
    sleep 10
    if grep -q "成功生成并添加区块" node1/stdout.log; then
        echo "SUCCESS: Blocks are being generated after sync!"
    else
        echo "WARNING: No block generation detected even after sync completed"
    fi
else
    echo "WARNING: No block generation detected in logs"
fi

# 停止节点
echo "Stopping node..."
stop_node 1

echo "=== Single Node Test Completed ==="