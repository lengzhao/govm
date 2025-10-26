package crypto

import (
	"crypto/sha256"
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/ecdsa"
	"github.com/lengzhao/govm/types"
)

// secp256k1PrivateKey implements the PrivateKey interface for secp256k1 keys.
type secp256k1PrivateKey struct {
	key *btcec.PrivateKey
}

// Bytes returns the secp256k1 private key as bytes.
func (k *secp256k1PrivateKey) Bytes() []byte {
	return k.key.Serialize()
}

// FromBytes creates a secp256k1 private key from bytes.
func (k *secp256k1PrivateKey) FromBytes(data []byte) (PrivateKey, error) {
	privKey, _ := btcec.PrivKeyFromBytes(data)
	return &secp256k1PrivateKey{key: privKey}, nil
}

// PublicKey returns the corresponding secp256k1 public key.
func (k *secp256k1PrivateKey) PublicKey() PublicKey {
	return &secp256k1PublicKey{key: k.key.PubKey()}
}

// Sign signs data with the secp256k1 private key using ECDSA.
func (k *secp256k1PrivateKey) Sign(data []byte) ([]byte, error) {
	hash := sha256.Sum256(data)
	signature := ecdsa.Sign(k.key, hash[:])

	// Serialize the signature
	return signature.Serialize(), nil
}

// Type returns the key type.
func (k *secp256k1PrivateKey) Type() KeyType {
	return "secp256k1"
}

// secp256k1PublicKey implements the PublicKey interface for secp256k1 keys.
type secp256k1PublicKey struct {
	key *btcec.PublicKey
}

// Bytes returns the secp256k1 public key as bytes (compressed format).
func (k *secp256k1PublicKey) Bytes() []byte {
	// Compress the public key (33 bytes)
	return k.key.SerializeCompressed()
}

// FromBytes creates a secp256k1 public key from bytes.
func (k *secp256k1PublicKey) FromBytes(data []byte) (PublicKey, error) {
	// Parse the compressed public key
	pubKey, err := btcec.ParsePubKey(data)
	if err != nil {
		return nil, fmt.Errorf("invalid public key data: %w", err)
	}

	return &secp256k1PublicKey{key: pubKey}, nil
}

// Address generates an address from the secp256k1 public key.
func (k *secp256k1PublicKey) Address() types.Address {
	// Hash the public key and take the last 20 bytes as the address
	pubKeyBytes := k.key.SerializeCompressed()
	hash := sha256.Sum256(pubKeyBytes)
	var addr types.Address
	copy(addr[:], hash[len(hash)-20:])
	return addr
}

// Verify verifies a signature with the secp256k1 public key.
func (k *secp256k1PublicKey) Verify(data []byte, signature []byte) bool {
	hash := sha256.Sum256(data)

	// Parse the signature
	sig, err := ecdsa.ParseSignature(signature)
	if err != nil {
		return false
	}

	return sig.Verify(hash[:], k.key)
}

// Type returns the key type.
func (k *secp256k1PublicKey) Type() KeyType {
	return "secp256k1"
}

// GenerateSecp256k1KeyPair generates a new secp256k1 key pair.
func GenerateSecp256k1KeyPair() (PrivateKey, PublicKey, error) {
	privKey, err := btcec.NewPrivateKey()
	if err != nil {
		return nil, nil, err
	}

	priv := &secp256k1PrivateKey{key: privKey}
	pub := &secp256k1PublicKey{key: privKey.PubKey()}
	return priv, pub, nil
}
