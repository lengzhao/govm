package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lengzhao/govm/consensus"
	"github.com/lengzhao/govm/core"
	"github.com/lengzhao/govm/generator"
	"github.com/lengzhao/govm/storage"
	"github.com/lengzhao/govm/sync"
	"github.com/lengzhao/govm/txpool"
	"github.com/lengzhao/govm/types"
	"github.com/lengzhao/network"
)

// Validator 验证节点配置
type Validator struct {
	ID        int    `json:"id"`
	Address   string `json:"address"`
	PublicKey string `json:"public_key"`
}

// ValidatorsConfig 验证节点配置文件结构
type ValidatorsConfig struct {
	Validators []Validator `json:"validators"`
}

// GenesisConfig 创世区块配置文件结构
type GenesisConfig struct {
	Genesis types.GenesisConfig `json:"genesis"`
}

// 命令行参数
var (
	nodeID      = flag.Int("node-id", 1, "Node ID")
	port        = flag.Int("port", 8000, "Port to listen on")
	dataDir     = flag.String("data-dir", "./data", "Data directory")
	configFile  = flag.String("config", "./config/validators.json", "Validators configuration file")
	genesisFile = flag.String("genesis", "./config/genesis.json", "Genesis configuration file")
)

// loadValidatorsFromConfig 从配置文件加载验证节点
func loadValidatorsFromConfig(configFile string) ([]types.Address, error) {
	// 读取配置文件
	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// 解析JSON配置
	var config ValidatorsConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// 转换为验证节点地址列表
	validators := make([]types.Address, len(config.Validators))
	for i, v := range config.Validators {
		addr := types.Address{}
		// 简化实现：使用地址字符串的前20个字节
		copy(addr[:], []byte(v.Address))
		validators[i] = addr
	}

	return validators, nil
}

// loadGenesisConfig 从配置文件加载创世区块配置
func loadGenesisConfig(genesisFile string) (*types.GenesisConfig, error) {
	// 检查配置文件是否存在
	if _, err := os.Stat(genesisFile); os.IsNotExist(err) {
		// 如果文件不存在，使用默认配置
		fmt.Printf("Genesis config file not found, using default timestamp\n")
		return &types.GenesisConfig{
			Timestamp: uint64(time.Now().Unix()),
		}, nil
	}

	// 读取配置文件
	data, err := os.ReadFile(genesisFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read genesis config file: %w", err)
	}

	// 解析JSON配置
	var config GenesisConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse genesis config file: %w", err)
	}

	return &config.Genesis, nil
}

func main() {
	// 解析命令行参数
	flag.Parse()

	fmt.Printf("govm - 高性能分片区块链平台 (Node ID: %d)\n", *nodeID)
	fmt.Printf("项目启动中... (Port: %d, DataDir: %s)\n", *port, *dataDir)

	// 初始化存储模块
	store := storage.NewMemoryStorage("")
	if err := store.Start(); err != nil {
		fmt.Printf("存储模块启动失败: %v\n", err)
		return
	}
	defer store.Stop()

	// 初始化网络模块
	netConfig := network.NewNetworkConfig()
	netConfig.Host = "0.0.0.0" // 监听所有接口
	netConfig.Port = *port
	netConfig.MaxPeers = 50
	netConfig.PrivateKeyPath = fmt.Sprintf("./node%d/private.key", *nodeID)

	net, err := network.New(netConfig)
	if err != nil {
		fmt.Printf("网络模块创建失败: %v\n", err)
		return
	}

	// 加载验证节点配置
	validators, err := loadValidatorsFromConfig(*configFile)
	if err != nil {
		fmt.Printf("加载验证节点配置失败: %v\n", err)
		return
	}

	fmt.Printf("加载了 %d 个验证节点\n", len(validators))
	for i, addr := range validators {
		fmt.Printf("验证节点 %d: %x\n", i+1, addr)
	}

	// 加载创世区块配置
	genesisConfig, err := loadGenesisConfig(*genesisFile)
	if err != nil {
		fmt.Printf("加载创世区块配置失败: %v\n", err)
		return
	}

	fmt.Printf("创世区块时间戳: %d\n", genesisConfig.Timestamp)

	// 初始化共识模块
	config := &consensus.PoAConfig{
		Validators:    validators,
		BlockInterval: types.BlockInterval,     // 2秒区块间隔
		RoundLength:   uint64(len(validators)), // 轮次长度等于验证节点数量
	}
	cons := consensus.NewPoAConsensus(config, store)

	// 初始化核心模块
	coreConfig := &core.CoreConfig{
		ShardID: types.DefaultShardID,
		DataDir: *dataDir,
		Genesis: genesisConfig, // 传递创世区块配置
	}
	coreModule, err := core.NewCore(coreConfig, cons, store)
	if err != nil {
		fmt.Printf("核心模块创建失败: %v\n", err)
		return
	}

	// 启动核心模块
	if err := coreModule.Start(); err != nil {
		fmt.Printf("核心模块启动失败: %v\n", err)
		return
	}
	defer coreModule.Stop()

	// 设置网络接口并注册消息处理器
	if err := coreModule.SetNetwork(net, *nodeID); err != nil {
		fmt.Printf("设置网络接口失败: %v\n", err)
		return
	}

	// 初始化交易池模块
	txPool := txpool.NewTxPool(coreModule, store)
	if err := txPool.Start(); err != nil {
		fmt.Printf("交易池模块启动失败: %v\n", err)
		return
	}
	defer txPool.Stop()

	// 初始化同步模块
	syncer := sync.NewSyncer(coreModule, net, store)

	// 启动同步模块
	if err := syncer.StartSync(); err != nil {
		fmt.Printf("同步模块启动失败: %v\n", err)
		// 不中断程序执行，继续启动其他模块
	} else {
		fmt.Println("govm 同步模块已启动")
	}
	defer syncer.StopSync()

	// 启动网络模块
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := net.Run(ctx); err != nil {
			fmt.Printf("网络模块启动失败: %v\n", err)
		}
	}()

	fmt.Println("govm 核心模块已启动")

	// 初始化区块生成器
	blockGenerator := generator.NewBlockGenerator(cons, store, txPool)

	// 获取当前节点地址（基于节点ID分配验证者地址）
	var currentNodeAddr types.Address
	if *nodeID <= len(validators) {
		currentNodeAddr = validators[*nodeID-1]
	} else {
		// 如果节点ID超出验证节点数量，使用循环分配
		currentNodeAddr = validators[(*nodeID-1)%len(validators)]
	}

	fmt.Printf("当前节点地址: %x\n", currentNodeAddr)

	// 启动出块循环
	go startBlockGeneration(coreModule, blockGenerator, cons, currentNodeAddr, syncer)

	fmt.Printf("govm 服务已启动 (Node ID: %d, Port: %d)\n", *nodeID, *port)

	// 等待中断信号以优雅地关闭所有服务
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	fmt.Println("govm 服务正在关闭...")
}

// startBlockGeneration 启动出块循环
func startBlockGeneration(coreModule core.Core, generator generator.BlockGenerator, cons consensus.PoAConsensus, nodeAddr types.Address, syncChecker types.SyncChecker) {
	// 启动区块生成器的出块循环
	if err := generator.StartBlockGeneration(coreModule, cons, nodeAddr, syncChecker); err != nil {
		fmt.Printf("启动区块生成循环失败: %v\n", err)
	}
}
