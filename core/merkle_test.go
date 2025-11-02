package core

import (
	"testing"

	"github.com/lengzhao/govm/types"
	"github.com/stretchr/testify/assert"
)

func TestNewMerkleTree(t *testing.T) {
	// 测试空数据
	tree := NewMerkleTree([]types.Hash{})
	assert.NotNil(t, tree)
	assert.Equal(t, types.Hash{}, tree.GetRootHash())

	// 测试单个数据
	hash1 := types.Hash{1}
	tree = NewMerkleTree([]types.Hash{hash1})
	assert.NotNil(t, tree)
	assert.Equal(t, hash1, tree.GetRootHash())

	// 测试两个数据
	hash2 := types.Hash{2}
	tree = NewMerkleTree([]types.Hash{hash1, hash2})
	assert.NotNil(t, tree)
	// 验证根哈希不为空且不等于任何一个叶子节点
	assert.NotEqual(t, types.Hash{}, tree.GetRootHash())
	assert.NotEqual(t, hash1, tree.GetRootHash())
	assert.NotEqual(t, hash2, tree.GetRootHash())

	// 测试三个数据
	hash3 := types.Hash{3}
	tree = NewMerkleTree([]types.Hash{hash1, hash2, hash3})
	assert.NotNil(t, tree)
	// 验证根哈希不为空
	assert.NotEqual(t, types.Hash{}, tree.GetRootHash())
}

func TestMerkleTree_GetRootHash(t *testing.T) {
	// 测试空树
	tree := NewMerkleTree([]types.Hash{})
	rootHash := tree.GetRootHash()
	assert.Equal(t, types.Hash{}, rootHash)

	// 测试单节点树
	hash := types.Hash{1, 2, 3}
	tree = NewMerkleTree([]types.Hash{hash})
	rootHash = tree.GetRootHash()
	assert.Equal(t, hash, rootHash)
}