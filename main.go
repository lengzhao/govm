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

	"github.com/lengzhao/binary"
	"github.com/lengzhao/govm/consensus"
	"github.com/lengzhao/govm/core"
	"github.com/lengzhao/govm/generator"
	"github.com/lengzhao/govm/storage"
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

// 命令行参数
var (
	nodeID     = flag.Int("node-id", 1, "Node ID")
	port       = flag.Int("port", 8000, "Port to listen on")
	dataDir    = flag.String("data-dir", "./data", "Data directory")
	configFile = flag.String("config", "./config/validators.json", "Validators configuration file")
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
	blockGenerator := generator.NewBlockGenerator(cons, store)

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
	go startBlockGeneration(coreModule, blockGenerator, cons, net, currentNodeAddr)

	fmt.Printf("govm 服务已启动 (Node ID: %d, Port: %d)\n", *nodeID, *port)

	// 等待中断信号以优雅地关闭所有服务
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	fmt.Println("govm 服务正在关闭...")
}

// startBlockGeneration 启动出块循环
func startBlockGeneration(core core.Core, generator generator.BlockGenerator, cons consensus.PoAConsensus, net network.NetworkInterface, nodeAddr types.Address) {
	ticker := time.NewTicker(time.Duration(types.BlockInterval) * time.Millisecond)
	defer ticker.Stop()

	fmt.Println("出块循环已启动...")

	// 注册新区块消息处理器
	net.RegisterMessageHandler("new_block", func(from string, topic string, data []byte) error {
		fmt.Printf("Received new block from %s\n", from)
		var block types.Block
		if err := binary.Unmarshal(data, &block); err != nil {
			fmt.Printf("反序列化区块失败: %v\n", err)
			return err
		}
		fmt.Printf("Received block: %d\n", block.Header.BlockNumber)
		err := core.AddBlock(&block)
		if err != nil {
			fmt.Printf("添加区块失败: %v\n", err)
			return err
		}

		// 这里应该处理接收到的新区块
		return nil
	})

	fmt.Println("等待10秒以确保网络模块启动...")
	time.Sleep(10 * time.Second)
	fmt.Println("开始生成区块")

	for range ticker.C {
		// 获取最新的区块作为前一个区块
		lastBlock := core.GetLastBlock()

		// 计算下一个区块的高度
		nextBlockHeight := lastBlock.Header.BlockNumber + 1

		// 检查当前是否轮到本节点出块
		if !cons.IsMyTurn(nextBlockHeight, nodeAddr) {
			// 如果不是轮到本节点出块，跳过
			fmt.Printf("不是轮到本节点出块，跳过 (高度: %d)\n", nextBlockHeight)
			continue
		}

		fmt.Printf("轮到本节点出块 (高度: %d)\n", nextBlockHeight)

		// 生成新区块
		block, err := generator.GenerateBlock(lastBlock)
		if err != nil {
			fmt.Printf("生成区块失败: %v\n", err)
			continue
		}

		data, err := binary.Marshal(block)
		if err != nil {
			fmt.Printf("序列化区块失败: %v\n", err)
			continue
		}

		// 添加区块到区块链
		if err := core.AddBlock(block); err != nil {
			fmt.Printf("添加区块失败: %v\n", err)
			continue
		}

		fmt.Printf("成功生成并添加区块，高度: %d\n", block.Header.BlockNumber)

		// 广播新区块到网络
		net.BroadcastMessage("new_block", data)
	}

}
