# Crypto Package

This package provides cryptographic functions for the govm blockchain, including support for multiple signature schemes.

## Supported Signature Schemes

1. **Ed25519** - Default signature scheme
2. **ECDSA** - Elliptic Curve Digital Signature Algorithm using P-256 curve
3. **Secp256k1** - Elliptic Curve Digital Signature Algorithm using secp256k1 curve (Bitcoin standard)
4. **Schnorr** - Schnorr signatures using secp256k1 curve (Bitcoin Taproot standard)

## Key Generation

```go
// Create a new crypto instance
crypto := NewCrypto()

// Generate different types of key pairs using the unified method
ed25519Priv, ed25519Pub, _ := crypto.GenerateKeyPair(Ed25519)
ecdsaPriv, ecdsaPub, _ := crypto.GenerateKeyPair(ECDSA)
secp256k1Priv, secp256k1Pub, _ := crypto.GenerateKeyPair(Secp256k1)
schnorrPriv, schnorrPub, _ := crypto.GenerateKeyPair(Schnorr)
```

## Signing and Verification

```go
// Sign data
data := []byte("Hello, World!")
signature, _ := crypto.Sign(data, privateKey)

// Verify signature
valid := crypto.Verify(data, signature, publicKey)
```

## Address Generation

```go
// Generate address from public key
address := crypto.GenerateAddress(publicKey)
```

## Key Serialization

```go
// Serialize keys to bytes
privBytes := privateKey.Bytes()
pubBytes := publicKey.Bytes()

// Deserialize keys from bytes
reconstructedPriv, _ := privateKey.FromBytes(privBytes)
reconstructedPub, _ := publicKey.FromBytes(pubBytes)
```

## Key Storage

```go
// Save private key to encrypted file
err := SaveToFile(privateKey, "key.json", "password")

// Load private key from encrypted file
loadedPriv, err := LoadFromFile("key.json", "password")
```

## Secp256k1 and Schnorr Support

The package now includes support for Bitcoin-standard secp256k1 ECDSA and Schnorr signatures through the btcec library.

### Secp256k1 Features:
- Uses the same elliptic curve as Bitcoin (secp256k1)
- Compatible with Bitcoin wallets and tools
- Compressed public key format (33 bytes)
- Standard ECDSA signatures

### Schnorr Features:
- Implements BIP-340 Schnorr signatures
- 32-byte public keys
- 64-byte signatures
- Better privacy and efficiency than ECDSA
- Supports signature aggregation (for future multi-signature support)

## Module Structure

The crypto package is organized into separate files for each signature scheme:

- `ed25519.go` - Ed25519 signature implementation
- `ecdsa.go` - ECDSA signature implementation using P-256 curve
- `secp256k1.go` - ECDSA signature implementation using secp256k1 curve
- `schnorr.go` - Schnorr signature implementation using secp256k1 curve

## Usage Examples

See `example_secp256k1.go` for detailed examples of how to use secp256k1 and Schnorr signatures.
See `example_unified.go` for examples of how to use the unified GenerateKeyPair method.

## Installation

```bash
go get github.com/lengzhao/govm/crypto
```

## Usage

### Basic Usage

```go
package main

import (
    "fmt"
    "github.com/lengzhao/govm/crypto"
)

func main() {
    // Create a new crypto instance
    c := crypto.NewCrypto()
    
    // Generate a key pair
    priv, pub, err := c.GenerateKeyPair()
    if err != nil {
        panic(err)
    }
    
    // Sign data
    data := []byte("Hello, World!")
    signature, err := c.Sign(data, priv)
    if err != nil {
        panic(err)
    }
    
    // Verify signature
    valid := c.Verify(data, signature, pub)
    fmt.Printf("Signature valid: %v\n", valid)
    
    // Generate address
    address := c.GenerateAddress(pub)
    fmt.Printf("Address: %x\n", address)
    
    // Hash data
    hash := c.Hash(data)
    fmt.Printf("SHA256 Hash: %x\n", hash)
    
    // Hash hash
    keccakHash := c.Hash(data)
    fmt.Printf("Hash Hash: %x\n", keccakHash)
}
```

### Ed25519 Specific Usage

```go
// Generate Ed25519 key pair
priv, pub, err := c.GenerateEd25519KeyPair()
if err != nil {
    panic(err)
}

// Sign and verify
data := []byte("Hello, Ed25519!")
signature, err := c.Sign(data, priv)
if err != nil {
    panic(err)
}

valid := c.Verify(data, signature, pub)
fmt.Printf("Ed25519 signature valid: %v\n", valid)
```

### ECDSA Specific Usage

```go
// Generate ECDSA key pair
priv, pub, err := c.GenerateECDSAKeyPair()
if err != nil {
    panic(err)
}

// Sign and verify
data := []byte("Hello, ECDSA!")
signature, err := c.Sign(data, priv)
if err != nil {
    panic(err)
}

valid := c.Verify(data, signature, pub)
fmt.Printf("ECDSA signature valid: %v\n", valid)
```

### Key Encryption with Password Protection

```go
// Generate a key pair
priv, _, err := c.GenerateKeyPair()
if err != nil {
    panic(err)
}

// Save key with custom password
password := "your-secure-password"
err = crypto.SaveToFile(priv, "encrypted_key.json", password)
if err != nil {
    panic(err)
}

// Load key with custom password
loadedPriv, err := crypto.LoadFromFile("encrypted_key.json", password)
if err != nil {
    panic(err)
}

// Verify loaded key
fmt.Printf("Loaded key type: %s\n", loadedPriv.Type())
```

### Using Default Password

```go
// Save key with default password (empty string will use "govm-password")
err = crypto.SaveToFile(priv, "encrypted_key.json", "")
if err != nil {
    panic(err)
}

// Load key with default password (empty string will use "govm-password")
loadedPriv, err := crypto.LoadFromFile("encrypted_key.json", "")
if err != nil {
    panic(err)
}
```

## Code Structure

The crypto module is organized into several files:

- `crypto.go`: Main interface definitions and core logic
- `ed25519.go`: Ed25519 algorithm implementation
- `ecdsa.go`: ECDSA algorithm implementation
- `crypto_test.go`: Unit tests for all functionality

## API Reference

### Crypto Interface

```go
type Crypto interface {
    GenerateKeyPair(keyType KeyType) (PrivateKey, PublicKey, error)
    Sign(data []byte, privateKey PrivateKey) ([]byte, error)
    Verify(data []byte, signature []byte, publicKey PublicKey) bool
    Hash(data []byte) types.Hash
    Hash(data []byte) types.Hash
    GenerateAddress(publicKey PublicKey) types.Address
}
```

### PrivateKey Interface

```go
type PrivateKey interface {
    Bytes() []byte
    PublicKey() PublicKey
    Sign(data []byte) ([]byte, error)
    Type() KeyType
    FromBytes(data []byte) (PrivateKey, error)
}
```

### PublicKey Interface

```go
type PublicKey interface {
    Bytes() []byte
    Address() types.Address
    Verify(data []byte, signature []byte) bool
    Type() KeyType
    FromBytes(data []byte) (PublicKey, error)
}
```

### Key Storage Functions

```go
// Save private key to file (encrypted with password, or default password if empty)
func SaveToFile(privateKey PrivateKey, filename string, password string) error

// Load private key from file (using provided password, or default password if empty)
func LoadFromFile(filename string, password string) (PrivateKey, error)
```

## Security Considerations

1. **Key Storage**: Private keys should be stored securely and never exposed in plaintext.
2. **Random Number Generation**: The module uses crypto/rand for secure random number generation.
3. **Side-channel Attacks**: The implementation uses constant-time algorithms where possible.
4. **Signature Malleability**: ECDSA signatures are not normalized, which may lead to malleability issues.
5. **Password Security**: When using password protection, use strong passwords to prevent brute-force attacks.
6. **Key Derivation**: The module uses scrypt with strong parameters (N=32768, r=8, p=1) for key derivation.
7. **Encryption**: AES-GCM is used for encryption, providing both confidentiality and authenticity.
8. **Default Password**: When no password is provided, the default password "govm-password" is used.

## Testing

To run the tests:

```bash
go test -v ./crypto
```

## Dependencies

- `golang.org/x/crypto/sha3` for Hash hash function
- `golang.org/x/crypto/scrypt` for password-based key derivation
- `github.com/btcsuite/btcd/btcec/v2` for secp256k1 and Schnorr signatures

## License

MIT