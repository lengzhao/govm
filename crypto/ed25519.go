package crypto

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"

	"github.com/lengzhao/govm/types"
)

// ed25519PrivateKey implements the PrivateKey interface for Ed25519 keys.
type ed25519PrivateKey struct {
	key ed25519.PrivateKey
}

// Bytes returns the Ed25519 private key as bytes.
func (k *ed25519PrivateKey) Bytes() []byte {
	return k.key
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
func (k *ed25519PrivateKey) Type() KeyType {
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
func (k *ed25519PublicKey) Type() KeyType {
	return Ed25519
}

// GenerateEd25519KeyPair generates a new Ed25519 key pair.
func GenerateEd25519KeyPair() (PrivateKey, PublicKey, error) {
	_, privKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, nil, err
	}

	priv := &ed25519PrivateKey{key: privKey}
	pub := &ed25519PublicKey{key: privKey.Public().(ed25519.PublicKey)}
	return priv, pub, nil
}

// Ed25519PrivateKeyFromBytes creates an Ed25519 private key from bytes.
func Ed25519PrivateKeyFromBytes(data []byte) (PrivateKey, error) {
	// For Ed25519, the private key bytes can be used directly
	privKey := ed25519.PrivateKey(data)
	return &ed25519PrivateKey{key: privKey}, nil
}
