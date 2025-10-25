# Crypto Module

The crypto module provides cryptographic functions for the govm blockchain platform. It supports both Ed25519 and ECDSA algorithms for digital signatures, along with hash functions and address generation.

## Features

- **Ed25519**: Fast and secure digital signature algorithm
- **ECDSA**: Elliptic Curve Digital Signature Algorithm
- **Hash Functions**: SHA256 and Keccak256
- **Address Generation**: Generate blockchain addresses from public keys
- **Key Encryption**: Save and load encrypted private keys with password protection

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
    
    // Keccak256 hash
    keccakHash := c.Keccak256(data)
    fmt.Printf("Keccak256 Hash: %x\n", keccakHash)
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
    GenerateKeyPair() (PrivateKey, PublicKey, error)
    GenerateEd25519KeyPair() (PrivateKey, PublicKey, error)
    GenerateECDSAKeyPair() (PrivateKey, PublicKey, error)
    Sign(data []byte, privateKey PrivateKey) ([]byte, error)
    Verify(data []byte, signature []byte, publicKey PublicKey) bool
    Hash(data []byte) types.Hash
    Keccak256(data []byte) types.Hash
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
}
```

### PublicKey Interface

```go
type PublicKey interface {
    Bytes() []byte
    Address() types.Address
    Verify(data []byte, signature []byte) bool
    Type() KeyType
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

- `golang.org/x/crypto/sha3` for Keccak256 hash function
- `golang.org/x/crypto/scrypt` for password-based key derivation

## License

MIT