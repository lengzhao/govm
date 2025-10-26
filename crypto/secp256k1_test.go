package crypto

import (
	"testing"
)

func TestSecp256k1KeyGeneration(t *testing.T) {
	crypto := NewCrypto()

	priv, pub, err := crypto.GenerateKeyPair(Secp256k1)
	if err != nil {
		t.Fatalf("Failed to generate secp256k1 key pair: %v", err)
	}

	if priv == nil {
		t.Error("Private key is nil")
	}

	if pub == nil {
		t.Error("Public key is nil")
	}

	if priv.Type() != Secp256k1 {
		t.Errorf("Expected private key type Secp256k1, got %s", priv.Type())
	}

	if pub.Type() != Secp256k1 {
		t.Errorf("Expected public key type Secp256k1, got %s", pub.Type())
	}
}

func TestSecp256k1SigningAndVerification(t *testing.T) {
	crypto := NewCrypto()

	priv, pub, err := crypto.GenerateKeyPair(Secp256k1)
	if err != nil {
		t.Fatalf("Failed to generate secp256k1 key pair: %v", err)
	}

	data := []byte("test message")
	signature, err := crypto.Sign(data, priv)
	if err != nil {
		t.Fatalf("Failed to sign data: %v", err)
	}

	if signature == nil {
		t.Error("Signature is nil")
	}

	valid := crypto.Verify(data, signature, pub)
	if !valid {
		t.Error("Failed to verify signature")
	}

	// Test with wrong data
	wrongData := []byte("wrong message")
	valid = crypto.Verify(wrongData, signature, pub)
	if valid {
		t.Error("Verification should fail with wrong data")
	}
}

func TestSecp256k1AddressGeneration(t *testing.T) {
	crypto := NewCrypto()

	_, pub, err := crypto.GenerateKeyPair(Secp256k1)
	if err != nil {
		t.Fatalf("Failed to generate secp256k1 key pair: %v", err)
	}

	address := crypto.GenerateAddress(pub)
	if address == [20]byte{} {
		t.Error("Address is empty")
	}

	// Verify that the address can be derived from the public key directly
	derivedAddress := pub.Address()
	if address != derivedAddress {
		t.Error("Address mismatch between crypto.GenerateAddress and pub.Address()")
	}
}

func TestSecp256k1PrivateKeyBytes(t *testing.T) {
	crypto := NewCrypto()

	// Test secp256k1 private key bytes
	secpPriv, _, err := crypto.GenerateKeyPair(Secp256k1)
	if err != nil {
		t.Fatalf("Failed to generate secp256k1 key pair: %v", err)
	}

	secpBytes := secpPriv.Bytes()
	if len(secpBytes) == 0 {
		t.Error("secp256k1 private key bytes are empty")
	}
}

func TestSecp256k1PublicKeyBytes(t *testing.T) {
	crypto := NewCrypto()

	// Test secp256k1 public key bytes
	_, secpPub, err := crypto.GenerateKeyPair(Secp256k1)
	if err != nil {
		t.Fatalf("Failed to generate secp256k1 key pair: %v", err)
	}

	secpBytes := secpPub.Bytes()
	if len(secpBytes) == 0 {
		t.Error("secp256k1 public key bytes are empty")
	}
}

func TestSecp256k1PublicKeyFromPrivateKey(t *testing.T) {
	crypto := NewCrypto()

	// Test secp256k1
	secpPriv, secpPub, err := crypto.GenerateKeyPair(Secp256k1)
	if err != nil {
		t.Fatalf("Failed to generate secp256k1 key pair: %v", err)
	}

	derivedPub := secpPriv.PublicKey()
	if derivedPub.Type() != secpPub.Type() {
		t.Error("Derived public key type mismatch")
	}
}

// Test FromBytes methods
func TestSecp256k1PrivateKeyFromBytes(t *testing.T) {
	crypto := NewCrypto()

	// Test secp256k1 private key FromBytes
	secpPriv, _, err := crypto.GenerateKeyPair(Secp256k1)
	if err != nil {
		t.Fatalf("Failed to generate secp256k1 key pair: %v", err)
	}

	secpBytes := secpPriv.Bytes()
	reconstructedSecpPriv, err := secpPriv.FromBytes(secpBytes)
	if err != nil {
		t.Fatalf("Failed to reconstruct secp256k1 private key from bytes: %v", err)
	}

	if reconstructedSecpPriv.Type() != secpPriv.Type() {
		t.Errorf("Expected reconstructed key type %s, got %s", secpPriv.Type(), reconstructedSecpPriv.Type())
	}

	if string(reconstructedSecpPriv.Bytes()) != string(secpPriv.Bytes()) {
		t.Error("Reconstructed secp256k1 private key bytes do not match original")
	}
}

func TestSecp256k1PublicKeyFromBytesTest(t *testing.T) {
	crypto := NewCrypto()

	// Test secp256k1 public key FromBytes
	_, secpPub, err := crypto.GenerateKeyPair(Secp256k1)
	if err != nil {
		t.Fatalf("Failed to generate secp256k1 key pair: %v", err)
	}

	secpBytes := secpPub.Bytes()
	reconstructedSecpPub, err := secpPub.FromBytes(secpBytes)
	if err != nil {
		t.Fatalf("Failed to reconstruct secp256k1 public key from bytes: %v", err)
	}

	if reconstructedSecpPub.Type() != secpPub.Type() {
		t.Errorf("Expected reconstructed key type %s, got %s", secpPub.Type(), reconstructedSecpPub.Type())
	}

	if string(reconstructedSecpPub.Bytes()) != string(secpPub.Bytes()) {
		t.Error("Reconstructed secp256k1 public key bytes do not match original")
	}

	// Verify that both keys produce the same address
	originalAddr := secpPub.Address()
	reconstructedAddr := reconstructedSecpPub.Address()
	if originalAddr != reconstructedAddr {
		t.Error("Reconstructed secp256k1 public key produces different address")
	}
}
