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
	"golang.org/x/crypto/sha3"
)

// Crypto defines the interface for cryptographic operations.
type Crypto interface {
	// GenerateKeyPair generates a new key pair.
	GenerateKeyPair() (PrivateKey, PublicKey, error)

	// GenerateEd25519KeyPair generates a new Ed25519 key pair.
	GenerateEd25519KeyPair() (PrivateKey, PublicKey, error)

	// GenerateECDSAKeyPair generates a new ECDSA key pair.
	GenerateECDSAKeyPair() (PrivateKey, PublicKey, error)

	// Sign signs data with the given private key.
	Sign(data []byte, privateKey PrivateKey) ([]byte, error)

	// Verify verifies a signature with the given public key.
	Verify(data []byte, signature []byte, publicKey PublicKey) bool

	// Hash computes the hash of the given data.
	Hash(data []byte) types.Hash

	// Keccak256 computes the Keccak256 hash of the given data.
	Keccak256(data []byte) types.Hash

	// GenerateAddress generates an address from the given public key.
	GenerateAddress(publicKey PublicKey) types.Address
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
	Type() KeyType
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
	Type() KeyType
}

// KeyType represents the type of cryptographic key.
type KeyType string

const (
	// Ed25519 key type
	Ed25519 KeyType = "ed25519"
	// ECDSA key type
	ECDSA KeyType = "ecdsa"
)

// KeyFile represents the structure of a saved key file.
type KeyFile struct {
	Type       KeyType `json:"type"`
	PublicKey  []byte  `json:"public_key"`
	PrivateKey []byte  `json:"private_key"`
}

// EncryptedKeyFile represents the structure of an encrypted saved key file.
type EncryptedKeyFile struct {
	Type      KeyType `json:"type"`
	PublicKey []byte  `json:"public_key"`
	// EncryptedPrivateKey is the AES-GCM encrypted private key
	EncryptedPrivateKey []byte `json:"encrypted_private_key"`
	// Salt used for key derivation
	Salt []byte `json:"salt"`
	// Nonce used for AES-GCM
	Nonce []byte `json:"nonce"`
}

// DefaultCrypto implements the Crypto interface.
type DefaultCrypto struct{}

// GenerateKeyPair generates a new key pair of the specified type.
func (c *DefaultCrypto) GenerateKeyPair() (PrivateKey, PublicKey, error) {
	// For now, we'll generate an Ed25519 key pair by default
	return c.GenerateEd25519KeyPair()
}

// GenerateEd25519KeyPair generates a new Ed25519 key pair.
func (c *DefaultCrypto) GenerateEd25519KeyPair() (PrivateKey, PublicKey, error) {
	return GenerateEd25519KeyPair()
}

// GenerateECDSAKeyPair generates a new ECDSA key pair using the secp256k1 curve.
func (c *DefaultCrypto) GenerateECDSAKeyPair() (PrivateKey, PublicKey, error) {
	return GenerateECDSAKeyPair()
}

// Sign signs data with the given private key.
func (c *DefaultCrypto) Sign(data []byte, privateKey PrivateKey) ([]byte, error) {
	return privateKey.Sign(data)
}

// Verify verifies a signature with the given public key.
func (c *DefaultCrypto) Verify(data []byte, signature []byte, publicKey PublicKey) bool {
	return publicKey.Verify(data, signature)
}

// Hash computes the SHA256 hash of the given data.
func (c *DefaultCrypto) Hash(data []byte) types.Hash {
	hash := sha256.Sum256(data)
	return types.Hash(hash)
}

// Keccak256 computes the Keccak256 hash of the given data.
func (c *DefaultCrypto) Keccak256(data []byte) types.Hash {
	hasher := sha3.NewLegacyKeccak256()
	hasher.Write(data)
	var result types.Hash
	copy(result[:], hasher.Sum(nil))
	return result
}

// GenerateAddress generates an address from the given public key.
func (c *DefaultCrypto) GenerateAddress(publicKey PublicKey) types.Address {
	return publicKey.Address()
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
	switch encryptedKeyFile.Type {
	case Ed25519:
		return Ed25519PrivateKeyFromBytes(privateKeyBytes)
	case ECDSA:
		return ECDSAPrivateKeyFromBytes(privateKeyBytes)
	default:
		return nil, fmt.Errorf("unsupported key type: %s", encryptedKeyFile.Type)
	}
}

// NewCrypto creates a new instance of the default crypto implementation.
func NewCrypto() Crypto {
	return &DefaultCrypto{}
}
