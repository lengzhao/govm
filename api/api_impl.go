package api

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
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

	server *http.Server
	port   string

	// 钱包相关
	accounts  map[types.Address][]byte // 地址到私钥的映射
	walletDir string                   // 钱包数据目录
	mutex     sync.RWMutex             // 保护accounts的读写锁
}

// loadAccountsFromFiles 从文件加载账户
func (a *DefaultAPI) loadAccountsFromFiles() error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	// 检查钱包目录是否存在
	if _, err := os.Stat(a.walletDir); os.IsNotExist(err) {
		return nil // 目录不存在，无需加载
	}

	// 读取目录中的所有文件
	entries, err := os.ReadDir(a.walletDir)
	if err != nil {
		return fmt.Errorf("failed to read wallet directory: %w", err)
	}

	// 遍历所有.key文件
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if filepath.Ext(name) == ".key" {
			// 从文件名提取地址
			addressStr := strings.TrimSuffix(name, ".key")

			// 解析地址
			var address types.Address
			addrBytes, err := hex.DecodeString(addressStr)
			if err != nil {
				fmt.Printf("Warning: invalid address in filename %s: %v\n", name, err)
				continue
			}

			if len(addrBytes) != len(address) {
				fmt.Printf("Warning: invalid address length in filename %s\n", name)
				continue
			}

			copy(address[:], addrBytes)

			// 加载私钥
			filename := filepath.Join(a.walletDir, name)
			privateKey, err := crypto.LoadFromFile(filename, "")
			if err != nil {
				fmt.Printf("Warning: failed to load key from %s: %v\n", name, err)
				continue
			}

			// 保存到内存
			a.accounts[address] = privateKey.Bytes()
		}
	}

	return nil
}

// NewAPI 创建新的API实例
func NewAPI(core core.Core, txPool txpool.TxPool, storage storage.Storage, network network.NetworkInterface) API {
	walletDir := "./wallet" // 默认钱包目录

	// 创建钱包目录（如果不存在）
	os.MkdirAll(walletDir, 0755)

	api := &DefaultAPI{
		core:      core,
		txPool:    txPool,
		storage:   storage,
		network:   network,
		port:      ":8080", // 默认端口
		accounts:  make(map[types.Address][]byte),
		walletDir: walletDir,
	}

	// 从文件加载已保存的账户
	if err := api.loadAccountsFromFiles(); err != nil {
		fmt.Printf("Warning: failed to load accounts from files: %v\n", err)
	}

	return api
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
	mux.HandleFunc("/blockchain/height", a.getBlockchainHeightHandler)
	mux.HandleFunc("/blockchain/latest/hash", a.getLatestBlockHashHandler)
	mux.HandleFunc("/transaction/", a.getTransactionHandler)
	mux.HandleFunc("/transaction/send", a.sendTransactionHandler)
	mux.HandleFunc("/account/balance/", a.getBalanceHandler)
	mux.HandleFunc("/account/create", a.createAccountHandler)
	mux.HandleFunc("/account/info/", a.getAccountHandler)
	mux.HandleFunc("/node/info", a.getNodeInfoHandler)
	mux.HandleFunc("/node/peers", a.getPeersHandler)

	// 钱包相关端点
	mux.HandleFunc("/wallet/accounts", a.getAccountsHandler)
	mux.HandleFunc("/wallet/account/import", a.importAccountHandler)
	mux.HandleFunc("/wallet/account/export/", a.exportAccountHandler)
	mux.HandleFunc("/wallet/transaction/sign", a.signTransactionHandler)

	// 管理相关端点
	mux.HandleFunc("/admin/peer/add", a.addPeerHandler)
	mux.HandleFunc("/admin/peer/remove", a.removePeerHandler)
	mux.HandleFunc("/admin/mining/start", a.startMiningHandler)
	mux.HandleFunc("/admin/mining/stop", a.stopMiningHandler)
	mux.HandleFunc("/admin/metrics", a.getMetricsHandler)
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

// createAccountHandler 创建账户
func (a *DefaultAPI) createAccountHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 创建新账户
	address, err := a.CreateAccount()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 返回账户地址
	response := map[string]interface{}{
		"address": fmt.Sprintf("%x", address),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// getAccountHandler 获取账户信息
func (a *DefaultAPI) getAccountHandler(w http.ResponseWriter, r *http.Request) {
	// 从URL中提取地址
	addressStr := r.URL.Path[len("/account/info/"):]
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

	// 获取账户信息
	account, err := a.GetAccount(address)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// 返回账户信息
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(account)
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

// getBlockchainHeightHandler 获取区块链高度
func (a *DefaultAPI) getBlockchainHeightHandler(w http.ResponseWriter, r *http.Request) {
	// 获取当前区块高度
	height := a.core.GetHeight()

	// 返回区块高度
	response := map[string]interface{}{
		"height": height,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// getLatestBlockHashHandler 获取最新区块哈希
func (a *DefaultAPI) getLatestBlockHashHandler(w http.ResponseWriter, r *http.Request) {
	// 获取最新区块
	latestBlock := a.core.GetLastBlock()
	if latestBlock == nil {
		http.Error(w, "no blocks found", http.StatusNotFound)
		return
	}

	// 计算区块哈希
	data, err := json.Marshal(latestBlock.Header.BlockHeader)
	if err != nil {
		http.Error(w, "failed to marshal block header", http.StatusInternalServerError)
		return
	}
	blockHash := crypto.Hash(data)

	// 返回区块哈希
	response := map[string]interface{}{
		"hash": fmt.Sprintf("%x", blockHash),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
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
	hash := crypto.Hash(data)

	return hash, nil
}

// GetBalance 获取余额
func (a *DefaultAPI) GetBalance(address types.Address) (uint64, error) {
	// 简化实现，实际应该从状态中获取余额
	// 这里返回一个默认值
	return 1000000, nil
}

// GetAccount 获取账户信息
func (a *DefaultAPI) GetAccount(address types.Address) (*types.Account, error) {
	// 从核心模块获取账户信息
	// 这里简化实现，实际应该从状态中获取完整的账户信息
	account := &types.Account{
		Address: address,
		Balance: 1000000, // 默认余额
		Nonce:   0,
	}
	return account, nil
}

// saveAccountToFile 将账户保存到文件
func (a *DefaultAPI) saveAccountToFile(address types.Address, privKey []byte) error {
	// 创建文件路径
	filename := fmt.Sprintf("%s/%x.key", a.walletDir, address)

	// 从字节数据创建私钥对象
	algorithm, err := crypto.AlgorithmFactory(crypto.Ed25519)
	if err != nil {
		return fmt.Errorf("failed to get algorithm: %w", err)
	}

	privateKey, err := algorithm.PrivateKeyFromBytes(privKey)
	if err != nil {
		return fmt.Errorf("failed to create private key from bytes: %w", err)
	}

	// 使用默认密码保存私钥到文件
	return crypto.SaveToFile(privateKey, filename, "")
}

// CreateAccount 创建新账户
func (a *DefaultAPI) CreateAccount() (types.Address, error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	// 生成新的密钥对
	privKey, pubKey, err := crypto.GenerateKeyPair(crypto.Ed25519)
	if err != nil {
		return types.Address{}, fmt.Errorf("failed to generate key pair: %w", err)
	}

	// 从公钥生成地址
	address := crypto.GenerateAddress(pubKey, crypto.Ed25519)

	// 保存私钥到内存
	a.accounts[address] = privKey

	// 保存私钥到文件
	if err := a.saveAccountToFile(address, privKey); err != nil {
		return types.Address{}, fmt.Errorf("failed to save account to file: %w", err)
	}

	return address, nil
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

// WalletAPI 获取钱包API接口
func (a *DefaultAPI) WalletAPI() WalletAPI {
	return a
}

// AdminAPI 获取管理API接口
func (a *DefaultAPI) AdminAPI() AdminAPI {
	return a
}

// WalletAPI接口实现

// ImportAccount 导入账户
func (a *DefaultAPI) ImportAccount(privateKey []byte) (types.Address, error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	// 从私钥获取公钥
	algorithm, err := crypto.AlgorithmFactory(crypto.Ed25519)
	if err != nil {
		return types.Address{}, fmt.Errorf("failed to get algorithm: %w", err)
	}

	privKey, err := algorithm.PrivateKeyFromBytes(privateKey)
	if err != nil {
		return types.Address{}, fmt.Errorf("failed to create private key from bytes: %w", err)
	}

	pubKey := privKey.PublicKey()
	address := pubKey.Address()

	// 保存私钥到内存
	a.accounts[address] = privateKey

	// 保存私钥到文件
	if err := a.saveAccountToFile(address, privateKey); err != nil {
		return types.Address{}, fmt.Errorf("failed to save account to file: %w", err)
	}

	return address, nil
}

// ExportAccount 导出账户
func (a *DefaultAPI) ExportAccount(address types.Address) ([]byte, error) {
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	privKey, exists := a.accounts[address]
	if !exists {
		return nil, fmt.Errorf("account not found")
	}

	return privKey, nil
}

// SignTransaction 签名交易
func (a *DefaultAPI) SignTransaction(tx *types.Transaction, address types.Address) (*types.TransactionWithSign, error) {
	a.mutex.RLock()
	privKey, exists := a.accounts[address]
	a.mutex.RUnlock()

	if !exists {
		return nil, fmt.Errorf("account not found")
	}

	// 序列化交易数据用于签名
	data, err := json.Marshal(tx)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal transaction: %w", err)
	}

	// 签名交易
	signature, err := crypto.Sign(data, privKey, crypto.Ed25519)
	if err != nil {
		return nil, fmt.Errorf("failed to sign transaction: %w", err)
	}

	// 创建带签名的交易
	txWithSign := &types.TransactionWithSign{
		Transaction: *tx,
		Signature:   signature,
	}

	return txWithSign, nil
}

// GetAccounts 获取账户列表
func (a *DefaultAPI) GetAccounts() ([]types.Address, error) {
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	addresses := make([]types.Address, 0, len(a.accounts))
	for addr := range a.accounts {
		addresses = append(addresses, addr)
	}

	return addresses, nil
}

// AdminAPI接口实现

// AddPeer 添加节点
func (a *DefaultAPI) AddPeer(peerAddr string) error {
	return a.network.ConnectToPeer(peerAddr)
}

// RemovePeer 移除节点
func (a *DefaultAPI) RemovePeer(peerAddr string) error {
	// 网络库可能没有直接的移除节点方法
	// 这里我们简单返回nil表示未实现
	return nil
}

// StartMining 开始挖矿
func (a *DefaultAPI) StartMining() error {
	// 对于PoA共识，挖矿是由验证者自动进行的
	// 这里我们简单返回nil表示未实现
	return nil
}

// StopMining 停止挖矿
func (a *DefaultAPI) StopMining() error {
	// 对于PoA共识，挖矿是由验证者自动进行的
	// 这里我们简单返回nil表示未实现
	return nil
}

// GetMetrics 获取节点指标
func (a *DefaultAPI) GetMetrics() (Metrics, error) {
	metrics := Metrics{
		BlockHeight:      a.core.GetHeight(),
		TransactionCount: uint64(a.txPool.GetTransactionCount()),
		PeerCount:        len(a.network.GetPeers()),
	}
	return metrics, nil
}

// HTTP处理方法

// getAccountsHandler 获取账户列表
func (a *DefaultAPI) getAccountsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	addresses, err := a.GetAccounts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"accounts": addresses,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// importAccountHandler 导入账户
func (a *DefaultAPI) importAccountHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		PrivateKey string `json:"private_key"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		return
	}

	// 解码私钥（假设是十六进制编码）
	privateKey, err := hex.DecodeString(req.PrivateKey)
	if err != nil {
		http.Error(w, "Invalid private key format", http.StatusBadRequest)
		return
	}

	address, err := a.ImportAccount(privateKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"address": fmt.Sprintf("%x", address),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// exportAccountHandler 导出账户
func (a *DefaultAPI) exportAccountHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 从URL中提取地址
	addressStr := r.URL.Path[len("/wallet/account/export/"):]
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
	addrBytes, err := hex.DecodeString(addressStr)
	if err != nil {
		http.Error(w, "Invalid address format", http.StatusBadRequest)
		return
	}
	copy(address[:], addrBytes)

	privateKey, err := a.ExportAccount(address)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	response := map[string]interface{}{
		"private_key": fmt.Sprintf("%x", privateKey),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// signTransactionHandler 签名交易
func (a *DefaultAPI) signTransactionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Transaction types.Transaction `json:"transaction"`
		Address     string            `json:"address"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		return
	}

	// 解析地址
	var address types.Address
	if len(req.Address) != 40 {
		http.Error(w, "Invalid address length", http.StatusBadRequest)
		return
	}

	addrBytes, err := hex.DecodeString(req.Address)
	if err != nil {
		http.Error(w, "Invalid address format", http.StatusBadRequest)
		return
	}
	copy(address[:], addrBytes)

	signedTx, err := a.SignTransaction(&req.Transaction, address)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(signedTx)
}

// addPeerHandler 添加节点
func (a *DefaultAPI) addPeerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		PeerAddr string `json:"peer_addr"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		return
	}

	if err := a.AddPeer(req.PeerAddr); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// removePeerHandler 移除节点
func (a *DefaultAPI) removePeerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		PeerAddr string `json:"peer_addr"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		return
	}

	if err := a.RemovePeer(req.PeerAddr); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// startMiningHandler 开始挖矿
func (a *DefaultAPI) startMiningHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := a.StartMining(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// stopMiningHandler 停止挖矿
func (a *DefaultAPI) stopMiningHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := a.StopMining(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// getMetricsHandler 获取节点指标
func (a *DefaultAPI) getMetricsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	metrics, err := a.GetMetrics()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metrics)
}
