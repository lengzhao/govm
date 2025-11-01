#!/bin/bash

# 停止3节点网络脚本
# 该脚本用于停止所有正在运行的节点

echo "Stopping 3-node blockchain network..."

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

# 停止所有节点
for i in {1..3}; do
    stop_node $i
done

echo "All nodes stopped!"