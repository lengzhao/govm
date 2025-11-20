// Package crypto provides examples of using secp256k1 and Schnorr signatures.
package crypto

import (
	"fmt"
)

// ExampleSecp256k1 demonstrates how to use secp256k1 key pairs for signing and verification.
func ExampleSecp256k1() {
	// Create a new crypto instance
	crypto := NewCrypto()

	// Generate a secp256k1 key pair
	privKey, pubKey, err := crypto.GenerateKeyPair(Secp256k1)
	if err != nil {
		fmt.Printf("Error generating secp256k1 key pair: %v\n", err)
		return
	}

	fmt.Printf("Generated secp256k1 key pair\n")
	fmt.Printf("Private key type: %s\n", privKey.Type())
	fmt.Printf("Public key type: %s\n", pubKey.Type())

	// Sign some data
	data := []byte("Hello, secp256k1!")
	signature, err := crypto.Sign(data, privKey.Bytes(), Secp256k1)
	if err != nil {
		fmt.Printf("Error signing data: %v\n", err)
		return
	}

	fmt.Printf("Data signed successfully, signature length: %d bytes\n", len(signature))

	// Verify the signature
	valid := crypto.Verify(data, signature, pubKey.Bytes(), Secp256k1)
	fmt.Printf("Signature verification result: %t\n", valid)

	// Try to verify with wrong data
	wrongData := []byte("Wrong data")
	valid = crypto.Verify(wrongData, signature, pubKey.Bytes(), Secp256k1)
	fmt.Printf("Signature verification with wrong data: %t\n", valid)

	// Generate an address from the public key
	address := crypto.GenerateAddress(pubKey.Bytes(), Secp256k1)
	fmt.Printf("Generated address: %x\n", address)

	// Demonstrate key serialization and deserialization
	privBytes := privKey.Bytes()
	pubBytes := pubKey.Bytes()

	fmt.Printf("Private key bytes length: %d\n", len(privBytes))
	fmt.Printf("Public key bytes length: %d\n", len(pubBytes))

	// Reconstruct keys from bytes
	reconstructedPriv, err := privKey.FromBytes(privBytes)
	if err != nil {
		fmt.Printf("Error reconstructing private key: %v\n", err)
		return
	}

	reconstructedPub, err := pubKey.FromBytes(pubBytes)
	if err != nil {
		fmt.Printf("Error reconstructing public key: %v\n", err)
		return
	}

	fmt.Printf("Reconstructed private key type: %s\n", reconstructedPriv.Type())
	fmt.Printf("Reconstructed public key type: %s\n", reconstructedPub.Type())

	// Verify that reconstructed keys work
	reconstructedSignature, err := crypto.Sign(data, reconstructedPriv.Bytes(), Secp256k1)
	if err != nil {
		fmt.Printf("Error signing with reconstructed private key: %v\n", err)
		return
	}

	valid = crypto.Verify(data, reconstructedSignature, reconstructedPub.Bytes(), Secp256k1)
	fmt.Printf("Signature verification with reconstructed keys: %t\n", valid)
}

// ExampleSchnorr demonstrates how to use Schnorr key pairs for signing and verification.
func ExampleSchnorr() {
	// Create a new crypto instance
	crypto := NewCrypto()

	// Generate a Schnorr key pair
	privKey, pubKey, err := crypto.GenerateKeyPair(Schnorr)
	if err != nil {
		fmt.Printf("Error generating Schnorr key pair: %v\n", err)
		return
	}

	fmt.Printf("Generated Schnorr key pair\n")
	fmt.Printf("Private key type: %s\n", privKey.Type())
	fmt.Printf("Public key type: %s\n", pubKey.Type())

	// Sign some data
	data := []byte("Hello, Schnorr!")
	signature, err := crypto.Sign(data, privKey.Bytes(), Schnorr)
	if err != nil {
		fmt.Printf("Error signing data: %v\n", err)
		return
	}

	fmt.Printf("Data signed successfully, signature length: %d bytes\n", len(signature))

	// Verify the signature
	valid := crypto.Verify(data, signature, pubKey.Bytes(), Schnorr)
	fmt.Printf("Signature verification result: %t\n", valid)

	// Try to verify with wrong data
	wrongData := []byte("Wrong data")
	valid = crypto.Verify(wrongData, signature, pubKey.Bytes(), Schnorr)
	fmt.Printf("Signature verification with wrong data: %t\n", valid)

	// Generate an address from the public key
	address := crypto.GenerateAddress(pubKey.Bytes(), Schnorr)
	fmt.Printf("Generated address: %x\n", address)

	// Demonstrate key serialization and deserialization
	privBytes := privKey.Bytes()
	pubBytes := pubKey.Bytes()

	fmt.Printf("Private key bytes length: %d\n", len(privBytes))
	fmt.Printf("Public key bytes length: %d\n", len(pubBytes))

	// Reconstruct keys from bytes
	reconstructedPriv, err := privKey.FromBytes(privBytes)
	if err != nil {
		fmt.Printf("Error reconstructing private key: %v\n", err)
		return
	}

	reconstructedPub, err := pubKey.FromBytes(pubBytes)
	if err != nil {
		fmt.Printf("Error reconstructing public key: %v\n", err)
		return
	}

	fmt.Printf("Reconstructed private key type: %s\n", reconstructedPriv.Type())
	fmt.Printf("Reconstructed public key type: %s\n", reconstructedPub.Type())

	// Verify that reconstructed keys work
	reconstructedSignature, err := crypto.Sign(data, reconstructedPriv.Bytes(), Schnorr)
	if err != nil {
		fmt.Printf("Error signing with reconstructed private key: %v\n", err)
		return
	}

	valid = crypto.Verify(data, reconstructedSignature, reconstructedPub.Bytes(), Schnorr)
	fmt.Printf("Signature verification with reconstructed keys: %t\n", valid)
}
