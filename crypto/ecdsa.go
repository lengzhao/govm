package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"math/big"

	"github.com/lengzhao/govm/types"
)

// ecdsaPrivateKey implements the PrivateKey interface for ECDSA keys.
type ecdsaPrivateKey struct {
	key *ecdsa.PrivateKey
}

// Bytes returns the ECDSA private key as bytes.
func (k *ecdsaPrivateKey) Bytes() []byte {
	return k.key.D.Bytes()
}

// FromBytes creates an ECDSA private key from bytes.
func (k *ecdsaPrivateKey) FromBytes(data []byte) (PrivateKey, error) {
	// Create a new private key from the D value
	d := new(big.Int).SetBytes(data)

	// Create the private key with the P256 curve
	privKey := &ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{
			Curve: elliptic.P256(),
		},
		D: d,
	}

	// Calculate the public key coordinates
	privKey.PublicKey.X, privKey.PublicKey.Y = privKey.PublicKey.Curve.ScalarBaseMult(data)

	return &ecdsaPrivateKey{key: privKey}, nil
}

// PublicKey returns the corresponding ECDSA public key.
func (k *ecdsaPrivateKey) PublicKey() PublicKey {
	return &ecdsaPublicKey{key: &k.key.PublicKey}
}

// Sign signs data with the ECDSA private key.
func (k *ecdsaPrivateKey) Sign(data []byte) ([]byte, error) {
	hash := sha256.Sum256(data)
	return k.signHash(hash[:])
}

// signHash signs a hash with the ECDSA private key.
func (k *ecdsaPrivateKey) signHash(hash []byte) ([]byte, error) {
	r, s, err := ecdsa.Sign(rand.Reader, k.key, hash)
	if err != nil {
		return nil, err
	}

	// Serialize the signature (R || S)
	signature := make([]byte, 64)
	r.FillBytes(signature[:32])
	s.FillBytes(signature[32:])
	return signature, nil
}

// Type returns the key type.
func (k *ecdsaPrivateKey) Type() KeyType {
	return ECDSA
}

// ecdsaPublicKey implements the PublicKey interface for ECDSA keys.
type ecdsaPublicKey struct {
	key *ecdsa.PublicKey
}

// Bytes returns the ECDSA public key as bytes (compressed format).
func (k *ecdsaPublicKey) Bytes() []byte {
	// Compress the public key (33 bytes)
	return elliptic.MarshalCompressed(k.key.Curve, k.key.X, k.key.Y)
}

// FromBytes creates an ECDSA public key from bytes.
func (k *ecdsaPublicKey) FromBytes(data []byte) (PublicKey, error) {
	// Parse the compressed public key
	x, y := elliptic.UnmarshalCompressed(elliptic.P256(), data)
	if x == nil || y == nil {
		return nil, fmt.Errorf("invalid public key data")
	}

	pubKey := &ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     x,
		Y:     y,
	}

	return &ecdsaPublicKey{key: pubKey}, nil
}

// Address generates an address from the ECDSA public key.
func (k *ecdsaPublicKey) Address() types.Address {
	// Hash the public key and take the last 20 bytes as the address
	pubKeyBytes := elliptic.MarshalCompressed(k.key.Curve, k.key.X, k.key.Y)
	hash := sha256.Sum256(pubKeyBytes)
	var addr types.Address
	copy(addr[:], hash[len(hash)-20:])
	return addr
}

// Verify verifies a signature with the ECDSA public key.
func (k *ecdsaPublicKey) Verify(data []byte, signature []byte) bool {
	if len(signature) != 64 {
		return false
	}

	hash := sha256.Sum256(data)
	r := new(big.Int).SetBytes(signature[:32])
	s := new(big.Int).SetBytes(signature[32:])
	return ecdsa.Verify(k.key, hash[:], r, s)
}

// Type returns the key type.
func (k *ecdsaPublicKey) Type() KeyType {
	return ECDSA
}

// GenerateECDSAKeyPair generates a new ECDSA key pair using the secp256k1 curve.
func GenerateECDSAKeyPair() (PrivateKey, PublicKey, error) {
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, nil, err
	}

	priv := &ecdsaPrivateKey{key: privKey}
	pub := &ecdsaPublicKey{key: &privKey.PublicKey}
	return priv, pub, nil
}

// ECDSAPrivateKeyFromBytes creates an ECDSA private key from bytes.
func ECDSAPrivateKeyFromBytes(data []byte) (PrivateKey, error) {
	// Create a new private key from the D value
	d := new(big.Int).SetBytes(data)

	// Create the private key with the secp256k1 curve
	privKey := &ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{
			Curve: elliptic.P256(),
		},
		D: d,
	}

	// Calculate the public key coordinates
	privKey.PublicKey.X, privKey.PublicKey.Y = privKey.PublicKey.Curve.ScalarBaseMult(data)

	return &ecdsaPrivateKey{key: privKey}, nil
}
