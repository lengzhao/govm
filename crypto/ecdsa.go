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

func init() {
	RegisterAlgorithm(&ecdsaAlgorithm{})
}

// ecdsaPrivateKey implements the PrivateKey interface for ECDSA keys.
type ecdsaPrivateKey struct {
	key *ecdsa.PrivateKey
}

// ecdsaAlgorithm implements the Algorithm interface for ECDSA.
type ecdsaAlgorithm struct{}

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
func (k *ecdsaPrivateKey) Type() string {
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
func (k *ecdsaPublicKey) Type() string {
	return ECDSA
}

// GenerateKeyPair 生成新的密钥对
func (a *ecdsaAlgorithm) GenerateKeyPair() ([]byte, []byte, error) {
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, nil, err
	}

	// 返回私钥的D值和压缩的公钥
	privKeyBytes := privKey.D.Bytes()
	pubKeyBytes := elliptic.MarshalCompressed(privKey.PublicKey.Curve, privKey.PublicKey.X, privKey.PublicKey.Y)

	return privKeyBytes, pubKeyBytes, nil
}

// Sign 使用私钥对数据进行签名
func (a *ecdsaAlgorithm) Sign(data []byte, privateKey []byte) ([]byte, error) {
	// Create a new private key from the D value
	d := new(big.Int).SetBytes(privateKey)

	// Create the private key with the P256 curve
	privKey := &ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{
			Curve: elliptic.P256(),
		},
		D: d,
	}

	// Calculate the public key coordinates
	privKey.PublicKey.X, privKey.PublicKey.Y = privKey.PublicKey.Curve.ScalarBaseMult(privateKey)

	hash := sha256.Sum256(data)
	return a.signHash(hash[:], privKey)
}

// signHash signs a hash with the ECDSA private key.
func (a *ecdsaAlgorithm) signHash(hash []byte, privKey *ecdsa.PrivateKey) ([]byte, error) {
	r, s, err := ecdsa.Sign(rand.Reader, privKey, hash)
	if err != nil {
		return nil, err
	}

	// Serialize the signature (R || S)
	signature := make([]byte, 64)
	r.FillBytes(signature[:32])
	s.FillBytes(signature[32:])
	return signature, nil
}

// Verify 使用公钥验证签名
func (a *ecdsaAlgorithm) Verify(data []byte, signature []byte, publicKey []byte) bool {
	if len(signature) != 64 {
		return false
	}

	// Parse the compressed public key
	x, y := elliptic.UnmarshalCompressed(elliptic.P256(), publicKey)
	if x == nil || y == nil {
		return false
	}

	pubKey := &ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     x,
		Y:     y,
	}

	hash := sha256.Sum256(data)
	r := new(big.Int).SetBytes(signature[:32])
	s := new(big.Int).SetBytes(signature[32:])
	return ecdsa.Verify(pubKey, hash[:], r, s)
}

// GenerateAddress 从公钥生成地址
func (a *ecdsaAlgorithm) GenerateAddress(publicKey []byte) types.Address {
	// Parse the compressed public key
	x, y := elliptic.UnmarshalCompressed(elliptic.P256(), publicKey)
	if x == nil || y == nil {
		// Return empty address if parsing fails
		return types.Address{}
	}

	pubKey := &ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     x,
		Y:     y,
	}

	// Hash the public key and take the last 20 bytes as the address
	pubKeyBytes := elliptic.MarshalCompressed(pubKey.Curve, pubKey.X, pubKey.Y)
	hash := sha256.Sum256(pubKeyBytes)
	var addr types.Address
	copy(addr[:], hash[len(hash)-20:])
	return addr
}

// Type 返回算法类型
func (a *ecdsaAlgorithm) Type() string {
	return ECDSA
}

// PrivateKeyFromBytes 从字节数据创建私钥
func (a *ecdsaAlgorithm) PrivateKeyFromBytes(data []byte) (PrivateKey, error) {
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

// PublicKeyFromBytes 从字节数据创建公钥
func (a *ecdsaAlgorithm) PublicKeyFromBytes(data []byte) (PublicKey, error) {
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
