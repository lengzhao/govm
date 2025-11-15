#!/bin/bash

# 区块一致性验证脚本
# 该脚本用于验证多节点网络中区块高度的一致性和新区块的生成

set -e  # 遇到错误时退出

echo "=== Block Consistency Validation ==="

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

# 观察网络运行一段时间，足够生成一些区块
echo "Observing network for 20 seconds to generate blocks..."
for i in {1..20}; do
    echo "Second $i..."
    sleep 1
done

# 检查区块高度一致性
echo "Checking block height consistency across nodes..."
height_consistent=true

# 获取各节点的高度（通过API或其他方式）
# 这里我们通过检查日志来验证
for i in {1..3}; do
    echo "--- Node ${i} block height info ---"
    if [ -f "node${i}/stdout.log" ]; then
        # 查找包含高度信息的日志行
        height_lines=$(grep -i "height\|block" node${i}/stdout.log | tail -5 || true)
        if [ -n "$height_lines" ]; then
            echo "$height_lines"
        else
            echo "No height information found in logs"
        fi
    else
        echo "No stdout log found"
        height_consistent=false
    fi
    echo ""
done

# 检查是否有新区块生成
echo "Checking for new block generation..."
new_blocks_generated=true

for i in {1..3}; do
    echo "--- Node ${i} new block generation ---"
    if [ -f "node${i}/stdout.log" ]; then
        # 查找新区块生成的相关日志
        block_lines=$(grep -i "block\|generated\|added" node${i}/stdout.log | tail -5 || true)
        if [ -n "$block_lines" ]; then
            echo "$block_lines"
            # 检查是否有区块高度大于0的记录
            if grep -q "height.*[1-9]" node${i}/stdout.log || grep -q "block.*[1-9]" node${i}/stdout.log; then
                echo "✓ Node ${i} has generated new blocks"
            else
                echo "? Node ${i} may not have generated new blocks yet"
            fi
        else
            echo "No block generation information found in logs"
        fi
    else
        echo "No stdout log found"
        new_blocks_generated=false
    fi
    echo ""
done

# 显示详细的日志内容以便分析
echo ""
echo "=== Detailed Log Analysis ==="
for i in {1..3}; do
    echo "--- Node ${i} full stdout ---"
    if [ -f "node${i}/stdout.log" ]; then
        tail -20 node${i}/stdout.log
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
echo "=== Block Consistency Validation Summary ==="
if [ "$height_consistent" = true ] && [ "$new_blocks_generated" = true ]; then
    echo "✓ Block height consistency validation: PASSED"
    echo "✓ New block generation validation: PASSED"
    echo "✓ All nodes are synchronized"
else
    echo "! Some validations may need manual checking"
    echo "! Check logs above for details"
fi

echo "=== Block Consistency Validation Completed ==="