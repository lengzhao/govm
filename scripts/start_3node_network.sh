#!/bin/bash

# 3节点网络验证脚本
# 该脚本用于启动和验证一个包含3个节点的区块链网络

set -e  # 遇到错误时退出

echo "Starting 3-node blockchain network..."

# 创建数据目录
mkdir -p node1/data node2/data node3/data config

# 创建创世区块配置文件，设置时间为当前时间+10秒
GENESIS_TIME=$(($(date +%s)*1000 + 10000))
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
    
    # 启动节点（后台运行）
    ./govm --node-id=${node_id} --port=${port} --data-dir=${data_dir} --config=./config/validators.json --genesis=./config/genesis.json &
    
    # 保存进程ID
    echo $! > node${node_id}/govm.pid
    
    # 等待节点启动
    sleep 2
}

# 启动3个节点
start_node 1 8001
start_node 2 8002
start_node 3 8003

echo "All nodes started successfully!"

# 等待一段时间让网络稳定
echo "Waiting for network to stabilize..."
sleep 10

# 验证节点是否正常运行
echo "Checking node status..."
for i in {1..3}; do
    if [ -f "node${i}/govm.pid" ]; then
        pid=$(cat node${i}/govm.pid)
        if ps -p $pid > /dev/null; then
            echo "Node ${i} is running (PID: $pid)"
        else
            echo "Node ${i} is not running"
        fi
    else
        echo "Node ${i} PID file not found"
    fi
done

echo "3-node network is now running!"
echo "To stop the network, run: ./scripts/stop_3node_network.sh"