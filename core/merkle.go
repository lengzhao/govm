package core

import (
	"github.com/lengzhao/govm/crypto"
	"github.com/lengzhao/govm/types"
)

// MerkleTree Merkle树结构
type MerkleTree struct {
	Root   *MerkleNode
	Leaves []*MerkleNode
	crypto crypto.Crypto
}

// MerkleNode Merkle树节点
type MerkleNode struct {
	Left  *MerkleNode
	Right *MerkleNode
	Hash  types.Hash
}

// NewMerkleTree 创建新的Merkle树
func NewMerkleTree(data []types.Hash) *MerkleTree {
	if len(data) == 0 {
		return &MerkleTree{
			Root:   &MerkleNode{Hash: types.Hash{}},
			Leaves: []*MerkleNode{},
			crypto: crypto.NewCrypto(),
		}
	}

	// 创建叶子节点
	leaves := make([]*MerkleNode, len(data))
	for i, datum := range data {
		leaves[i] = &MerkleNode{Hash: datum}
	}

	// 构建树
	nodes := leaves
	for len(nodes) > 1 {
		var newLevel []*MerkleNode
		for i := 0; i < len(nodes); i += 2 {
			left := nodes[i]
			var right *MerkleNode
			if i+1 < len(nodes) {
				right = nodes[i+1]
			} else {
				// 如果节点数为奇数，复制最后一个节点
				right = left
			}
			newNode := &MerkleNode{
				Left:  left,
				Right: right,
				Hash:  combineAndHash(left.Hash, right.Hash),
			}
			newLevel = append(newLevel, newNode)
		}
		nodes = newLevel
	}

	return &MerkleTree{
		Root:   nodes[0],
		Leaves: leaves,
		crypto: crypto.NewCrypto(),
	}
}

// GetRootHash 获取Merkle根哈希
func (mt *MerkleTree) GetRootHash() types.Hash {
	if mt.Root == nil {
		return types.Hash{}
	}
	return mt.Root.Hash
}

// combineAndHash 组合两个哈希并计算新的哈希
func combineAndHash(left, right types.Hash) types.Hash {
	combined := make([]byte, len(left)+len(right))
	copy(combined[:len(left)], left[:])
	copy(combined[len(left):], right[:])

	cryptoInstance := crypto.NewCrypto()
	return cryptoInstance.Hash(combined)
}
