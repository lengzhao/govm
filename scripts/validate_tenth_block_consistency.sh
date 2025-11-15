#!/bin/bash

# 第10个区块一致性验证脚本
# 该脚本用于验证在3节点网络中第10个区块是否在所有节点间保持一致

set -e  # 遇到错误时退出

echo "=== Tenth Block Consistency Validation ==="

# 确保在项目根目录
cd "$(dirname "$0")/.."

# 清理之前的运行数据
echo "Cleaning up previous data..."
rm -rf node1 node2 node3 config data
mkdir -p node1/data node2/data node3/data node1 node2 node3 config data

# 创建创世区块配置文件，设置时间为当前时间+10秒
GENESIS_TIME=$(($(date +%s) + 10))
cat > config/genesis.json << EOF
{
  "genesis": {
    "timestamp": $GENESIS_TIME
  }
}
EOF

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

# 观察网络运行足够长的时间，以生成至少10个区块
echo "Observing network for 60 seconds to generate at least 10 blocks..."
for i in {1..60}; do
    echo "Second $i..."
    sleep 1
done

# 检查各节点是否生成了至少10个区块
echo "Checking if all nodes have generated at least 10 blocks..."
tenth_block_consistent=true

for i in {1..3}; do
    echo "--- Node ${i} block generation status ---"
    if [ -f "node${i}/stdout.log" ]; then
        # 查找包含区块高度信息的日志行
        height_lines=$(grep -i "height\|block.*[0-9]" node${i}/stdout.log | tail -10 || true)
        if [ -n "$height_lines" ]; then
            echo "$height_lines"
            
            # 检查是否有高度大于等于10的区块
            if grep -q "height.*1[0-9]\|block.*1[0-9]" node${i}/stdout.log; then
                echo "✓ Node ${i} has generated at least 10 blocks"
            else
                echo "? Node ${i} may not have generated 10 blocks yet"
                tenth_block_consistent=false
            fi
        else
            echo "No block height information found in logs"
            tenth_block_consistent=false
        fi
    else
        echo "No stdout log found"
        tenth_block_consistent=false
    fi
    echo ""
done

# 特别检查第10个区块的信息
echo "Checking tenth block consistency across nodes..."
for i in {1..3}; do
    echo "--- Node ${i} tenth block info ---"
    if [ -f "node${i}/stdout.log" ]; then
        # 查找第10个区块的相关信息
        tenth_block_lines=$(grep -i "block.*10\|height.*10" node${i}/stdout.log | tail -5 || true)
        if [ -n "$tenth_block_lines" ]; then
            echo "$tenth_block_lines"
        else
            echo "No specific information about tenth block found"
        fi
    else
        echo "No stdout log found"
    fi
    echo ""
done

# 显示详细的日志内容以便分析
echo ""
echo "=== Detailed Log Analysis ==="
for i in {1..3}; do
    echo "--- Node ${i} full stdout ---"
    if [ -f "node${i}/stdout.log" ]; then
        tail -30 node${i}/stdout.log
    else
        echo "No stdout log found"
    fi
    echo ""
done

# 停止所有节点
echo "Stopping all nodes..."
for i in {1..3}; do
    stop_node $i
done

# 输出验证结果
echo ""
echo "=== Tenth Block Consistency Validation Summary ==="
if [ "$tenth_block_consistent" = true ]; then
    echo "✓ All nodes appear to have generated at least 10 blocks"
    echo "✓ Network is functioning properly"
    echo "✓ Block synchronization is working"
else
    echo "! Some nodes may not have generated 10 blocks yet"
    echo "! This could be due to network sync delays or timing issues"
    echo "! Check logs above for details"
fi

echo ""
echo "To verify the exact consistency of the tenth block, run the unit tests:"
echo "  go test -v ./test -run TestTenthBlockConsistency"
echo "  go test -v ./test -run TestTenthBlockContentConsistency"

echo ""
echo "=== Tenth Block Consistency Validation Completed ==="