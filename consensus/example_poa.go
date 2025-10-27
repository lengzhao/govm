package consensus

import (
	"fmt"
	"time"

	"github.com/lengzhao/govm/storage"
	"github.com/lengzhao/govm/types"
)

// ExamplePoAUsage 演示如何使用PoA共识机制
func ExamplePoAUsage() {
	// 创建验证节点地址列表
	validators := make([]types.Address, 21)
	for i := 0; i < 21; i++ {
		var addr types.Address
		// 简化实现，实际应该使用真实的地址
		copy(addr[:], []byte(fmt.Sprintf("validator-%02d", i)))
		validators[i] = addr
	}

	// 创建PoA配置
	config := &PoAConfig{
		Validators:    validators,
		BlockInterval: 2000, // 2秒
		RoundLength:   21,   // 21个区块为一轮
	}

	// 创建PoA共识实例
	poa := NewPoAConsensus(config, storage.NewMemoryStorage(""))

	// 获取当前验证者
	currentValidator := poa.GetValidator()
	fmt.Printf("Current validator: %+v\n", currentValidator)

	// 获取所有验证者
	allValidators := poa.GetValidators()
	fmt.Printf("Total validators: %d\n", len(allValidators))

	// 获取当前轮次
	round := poa.GetRound()
	fmt.Printf("Current round: %d\n", round)

	// 获取当前轮值
	turn := poa.GetTurn()
	fmt.Printf("Current turn: %d\n", turn)

	// 检查特定节点是否为验证者
	var testAddr types.Address
	copy(testAddr[:], []byte("validator-05"))
	isValidator := poa.IsValidator(testAddr)
	fmt.Printf("Is 'validator-05' a validator: %t\n", isValidator)

	// 演示空区块功能
	fmt.Println("\n--- 空区块功能演示 ---")

	// 创建一个空区块（使用空签名）
	emptyBlock := &types.Block{
		Header: types.BlockHeaderWithSign{
			BlockHeader: types.BlockHeader{
				ShardID:     1,
				BlockNumber: 1,
				Timestamp:   uint64(time.Now().UnixNano() / 1000000),
				PrevHash:    types.Hash{},
				MerkleRoot:  types.Hash{},
				OtherShards: [3]types.Hash{},
			},
			Signature: []byte{}, // 空签名
		},
		Transactions: []types.Hash{},
	}

	// 验证空区块（由于简化实现，这里可能会失败，但我们主要测试接口）
	err := poa.ValidateBlock(emptyBlock)
	if err != nil {
		fmt.Printf("Empty block validation failed: %v\n", err)
	} else {
		fmt.Println("Empty block validation passed")
	}

	// 演示节点管理功能
	fmt.Println("\n--- 节点管理功能演示 ---")

	// 设置新的验证节点列表
	newValidators := make([]types.Address, 21)
	for i := 0; i < 21; i++ {
		var addr types.Address
		copy(addr[:], []byte(fmt.Sprintf("new-validator-%02d", i)))
		newValidators[i] = addr
	}

	err = poa.UpdateValidators(newValidators)
	if err != nil {
		fmt.Printf("Failed to update validators: %v\n", err)
	} else {
		fmt.Println("Successfully updated validators")
		fmt.Printf("New validators count: %d\n", len(newValidators))
	}

	// 检查新验证节点
	var newValidatorAddr types.Address
	copy(newValidatorAddr[:], []byte("new-validator-05"))
	isNewValidator := poa.IsValidator(newValidatorAddr)
	fmt.Printf("Is 'new-validator-05' a validator: %t\n", isNewValidator)
}
