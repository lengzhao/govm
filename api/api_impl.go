package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/lengzhao/govm/core"
	"github.com/lengzhao/govm/crypto"
	"github.com/lengzhao/govm/storage"
	"github.com/lengzhao/govm/txpool"
	"github.com/lengzhao/govm/types"
	"github.com/lengzhao/network"
)

// DefaultAPI 默认API实现
type DefaultAPI struct {
	core    core.Core
	txPool  txpool.TxPool
	storage storage.Storage
	network network.NetworkInterface
	crypto  crypto.Crypto

	server *http.Server
	port   string
}

// NewAPI 创建新的API实例
func NewAPI(core core.Core, txPool txpool.TxPool, storage storage.Storage, network network.NetworkInterface) API {
	return &DefaultAPI{
		core:    core,
		txPool:  txPool,
		storage: storage,
		network: network,
		crypto:  crypto.NewCrypto(),
		port:    ":8080", // 默认端口
	}
}

// Start 启动API服务
func (a *DefaultAPI) Start() error {
	mux := http.NewServeMux()

	// 注册API端点
	a.RegisterEndpoints(mux)

	// 创建HTTP服务器
	a.server = &http.Server{
		Addr:    a.port,
		Handler: mux,
	}

	// 启动HTTP服务器
	go func() {
		fmt.Printf("API server starting on %s\n", a.port)
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("API server error: %v\n", err)
		}
	}()

	return nil
}

// Stop 停止API服务
func (a *DefaultAPI) Stop() error {
	if a.server != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return a.server.Shutdown(ctx)
	}
	return nil
}

// SetPort 设置API服务端口
func (a *DefaultAPI) SetPort(port string) {
	a.port = port
}

// RegisterEndpoints 注册API端点
func (a *DefaultAPI) RegisterEndpoints(mux *http.ServeMux) {
	mux.HandleFunc("/block/hash/", a.getBlockByHashHandler)
	mux.HandleFunc("/block/number/", a.getBlockByNumberHandler)
	mux.HandleFunc("/transaction/", a.getTransactionHandler)
	mux.HandleFunc("/transaction/send", a.sendTransactionHandler)
	mux.HandleFunc("/account/balance/", a.getBalanceHandler)
	mux.HandleFunc("/node/info", a.getNodeInfoHandler)
	mux.HandleFunc("/node/peers", a.getPeersHandler)
}

// getBlockByHashHandler 根据哈希获取区块
func (a *DefaultAPI) getBlockByHashHandler(w http.ResponseWriter, r *http.Request) {
	// 从URL中提取哈希值
	hashStr := r.URL.Path[len("/block/hash/"):]
	if hashStr == "" {
		http.Error(w, "Missing block hash", http.StatusBadRequest)
		return
	}

	// 解析哈希
	var hash types.Hash
	if len(hashStr) != 64 {
		http.Error(w, "Invalid hash length", http.StatusBadRequest)
		return
	}

	// 将十六进制字符串转换为哈希
	// 简化实现，实际应该使用hex.DecodeString
	copy(hash[:], []byte(hashStr)[:32])

	// 获取区块
	block, err := a.GetBlockByHash(hash)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// 返回JSON响应
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(block)
}

// getBlockByNumberHandler 根据区块号获取区块
func (a *DefaultAPI) getBlockByNumberHandler(w http.ResponseWriter, r *http.Request) {
	// 从URL中提取区块号
	numberStr := r.URL.Path[len("/block/number/"):]
	if numberStr == "" {
		http.Error(w, "Missing block number", http.StatusBadRequest)
		return
	}

	// 解析区块号
	number, err := strconv.ParseUint(numberStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid block number", http.StatusBadRequest)
		return
	}

	// 获取区块
	block, err := a.GetBlockByNumber(number)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// 返回JSON响应
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(block)
}

// getTransactionHandler 获取交易
func (a *DefaultAPI) getTransactionHandler(w http.ResponseWriter, r *http.Request) {
	// 从URL中提取交易哈希
	hashStr := r.URL.Path[len("/transaction/"):]
	if hashStr == "" {
		http.Error(w, "Missing transaction hash", http.StatusBadRequest)
		return
	}

	// 解析哈希
	var hash types.Hash
	if len(hashStr) != 64 {
		http.Error(w, "Invalid hash length", http.StatusBadRequest)
		return
	}

	// 将十六进制字符串转换为哈希
	// 简化实现，实际应该使用hex.DecodeString
	copy(hash[:], []byte(hashStr)[:32])

	// 获取交易
	tx, err := a.GetTransaction(hash)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// 返回JSON响应
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tx)
}

// sendTransactionHandler 发送交易
func (a *DefaultAPI) sendTransactionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 解析请求体中的交易数据
	var tx types.Transaction
	if err := json.NewDecoder(r.Body).Decode(&tx); err != nil {
		http.Error(w, "Invalid transaction data", http.StatusBadRequest)
		return
	}

	// 发送交易
	txHash, err := a.SendTransaction(&tx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 返回交易哈希
	response := map[string]interface{}{
		"hash": txHash,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// getBalanceHandler 获取余额
func (a *DefaultAPI) getBalanceHandler(w http.ResponseWriter, r *http.Request) {
	// 从URL中提取地址
	addressStr := r.URL.Path[len("/account/balance/"):]
	if addressStr == "" {
		http.Error(w, "Missing address", http.StatusBadRequest)
		return
	}

	// 解析地址
	var address types.Address
	if len(addressStr) != 40 {
		http.Error(w, "Invalid address length", http.StatusBadRequest)
		return
	}

	// 将十六进制字符串转换为地址
	// 简化实现，实际应该使用hex.DecodeString
	copy(address[:], []byte(addressStr)[:20])

	// 获取余额
	balance, err := a.GetBalance(address)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// 返回余额
	response := map[string]interface{}{
		"balance": balance,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// getNodeInfoHandler 获取节点信息
func (a *DefaultAPI) getNodeInfoHandler(w http.ResponseWriter, r *http.Request) {
	info, err := a.GetNodeInfo()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(info)
}

// getPeersHandler 获取节点列表
func (a *DefaultAPI) getPeersHandler(w http.ResponseWriter, r *http.Request) {
	peers, err := a.GetPeers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(peers)
}

// GetBlockByHash 根据哈希获取区块
func (a *DefaultAPI) GetBlockByHash(hash types.Hash) (*types.Block, error) {
	return a.core.GetBlockByHash(hash)
}

// GetBlockByNumber 根据区块号获取区块
func (a *DefaultAPI) GetBlockByNumber(number uint64) (*types.Block, error) {
	return a.core.GetBlockByHeight(number)
}

// GetTransaction 获取交易
func (a *DefaultAPI) GetTransaction(hash types.Hash) (*types.Transaction, error) {
	tx, err := a.txPool.GetTransaction(hash)
	if err != nil {
		return nil, err
	}
	return &tx.Transaction, nil
}

// SendTransaction 发送交易
func (a *DefaultAPI) SendTransaction(tx *types.Transaction) (types.Hash, error) {
	// 创建带签名的交易
	txWithSign := &types.TransactionWithSign{
		Transaction: *tx,
		Signature:   nil, // 简化实现，实际应该进行签名
	}

	// 验证交易
	if err := a.txPool.ValidateTransaction(txWithSign); err != nil {
		return types.Hash{}, fmt.Errorf("transaction validation failed: %w", err)
	}

	// 添加到交易池
	if err := a.txPool.AddTransaction(txWithSign); err != nil {
		return types.Hash{}, fmt.Errorf("failed to add transaction to pool: %w", err)
	}

	// 计算交易哈希
	data, err := json.Marshal(tx)
	if err != nil {
		return types.Hash{}, fmt.Errorf("failed to marshal transaction: %w", err)
	}
	hash := a.crypto.Hash(data)

	return hash, nil
}

// GetBalance 获取余额
func (a *DefaultAPI) GetBalance(address types.Address) (uint64, error) {
	// 简化实现，实际应该从状态中获取余额
	// 这里返回一个默认值
	return 1000000, nil
}

// GetNodeInfo 获取节点信息
func (a *DefaultAPI) GetNodeInfo() (NodeInfo, error) {
	info := NodeInfo{
		ID:      "node-1",
		Address: "localhost:8080",
		Status:  "running",
	}
	return info, nil
}

// GetPeers 获取节点列表
func (a *DefaultAPI) GetPeers() ([]NodeInfo, error) {
	// 获取连接的节点列表
	peerIDs := a.network.GetPeers()

	peers := make([]NodeInfo, len(peerIDs))
	for i, peerID := range peerIDs {
		peers[i] = NodeInfo{
			ID:      peerID,
			Address: fmt.Sprintf("peer-%d", i),
			Status:  "connected",
		}
	}

	return peers, nil
}
