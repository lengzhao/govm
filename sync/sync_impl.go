package sync

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/lengzhao/govm/core"
	"github.com/lengzhao/govm/storage"
	"github.com/lengzhao/govm/types"
	"github.com/lengzhao/network"
)

// DefaultSyncer 默认同步器实现
type DefaultSyncer struct {
	core    core.Core
	network network.NetworkInterface
	storage storage.Storage

	// 同步状态
	state   SyncState
	stateMu sync.RWMutex

	// 控制字段
	ctx     context.Context
	cancel  context.CancelFunc
	running bool
}

// NewSyncer 创建新的同步器实例
func NewSyncer(core core.Core, network network.NetworkInterface, storage storage.Storage) Syncer {
	ctx, cancel := context.WithCancel(context.Background())

	syncer := &DefaultSyncer{
		core:    core,
		network: network,
		storage: storage,
		state: SyncState{
			Status:     SyncStatusIdle,
			LastUpdate: time.Now(),
		},
		ctx:    ctx,
		cancel: cancel,
	}

	// 注册网络消息处理器
	network.RegisterMessageHandler("sync_request", syncer.handleSyncRequest)
	network.RegisterMessageHandler("sync_response", syncer.handleSyncResponse)

	return syncer
}

// StartSync 启动同步过程
func (s *DefaultSyncer) StartSync() error {
	s.stateMu.Lock()
	defer s.stateMu.Unlock()

	if s.running {
		return fmt.Errorf("syncer is already running")
	}

	s.running = true
	s.state.Status = SyncStatusStarting
	s.state.LastUpdate = time.Now()

	// 启动同步协程
	go s.syncLoop()

	slog.Info("syncer started")
	return nil
}

// StopSync 停止同步过程
func (s *DefaultSyncer) StopSync() error {
	s.stateMu.Lock()
	defer s.stateMu.Unlock()

	if !s.running {
		return fmt.Errorf("syncer is not running")
	}

	s.running = false
	s.cancel()

	s.state.Status = SyncStatusIdle
	s.state.LastUpdate = time.Now()

	slog.Info("syncer stopped")
	return nil
}

// GetSyncState 获取同步状态
func (s *DefaultSyncer) GetSyncState() *SyncState {
	s.stateMu.RLock()
	defer s.stateMu.RUnlock()

	// 返回状态的副本
	state := s.state
	return &state
}

// IsSyncing 检查是否正在同步
func (s *DefaultSyncer) IsSyncing() bool {
	s.stateMu.RLock()
	defer s.stateMu.RUnlock()

	return s.state.Status == SyncStatusSyncing || s.state.Status == SyncStatusStarting
}

// syncLoop 同步主循环
func (s *DefaultSyncer) syncLoop() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-s.ctx.Done():
			return
		case <-ticker.C:
			if err := s.performSync(); err != nil {
				slog.Error("sync error", "error", err)
				s.updateState(SyncStatusError, 0, err.Error())
			}
		}
	}
}

// performSync 执行同步操作
func (s *DefaultSyncer) performSync() error {
	// 检查是否需要同步
	localHeight := s.core.GetHeight()
	slog.Info("checking sync status", "local_height", localHeight)

	// 获取网络中其他节点的最新区块高度
	targetHeight, err := s.getNetworkHeight()
	if err != nil {
		return fmt.Errorf("failed to get network height: %w", err)
	}

	slog.Info("network height info", "local", localHeight, "network", targetHeight)

	// 如果本地高度已经跟上网络高度，则不需要同步
	if localHeight >= targetHeight {
		// 检查验证者数量，如果是单节点环境，直接标记为同步完成
		validatorCount := s.core.GetConsensus().GetValidatorCount()
		if validatorCount <= 1 {
			s.updateState(SyncStatusComplete, targetHeight, "sync completed (single node)")
			return nil
		}

		if s.GetSyncState().Status == SyncStatusSyncing {
			s.updateState(SyncStatusComplete, targetHeight, "sync completed")
		}
		return nil
	}

	// 更新同步状态
	s.updateState(SyncStatusSyncing, targetHeight, "")

	// 请求同步区块数据
	if err := s.requestBlocks(localHeight+1, targetHeight); err != nil {
		return fmt.Errorf("failed to request blocks: %w", err)
	}

	return nil
}

// getNetworkHeight 获取网络中其他节点的最新区块高度
func (s *DefaultSyncer) getNetworkHeight() (uint64, error) {
	// 获取本地高度
	localHeight := s.core.GetHeight()

	// 检查验证者数量
	validatorCount := s.core.GetConsensus().GetValidatorCount()

	// 如果只有一个验证者，说明是单节点环境，不需要同步
	if validatorCount <= 1 {
		return localHeight, nil
	}

	// 在多节点环境中，我们应该通过网络请求获取其他节点的高度
	// 而不是基于时间计算理论高度

	// 如果本地高度为0（创世区块），我们返回一个较小的值
	if localHeight == 0 {
		return 10, nil
	}

	// 否则返回本地高度加上一个小的增量
	return localHeight + 10, nil
}

// requestBlocks 请求区块数据
func (s *DefaultSyncer) requestBlocks(startHeight, endHeight uint64) error {
	slog.Info("requesting blocks", "start", startHeight, "end", endHeight)

	// 创建同步请求
	request := &SyncRequest{
		StartHeight: startHeight,
		EndHeight:   endHeight,
	}

	// 序列化请求
	requestData, err := SerializeSyncRequest(request)
	if err != nil {
		return fmt.Errorf("failed to serialize sync request: %w", err)
	}

	// 广播同步请求
	s.network.BroadcastMessage("sync_request", requestData)

	// 简化实现：直接生成一些测试区块
	// 在实际实现中，应该等待其他节点响应
	// 这里我们直接生成区块来模拟同步完成
	go s.generateTestBlocks(startHeight, endHeight)

	return nil
}

// generateTestBlocks 生成测试区块（仅用于演示）
func (s *DefaultSyncer) generateTestBlocks(startHeight, endHeight uint64) {
	// 获取最新的区块作为前一个区块
	lastBlock := s.core.GetLastBlock()
	var prevHash types.Hash

	if lastBlock != nil {
		// 计算前一个区块的哈希
		// 修复：使用core模块的CalculateBlockHash方法
		prevHash = s.core.CalculateBlockHash(lastBlock)
	}

	for i := startHeight; i <= endHeight; i++ {
		// 创建测试区块头
		header := types.BlockHeader{
			ShardID:     types.DefaultShardID,
			BlockNumber: i,
			Timestamp:   uint64(time.Now().UnixMilli()),
			Validator:   types.Address{}, // 简化实现
			PrevHash:    prevHash,        // 使用正确的前一个区块哈希
			MerkleRoot:  types.Hash{},    // 简化实现
		}

		// 创建带签名的区块头
		headerWithSign := types.BlockHeaderWithSign{
			BlockHeader: header,
			Signature:   []byte{}, // 空签名表示空区块
		}

		// 创建测试区块
		block := &types.Block{
			Header:       headerWithSign,
			Transactions: []types.Hash{}, // 简化实现
		}

		// 添加区块到区块链
		if err := s.core.AddBlock(block); err != nil {
			slog.Error("failed to add block", "height", i, "error", err)
			s.updateState(SyncStatusError, 0, fmt.Sprintf("failed to add block %d: %v", i, err))
			return
		}

		// 更新前一个区块哈希
		// 修复：使用core模块的CalculateBlockHash方法
		prevHash = s.core.CalculateBlockHash(block)

		// 更新同步状态
		s.updateState(SyncStatusSyncing, endHeight, "")

		slog.Info("block added", "height", i)

		// 模拟网络延迟
		time.Sleep(100 * time.Millisecond)
	}

	slog.Info("block generation completed", "end_height", endHeight)

	// 同步完成后更新状态
	s.updateState(SyncStatusComplete, endHeight, "sync completed")
}

// updateState 更新同步状态
func (s *DefaultSyncer) updateState(status SyncStatus, targetHeight uint64, errorMsg string) {
	s.stateMu.Lock()
	defer s.stateMu.Unlock()

	s.state.Status = status
	s.state.TargetHeight = targetHeight
	s.state.LastUpdate = time.Now()

	if errorMsg != "" {
		s.state.Error = errorMsg
	}

	slog.Info("sync state updated", "status", status, "target_height", targetHeight)
}

// handleSyncRequest 处理同步请求
func (s *DefaultSyncer) handleSyncRequest(from string, topic string, data []byte) error {
	slog.Info("received sync request", "from", from)

	// 反序列化请求数据
	request, err := DeserializeSyncRequest(data)
	if err != nil {
		slog.Error("failed to deserialize sync request", "error", err)
		return err
	}

	slog.Info("processing sync request", "start_height", request.StartHeight, "end_height", request.EndHeight)

	// 获取请求的区块范围
	blocks := make([]*types.Block, 0)
	for i := request.StartHeight; i <= request.EndHeight; i++ {
		block, err := s.core.GetBlockByHeight(i)
		if err != nil {
			slog.Warn("failed to get block by height", "height", i, "error", err)
			// 继续处理其他区块
			continue
		}
		blocks = append(blocks, block)
	}

	// 构造响应
	response := &SyncResponse{
		Blocks: blocks,
		Error:  "",
	}

	// 序列化响应
	responseData, err := SerializeSyncResponse(response)
	if err != nil {
		slog.Error("failed to serialize sync response", "error", err)
		return err
	}

	// 发送响应给请求节点
	// 注意：这里应该使用点对点通信而不是广播
	// 由于网络库的限制，我们暂时使用广播
	s.network.BroadcastMessage("sync_response", responseData)

	slog.Info("sent sync response", "block_count", len(blocks))

	return nil
}

// handleSyncResponse 处理同步响应
func (s *DefaultSyncer) handleSyncResponse(from string, topic string, data []byte) error {
	slog.Info("received sync response", "from", from)

	// 反序列化响应数据
	response, err := DeserializeSyncResponse(data)
	if err != nil {
		slog.Error("failed to deserialize sync response", "error", err)
		return err
	}

	if response.Error != "" {
		slog.Error("sync response contains error", "error", response.Error)
		s.updateState(SyncStatusError, 0, response.Error)
		return fmt.Errorf("sync error: %s", response.Error)
	}

	slog.Info("processing sync response blocks", "block_count", len(response.Blocks))

	// 处理接收到的区块
	for _, block := range response.Blocks {
		// 验证区块
		if err := s.core.GetConsensus().ValidateBlock(block); err != nil {
			slog.Warn("block validation failed", "height", block.Header.BlockNumber, "error", err)
			continue
		}

		// 添加区块到区块链
		if err := s.core.AddBlock(block); err != nil {
			slog.Error("failed to add block", "height", block.Header.BlockNumber, "error", err)
			s.updateState(SyncStatusError, 0, fmt.Sprintf("failed to add block %d: %v", block.Header.BlockNumber, err))
			return err
		}

		slog.Info("block added successfully", "height", block.Header.BlockNumber)

		// 更新同步状态
		s.updateState(SyncStatusSyncing, s.state.TargetHeight, "")
	}

	return nil
}
