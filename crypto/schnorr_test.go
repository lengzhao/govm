package crypto

import (
	"testing"
)

func TestSchnorrKeyGeneration(t *testing.T) {
	crypto := NewCrypto()

	priv, pub, err := crypto.GenerateKeyPair(Schnorr)
	if err != nil {
		t.Fatalf("Failed to generate Schnorr key pair: %v", err)
	}

	if priv == nil {
		t.Error("Private key is nil")
	}

	if pub == nil {
		t.Error("Public key is nil")
	}

	if priv.Type() != Schnorr {
		t.Errorf("Expected private key type Schnorr, got %s", priv.Type())
	}

	if pub.Type() != Schnorr {
		t.Errorf("Expected public key type Schnorr, got %s", pub.Type())
	}
}

func TestSchnorrSigningAndVerification(t *testing.T) {
	crypto := NewCrypto()

	priv, pub, err := crypto.GenerateKeyPair(Schnorr)
	if err != nil {
		t.Fatalf("Failed to generate Schnorr key pair: %v", err)
	}

	data := []byte("test message")
	signature, err := crypto.Sign(data, priv.Bytes(), Schnorr)
	if err != nil {
		t.Fatalf("Failed to sign data: %v", err)
	}

	if signature == nil {
		t.Error("Signature is nil")
	}

	valid := crypto.Verify(data, signature, pub.Bytes(), Schnorr)
	if !valid {
		t.Error("Failed to verify signature")
	}

	// Test with wrong data
	wrongData := []byte("wrong message")
	valid = crypto.Verify(wrongData, signature, pub.Bytes(), Schnorr)
	if valid {
		t.Error("Verification should fail with wrong data")
	}
}

func TestSchnorrAddressGeneration(t *testing.T) {
	crypto := NewCrypto()

	_, pub, err := crypto.GenerateKeyPair(Schnorr)
	if err != nil {
		t.Fatalf("Failed to generate Schnorr key pair: %v", err)
	}

	address := crypto.GenerateAddress(pub.Bytes(), Schnorr)
	if address == [20]byte{} {
		t.Error("Address is empty")
	}

	// Verify that the address can be derived from the public key directly
	derivedAddress := pub.Address()
	if address != derivedAddress {
		t.Error("Address mismatch between crypto.GenerateAddress and pub.Address()")
	}
}

func TestSchnorrPrivateKeyBytes(t *testing.T) {
	crypto := NewCrypto()

	// Test Schnorr private key bytes
	schnorrPriv, _, err := crypto.GenerateKeyPair(Schnorr)
	if err != nil {
		t.Fatalf("Failed to generate Schnorr key pair: %v", err)
	}

	schnorrBytes := schnorrPriv.Bytes()
	if len(schnorrBytes) == 0 {
		t.Error("Schnorr private key bytes are empty")
	}
}

func TestSchnorrPublicKeyBytes(t *testing.T) {
	crypto := NewCrypto()

	// Test Schnorr public key bytes
	_, schnorrPub, err := crypto.GenerateKeyPair(Schnorr)
	if err != nil {
		t.Fatalf("Failed to generate Schnorr key pair: %v", err)
	}

	schnorrBytes := schnorrPub.Bytes()
	if len(schnorrBytes) == 0 {
		t.Error("Schnorr public key bytes are empty")
	}
}

func TestSchnorrPublicKeyFromPrivateKey(t *testing.T) {
	crypto := NewCrypto()

	// Test Schnorr
	schnorrPriv, schnorrPub, err := crypto.GenerateKeyPair(Schnorr)
	if err != nil {
		t.Fatalf("Failed to generate Schnorr key pair: %v", err)
	}

	derivedSchnorrPub := schnorrPriv.PublicKey()
	if derivedSchnorrPub.Type() != schnorrPub.Type() {
		t.Error("Derived public key type mismatch")
	}
}

// Test FromBytes methods
func TestSchnorrPrivateKeyFromBytes(t *testing.T) {
	crypto := NewCrypto()

	// Test Schnorr private key FromBytes
	schnorrPriv, _, err := crypto.GenerateKeyPair(Schnorr)
	if err != nil {
		t.Fatalf("Failed to generate Schnorr key pair: %v", err)
	}

	schnorrBytes := schnorrPriv.Bytes()
	reconstructedSchnorrPriv, err := schnorrPriv.FromBytes(schnorrBytes)
	if err != nil {
		t.Fatalf("Failed to reconstruct Schnorr private key from bytes: %v", err)
	}

	if reconstructedSchnorrPriv.Type() != schnorrPriv.Type() {
		t.Errorf("Expected reconstructed key type %s, got %s", schnorrPriv.Type(), reconstructedSchnorrPriv.Type())
	}

	if string(reconstructedSchnorrPriv.Bytes()) != string(schnorrPriv.Bytes()) {
		t.Error("Reconstructed Schnorr private key bytes do not match original")
	}
}

func TestSchnorrPublicKeyFromBytesTest(t *testing.T) {
	crypto := NewCrypto()

	// Test Schnorr public key FromBytes
	_, schnorrPub, err := crypto.GenerateKeyPair(Schnorr)
	if err != nil {
		t.Fatalf("Failed to generate Schnorr key pair: %v", err)
	}

	schnorrBytes := schnorrPub.Bytes()
	reconstructedSchnorrPub, err := schnorrPub.FromBytes(schnorrBytes)
	if err != nil {
		t.Fatalf("Failed to reconstruct Schnorr public key from bytes: %v", err)
	}

	if reconstructedSchnorrPub.Type() != schnorrPub.Type() {
		t.Errorf("Expected reconstructed key type %s, got %s", schnorrPub.Type(), reconstructedSchnorrPub.Type())
	}

	if string(reconstructedSchnorrPub.Bytes()) != string(schnorrPub.Bytes()) {
		t.Error("Reconstructed Schnorr public key bytes do not match original")
	}

	// Verify that both keys produce the same address
	originalAddr := schnorrPub.Address()
	reconstructedAddr := reconstructedSchnorrPub.Address()
	if originalAddr != reconstructedAddr {
		t.Error("Reconstructed Schnorr public key produces different address")
	}
}
