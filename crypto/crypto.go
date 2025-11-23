// Package crypto provides cryptographic functions for the govm blockchain.
package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/lengzhao/govm/types"
	"golang.org/x/crypto/scrypt"
)

// GenerateKeyPair generates a new key pair of the specified type.
// If no key type is specified, it generates an Ed25519 key pair by default.
func GenerateKeyPair(keyType string) ([]byte, []byte, error) {
	algorithm, err := algorithmFactory(keyType)
	if err != nil {
		return nil, nil, err
	}
	return algorithm.GenerateKeyPair()
}

// Sign signs data with the given private key bytes.
func Sign(data []byte, privateKey []byte, keyType string) ([]byte, error) {
	algorithm, err := algorithmFactory(keyType)
	if err != nil {
		return nil, err
	}

	return algorithm.Sign(data, privateKey)
}

// Verify verifies a signature with the given public key bytes.
func Verify(data []byte, signature []byte, publicKey []byte, keyType string) bool {
	algorithm, err := algorithmFactory(keyType)
	if err != nil {
		return false
	}

	return algorithm.Verify(data, signature, publicKey)
}

// Hash computes the SHA256 hash of the given data.
func Hash(data []byte) types.Hash {
	hash := sha256.Sum256(data)
	return types.Hash(hash)
}

// GenerateAddress generates an address from the given public key bytes.
func GenerateAddress(publicKey []byte, keyType string) types.Address {
	algorithm, err := algorithmFactory(keyType)
	if err != nil {
		return types.Address{}
	}

	return algorithm.GenerateAddress(publicKey)
}

// GenerateRandomBytes generates cryptographically secure random bytes.
func GenerateRandomBytes(length int) ([]byte, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to generate random bytes: %w", err)
	}
	return bytes, nil
}

// EncryptConfig encrypts configuration data.
func EncryptConfig(data []byte, password string) ([]byte, error) {
	// Generate a random salt
	salt := make([]byte, 32)
	if _, err := rand.Read(salt); err != nil {
		return nil, fmt.Errorf("failed to generate salt: %w", err)
	}

	// Derive key from password
	key, err := deriveKey(password, salt)
	if err != nil {
		return nil, fmt.Errorf("failed to derive key: %w", err)
	}

	// Encrypt data
	ciphertext, nonce, err := encryptData(data, key)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt data: %w", err)
	}

	// Combine salt, nonce, and ciphertext
	result := make([]byte, 0, len(salt)+len(nonce)+len(ciphertext))
	result = append(result, salt...)
	result = append(result, nonce...)
	result = append(result, ciphertext...)

	return result, nil
}

// DecryptConfig decrypts configuration data.
func DecryptConfig(encryptedData []byte, password string) ([]byte, error) {
	// Extract salt, nonce, and ciphertext
	if len(encryptedData) < 32+12 {
		return nil, fmt.Errorf("encrypted data is too short")
	}

	salt := encryptedData[:32]
	nonce := encryptedData[32:44]
	ciphertext := encryptedData[44:]

	// Derive key from password and salt
	key, err := deriveKey(password, salt)
	if err != nil {
		return nil, fmt.Errorf("failed to derive key: %w", err)
	}

	// Decrypt data
	plaintext, err := decryptData(ciphertext, nonce, key)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt data: %w", err)
	}

	return plaintext, nil
}

// Algorithm 定义了统一的密码学算法接口
type Algorithm interface {
	// GenerateKeyPair 生成新的密钥对，直接返回私钥和公钥的字节
	GenerateKeyPair() (privateKey []byte, publicKey []byte, error error)

	// Sign 使用私钥对数据进行签名
	Sign(data []byte, privateKey []byte) ([]byte, error)

	// Verify 使用公钥验证签名
	Verify(data []byte, signature []byte, publicKey []byte) bool

	// GenerateAddress 从公钥生成地址
	GenerateAddress(publicKey []byte) types.Address

	// Type 返回算法类型
	Type() string

	// PrivateKeyFromBytes 从字节数据创建私钥
	PrivateKeyFromBytes(data []byte) (PrivateKey, error)

	// PublicKeyFromBytes 从字节数据创建公钥
	PublicKeyFromBytes(data []byte) (PublicKey, error)
}

// PrivateKey represents a private key interface.
type PrivateKey interface {
	// Bytes returns the private key as bytes.
	Bytes() []byte

	// PublicKey returns the corresponding public key.
	PublicKey() PublicKey

	// Sign signs data with this private key.
	Sign(data []byte) ([]byte, error)

	// Type returns the type of the private key.
	Type() string
}

// PublicKey represents a public key interface.
type PublicKey interface {
	// Bytes returns the public key as bytes.
	Bytes() []byte

	// Address returns the address derived from this public key.
	Address() types.Address

	// Verify verifies a signature with this public key.
	Verify(data []byte, signature []byte) bool

	// Type returns the type of the public key.
	Type() string
}

const (
	// Ed25519 key type
	Ed25519 = "ed25519"
	// ECDSA key type
	ECDSA = "ecdsa"
	// Secp256k1 key type
	Secp256k1 = "secp256k1"
	// Schnorr key type
	Schnorr = "schnorr"
)

// KeyFile represents the structure of a saved key file.
type KeyFile struct {
	Type       string `json:"type"`
	PublicKey  []byte `json:"public_key"`
	PrivateKey []byte `json:"private_key"`
}

// EncryptedKeyFile represents the structure of an encrypted saved key file.
type EncryptedKeyFile struct {
	Type      string `json:"type"`
	PublicKey []byte `json:"public_key"`
	// EncryptedPrivateKey is the AES-GCM encrypted private key
	EncryptedPrivateKey []byte `json:"encrypted_private_key"`
	// Salt used for key derivation
	Salt []byte `json:"salt"`
	// Nonce used for AES-GCM
	Nonce []byte `json:"nonce"`
}

// algorithmRegistry 算法注册表
var algorithmRegistry = make(map[string]Algorithm)

// RegisterAlgorithm 注册算法
func RegisterAlgorithm(algorithm Algorithm) error {
	if algorithm == nil {
		return fmt.Errorf("algorithm cannot be nil")
	}
	keyType := algorithm.Type()
	if _, exists := algorithmRegistry[keyType]; exists {
		return fmt.Errorf("algorithm for key type %s already registered", keyType)
	}
	algorithmRegistry[keyType] = algorithm
	return nil
}

// algorithmFactory 创建算法实例的工厂
func algorithmFactory(keyType string) (Algorithm, error) {
	if algorithm, exists := algorithmRegistry[keyType]; exists {
		return algorithm, nil
	}
	return nil, fmt.Errorf("unsupported key type: %s", keyType)
}

// AlgorithmFactory creates an algorithm instance by key type
func AlgorithmFactory(keyType string) (Algorithm, error) {
	return algorithmFactory(keyType)
}

// deriveKey derives a 32-byte key from a password using scrypt
func deriveKey(password string, salt []byte) ([]byte, error) {
	// Use scrypt with N=32768, r=8, p=1 for strong key derivation
	key, err := scrypt.Key([]byte(password), salt, 32768, 8, 1, 32)
	if err != nil {
		return nil, fmt.Errorf("failed to derive key: %w", err)
	}
	return key, nil
}

// encryptData encrypts data using AES-GCM with the provided key
func encryptData(data []byte, key []byte) ([]byte, []byte, error) {
	// Create cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	// Create GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	// Generate nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Encrypt data
	ciphertext := gcm.Seal(nil, nonce, data, nil)
	return ciphertext, nonce, nil
}

// decryptData decrypts data using AES-GCM with the provided key
func decryptData(ciphertext []byte, nonce []byte, key []byte) ([]byte, error) {
	// Create cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	// Create GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	// Decrypt data
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt data: %w", err)
	}

	return plaintext, nil
}

// SaveToFile saves a private key to a file.
// If password is empty, it will save as unencrypted key file.
// If password is provided, it will save as encrypted key file using that password.
// If password is empty, it will be set to default password "govm-password".
func SaveToFile(privateKey PrivateKey, filename string, password string) error {
	// If password is empty, set it to default password
	if password == "" {
		password = "govm-password"
	}

	// Generate salt for key derivation
	salt := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return fmt.Errorf("failed to generate salt: %w", err)
	}

	// Derive key from password
	key, err := deriveKey(password, salt)
	if err != nil {
		return fmt.Errorf("failed to derive key: %w", err)
	}

	// Encrypt private key
	ciphertext, nonce, err := encryptData(privateKey.Bytes(), key)
	if err != nil {
		return fmt.Errorf("failed to encrypt private key: %w", err)
	}

	// Create the encrypted key file structure
	encryptedKeyFile := EncryptedKeyFile{
		Type:                privateKey.Type(),
		PublicKey:           privateKey.PublicKey().Bytes(),
		EncryptedPrivateKey: ciphertext,
		Salt:                salt,
		Nonce:               nonce,
	}

	// Convert to JSON
	data, err := json.Marshal(encryptedKeyFile)
	if err != nil {
		return fmt.Errorf("failed to marshal encrypted key file: %w", err)
	}

	// Write to file
	return os.WriteFile(filename, data, 0600)
}

// LoadFromFile loads a private key from a file.
// If password is empty, it will be set to default password "govm-password".
func LoadFromFile(filename string, password string) (PrivateKey, error) {
	// If password is empty, set it to default password
	if password == "" {
		password = "govm-password"
	}

	// Read the file
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read key file: %w", err)
	}

	// Parse the JSON
	var encryptedKeyFile EncryptedKeyFile
	if err := json.Unmarshal(data, &encryptedKeyFile); err != nil {
		return nil, fmt.Errorf("failed to unmarshal encrypted key file: %w", err)
	}

	// Derive key from password and salt
	key, err := deriveKey(password, encryptedKeyFile.Salt)
	if err != nil {
		return nil, fmt.Errorf("failed to derive key: %w", err)
	}

	// Decrypt private key
	privateKeyBytes, err := decryptData(encryptedKeyFile.EncryptedPrivateKey, encryptedKeyFile.Nonce, key)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt private key: %w", err)
	}

	// Recreate the private key based on its type
	algorithm, err := algorithmFactory(encryptedKeyFile.Type)
	if err != nil {
		return nil, fmt.Errorf("failed to get algorithm for %s: %w", encryptedKeyFile.Type, err)
	}

	return algorithm.PrivateKeyFromBytes(privateKeyBytes)
}
