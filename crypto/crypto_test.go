package crypto

import (
	"os"
	"testing"
)

func TestEd25519KeyGeneration(t *testing.T) {
	crypto := NewCrypto()

	priv, pub, err := crypto.GenerateKeyPair(Ed25519)
	if err != nil {
		t.Fatalf("Failed to generate Ed25519 key pair: %v", err)
	}

	if priv == nil {
		t.Error("Private key is nil")
	}

	if pub == nil {
		t.Error("Public key is nil")
	}

	if priv.Type() != Ed25519 {
		t.Errorf("Expected private key type Ed25519, got %s", priv.Type())
	}

	if pub.Type() != Ed25519 {
		t.Errorf("Expected public key type Ed25519, got %s", pub.Type())
	}
}

func TestECDSAKeyGeneration(t *testing.T) {
	crypto := NewCrypto()

	priv, pub, err := crypto.GenerateKeyPair(ECDSA)
	if err != nil {
		t.Fatalf("Failed to generate ECDSA key pair: %v", err)
	}

	if priv == nil {
		t.Error("Private key is nil")
	}

	if pub == nil {
		t.Error("Public key is nil")
	}

	if priv.Type() != ECDSA {
		t.Errorf("Expected private key type ECDSA, got %s", priv.Type())
	}

	if pub.Type() != ECDSA {
		t.Errorf("Expected public key type ECDSA, got %s", pub.Type())
	}
}

func TestEd25519SigningAndVerification(t *testing.T) {
	crypto := NewCrypto()

	priv, pub, err := crypto.GenerateKeyPair(Ed25519)
	if err != nil {
		t.Fatalf("Failed to generate Ed25519 key pair: %v", err)
	}

	data := []byte("test message")
	signature, err := crypto.Sign(data, priv.Bytes(), Ed25519)
	if err != nil {
		t.Fatalf("Failed to sign data: %v", err)
	}

	if signature == nil {
		t.Error("Signature is nil")
	}

	valid := crypto.Verify(data, signature, pub.Bytes(), Ed25519)
	if !valid {
		t.Error("Failed to verify signature")
	}

	// Test with wrong data
	wrongData := []byte("wrong message")
	valid = crypto.Verify(wrongData, signature, pub.Bytes(), Ed25519)
	if valid {
		t.Error("Verification should fail with wrong data")
	}
}

func TestECDSASigningAndVerification(t *testing.T) {
	crypto := NewCrypto()

	priv, pub, err := crypto.GenerateKeyPair(ECDSA)
	if err != nil {
		t.Fatalf("Failed to generate ECDSA key pair: %v", err)
	}

	data := []byte("test message")
	signature, err := crypto.Sign(data, priv.Bytes(), ECDSA)
	if err != nil {
		t.Fatalf("Failed to sign data: %v", err)
	}

	if signature == nil {
		t.Error("Signature is nil")
	}

	if len(signature) != 64 {
		t.Errorf("Expected signature length 64, got %d", len(signature))
	}

	valid := crypto.Verify(data, signature, pub.Bytes(), ECDSA)
	if !valid {
		t.Error("Failed to verify signature")
	}

	// Test with wrong data
	wrongData := []byte("wrong message")
	valid = crypto.Verify(wrongData, signature, pub.Bytes(), ECDSA)
	if valid {
		t.Error("Verification should fail with wrong data")
	}

	// Test with wrong signature
	wrongSignature := make([]byte, 64)
	valid = crypto.Verify(data, wrongSignature, pub.Bytes(), ECDSA)
	if valid {
		t.Error("Verification should fail with wrong signature")
	}
}

func TestHash(t *testing.T) {
	crypto := NewCrypto()

	data := []byte("test data")
	hash := crypto.Hash(data)

	if hash == [32]byte{} {
		t.Error("Hash is empty")
	}
}

func TestAddressGeneration(t *testing.T) {
	crypto := NewCrypto()

	// Test Ed25519 address generation
	_, edPub, err := crypto.GenerateKeyPair(Ed25519)
	if err != nil {
		t.Fatalf("Failed to generate Ed25519 key pair: %v", err)
	}

	edAddress := crypto.GenerateAddress(edPub.Bytes(), Ed25519)
	if edAddress == [20]byte{} {
		t.Error("Ed25519 address is empty")
	}

	// Verify that the address can be derived from the public key directly
	derivedAddress := edPub.Address()
	if edAddress != derivedAddress {
		t.Error("Address mismatch between crypto.GenerateAddress and pub.Address()")
	}

	// Test ECDSA address generation
	_, ecPub, err := crypto.GenerateKeyPair(ECDSA)
	if err != nil {
		t.Fatalf("Failed to generate ECDSA key pair: %v", err)
	}

	ecAddress := crypto.GenerateAddress(ecPub.Bytes(), ECDSA)
	if ecAddress == [20]byte{} {
		t.Error("ECDSA address is empty")
	}

	// Verify that the address can be derived from the public key directly
	derivedECAddress := ecPub.Address()
	if ecAddress != derivedECAddress {
		t.Error("Address mismatch between crypto.GenerateAddress and pub.Address()")
	}

	// Verify that the two addresses are different
	if edAddress == ecAddress {
		t.Error("Ed25519 and ECDSA addresses should be different")
	}
}

func TestPrivateKeyBytes(t *testing.T) {
	crypto := NewCrypto()

	// Test Ed25519 private key bytes
	edPriv, _, err := crypto.GenerateKeyPair(Ed25519)
	if err != nil {
		t.Fatalf("Failed to generate Ed25519 key pair: %v", err)
	}

	edBytes := edPriv.Bytes()
	if len(edBytes) == 0 {
		t.Error("Ed25519 private key bytes are empty")
	}

	// Test ECDSA private key bytes
	ecPriv, _, err := crypto.GenerateKeyPair(ECDSA)
	if err != nil {
		t.Fatalf("Failed to generate ECDSA key pair: %v", err)
	}

	ecBytes := ecPriv.Bytes()
	if len(ecBytes) == 0 {
		t.Error("ECDSA private key bytes are empty")
	}
}

func TestPublicKeyBytes(t *testing.T) {
	crypto := NewCrypto()

	// Test Ed25519 public key bytes
	_, edPub, err := crypto.GenerateKeyPair(Ed25519)
	if err != nil {
		t.Fatalf("Failed to generate Ed25519 key pair: %v", err)
	}

	edBytes := edPub.Bytes()
	if len(edBytes) == 0 {
		t.Error("Ed25519 public key bytes are empty")
	}

	// Test ECDSA public key bytes
	_, ecPub, err := crypto.GenerateKeyPair(ECDSA)
	if err != nil {
		t.Fatalf("Failed to generate ECDSA key pair: %v", err)
	}

	ecBytes := ecPub.Bytes()
	if len(ecBytes) == 0 {
		t.Error("ECDSA public key bytes are empty")
	}

	// ECDSA public key should be 33 bytes (compressed)
	if len(ecBytes) != 33 {
		t.Errorf("Expected ECDSA public key length 33, got %d", len(ecBytes))
	}
}

func TestPublicKeyFromPrivateKey(t *testing.T) {
	crypto := NewCrypto()

	// Test Ed25519
	edPriv, edPub, err := crypto.GenerateKeyPair(Ed25519)
	if err != nil {
		t.Fatalf("Failed to generate Ed25519 key pair: %v", err)
	}

	derivedPub := edPriv.PublicKey()
	if derivedPub.Type() != edPub.Type() {
		t.Error("Derived public key type mismatch")
	}

	// Test ECDSA
	ecPriv, ecPub, err := crypto.GenerateKeyPair(ECDSA)
	if err != nil {
		t.Fatalf("Failed to generate ECDSA key pair: %v", err)
	}

	derivedECPub := ecPriv.PublicKey()
	if derivedECPub.Type() != ecPub.Type() {
		t.Error("Derived public key type mismatch")
	}
}

// Additional tests for edge cases and robustness

func TestEmptyDataSigning(t *testing.T) {
	crypto := NewCrypto()

	// Test Ed25519 with empty data
	edPriv, edPub, err := crypto.GenerateKeyPair(Ed25519)
	if err != nil {
		t.Fatalf("Failed to generate Ed25519 key pair: %v", err)
	}

	emptyData := []byte{}
	edSignature, err := crypto.Sign(emptyData, edPriv.Bytes(), Ed25519)
	if err != nil {
		t.Fatalf("Failed to sign empty data with Ed25519: %v", err)
	}

	if !crypto.Verify(emptyData, edSignature, edPub.Bytes(), Ed25519) {
		t.Error("Failed to verify signature for empty data with Ed25519")
	}

	// Test ECDSA with empty data
	ecPriv, ecPub, err := crypto.GenerateKeyPair(ECDSA)
	if err != nil {
		t.Fatalf("Failed to generate ECDSA key pair: %v", err)
	}

	ecSignature, err := crypto.Sign(emptyData, ecPriv.Bytes(), ECDSA)
	if err != nil {
		t.Fatalf("Failed to sign empty data with ECDSA: %v", err)
	}

	if !crypto.Verify(emptyData, ecSignature, ecPub.Bytes(), ECDSA) {
		t.Error("Failed to verify signature for empty data with ECDSA")
	}
}

func TestLargeDataSigning(t *testing.T) {
	crypto := NewCrypto()

	// Create large data
	largeData := make([]byte, 1024*1024) // 1MB of data
	for i := range largeData {
		largeData[i] = byte(i % 256)
	}

	// Test Ed25519 with large data
	edPriv, edPub, err := crypto.GenerateKeyPair(Ed25519)
	if err != nil {
		t.Fatalf("Failed to generate Ed25519 key pair: %v", err)
	}

	edSignature, err := crypto.Sign(largeData, edPriv.Bytes(), Ed25519)
	if err != nil {
		t.Fatalf("Failed to sign large data with Ed25519: %v", err)
	}

	if !crypto.Verify(largeData, edSignature, edPub.Bytes(), Ed25519) {
		t.Error("Failed to verify signature for large data with Ed25519")
	}

	// Test ECDSA with large data
	ecPriv, ecPub, err := crypto.GenerateKeyPair(ECDSA)
	if err != nil {
		t.Fatalf("Failed to generate ECDSA key pair: %v", err)
	}

	ecSignature, err := crypto.Sign(largeData, ecPriv.Bytes(), ECDSA)
	if err != nil {
		t.Fatalf("Failed to sign large data with ECDSA: %v", err)
	}

	if !crypto.Verify(largeData, ecSignature, ecPub.Bytes(), ECDSA) {
		t.Error("Failed to verify signature for large data with ECDSA")
	}
}

func TestDeterministicSignatures(t *testing.T) {
	crypto := NewCrypto()

	data := []byte("test data")

	// Test Ed25519 deterministic signatures
	edPriv, edPub, err := crypto.GenerateKeyPair(Ed25519)
	if err != nil {
		t.Fatalf("Failed to generate Ed25519 key pair: %v", err)
	}

	edSignature1, err := crypto.Sign(data, edPriv.Bytes(), Ed25519)
	if err != nil {
		t.Fatalf("Failed to sign data with Ed25519: %v", err)
	}

	// For Ed25519, signatures are deterministic for the same key and data
	edSignature2, err := crypto.Sign(data, edPriv.Bytes(), Ed25519)
	if err != nil {
		t.Fatalf("Failed to sign data with Ed25519: %v", err)
	}

	// Signatures should be identical
	if string(edSignature1) != string(edSignature2) {
		t.Error("Ed25519 signatures should be deterministic")
	}

	// Verify both signatures
	if !crypto.Verify(data, edSignature1, edPub.Bytes(), Ed25519) {
		t.Error("Failed to verify first Ed25519 signature")
	}

	if !crypto.Verify(data, edSignature2, edPub.Bytes(), Ed25519) {
		t.Error("Failed to verify second Ed25519 signature")
	}
}

// Test SaveToFile and LoadFromFile functionality
func TestSaveAndLoadKeys(t *testing.T) {
	crypto := NewCrypto()

	// Test with Ed25519
	edPriv, _, err := crypto.GenerateKeyPair(Ed25519)
	if err != nil {
		t.Fatalf("Failed to generate Ed25519 key pair: %v", err)
	}

	// Save the key with default password (empty string)
	edFilename := "test_ed25519_key.json"
	defer os.Remove(edFilename)
	err = SaveToFile(edPriv, edFilename, "")
	if err != nil {
		t.Fatalf("Failed to save Ed25519 key: %v", err)
	}

	// Load the key with default password (empty string)
	loadedEdPriv, err := LoadFromFile(edFilename, "")
	if err != nil {
		t.Fatalf("Failed to load Ed25519 key: %v", err)
	}

	// Verify the loaded key
	if loadedEdPriv.Type() != edPriv.Type() {
		t.Errorf("Expected loaded key type %s, got %s", edPriv.Type(), loadedEdPriv.Type())
	}

	// Verify that the keys have the same bytes
	if string(loadedEdPriv.Bytes()) != string(edPriv.Bytes()) {
		t.Error("Loaded Ed25519 key bytes do not match original")
	}

	// Test with ECDSA
	ecPriv, _, err := crypto.GenerateKeyPair(ECDSA)
	if err != nil {
		t.Fatalf("Failed to generate ECDSA key pair: %v", err)
	}

	// Save the key with default password (empty string)
	ecFilename := "test_ecdsa_key.json"
	defer os.Remove(ecFilename)
	err = SaveToFile(ecPriv, ecFilename, "")
	if err != nil {
		t.Fatalf("Failed to save ECDSA key: %v", err)
	}

	// Load the key with default password (empty string)
	loadedEcPriv, err := LoadFromFile(ecFilename, "")
	if err != nil {
		t.Fatalf("Failed to load ECDSA key: %v", err)
	}

	// Verify the loaded key
	if loadedEcPriv.Type() != ecPriv.Type() {
		t.Errorf("Expected loaded key type %s, got %s", ecPriv.Type(), loadedEcPriv.Type())
	}

	// Verify that the keys have the same bytes
	if string(loadedEcPriv.Bytes()) != string(ecPriv.Bytes()) {
		t.Error("Loaded ECDSA key bytes do not match original")
	}
}

// Test SaveToFile and LoadFromFile functionality with custom password
func TestSaveAndLoadKeysWithPassword(t *testing.T) {
	crypto := NewCrypto()
	password := "test-password-123"

	// Test with Ed25519
	edPriv, _, err := crypto.GenerateKeyPair(Ed25519)
	if err != nil {
		t.Fatalf("Failed to generate Ed25519 key pair: %v", err)
	}

	// Save the key with password
	edFilename := "test_ed25519_key_encrypted.json"
	defer os.Remove(edFilename)
	err = SaveToFile(edPriv, edFilename, password)
	if err != nil {
		t.Fatalf("Failed to save Ed25519 key with password: %v", err)
	}

	// Load the key with password
	loadedEdPriv, err := LoadFromFile(edFilename, password)
	if err != nil {
		t.Fatalf("Failed to load Ed25519 key with password: %v", err)
	}

	// Verify the loaded key
	if loadedEdPriv.Type() != edPriv.Type() {
		t.Errorf("Expected loaded key type %s, got %s", edPriv.Type(), loadedEdPriv.Type())
	}

	// Verify that the keys have the same bytes
	if string(loadedEdPriv.Bytes()) != string(edPriv.Bytes()) {
		t.Error("Loaded Ed25519 key bytes do not match original")
	}

	// Test with wrong password
	_, err = LoadFromFile(edFilename, "wrong-password")
	if err == nil {
		t.Error("Expected error when loading with wrong password, but got none")
	}

	// Test with ECDSA
	ecPriv, _, err := crypto.GenerateKeyPair(ECDSA)
	if err != nil {
		t.Fatalf("Failed to generate ECDSA key pair: %v", err)
	}

	// Save the key with password
	ecFilename := "test_ecdsa_key_encrypted.json"
	defer os.Remove(ecFilename)
	err = SaveToFile(ecPriv, ecFilename, password)
	if err != nil {
		t.Fatalf("Failed to save ECDSA key with password: %v", err)
	}

	// Load the key with password
	loadedEcPriv, err := LoadFromFile(ecFilename, password)
	if err != nil {
		t.Fatalf("Failed to load ECDSA key with password: %v", err)
	}

	// Verify the loaded key
	if loadedEcPriv.Type() != ecPriv.Type() {
		t.Errorf("Expected loaded key type %s, got %s", ecPriv.Type(), loadedEcPriv.Type())
	}

	// Verify that the keys have the same bytes
	if string(loadedEcPriv.Bytes()) != string(ecPriv.Bytes()) {
		t.Error("Loaded ECDSA key bytes do not match original")
	}

	// Test with wrong password
	_, err = LoadFromFile(ecFilename, "wrong-password")
	if err == nil {
		t.Error("Expected error when loading with wrong password, but got none")
	}
}

// Test backward compatibility - loading with default password
func TestBackwardCompatibility(t *testing.T) {
	crypto := NewCrypto()

	// Create a key file using the new method with default password
	edPriv, _, err := crypto.GenerateKeyPair(Ed25519)
	if err != nil {
		t.Fatalf("Failed to generate Ed25519 key pair: %v", err)
	}

	// Save the key using the new method with empty password (will use default)
	edFilename := "test_ed25519_key_default.json"
	defer os.Remove(edFilename)
	err = SaveToFile(edPriv, edFilename, "")
	if err != nil {
		t.Fatalf("Failed to save Ed25519 key: %v", err)
	}

	// Load the key using the new method with empty password (will use default)
	loadedEdPriv, err := LoadFromFile(edFilename, "")
	if err != nil {
		t.Fatalf("Failed to load Ed25519 key: %v", err)
	}

	// Verify the loaded key
	if loadedEdPriv.Type() != edPriv.Type() {
		t.Errorf("Expected loaded key type %s, got %s", edPriv.Type(), loadedEdPriv.Type())
	}

	// Verify that the keys have the same bytes
	if string(loadedEdPriv.Bytes()) != string(edPriv.Bytes()) {
		t.Error("Loaded Ed25519 key bytes do not match original")
	}
}

// Test FromBytes methods
func TestPrivateKeyFromBytes(t *testing.T) {
	crypto := NewCrypto()

	// Test Ed25519 private key FromBytes
	edPriv, _, err := crypto.GenerateKeyPair(Ed25519)
	if err != nil {
		t.Fatalf("Failed to generate Ed25519 key pair: %v", err)
	}

	edBytes := edPriv.Bytes()
	reconstructedEdPriv, err := edPriv.FromBytes(edBytes)
	if err != nil {
		t.Fatalf("Failed to reconstruct Ed25519 private key from bytes: %v", err)
	}

	if reconstructedEdPriv.Type() != edPriv.Type() {
		t.Errorf("Expected reconstructed key type %s, got %s", edPriv.Type(), reconstructedEdPriv.Type())
	}

	if string(reconstructedEdPriv.Bytes()) != string(edPriv.Bytes()) {
		t.Error("Reconstructed Ed25519 private key bytes do not match original")
	}

	// Test ECDSA private key FromBytes
	ecPriv, _, err := crypto.GenerateKeyPair(ECDSA)
	if err != nil {
		t.Fatalf("Failed to generate ECDSA key pair: %v", err)
	}

	ecBytes := ecPriv.Bytes()
	reconstructedEcPriv, err := ecPriv.FromBytes(ecBytes)
	if err != nil {
		t.Fatalf("Failed to reconstruct ECDSA private key from bytes: %v", err)
	}

	if reconstructedEcPriv.Type() != ecPriv.Type() {
		t.Errorf("Expected reconstructed key type %s, got %s", ecPriv.Type(), reconstructedEcPriv.Type())
	}

	if string(reconstructedEcPriv.Bytes()) != string(ecPriv.Bytes()) {
		t.Error("Reconstructed ECDSA private key bytes do not match original")
	}
}

func TestPublicKeyFromBytes(t *testing.T) {
	crypto := NewCrypto()

	// Test Ed25519 public key FromBytes
	_, edPub, err := crypto.GenerateKeyPair(Ed25519)
	if err != nil {
		t.Fatalf("Failed to generate Ed25519 key pair: %v", err)
	}

	edBytes := edPub.Bytes()
	reconstructedEdPub, err := edPub.FromBytes(edBytes)
	if err != nil {
		t.Fatalf("Failed to reconstruct Ed25519 public key from bytes: %v", err)
	}

	if reconstructedEdPub.Type() != edPub.Type() {
		t.Errorf("Expected reconstructed key type %s, got %s", edPub.Type(), reconstructedEdPub.Type())
	}

	if string(reconstructedEdPub.Bytes()) != string(edPub.Bytes()) {
		t.Error("Reconstructed Ed25519 public key bytes do not match original")
	}

	// Verify that both keys produce the same address
	originalAddr := edPub.Address()
	reconstructedAddr := reconstructedEdPub.Address()
	if originalAddr != reconstructedAddr {
		t.Error("Reconstructed Ed25519 public key produces different address")
	}

	// Test ECDSA public key FromBytes
	_, ecPub, err := crypto.GenerateKeyPair(ECDSA)
	if err != nil {
		t.Fatalf("Failed to generate ECDSA key pair: %v", err)
	}

	ecBytes := ecPub.Bytes()
	reconstructedEcPub, err := ecPub.FromBytes(ecBytes)
	if err != nil {
		t.Fatalf("Failed to reconstruct ECDSA public key from bytes: %v", err)
	}

	if reconstructedEcPub.Type() != ecPub.Type() {
		t.Errorf("Expected reconstructed key type %s, got %s", ecPub.Type(), reconstructedEcPub.Type())
	}

	if string(reconstructedEcPub.Bytes()) != string(ecPub.Bytes()) {
		t.Error("Reconstructed ECDSA public key bytes do not match original")
	}

	// Verify that both keys produce the same address
	originalECAddr := ecPub.Address()
	reconstructedECAddr := reconstructedEcPub.Address()
	if originalECAddr != reconstructedECAddr {
		t.Error("Reconstructed ECDSA public key produces different address")
	}
}
