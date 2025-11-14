package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestAPIInterface 测试API接口
func TestAPIInterface(t *testing.T) {
	// 测试NodeInfo结构
	nodeInfo := NodeInfo{
		ID:      "test-node",
		Address: "localhost:8080",
		Status:  "running",
	}
	assert.Equal(t, "test-node", nodeInfo.ID)
	assert.Equal(t, "localhost:8080", nodeInfo.Address)
	assert.Equal(t, "running", nodeInfo.Status)

	// 测试Metrics结构
	metrics := Metrics{
		BlockHeight:      100,
		TransactionCount: 500,
		PeerCount:        3,
	}
	assert.Equal(t, uint64(100), metrics.BlockHeight)
	assert.Equal(t, uint64(500), metrics.TransactionCount)
	assert.Equal(t, 3, metrics.PeerCount)
}
