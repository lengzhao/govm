// Package crypto provides examples of using the unified GenerateKeyPair method.
package crypto

import (
	"fmt"
)

// ExampleUnifiedKeyGeneration demonstrates how to use the unified GenerateKeyPair method.
func ExampleUnifiedKeyGeneration() {
	// Create a new crypto instance
	crypto := NewCrypto()

	// Generate different types of key pairs using the unified method
	fmt.Println("Generating key pairs using unified GenerateKeyPair method:")

	// Generate Ed25519 key pair
	edPriv, edPub, err := crypto.GenerateKeyPair(Ed25519)
	if err != nil {
		fmt.Printf("Error generating Ed25519 key pair: %v\n", err)
		return
	}
	fmt.Printf("Ed25519 - Private key type: %s, Public key type: %s\n", edPriv.Type(), edPub.Type())

	// Generate ECDSA key pair
	ecPriv, ecPub, err := crypto.GenerateKeyPair(ECDSA)
	if err != nil {
		fmt.Printf("Error generating ECDSA key pair: %v\n", err)
		return
	}
	fmt.Printf("ECDSA - Private key type: %s, Public key type: %s\n", ecPriv.Type(), ecPub.Type())

	// Generate Secp256k1 key pair
	secpPriv, secpPub, err := crypto.GenerateKeyPair(Secp256k1)
	if err != nil {
		fmt.Printf("Error generating Secp256k1 key pair: %v\n", err)
		return
	}
	fmt.Printf("Secp256k1 - Private key type: %s, Public key type: %s\n", secpPriv.Type(), secpPub.Type())

	// Generate Schnorr key pair
	schnorrPriv, schnorrPub, err := crypto.GenerateKeyPair(Schnorr)
	if err != nil {
		fmt.Printf("Error generating Schnorr key pair: %v\n", err)
		return
	}
	fmt.Printf("Schnorr - Private key type: %s, Public key type: %s\n", schnorrPriv.Type(), schnorrPub.Type())

	// Try to generate an unsupported key type
	_, _, err = crypto.GenerateKeyPair("unsupported")
	if err != nil {
		fmt.Printf("Expected error for unsupported key type: %v\n", err)
	}
}
