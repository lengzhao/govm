package crypto

import (
	"crypto/sha256"
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/lengzhao/govm/types"
)

// schnorrPrivateKey implements the PrivateKey interface for Schnorr keys.
type schnorrPrivateKey struct {
	key *btcec.PrivateKey
}

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
func (k *schnorrPrivateKey) Type() KeyType {
	return "schnorr"
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
func (k *schnorrPublicKey) Type() KeyType {
	return "schnorr"
}

// GenerateSchnorrKeyPair generates a new Schnorr key pair.
func GenerateSchnorrKeyPair() (PrivateKey, PublicKey, error) {
	privKey, err := btcec.NewPrivateKey()
	if err != nil {
		return nil, nil, err
	}

	priv := &schnorrPrivateKey{key: privKey}
	pub := &schnorrPublicKey{key: privKey.PubKey()}
	return priv, pub, nil
}
