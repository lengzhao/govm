package core

import (
	"fmt"

	lzbinary "github.com/lengzhao/binary"
	"github.com/lengzhao/govm/crypto"
	"github.com/lengzhao/govm/storage"
	"github.com/lengzhao/govm/types"
)

// TxProcessor 交易处理器接口
type TxProcessor interface {
	// ValidateTransaction 验证交易的有效性
	ValidateTransaction(tx *types.TransactionWithSign) error

	// ApplyTransaction 应用交易到状态
	ApplyTransaction(tx *types.TransactionWithSign) error

	// GetTransactionByHash 根据哈希获取交易
	GetTransactionByHash(hash types.Hash) (*types.TransactionWithSign, error)

	// GetBalance 获取账户余额
	GetBalance(addr types.Address) (uint64, error)

	// GetNonce 获取账户nonce值
	GetNonce(addr types.Address) (uint64, error)

	// GetAccount 获取账户信息
	GetAccount(addr types.Address) (*types.Account, error)

	// CreateAccount 创建新账户
	CreateAccount(addr types.Address, pubKey []byte) error

	// UpdateAccount 更新账户信息
	UpdateAccount(account *types.Account) error
}

// DefaultTxProcessor 默认交易处理器实现
type DefaultTxProcessor struct {
	crypto  crypto.Crypto
	storage storage.Storage
}

// NewTxProcessor 创建新的交易处理器
func NewTxProcessor(store storage.Storage) TxProcessor {
	return &DefaultTxProcessor{
		crypto:  crypto.NewCrypto(),
		storage: store,
	}
}

// ValidateTransaction 验证交易的有效性
func (tp *DefaultTxProcessor) ValidateTransaction(tx *types.TransactionWithSign) error {
	// 检查交易基本结构
	if tx == nil {
		return fmt.Errorf("transaction is nil")
	}

	// 检查签名是否存在
	if len(tx.Signature) == 0 {
		return fmt.Errorf("transaction signature is missing")
	}

	// 验证交易签名
	if err := tp.verifyTransactionSignature(&tx.Transaction, tx.Signature); err != nil {
		return fmt.Errorf("transaction signature verification failed: %w", err)
	}

	// 验证交易字段
	if err := tp.validateTransactionFields(&tx.Transaction); err != nil {
		return fmt.Errorf("transaction field validation failed: %w", err)
	}

	return nil
}

// verifyTransactionSignature 验证交易签名
func (tp *DefaultTxProcessor) verifyTransactionSignature(tx *types.Transaction, signature []byte) error {
	// 序列化交易用于签名验证
	data, err := lzbinary.Marshal(tx)
	if err != nil {
		return fmt.Errorf("failed to marshal transaction: %w", err)
	}

	// 检查公钥是否存在
	if len(tx.PublicKey) == 0 {
		return fmt.Errorf("public key is missing in transaction")
	}

	// 验证签名 (修改为使用新的接口)
	if !tp.crypto.Verify(data, signature, tx.PublicKey, crypto.Ed25519) {
		return fmt.Errorf("invalid transaction signature")
	}

	return nil
}

// validateTransactionFields 验证交易字段
func (tp *DefaultTxProcessor) validateTransactionFields(tx *types.Transaction) error {
	// 检查金额是否有效
	if tx.Amount <= 0 {
		return fmt.Errorf("transaction amount must be positive")
	}

	// 检查发送方和接收方地址是否相同
	if tx.From == tx.To {
		return fmt.Errorf("sender and receiver addresses cannot be the same")
	}

	// 检查分片ID是否有效
	if tx.ShardID == 0 {
		return fmt.Errorf("shard ID must be positive")
	}

	// 检查Gas相关字段
	if tx.GasPrice <= 0 {
		return fmt.Errorf("gas price must be positive")
	}

	if tx.GasLimit <= 0 {
		return fmt.Errorf("gas limit must be positive")
	}

	if tx.GasFeeCap <= 0 {
		return fmt.Errorf("gas fee cap must be positive")
	}

	if tx.GasPrice > tx.GasFeeCap {
		return fmt.Errorf("gas price cannot exceed gas fee cap")
	}

	// 检查公钥字段
	if len(tx.PublicKey) == 0 {
		return fmt.Errorf("public key is required")
	}

	// 可以添加更多验证规则

	return nil
}

// ApplyTransaction 应用交易到状态
func (tp *DefaultTxProcessor) ApplyTransaction(tx *types.TransactionWithSign) error {
	// 首先验证交易
	if err := tp.ValidateTransaction(tx); err != nil {
		return fmt.Errorf("transaction validation failed: %w", err)
	}

	// 计算Gas费用
	gasFee := tp.calculateGasFee(&tx.Transaction)

	// 应用交易到状态数据库
	// 1. 从发送方账户扣除金额和Gas费用
	fromBalance, err := tp.GetBalance(tx.From)
	if err != nil {
		fromBalance = 0 // 如果账户不存在，余额为0
	}

	totalDeduction := tx.Amount + gasFee
	if fromBalance < totalDeduction {
		return fmt.Errorf("insufficient balance: %d < %d (amount: %d + gas fee: %d)", fromBalance, totalDeduction, tx.Amount, gasFee)
	}

	if err := tp.setBalance(tx.From, fromBalance-totalDeduction); err != nil {
		return fmt.Errorf("failed to update sender balance: %w", err)
	}

	// 2. 向接收方账户增加金额
	toBalance, err := tp.GetBalance(tx.To)
	if err != nil {
		toBalance = 0 // 如果账户不存在，余额为0
	}

	if err := tp.setBalance(tx.To, toBalance+tx.Amount); err != nil {
		return fmt.Errorf("failed to update receiver balance: %w", err)
	}

	// 3. 更新账户nonce值
	if err := tp.updateNonce(tx.From, tx.Nonce); err != nil {
		return fmt.Errorf("failed to update nonce: %w", err)
	}

	// 4. 记录交易日志
	if err := tp.storeTransaction(tx); err != nil {
		return fmt.Errorf("failed to store transaction: %w", err)
	}

	return nil
}

// GetTransactionByHash 根据哈希获取交易
func (tp *DefaultTxProcessor) GetTransactionByHash(hash types.Hash) (*types.TransactionWithSign, error) {
	// 从存储中获取交易
	txStorage, err := tp.storage.NewStorage(types.SNTx)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction storage: %w", err)
	}

	data, err := txStorage.Get(hash[:])
	if err != nil {
		return nil, fmt.Errorf("transaction not found: %w", err)
	}

	var tx types.TransactionWithSign
	if err := lzbinary.Unmarshal(data, &tx); err != nil {
		return nil, fmt.Errorf("failed to unmarshal transaction: %w", err)
	}

	return &tx, nil
}

// GetBalance 获取账户余额
func (tp *DefaultTxProcessor) GetBalance(addr types.Address) (uint64, error) {
	// 从存储中获取账户余额
	statusStorage, err := tp.storage.NewStorage(types.SNStatus)
	if err != nil {
		return 0, fmt.Errorf("failed to create status storage: %w", err)
	}

	key := append([]byte("balance_"), addr[:]...)
	data, err := statusStorage.Get(key)
	if err != nil {
		return 0, fmt.Errorf("balance not found: %w", err)
	}

	if len(data) != 8 {
		return 0, fmt.Errorf("invalid balance data")
	}

	balance := uint64(0)
	for i, b := range data {
		balance |= uint64(b) << (8 * i)
	}

	return balance, nil
}

// GetNonce 获取账户nonce值
func (tp *DefaultTxProcessor) GetNonce(addr types.Address) (uint64, error) {
	// 从存储中获取账户nonce值
	statusStorage, err := tp.storage.NewStorage(types.SNStatus)
	if err != nil {
		return 0, fmt.Errorf("failed to create status storage: %w", err)
	}

	key := append([]byte("nonce_"), addr[:]...)
	data, err := statusStorage.Get(key)
	if err != nil {
		return 0, nil // 如果nonce不存在，返回0
	}

	if len(data) != 8 {
		return 0, fmt.Errorf("invalid nonce data")
	}

	nonce := uint64(0)
	for i, b := range data {
		nonce |= uint64(b) << (8 * i)
	}

	return nonce, nil
}

// setBalance 设置账户余额
func (tp *DefaultTxProcessor) setBalance(addr types.Address, balance uint64) error {
	// 存储账户余额
	statusStorage, err := tp.storage.NewStorage(types.SNStatus)
	if err != nil {
		return fmt.Errorf("failed to create status storage: %w", err)
	}

	key := append([]byte("balance_"), addr[:]...)
	data := make([]byte, 8)
	for i := 0; i < 8; i++ {
		data[i] = byte(balance >> (8 * i))
	}

	return statusStorage.Put(key, data)
}

// setNonce 设置账户nonce值
func (tp *DefaultTxProcessor) setNonce(addr types.Address, nonce uint64) error {
	// 存储账户nonce值
	statusStorage, err := tp.storage.NewStorage(types.SNStatus)
	if err != nil {
		return fmt.Errorf("failed to create status storage: %w", err)
	}

	key := append([]byte("nonce_"), addr[:]...)
	data := make([]byte, 8)
	for i := 0; i < 8; i++ {
		data[i] = byte(nonce >> (8 * i))
	}

	return statusStorage.Put(key, data)
}

// storeTransaction 存储交易
func (tp *DefaultTxProcessor) storeTransaction(tx *types.TransactionWithSign) error {
	// 存储交易
	txStorage, err := tp.storage.NewStorage(types.SNTx)
	if err != nil {
		return fmt.Errorf("failed to create transaction storage: %w", err)
	}

	data, err := lzbinary.Marshal(tx)
	if err != nil {
		return fmt.Errorf("failed to marshal transaction: %w", err)
	}

	txHash := tp.calculateTransactionHash(tx)
	return txStorage.Put(txHash[:], data)
}

// updateNonce 更新账户nonce值
func (tp *DefaultTxProcessor) updateNonce(addr types.Address, nonce uint64) error {
	// 获取当前nonce值
	currentNonce, err := tp.GetNonce(addr)
	if err != nil {
		currentNonce = 0 // 如果账户不存在，nonce为0
	}

	// 检查nonce是否有效
	if nonce <= currentNonce {
		return fmt.Errorf("invalid nonce: expected %d, got %d", currentNonce+1, nonce)
	}

	// 更新nonce值
	return tp.setNonce(addr, nonce)
}

// calculateTransactionHash 计算交易哈希
func (tp *DefaultTxProcessor) calculateTransactionHash(tx *types.TransactionWithSign) types.Hash {
	data, err := lzbinary.Marshal(tx.Transaction)
	if err != nil {
		return types.Hash{} // 返回空哈希
	}

	return tp.crypto.Hash(data)
}

// calculateGasFee 计算交易的Gas费用
func (tp *DefaultTxProcessor) calculateGasFee(tx *types.Transaction) uint64 {
	// 基础Gas费用
	baseGas := uint64(21000)

	// 数据Gas费用（每字节1单位Gas）
	dataGas := uint64(len(tx.Data))

	// 计算总Gas使用量
	gasUsed := baseGas + dataGas

	// 确保不超过Gas限制
	if gasUsed > tx.GasLimit {
		gasUsed = tx.GasLimit
	}

	// 计算Gas费用
	gasFee := gasUsed * tx.GasPrice

	// 确保不超过Gas费用上限
	if gasFee > tx.GasFeeCap {
		gasFee = tx.GasFeeCap
	}

	return gasFee
}

// GetAccount 获取账户信息
func (tp *DefaultTxProcessor) GetAccount(addr types.Address) (*types.Account, error) {
	// 获取账户余额
	balance, err := tp.GetBalance(addr)
	if err != nil {
		balance = 0
	}

	// 获取账户nonce值
	nonce, err := tp.GetNonce(addr)
	if err != nil {
		nonce = 0
	}

	// 创建账户对象
	account := &types.Account{
		Address:    addr,
		Balance:    balance,
		Nonce:      nonce,
		PublicKey:  nil, // 公钥需要从其他地方获取
		CodeHash:   types.Hash{},
		IsContract: false,
	}

	return account, nil
}

// CreateAccount 创建新账户
func (tp *DefaultTxProcessor) CreateAccount(addr types.Address, pubKey []byte) error {
	// 初始化账户余额为0
	if err := tp.setBalance(addr, 0); err != nil {
		return fmt.Errorf("failed to initialize account balance: %w", err)
	}

	// 初始化账户nonce为0
	if err := tp.setNonce(addr, 0); err != nil {
		return fmt.Errorf("failed to initialize account nonce: %w", err)
	}

	return nil
}

// UpdateAccount 更新账户信息
func (tp *DefaultTxProcessor) UpdateAccount(account *types.Account) error {
	// 更新账户余额
	if err := tp.setBalance(account.Address, account.Balance); err != nil {
		return fmt.Errorf("failed to update account balance: %w", err)
	}

	// 更新账户nonce值
	if err := tp.setNonce(account.Address, account.Nonce); err != nil {
		return fmt.Errorf("failed to update account nonce: %w", err)
	}

	return nil
}
