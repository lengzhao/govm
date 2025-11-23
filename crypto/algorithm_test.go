package crypto

import (
	"testing"

	"github.com/lengzhao/govm/types"
)

func TestEd25519Algorithm(t *testing.T) {
	algorithm := &ed25519Algorithm{}

	// Test Type
	if algorithm.Type() != string(Ed25519) {
		t.Errorf("Expected type Ed25519, got %s", algorithm.Type())
	}

	// Test GenerateKeyPair
	priv, pub, err := algorithm.GenerateKeyPair()
	if err != nil {
		t.Fatalf("Failed to generate key pair: %v", err)
	}

	if priv == nil || pub == nil {
		t.Fatal("Private or public key is nil")
	}

	// Test Sign and Verify
	data := []byte("test data")
	signature, err := algorithm.Sign(data, priv)
	if err != nil {
		t.Fatalf("Failed to sign data: %v", err)
	}

	if !algorithm.Verify(data, signature, pub) {
		t.Error("Failed to verify signature")
	}

	// Test GenerateAddress
	address := algorithm.GenerateAddress(pub)
	if address == [20]byte{} {
		t.Error("Generated address is empty")
	}
}

func TestECDSAAlgorithm(t *testing.T) {
	algorithm := &ecdsaAlgorithm{}

	// Test Type
	if algorithm.Type() != string(ECDSA) {
		t.Errorf("Expected type ECDSA, got %s", algorithm.Type())
	}

	// Test GenerateKeyPair
	priv, pub, err := algorithm.GenerateKeyPair()
	if err != nil {
		t.Fatalf("Failed to generate key pair: %v", err)
	}

	if priv == nil || pub == nil {
		t.Fatal("Private or public key is nil")
	}

	// Test Sign and Verify
	data := []byte("test data")
	signature, err := algorithm.Sign(data, priv)
	if err != nil {
		t.Fatalf("Failed to sign data: %v", err)
	}

	if !algorithm.Verify(data, signature, pub) {
		t.Error("Failed to verify signature")
	}

	// Test GenerateAddress
	address := algorithm.GenerateAddress(pub)
	if address == [20]byte{} {
		t.Error("Generated address is empty")
	}
}

func TestSecp256k1Algorithm(t *testing.T) {
	algorithm := &secp256k1Algorithm{}

	// Test Type
	if algorithm.Type() != string(Secp256k1) {
		t.Errorf("Expected type Secp256k1, got %s", algorithm.Type())
	}

	// Test GenerateKeyPair
	priv, pub, err := algorithm.GenerateKeyPair()
	if err != nil {
		t.Fatalf("Failed to generate key pair: %v", err)
	}

	if priv == nil || pub == nil {
		t.Fatal("Private or public key is nil")
	}

	// Test Sign and Verify
	data := []byte("test data")
	signature, err := algorithm.Sign(data, priv)
	if err != nil {
		t.Fatalf("Failed to sign data: %v", err)
	}

	if !algorithm.Verify(data, signature, pub) {
		t.Error("Failed to verify signature")
	}

	// Test GenerateAddress
	address := algorithm.GenerateAddress(pub)
	if address == [20]byte{} {
		t.Error("Generated address is empty")
	}
}

func TestSchnorrAlgorithm(t *testing.T) {
	algorithm := &schnorrAlgorithm{}

	// Test Type
	if algorithm.Type() != string(Schnorr) {
		t.Errorf("Expected type Schnorr, got %s", algorithm.Type())
	}

	// Test GenerateKeyPair
	priv, pub, err := algorithm.GenerateKeyPair()
	if err != nil {
		t.Fatalf("Failed to generate key pair: %v", err)
	}

	if priv == nil || pub == nil {
		t.Fatal("Private or public key is nil")
	}

	// Test Sign and Verify
	data := []byte("test data")
	signature, err := algorithm.Sign(data, priv)
	if err != nil {
		t.Fatalf("Failed to sign data: %v", err)
	}

	if !algorithm.Verify(data, signature, pub) {
		t.Error("Failed to verify signature")
	}

	// Test GenerateAddress
	address := algorithm.GenerateAddress(pub)
	if address == [20]byte{} {
		t.Error("Generated address is empty")
	}
}

func TestAlgorithmFactory(t *testing.T) {
	// Test Ed25519
	algorithm, err := algorithmFactory(Ed25519)
	if err != nil {
		t.Fatalf("Failed to create Ed25519 algorithm: %v", err)
	}

	if algorithm.Type() != string(Ed25519) {
		t.Errorf("Expected Ed25519 algorithm, got %s", algorithm.Type())
	}

	// Test ECDSA
	algorithm, err = algorithmFactory(ECDSA)
	if err != nil {
		t.Fatalf("Failed to create ECDSA algorithm: %v", err)
	}

	if algorithm.Type() != string(ECDSA) {
		t.Errorf("Expected ECDSA algorithm, got %s", algorithm.Type())
	}

	// Test Secp256k1
	algorithm, err = algorithmFactory(Secp256k1)
	if err != nil {
		t.Fatalf("Failed to create Secp256k1 algorithm: %v", err)
	}

	if algorithm.Type() != string(Secp256k1) {
		t.Errorf("Expected Secp256k1 algorithm, got %s", algorithm.Type())
	}

	// Test Schnorr
	algorithm, err = algorithmFactory(Schnorr)
	if err != nil {
		t.Fatalf("Failed to create Schnorr algorithm: %v", err)
	}

	if algorithm.Type() != string(Schnorr) {
		t.Errorf("Expected Schnorr algorithm, got %s", algorithm.Type())
	}

	// Test unsupported algorithm
	_, err = algorithmFactory("unsupported")
	if err == nil {
		t.Error("Expected error for unsupported algorithm, but got none")
	}
}

// mockAlgorithm is a mock implementation of the Algorithm interface for testing
type mockAlgorithm struct{}

func (m *mockAlgorithm) GenerateKeyPair() ([]byte, []byte, error) {
	return nil, nil, nil
}

func (m *mockAlgorithm) Sign(data []byte, privateKey []byte) ([]byte, error) {
	return nil, nil
}

func (m *mockAlgorithm) Verify(data []byte, signature []byte, publicKey []byte) bool {
	return false
}

func (m *mockAlgorithm) GenerateAddress(publicKey []byte) types.Address {
	return types.Address{}
}

func (m *mockAlgorithm) Type() string {
	return "mock"
}

func (m *mockAlgorithm) PrivateKeyFromBytes(data []byte) (PrivateKey, error) {
	return nil, nil
}

func (m *mockAlgorithm) PublicKeyFromBytes(data []byte) (PublicKey, error) {
	return nil, nil
}

func TestAlgorithmRegistration(t *testing.T) {
	// Test registering a new algorithm
	RegisterAlgorithm(&mockAlgorithm{})

	// Test that the algorithm can be created
	algorithm, err := algorithmFactory("mock")
	if err != nil {
		t.Fatalf("Failed to create mock algorithm: %v", err)
	}

	if algorithm.Type() != "mock" {
		t.Errorf("Expected mock algorithm, got %s", algorithm.Type())
	}
}

func TestAlgorithmPrivateKeyFromBytes(t *testing.T) {
	// Test Ed25519
	algorithm := &ed25519Algorithm{}
	priv, _, err := algorithm.GenerateKeyPair()
	if err != nil {
		t.Fatalf("Failed to generate Ed25519 key pair: %v", err)
	}

	reconstructedPriv, err := algorithm.PrivateKeyFromBytes(priv)
	if err != nil {
		t.Fatalf("Failed to reconstruct Ed25519 private key from bytes: %v", err)
	}

	if reconstructedPriv.Type() != Ed25519 {
		t.Errorf("Expected reconstructed key type %s, got %s", Ed25519, reconstructedPriv.Type())
	}

	// Test ECDSA
	ecdsaAlgorithm := &ecdsaAlgorithm{}
	ecPriv, _, err := ecdsaAlgorithm.GenerateKeyPair()
	if err != nil {
		t.Fatalf("Failed to generate ECDSA key pair: %v", err)
	}

	ecReconstructedPriv, err := ecdsaAlgorithm.PrivateKeyFromBytes(ecPriv)
	if err != nil {
		t.Fatalf("Failed to reconstruct ECDSA private key from bytes: %v", err)
	}

	if ecReconstructedPriv.Type() != ECDSA {
		t.Errorf("Expected reconstructed key type %s, got %s", ECDSA, ecReconstructedPriv.Type())
	}

	// Test Secp256k1
	secpAlgorithm := &secp256k1Algorithm{}
	secpPriv, _, err := secpAlgorithm.GenerateKeyPair()
	if err != nil {
		t.Fatalf("Failed to generate Secp256k1 key pair: %v", err)
	}

	secpReconstructedPriv, err := secpAlgorithm.PrivateKeyFromBytes(secpPriv)
	if err != nil {
		t.Fatalf("Failed to reconstruct Secp256k1 private key from bytes: %v", err)
	}

	if secpReconstructedPriv.Type() != Secp256k1 {
		t.Errorf("Expected reconstructed key type %s, got %s", Secp256k1, secpReconstructedPriv.Type())
	}

	// Test Schnorr
	schnorrAlgorithm := &schnorrAlgorithm{}
	schnorrPriv, _, err := schnorrAlgorithm.GenerateKeyPair()
	if err != nil {
		t.Fatalf("Failed to generate Schnorr key pair: %v", err)
	}

	schnorrReconstructedPriv, err := schnorrAlgorithm.PrivateKeyFromBytes(schnorrPriv)
	if err != nil {
		t.Fatalf("Failed to reconstruct Schnorr private key from bytes: %v", err)
	}

	if schnorrReconstructedPriv.Type() != Schnorr {
		t.Errorf("Expected reconstructed key type %s, got %s", Schnorr, schnorrReconstructedPriv.Type())
	}
}

func TestAlgorithmPublicKeyFromBytes(t *testing.T) {
	// Test Ed25519
	algorithm := &ed25519Algorithm{}
	_, pub, err := algorithm.GenerateKeyPair()
	if err != nil {
		t.Fatalf("Failed to generate Ed25519 key pair: %v", err)
	}

	reconstructedPub, err := algorithm.PublicKeyFromBytes(pub)
	if err != nil {
		t.Fatalf("Failed to reconstruct Ed25519 public key from bytes: %v", err)
	}

	if reconstructedPub.Type() != Ed25519 {
		t.Errorf("Expected reconstructed key type %s, got %s", Ed25519, reconstructedPub.Type())
	}

	// Test ECDSA
	ecdsaAlgorithm := &ecdsaAlgorithm{}
	_, ecPub, err := ecdsaAlgorithm.GenerateKeyPair()
	if err != nil {
		t.Fatalf("Failed to generate ECDSA key pair: %v", err)
	}

	ecReconstructedPub, err := ecdsaAlgorithm.PublicKeyFromBytes(ecPub)
	if err != nil {
		t.Fatalf("Failed to reconstruct ECDSA public key from bytes: %v", err)
	}

	if ecReconstructedPub.Type() != ECDSA {
		t.Errorf("Expected reconstructed key type %s, got %s", ECDSA, ecReconstructedPub.Type())
	}

	// Test Secp256k1
	secpAlgorithm := &secp256k1Algorithm{}
	_, secpPub, err := secpAlgorithm.GenerateKeyPair()
	if err != nil {
		t.Fatalf("Failed to generate Secp256k1 key pair: %v", err)
	}

	secpReconstructedPub, err := secpAlgorithm.PublicKeyFromBytes(secpPub)
	if err != nil {
		t.Fatalf("Failed to reconstruct Secp256k1 public key from bytes: %v", err)
	}

	if secpReconstructedPub.Type() != Secp256k1 {
		t.Errorf("Expected reconstructed key type %s, got %s", Secp256k1, secpReconstructedPub.Type())
	}

	// Test Schnorr
	schnorrAlgorithm := &schnorrAlgorithm{}
	_, schnorrPub, err := schnorrAlgorithm.GenerateKeyPair()
	if err != nil {
		t.Fatalf("Failed to generate Schnorr key pair: %v", err)
	}

	schnorrReconstructedPub, err := schnorrAlgorithm.PublicKeyFromBytes(schnorrPub)
	if err != nil {
		t.Fatalf("Failed to reconstruct Schnorr public key from bytes: %v", err)
	}

	if schnorrReconstructedPub.Type() != Schnorr {
		t.Errorf("Expected reconstructed key type %s, got %s", Schnorr, schnorrReconstructedPub.Type())
	}
}
