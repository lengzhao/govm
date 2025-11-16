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
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/scrypt"
)

// Crypto defines the interface for cryptographic operations.
type Crypto interface {
	// GenerateKeyPair generates a new key pair of the specified type.
	// If no key type is specified, it generates an Ed25519 key pair by default.
	GenerateKeyPair(keyType KeyType) (PrivateKey, PublicKey, error)

	// Sign signs data with the given private key.
	Sign(data []byte, privateKey PrivateKey) ([]byte, error)

	// Verify verifies a signature with the given public key.
	Verify(data []byte, signature []byte, publicKey PublicKey) bool

	// Hash computes the hash of the given data.
	Hash(data []byte) types.Hash

	// GenerateAddress generates an address from the given public key.
	GenerateAddress(publicKey PublicKey) types.Address

	// GenerateRandomBytes generates cryptographically secure random bytes.
	GenerateRandomBytes(length int) ([]byte, error)

	// HashPassword securely hashes a password using bcrypt.
	HashPassword(password string) (string, error)

	// VerifyPassword verifies a password against a hash.
	VerifyPassword(password, hash string) bool

	// EncryptConfig encrypts configuration data.
	EncryptConfig(data []byte, password string) ([]byte, error)

	// DecryptConfig decrypts configuration data.
	DecryptConfig(encryptedData []byte, password string) ([]byte, error)
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

	// FromBytes creates a private key from bytes.
	FromBytes(data []byte) (PrivateKey, error)
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

	// FromBytes creates a public key from bytes.
	FromBytes(data []byte) (PublicKey, error)
}

// KeyType represents the type of cryptographic key.
type KeyType string

const (
	// Ed25519 key type
	Ed25519 KeyType = "ed25519"
	// ECDSA key type
	ECDSA KeyType = "ecdsa"
	// Secp256k1 key type
	Secp256k1 KeyType = "secp256k1"
	// Schnorr key type
	Schnorr KeyType = "schnorr"
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
func (c *DefaultCrypto) GenerateKeyPair(keyType KeyType) (PrivateKey, PublicKey, error) {
	switch keyType {
	case Ed25519:
		return GenerateEd25519KeyPair()
	case ECDSA:
		return GenerateECDSAKeyPair()
	case Secp256k1:
		return GenerateSecp256k1KeyPair()
	case Schnorr:
		return GenerateSchnorrKeyPair()
	default:
		return nil, nil, fmt.Errorf("unsupported key type: %s", keyType)
	}
}

// Sign signs data with the given private key.
func (c *DefaultCrypto) Sign(data []byte, privateKey PrivateKey) ([]byte, error) {
	return privateKey.Sign(data)
}

// Verify verifies a signature with the given public key.
func (c *DefaultCrypto) Verify(data []byte, signature []byte, publicKey PublicKey) bool {
	if publicKey != nil {
		return publicKey.Verify(data, signature)
	}
	return false
}

// Hash computes the SHA256 hash of the given data.
func (c *DefaultCrypto) Hash(data []byte) types.Hash {
	hash := sha256.Sum256(data)
	return types.Hash(hash)
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
	priv, _, err := NewCrypto().GenerateKeyPair(encryptedKeyFile.Type)
	if err != nil {
		return nil, fmt.Errorf("failed to generate %s private key: %w", encryptedKeyFile.Type, err)
	}
	return priv.FromBytes(privateKeyBytes)
}

// GenerateRandomBytes generates cryptographically secure random bytes.
func (c *DefaultCrypto) GenerateRandomBytes(length int) ([]byte, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to generate random bytes: %w", err)
	}
	return bytes, nil
}

// HashPassword securely hashes a password using bcrypt.
func (c *DefaultCrypto) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hash), nil
}

// VerifyPassword verifies a password against a hash.
func (c *DefaultCrypto) VerifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// EncryptConfig encrypts configuration data.
func (c *DefaultCrypto) EncryptConfig(data []byte, password string) ([]byte, error) {
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
func (c *DefaultCrypto) DecryptConfig(encryptedData []byte, password string) ([]byte, error) {
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

// NewCrypto creates a new instance of the default crypto implementation.
func NewCrypto() Crypto {
	return &DefaultCrypto{}
}
