package sync

import (
	"context"
	"testing"
	"time"

	"github.com/lengzhao/govm/consensus"
	"github.com/lengzhao/govm/core"
	"github.com/lengzhao/govm/storage"
	"github.com/lengzhao/govm/types"
	"github.com/lengzhao/network"
	"github.com/stretchr/testify/assert"
)

func TestSyncWithRealComponents(t *testing.T) {
	// 创建真实的存储实例
	store := storage.NewMemoryStorage("")
	err := store.Start()
	assert.NoError(t, err)
	defer store.Stop()

	// 创建网络配置
	netConfig := network.NewNetworkConfig()
	netConfig.Port = 0 // 使用随机端口

	// 创建网络实例
	net, err := network.New(netConfig)
	assert.NoError(t, err)

	// 启动网络
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		_ = net.Run(ctx)
	}()

	// 等待网络启动
	time.Sleep(100 * time.Millisecond)

	// 创建验证节点列表
	validators := make([]types.Address, 3)
	for i := 0; i < 3; i++ {
		var addr types.Address
		copy(addr[:], []byte{byte(i + 1)})
		validators[i] = addr
	}

	// 创建共识配置
	config := &consensus.PoAConfig{
		Validators:    validators,
		BlockInterval: 2000,
		RoundLength:   3,
	}

	// 创建共识实例
	cons := consensus.NewPoAConsensus(config, store)

	// 创建核心配置
	coreConfig := &core.CoreConfig{
		ShardID: types.DefaultShardID,
		DataDir: "./test_data",
	}

	// 创建核心实例
	coreModule, err := core.NewCore(coreConfig, cons, store)
	assert.NoError(t, err)

	// 启动核心模块
	err = coreModule.Start()
	assert.NoError(t, err)
	defer coreModule.Stop()

	// 创建同步器
	syncer := NewSyncer(coreModule, net, store)
	assert.NotNil(t, syncer)

	// 验证初始状态
	state := syncer.GetSyncState()
	assert.Equal(t, SyncStatusIdle, state.Status)

	// 启动同步器
	err = syncer.StartSync()
	assert.NoError(t, err)

	// 验证状态
	state = syncer.GetSyncState()
	assert.Equal(t, SyncStatusStarting, state.Status)

	// 等待一点时间让状态更新
	time.Sleep(100 * time.Millisecond)

	// 检查是否正在同步
	isSyncing := syncer.IsSyncing()
	assert.True(t, isSyncing)

	// 停止同步器
	err = syncer.StopSync()
	assert.NoError(t, err)

	// 验证状态
	state = syncer.GetSyncState()
	assert.Equal(t, SyncStatusIdle, state.Status)
}