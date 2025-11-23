# Crypto Package

This package provides cryptographic functions for the govm blockchain, including support for multiple signature schemes.

## Supported Signature Schemes

1. **Ed25519** - Default signature scheme
2. **ECDSA** - Elliptic Curve Digital Signature Algorithm using P-256 curve
3. **Secp256k1** - Elliptic Curve Digital Signature Algorithm using secp256k1 curve (Bitcoin standard)
4. **Schnorr** - Schnorr signatures using secp256k1 curve (Bitcoin Taproot standard)

## Key Generation

```go
// Generate different types of key pairs using the unified method
ed25519Priv, ed25519Pub, _ := crypto.GenerateKeyPair(crypto.Ed25519)
ecdsaPriv, ecdsaPub, _ := crypto.GenerateKeyPair(crypto.ECDSA)
secp256k1Priv, secp256k1Pub, _ := crypto.GenerateKeyPair(crypto.Secp256k1)
schnorrPriv, schnorrPub, _ := crypto.GenerateKeyPair(crypto.Schnorr)
```

## Signing and Verification

```go
// Sign data
data := []byte("Hello, World!")
signature, _ := crypto.Sign(data, privateKey, keyType)

// Verify signature
valid := crypto.Verify(data, signature, publicKey, keyType)
```

## Address Generation

```go
// Generate address from public key
address := crypto.GenerateAddress(publicKey, keyType)
```

## Key Serialization

Keys can be serialized and deserialized through the Algorithm interface:

```go
// Get algorithm instance
algorithm, _ := crypto.AlgorithmFactory(keyType)

// Serialize keys to bytes
privBytes := privateKey.Bytes()
pubBytes := publicKey.Bytes()

// Deserialize keys from bytes
reconstructedPriv, _ := algorithm.PrivateKeyFromBytes(privBytes)
reconstructedPub, _ := algorithm.PublicKeyFromBytes(pubBytes)
```

## Key Storage

```go
// Save private key to encrypted file
err := crypto.SaveToFile(privateKey, "key.json", "password")

// Load private key from encrypted file
loadedPriv, err := crypto.LoadFromFile("key.json", "password")
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

See the examples directory for detailed usage examples.

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
    // Generate a key pair
    priv, pub, err := crypto.GenerateKeyPair(crypto.Ed25519)
    if err != nil {
        panic(err)
    }
    
    // Sign data
    data := []byte("Hello, World!")
    signature, err := crypto.Sign(data, priv, crypto.Ed25519)
    if err != nil {
        panic(err)
    }
    
    // Verify signature
    valid := crypto.Verify(data, signature, pub, crypto.Ed25519)
    fmt.Printf("Signature valid: %v\n", valid)
    
    // Generate address
    address := crypto.GenerateAddress(pub, crypto.Ed25519)
    fmt.Printf("Address: %x\n", address)
    
    // Hash data
    hash := crypto.Hash(data)
    fmt.Printf("SHA256 Hash: %x\n", hash)
}
```

### Key Encryption with Password Protection

```go
// Generate a key pair
priv, _, err := crypto.GenerateKeyPair(crypto.Ed25519)
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
- `secp256k1.go`: Secp256k1 algorithm implementation
- `schnorr.go`: Schnorr algorithm implementation
- `crypto_test.go`: Unit tests for all functionality

## API Reference

All cryptographic operations are accessed through the unified functions in `crypto.go`. Direct algorithm-specific functions have been removed to simplify the API and reduce module usage complexity.

### Unified Functions

```go
func GenerateKeyPair(keyType string) ([]byte, []byte, error)
func Sign(data []byte, privateKey []byte, keyType string) ([]byte, error)
func Verify(data []byte, signature []byte, publicKey []byte, keyType string) bool
func GenerateAddress(publicKey []byte, keyType string) types.Address
func Hash(data []byte) types.Hash
```

### Algorithm Interface

```go
type Algorithm interface {
    GenerateKeyPair() (privateKey []byte, publicKey []byte, error error)
    Sign(data []byte, privateKey []byte) ([]byte, error)
    Verify(data []byte, signature []byte, publicKey []byte) bool
    GenerateAddress(publicKey []byte) types.Address
    Type() string
    PrivateKeyFromBytes(data []byte) (PrivateKey, error)
    PublicKeyFromBytes(data []byte) (PublicKey, error)
}
```

### PrivateKey Interface

```go
type PrivateKey interface {
    Bytes() []byte
    PublicKey() PublicKey
    Sign(data []byte) ([]byte, error)
    Type() string
}
```

### PublicKey Interface

```go
type PublicKey interface {
    Bytes() []byte
    Address() types.Address
    Verify(data []byte, signature []byte) bool
    Type() string
}
```

### Key Storage Functions

```go
// Save private key to file (encrypted with password, or default password if empty)
func SaveToFile(privateKey PrivateKey, filename string, password string) error

// Load private key from file (using provided password, or default password if empty)
func LoadFromFile(filename string, password string) (PrivateKey, error)
```

## Algorithm Registration

The crypto package uses a factory pattern for algorithm registration. Each algorithm is registered at initialization time and can be accessed through the factory:

```go
// Register an algorithm (typically done in the init() function of each algorithm file)
RegisterAlgorithm(algorithm Algorithm) error

// Get an algorithm instance
algorithm, err := AlgorithmFactory(keyType string) (Algorithm, error)
```

For stateless algorithms (which is the case for all current implementations), the same instance is reused across all operations, which reduces memory allocation and improves performance. The algorithm type is determined by calling the `Type()` method on the algorithm instance, so there's no need to pass the key type separately during registration.

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