package api

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/lengzhao/govm/consensus"
	"github.com/lengzhao/govm/core"
	"github.com/lengzhao/govm/storage"
	"github.com/lengzhao/govm/txpool"
	"github.com/lengzhao/govm/types"
	"github.com/stretchr/testify/assert"
)

// TestAccountPersistence 测试账号持久化功能
func TestAccountPersistence(t *testing.T) {
	// 创建临时钱包目录
	tmpWalletDir := "./test_persistence_wallet_" + fmt.Sprintf("%d", time.Now().UnixNano())
	os.MkdirAll(tmpWalletDir, 0755)
	defer os.RemoveAll(tmpWalletDir)

	// 创建存储
	store := storage.NewMemoryStorage("")
	err := store.Start()
	assert.NoError(t, err)
	defer store.Stop()

	// 创建模拟网络
	mockNetwork := &MockNetwork{
		peers: []string{"peer1", "peer2"},
	}

	// 创建共识模块
	config := &consensus.PoAConfig{
		Validators:    []types.Address{{1}, {2}, {3}},
		BlockInterval: 2000,
		RoundLength:   3,
	}
	cons := consensus.NewPoAConsensus(config, store)

	// 创建核心模块
	coreConfig := &core.CoreConfig{
		ShardID: types.DefaultShardID,
		DataDir: "./test_data",
		Genesis: &types.GenesisConfig{
			Timestamp: uint64(time.Now().Unix()),
		},
	}
	coreModule, err := core.NewCore(coreConfig, cons, store)
	assert.NoError(t, err)

	err = coreModule.Start()
	assert.NoError(t, err)
	defer coreModule.Stop()

	// 创建交易池
	txPool := txpool.NewTxPool(coreModule, store)
	err = txPool.Start()
	assert.NoError(t, err)
	defer txPool.Stop()

	// 创建第一个API实例
	apiServer1 := &DefaultAPI{
		core:      coreModule,
		txPool:    txPool,
		storage:   store,
		network:   mockNetwork,
		port:      ":8080",
		accounts:  make(map[types.Address][]byte),
		walletDir: tmpWalletDir,
	}

	// 创建账户
	address1, err := apiServer1.CreateAccount()
	assert.NoError(t, err)
	assert.NotEqual(t, types.Address{}, address1)

	// 获取账户列表
	accounts1, err := apiServer1.WalletAPI().GetAccounts()
	assert.NoError(t, err)
	assert.Len(t, accounts1, 1)
	assert.Contains(t, accounts1, address1)

	// 创建第二个API实例（模拟服务重启）
	apiServer2 := &DefaultAPI{
		core:      coreModule,
		txPool:    txPool,
		storage:   store,
		network:   mockNetwork,
		port:      ":8080",
		accounts:  make(map[types.Address][]byte),
		walletDir: tmpWalletDir,
	}

	// 从文件加载账户
	err = apiServer2.loadAccountsFromFiles()
	assert.NoError(t, err)

	// 验证账户已从文件加载
	accounts2, err := apiServer2.WalletAPI().GetAccounts()
	assert.NoError(t, err)
	assert.Len(t, accounts2, 1)
	assert.Contains(t, accounts2, address1)
	assert.Contains(t, accounts2, address1)

	// 创建另一个账户
	address2, err := apiServer2.CreateAccount()
	assert.NoError(t, err)
	assert.NotEqual(t, types.Address{}, address2)
	assert.NotEqual(t, address1, address2)

	// 验证两个账户都存在
	accounts3, err := apiServer2.WalletAPI().GetAccounts()
	assert.NoError(t, err)
	assert.Len(t, accounts3, 2)
	assert.Contains(t, accounts3, address1)
	assert.Contains(t, accounts3, address2)

	// 创建第三个API实例（再次模拟服务重启）
	apiServer3 := &DefaultAPI{
		core:      coreModule,
		txPool:    txPool,
		storage:   store,
		network:   mockNetwork,
		port:      ":8080",
		accounts:  make(map[types.Address][]byte),
		walletDir: tmpWalletDir,
	}

	// 从文件加载账户
	err = apiServer3.loadAccountsFromFiles()
	assert.NoError(t, err)

	// 验证所有账户都已从文件加载
	accounts4, err := apiServer3.WalletAPI().GetAccounts()
	assert.NoError(t, err)
	assert.Len(t, accounts4, 2)
	assert.Contains(t, accounts4, address1)
	assert.Contains(t, accounts4, address2)
	assert.Contains(t, accounts4, address1)
	assert.Contains(t, accounts4, address2)
}
