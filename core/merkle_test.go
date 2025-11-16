package core

import (
	"testing"

	"github.com/lengzhao/govm/types"
	"github.com/stretchr/testify/assert"
)

func TestMerkleTree_NewMerkleTree(t *testing.T) {
	// 测试空数据
	emptyTree := NewMerkleTree([]types.Hash{})
	assert.NotNil(t, emptyTree)
	assert.NotNil(t, emptyTree.Root)
	assert.Equal(t, types.Hash{}, emptyTree.Root.Hash)
	assert.Len(t, emptyTree.Leaves, 0)

	// 测试单个数据
	singleData := []types.Hash{{1}}
	singleTree := NewMerkleTree(singleData)
	assert.NotNil(t, singleTree)
	assert.NotNil(t, singleTree.Root)
	assert.NotEqual(t, types.Hash{}, singleTree.Root.Hash)
	assert.Len(t, singleTree.Leaves, 1)
	assert.Equal(t, singleData[0], singleTree.Leaves[0].Hash)

	// 测试多个数据
	multipleData := []types.Hash{{1}, {2}, {3}, {4}}
	multipleTree := NewMerkleTree(multipleData)
	assert.NotNil(t, multipleTree)
	assert.NotNil(t, multipleTree.Root)
	assert.NotEqual(t, types.Hash{}, multipleTree.Root.Hash)
	assert.Len(t, multipleTree.Leaves, 4)
}

func TestMerkleTree_GetRootHash(t *testing.T) {
	// 测试空树的根哈希
	emptyTree := NewMerkleTree([]types.Hash{})
	rootHash := emptyTree.GetRootHash()
	assert.Equal(t, types.Hash{}, rootHash)

	// 测试单个叶子节点的根哈希
	singleData := []types.Hash{{1}}
	singleTree := NewMerkleTree(singleData)
	rootHash = singleTree.GetRootHash()
	assert.NotEqual(t, types.Hash{}, rootHash)

	// 测试多个叶子节点的根哈希
	multipleData := []types.Hash{{1}, {2}, {3}, {4}}
	multipleTree := NewMerkleTree(multipleData)
	rootHash = multipleTree.GetRootHash()
	assert.NotEqual(t, types.Hash{}, rootHash)

	// 验证相同数据产生相同根哈希
	tree1 := NewMerkleTree(multipleData)
	tree2 := NewMerkleTree(multipleData)
	assert.Equal(t, tree1.GetRootHash(), tree2.GetRootHash())
}

func TestMerkleTree_Structure(t *testing.T) {
	// 测试4个叶子节点的树结构
	data := []types.Hash{{1}, {2}, {3}, {4}}
	tree := NewMerkleTree(data)

	// 验证叶子节点数量
	assert.Len(t, tree.Leaves, 4)

	// 验证根节点存在
	assert.NotNil(t, tree.Root)

	// 验证叶子节点数据正确
	assert.Equal(t, data[0], tree.Leaves[0].Hash)
	assert.Equal(t, data[1], tree.Leaves[1].Hash)
	assert.Equal(t, data[2], tree.Leaves[2].Hash)
	assert.Equal(t, data[3], tree.Leaves[3].Hash)
}

func TestMerkleTree_OddNumberOfLeaves(t *testing.T) {
	// 测试奇数个叶子节点
	data := []types.Hash{{1}, {2}, {3}}
	tree := NewMerkleTree(data)

	// 验证叶子节点数量
	assert.Len(t, tree.Leaves, 3)

	// 验证根节点存在
	assert.NotNil(t, tree.Root)

	// 验证根哈希不为空
	rootHash := tree.GetRootHash()
	assert.NotEqual(t, types.Hash{}, rootHash)
}
