package crypto

import (
	"crypto/sha256"
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/lengzhao/govm/types"
)

func init() {
	RegisterAlgorithm(&schnorrAlgorithm{})
}

// schnorrPrivateKey implements the PrivateKey interface for Schnorr keys.
type schnorrPrivateKey struct {
	key *btcec.PrivateKey
}

// schnorrAlgorithm implements the Algorithm interface for Schnorr.
type schnorrAlgorithm struct{}

// Bytes returns the Schnorr private key as bytes.
func (k *schnorrPrivateKey) Bytes() []byte {
	return k.key.Serialize()
}

// FromBytes creates a Schnorr private key from bytes.
func (k *schnorrPrivateKey) FromBytes(data []byte) (PrivateKey, error) {
	privKey, _ := btcec.PrivKeyFromBytes(data)
	return &schnorrPrivateKey{key: privKey}, nil
}

// PublicKey returns the corresponding Schnorr public key.
func (k *schnorrPrivateKey) PublicKey() PublicKey {
	return &schnorrPublicKey{key: k.key.PubKey()}
}

// Sign signs data with the Schnorr private key.
func (k *schnorrPrivateKey) Sign(data []byte) ([]byte, error) {
	// Hash the data
	hash := chainhash.HashB(data)

	// Create the Schnorr signature
	signature, err := schnorr.Sign(k.key, hash)
	if err != nil {
		return nil, err
	}

	return signature.Serialize(), nil
}

// Type returns the key type.
func (k *schnorrPrivateKey) Type() string {
	return Schnorr
}

// schnorrPublicKey implements the PublicKey interface for Schnorr keys.
type schnorrPublicKey struct {
	key *btcec.PublicKey
}

// Bytes returns the Schnorr public key as bytes (compressed format).
func (k *schnorrPublicKey) Bytes() []byte {
	// Compress the public key (32 bytes for Schnorr)
	return schnorr.SerializePubKey(k.key)
}

// FromBytes creates a Schnorr public key from bytes.
func (k *schnorrPublicKey) FromBytes(data []byte) (PublicKey, error) {
	// Parse the Schnorr public key
	pubKey, err := schnorr.ParsePubKey(data)
	if err != nil {
		return nil, fmt.Errorf("invalid Schnorr public key data: %w", err)
	}

	return &schnorrPublicKey{key: pubKey}, nil
}

// Address generates an address from the Schnorr public key.
func (k *schnorrPublicKey) Address() types.Address {
	// Hash the public key and take the last 20 bytes as the address
	pubKeyBytes := schnorr.SerializePubKey(k.key)
	hash := sha256.Sum256(pubKeyBytes)
	var addr types.Address
	copy(addr[:], hash[len(hash)-20:])
	return addr
}

// Verify verifies a signature with the Schnorr public key.
func (k *schnorrPublicKey) Verify(data []byte, signature []byte) bool {
	// Hash the data
	hash := chainhash.HashB(data)

	// Parse the signature
	sig, err := schnorr.ParseSignature(signature)
	if err != nil {
		return false
	}

	return sig.Verify(hash, k.key)
}

// Type returns the key type.
func (k *schnorrPublicKey) Type() string {
	return Schnorr
}

// GenerateKeyPair 生成新的密钥对
func (a *schnorrAlgorithm) GenerateKeyPair() ([]byte, []byte, error) {
	privKey, err := btcec.NewPrivateKey()
	if err != nil {
		return nil, nil, err
	}

	return privKey.Serialize(), schnorr.SerializePubKey(privKey.PubKey()), nil
}

// Sign 使用私钥对数据进行签名
func (a *schnorrAlgorithm) Sign(data []byte, privateKey []byte) ([]byte, error) {
	privKey, _ := btcec.PrivKeyFromBytes(privateKey)
	// Hash the data
	hash := chainhash.HashB(data)

	// Create the Schnorr signature
	signature, err := schnorr.Sign(privKey, hash)
	if err != nil {
		return nil, err
	}

	return signature.Serialize(), nil
}

// Verify 使用公钥验证签名
func (a *schnorrAlgorithm) Verify(data []byte, signature []byte, publicKey []byte) bool {
	// Hash the data
	hash := chainhash.HashB(data)

	// Parse the Schnorr public key
	pubKey, err := schnorr.ParsePubKey(publicKey)
	if err != nil {
		return false
	}

	// Parse the signature
	sig, err := schnorr.ParseSignature(signature)
	if err != nil {
		return false
	}

	return sig.Verify(hash, pubKey)
}

// GenerateAddress 从公钥生成地址
func (a *schnorrAlgorithm) GenerateAddress(publicKey []byte) types.Address {
	// Parse the Schnorr public key
	pubKey, err := schnorr.ParsePubKey(publicKey)
	if err != nil {
		// Return empty address if parsing fails
		return types.Address{}
	}

	// Hash the public key and take the last 20 bytes as the address
	pubKeyBytes := schnorr.SerializePubKey(pubKey)
	hash := sha256.Sum256(pubKeyBytes)
	var addr types.Address
	copy(addr[:], hash[len(hash)-20:])
	return addr
}

// Type 返回算法类型
func (a *schnorrAlgorithm) Type() string {
	return Schnorr
}

// PrivateKeyFromBytes 从字节数据创建私钥
func (a *schnorrAlgorithm) PrivateKeyFromBytes(data []byte) (PrivateKey, error) {
	privKey, _ := btcec.PrivKeyFromBytes(data)
	return &schnorrPrivateKey{key: privKey}, nil
}

// PublicKeyFromBytes 从字节数据创建公钥
func (a *schnorrAlgorithm) PublicKeyFromBytes(data []byte) (PublicKey, error) {
	// Parse the Schnorr public key
	pubKey, err := schnorr.ParsePubKey(data)
	if err != nil {
		return nil, fmt.Errorf("invalid Schnorr public key data: %w", err)
	}

	return &schnorrPublicKey{key: pubKey}, nil
}
