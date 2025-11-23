package crypto

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"

	"github.com/lengzhao/govm/types"
)

func init() {
	RegisterAlgorithm(&ed25519Algorithm{})
}

// ed25519PrivateKey implements the PrivateKey interface for Ed25519 keys.
type ed25519PrivateKey struct {
	key ed25519.PrivateKey
}

// ed25519Algorithm implements the Algorithm interface for Ed25519.
type ed25519Algorithm struct{}

// Bytes returns the Ed25519 private key as bytes.
func (k *ed25519PrivateKey) Bytes() []byte {
	return k.key
}

// FromBytes creates an Ed25519 private key from bytes.
func (k *ed25519PrivateKey) FromBytes(data []byte) (PrivateKey, error) {
	// For Ed25519, the private key bytes can be used directly
	privKey := ed25519.PrivateKey(data)
	return &ed25519PrivateKey{key: privKey}, nil
}

// PublicKey returns the corresponding Ed25519 public key.
func (k *ed25519PrivateKey) PublicKey() PublicKey {
	pubKey := k.key.Public().(ed25519.PublicKey)
	return &ed25519PublicKey{key: pubKey}
}

// Sign signs data with the Ed25519 private key.
func (k *ed25519PrivateKey) Sign(data []byte) ([]byte, error) {
	return ed25519.Sign(k.key, data), nil
}

// Type returns the key type.
func (k *ed25519PrivateKey) Type() string {
	return Ed25519
}

// ed25519PublicKey implements the PublicKey interface for Ed25519 keys.
type ed25519PublicKey struct {
	key ed25519.PublicKey
}

// Bytes returns the Ed25519 public key as bytes.
func (k *ed25519PublicKey) Bytes() []byte {
	return k.key
}

// FromBytes creates an Ed25519 public key from bytes.
func (k *ed25519PublicKey) FromBytes(data []byte) (PublicKey, error) {
	// For Ed25519, the public key bytes can be used directly
	pubKey := ed25519.PublicKey(data)
	return &ed25519PublicKey{key: pubKey}, nil
}

// Address generates an address from the Ed25519 public key.
func (k *ed25519PublicKey) Address() types.Address {
	// Hash the public key and take the last 20 bytes as the address
	hash := sha256.Sum256(k.key)
	var addr types.Address
	copy(addr[:], hash[len(hash)-20:])
	return addr
}

// Verify verifies a signature with the Ed25519 public key.
func (k *ed25519PublicKey) Verify(data []byte, signature []byte) bool {
	return ed25519.Verify(k.key, data, signature)
}

// Type returns the key type.
func (k *ed25519PublicKey) Type() string {
	return Ed25519
}

// GenerateKeyPair 生成新的密钥对
func (a *ed25519Algorithm) GenerateKeyPair() ([]byte, []byte, error) {
	_, privKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, nil, err
	}

	pubKey := privKey.Public().(ed25519.PublicKey)
	return privKey, pubKey, nil
}

// Sign 使用私钥对数据进行签名
func (a *ed25519Algorithm) Sign(data []byte, privateKey []byte) ([]byte, error) {
	privKey := ed25519.PrivateKey(privateKey)
	return ed25519.Sign(privKey, data), nil
}

// Verify 使用公钥验证签名
func (a *ed25519Algorithm) Verify(data []byte, signature []byte, publicKey []byte) bool {
	pubKey := ed25519.PublicKey(publicKey)
	return ed25519.Verify(pubKey, data, signature)
}

// GenerateAddress 从公钥生成地址
func (a *ed25519Algorithm) GenerateAddress(publicKey []byte) types.Address {
	pubKey := ed25519.PublicKey(publicKey)
	// Hash the public key and take the last 20 bytes as the address
	hash := sha256.Sum256(pubKey)
	var addr types.Address
	copy(addr[:], hash[len(hash)-20:])
	return addr
}

// Type 返回算法类型
func (a *ed25519Algorithm) Type() string {
	return Ed25519
}

// PrivateKeyFromBytes 从字节数据创建私钥
func (a *ed25519Algorithm) PrivateKeyFromBytes(data []byte) (PrivateKey, error) {
	privKey := ed25519.PrivateKey(data)
	return &ed25519PrivateKey{key: privKey}, nil
}

// PublicKeyFromBytes 从字节数据创建公钥
func (a *ed25519Algorithm) PublicKeyFromBytes(data []byte) (PublicKey, error) {
	pubKey := ed25519.PublicKey(data)
	return &ed25519PublicKey{key: pubKey}, nil
}
