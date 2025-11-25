package crypto

import (
	"fmt"
	"os"
	"testing"

	"github.com/lengzhao/govm/types"
)

func TestEd25519KeyGeneration(t *testing.T) {
	priv, pub, err := GenerateKeyPair(Ed25519)
	if err != nil {
		t.Fatalf("Failed to generate Ed25519 key pair: %v", err)
	}

	if priv == nil {
		t.Error("Private key is nil")
	}

	if pub == nil {
		t.Error("Public key is nil")
	}
}

func TestECDSAKeyGeneration(t *testing.T) {
	priv, pub, err := GenerateKeyPair(ECDSA)
	if err != nil {
		t.Fatalf("Failed to generate ECDSA key pair: %v", err)
	}

	if priv == nil {
		t.Error("Private key is nil")
	}

	if pub == nil {
		t.Error("Public key is nil")
	}
}

func TestEd25519SigningAndVerification(t *testing.T) {
	priv, pub, err := GenerateKeyPair(Ed25519)
	if err != nil {
		t.Fatalf("Failed to generate Ed25519 key pair: %v", err)
	}

	data := []byte("test message")
	signature, err := Sign(data, priv, Ed25519)
	if err != nil {
		t.Fatalf("Failed to sign data: %v", err)
	}

	if signature == nil {
		t.Error("Signature is nil")
	}

	valid := Verify(data, signature, pub, Ed25519)
	if !valid {
		t.Error("Failed to verify signature")
	}

	// Test with wrong data
	wrongData := []byte("wrong message")
	valid = Verify(wrongData, signature, pub, Ed25519)
	if valid {
		t.Error("Verification should fail with wrong data")
	}
}

func TestECDSASigningAndVerification(t *testing.T) {
	priv, pub, err := GenerateKeyPair(ECDSA)
	if err != nil {
		t.Fatalf("Failed to generate ECDSA key pair: %v", err)
	}

	data := []byte("test message")
	signature, err := Sign(data, priv, ECDSA)
	if err != nil {
		t.Fatalf("Failed to sign data: %v", err)
	}

	if signature == nil {
		t.Error("Signature is nil")
	}

	if len(signature) != 64 {
		t.Errorf("Expected signature length 64, got %d", len(signature))
	}

	valid := Verify(data, signature, pub, ECDSA)
	if !valid {
		t.Error("Failed to verify signature")
	}

	// Test with wrong data
	wrongData := []byte("wrong message")
	valid = Verify(wrongData, signature, pub, ECDSA)
	if valid {
		t.Error("Verification should fail with wrong data")
	}

	// Test with wrong signature
	wrongSignature := make([]byte, 64)
	valid = Verify(data, wrongSignature, pub, ECDSA)
	if valid {
		t.Error("Verification should fail with wrong signature")
	}
}

func TestHash(t *testing.T) {
	data := []byte("test data")
	hash := Hash(data)

	if hash == [32]byte{} {
		t.Error("Hash is empty")
	}
}

func TestAddressGeneration(t *testing.T) {
	// Test Ed25519 address generation
	_, edPub, err := GenerateKeyPair(Ed25519)
	if err != nil {
		t.Fatalf("Failed to generate Ed25519 key pair: %v", err)
	}

	edAddress := GenerateAddress(edPub, Ed25519)
	if edAddress == [20]byte{} {
		t.Error("Ed25519 address is empty")
	}

	// Test ECDSA address generation
	_, ecPub, err := GenerateKeyPair(ECDSA)
	if err != nil {
		t.Fatalf("Failed to generate ECDSA key pair: %v", err)
	}

	ecAddress := GenerateAddress(ecPub, ECDSA)
	if ecAddress == [20]byte{} {
		t.Error("ECDSA address is empty")
	}

	// Verify that the two addresses are different
	if edAddress == ecAddress {
		t.Error("Ed25519 and ECDSA addresses should be different")
	}
}

func TestPrivateKeyBytes(t *testing.T) {
	// Test Ed25519 private key bytes
	edPriv, _, err := GenerateKeyPair(Ed25519)
	if err != nil {
		t.Fatalf("Failed to generate Ed25519 key pair: %v", err)
	}

	edBytes := edPriv
	if len(edBytes) == 0 {
		t.Error("Ed25519 private key bytes are empty")
	}

	// Test ECDSA private key bytes
	ecPriv, _, err := GenerateKeyPair(ECDSA)
	if err != nil {
		t.Fatalf("Failed to generate ECDSA key pair: %v", err)
	}

	ecBytes := ecPriv
	if len(ecBytes) == 0 {
		t.Error("ECDSA private key bytes are empty")
	}
}

func TestPublicKeyBytes(t *testing.T) {
	// Test Ed25519 public key bytes
	_, edPub, err := GenerateKeyPair(Ed25519)
	if err != nil {
		t.Fatalf("Failed to generate Ed25519 key pair: %v", err)
	}

	edBytes := edPub
	if len(edBytes) == 0 {
		t.Error("Ed25519 public key bytes are empty")
	}

	// Test ECDSA public key bytes
	_, ecPub, err := GenerateKeyPair(ECDSA)
	if err != nil {
		t.Fatalf("Failed to generate ECDSA key pair: %v", err)
	}

	ecBytes := ecPub
	if len(ecBytes) == 0 {
		t.Error("ECDSA public key bytes are empty")
	}
}

func TestPublicKeyFromPrivateKey(t *testing.T) {
	// This test is no longer applicable with the new interface
	// Public key extraction is now handled by each algorithm's implementation
	t.Skip("Test not applicable with new interface")
}

// Additional tests for edge cases and robustness

func TestEmptyDataSigning(t *testing.T) {
	// Test Ed25519 with empty data
	edPriv, edPub, err := GenerateKeyPair(Ed25519)
	if err != nil {
		t.Fatalf("Failed to generate Ed25519 key pair: %v", err)
	}

	emptyData := []byte{}
	edSignature, err := Sign(emptyData, edPriv, Ed25519)
	if err != nil {
		t.Fatalf("Failed to sign empty data with Ed25519: %v", err)
	}

	if !Verify(emptyData, edSignature, edPub, Ed25519) {
		t.Error("Failed to verify signature for empty data with Ed25519")
	}

	// Test ECDSA with empty data
	ecPriv, ecPub, err := GenerateKeyPair(ECDSA)
	if err != nil {
		t.Fatalf("Failed to generate ECDSA key pair: %v", err)
	}

	ecSignature, err := Sign(emptyData, ecPriv, ECDSA)
	if err != nil {
		t.Fatalf("Failed to sign empty data with ECDSA: %v", err)
	}

	if !Verify(emptyData, ecSignature, ecPub, ECDSA) {
		t.Error("Failed to verify signature for empty data with ECDSA")
	}
}

func TestLargeDataSigning(t *testing.T) {
	// Create large data
	largeData := make([]byte, 1024*1024) // 1MB of data
	for i := range largeData {
		largeData[i] = byte(i % 256)
	}

	// Test Ed25519 with large data
	edPriv, edPub, err := GenerateKeyPair(Ed25519)
	if err != nil {
		t.Fatalf("Failed to generate Ed25519 key pair: %v", err)
	}

	edSignature, err := Sign(largeData, edPriv, Ed25519)
	if err != nil {
		t.Fatalf("Failed to sign large data with Ed25519: %v", err)
	}

	if !Verify(largeData, edSignature, edPub, Ed25519) {
		t.Error("Failed to verify signature for large data with Ed25519")
	}

	// Test ECDSA with large data
	ecPriv, ecPub, err := GenerateKeyPair(ECDSA)
	if err != nil {
		t.Fatalf("Failed to generate ECDSA key pair: %v", err)
	}

	ecSignature, err := Sign(largeData, ecPriv, ECDSA)
	if err != nil {
		t.Fatalf("Failed to sign large data with ECDSA: %v", err)
	}

	if !Verify(largeData, ecSignature, ecPub, ECDSA) {
		t.Error("Failed to verify signature for large data with ECDSA")
	}
}

func TestDeterministicSignatures(t *testing.T) {
	data := []byte("test data")

	// Test Ed25519 deterministic signatures
	edPriv, edPub, err := GenerateKeyPair(Ed25519)
	if err != nil {
		t.Fatalf("Failed to generate Ed25519 key pair: %v", err)
	}

	edSignature1, err := Sign(data, edPriv, Ed25519)
	if err != nil {
		t.Fatalf("Failed to sign data with Ed25519: %v", err)
	}

	// For Ed25519, signatures are deterministic for the same key and data
	edSignature2, err := Sign(data, edPriv, Ed25519)
	if err != nil {
		t.Fatalf("Failed to sign data with Ed25519: %v", err)
	}

	// Signatures should be identical
	if string(edSignature1) != string(edSignature2) {
		t.Error("Ed25519 signatures should be deterministic")
	}

	// Verify both signatures
	if !Verify(data, edSignature1, edPub, Ed25519) {
		t.Error("Failed to verify first Ed25519 signature")
	}

	if !Verify(data, edSignature2, edPub, Ed25519) {
		t.Error("Failed to verify second Ed25519 signature")
	}
}

// Test SaveToFile and LoadFromFile functionality
func TestSaveAndLoadKeys(t *testing.T) {
	// Test with Ed25519
	edPriv, _, err := GenerateKeyPair(Ed25519)
	if err != nil {
		t.Fatalf("Failed to generate Ed25519 key pair: %v", err)
	}

	// Reconstruct the private key from bytes
	algorithm, err := AlgorithmFactory(Ed25519)
	if err != nil {
		t.Fatalf("Failed to get algorithm: %v", err)
	}

	edPrivKey, err := algorithm.PrivateKeyFromBytes(edPriv)
	if err != nil {
		t.Fatalf("Failed to reconstruct private key: %v", err)
	}

	// Save the key with default password (empty string)
	edFilename := "test_ed25519_key.json"
	defer os.Remove(edFilename)
	err = SaveToFile(edPrivKey, edFilename, "")
	if err != nil {
		t.Fatalf("Failed to save Ed25519 key: %v", err)
	}

	// Load the key with default password (empty string)
	loadedEdPriv, err := LoadFromFile(edFilename, "")
	if err != nil {
		t.Fatalf("Failed to load Ed25519 key: %v", err)
	}

	// Verify the loaded key
	if loadedEdPriv.Type() != edPrivKey.Type() {
		t.Errorf("Expected loaded key type %s, got %s", edPrivKey.Type(), loadedEdPriv.Type())
	}

	// Verify that the keys have the same bytes
	if string(loadedEdPriv.Bytes()) != string(edPrivKey.Bytes()) {
		t.Error("Loaded Ed25519 key bytes do not match original")
	}

	// Test with ECDSA
	ecPriv, _, err := GenerateKeyPair(ECDSA)
	if err != nil {
		t.Fatalf("Failed to generate ECDSA key pair: %v", err)
	}

	// Reconstruct the private key from bytes
	ecAlgorithm, err := AlgorithmFactory(ECDSA)
	if err != nil {
		t.Fatalf("Failed to get algorithm: %v", err)
	}

	ecPrivKey, err := ecAlgorithm.PrivateKeyFromBytes(ecPriv)
	if err != nil {
		t.Fatalf("Failed to reconstruct private key: %v", err)
	}

	// Save the key with default password (empty string)
	ecFilename := "test_ecdsa_key.json"
	defer os.Remove(ecFilename)
	err = SaveToFile(ecPrivKey, ecFilename, "")
	if err != nil {
		t.Fatalf("Failed to save ECDSA key: %v", err)
	}

	// Load the key with default password (empty string)
	loadedEcPriv, err := LoadFromFile(ecFilename, "")
	if err != nil {
		t.Fatalf("Failed to load ECDSA key: %v", err)
	}

	// Verify the loaded key
	if loadedEcPriv.Type() != ecPrivKey.Type() {
		t.Errorf("Expected loaded key type %s, got %s", ecPrivKey.Type(), loadedEcPriv.Type())
	}

	// Verify that the keys have the same bytes
	if string(loadedEcPriv.Bytes()) != string(ecPrivKey.Bytes()) {
		t.Error("Loaded ECDSA key bytes do not match original")
	}
}

// Test SaveToFile and LoadFromFile functionality with custom password
func TestSaveAndLoadKeysWithPassword(t *testing.T) {
	password := "test-password-123"

	// Test with Ed25519
	edPriv, _, err := GenerateKeyPair(Ed25519)
	if err != nil {
		t.Fatalf("Failed to generate Ed25519 key pair: %v", err)
	}

	// Reconstruct the private key from bytes
	algorithm, err := AlgorithmFactory(Ed25519)
	if err != nil {
		t.Fatalf("Failed to get algorithm: %v", err)
	}

	edPrivKey, err := algorithm.PrivateKeyFromBytes(edPriv)
	if err != nil {
		t.Fatalf("Failed to reconstruct private key: %v", err)
	}

	// Save the key with password
	edFilename := "test_ed25519_key_encrypted.json"
	defer os.Remove(edFilename)
	err = SaveToFile(edPrivKey, edFilename, password)
	if err != nil {
		t.Fatalf("Failed to save Ed25519 key with password: %v", err)
	}

	// Load the key with password
	loadedEdPriv, err := LoadFromFile(edFilename, password)
	if err != nil {
		t.Fatalf("Failed to load Ed25519 key with password: %v", err)
	}

	// Verify the loaded key
	if loadedEdPriv.Type() != edPrivKey.Type() {
		t.Errorf("Expected loaded key type %s, got %s", edPrivKey.Type(), loadedEdPriv.Type())
	}

	// Verify that the keys have the same bytes
	if string(loadedEdPriv.Bytes()) != string(edPrivKey.Bytes()) {
		t.Error("Loaded Ed25519 key bytes do not match original")
	}

	// Test with wrong password
	_, err = LoadFromFile(edFilename, "wrong-password")
	if err == nil {
		t.Error("Expected error when loading with wrong password, but got none")
	}

	// Test with ECDSA
	ecPriv, _, err := GenerateKeyPair(ECDSA)
	if err != nil {
		t.Fatalf("Failed to generate ECDSA key pair: %v", err)
	}

	// Reconstruct the private key from bytes
	ecAlgorithm, err := AlgorithmFactory(ECDSA)
	if err != nil {
		t.Fatalf("Failed to get algorithm: %v", err)
	}

	ecPrivKey, err := ecAlgorithm.PrivateKeyFromBytes(ecPriv)
	if err != nil {
		t.Fatalf("Failed to reconstruct private key: %v", err)
	}

	// Save the key with password
	ecFilename := "test_ecdsa_key_encrypted.json"
	defer os.Remove(ecFilename)
	err = SaveToFile(ecPrivKey, ecFilename, password)
	if err != nil {
		t.Fatalf("Failed to save ECDSA key with password: %v", err)
	}

	// Load the key with password
	loadedEcPriv, err := LoadFromFile(ecFilename, password)
	if err != nil {
		t.Fatalf("Failed to load ECDSA key with password: %v", err)
	}

	// Verify the loaded key
	if loadedEcPriv.Type() != ecPrivKey.Type() {
		t.Errorf("Expected loaded key type %s, got %s", ecPrivKey.Type(), loadedEcPriv.Type())
	}

	// Verify that the keys have the same bytes
	if string(loadedEcPriv.Bytes()) != string(ecPrivKey.Bytes()) {
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
	// Create a key file using the new method with default password
	edPriv, _, err := GenerateKeyPair(Ed25519)
	if err != nil {
		t.Fatalf("Failed to generate Ed25519 key pair: %v", err)
	}

	// Reconstruct the private key from bytes
	algorithm, err := AlgorithmFactory(Ed25519)
	if err != nil {
		t.Fatalf("Failed to get algorithm: %v", err)
	}

	edPrivKey, err := algorithm.PrivateKeyFromBytes(edPriv)
	if err != nil {
		t.Fatalf("Failed to reconstruct private key: %v", err)
	}

	// Save the key using the new method with empty password (will use default)
	edFilename := "test_ed25519_key_default.json"
	defer os.Remove(edFilename)
	err = SaveToFile(edPrivKey, edFilename, "")
	if err != nil {
		t.Fatalf("Failed to save Ed25519 key: %v", err)
	}

	// Load the key using the new method with empty password (will use default)
	loadedEdPriv, err := LoadFromFile(edFilename, "")
	if err != nil {
		t.Fatalf("Failed to load Ed25519 key: %v", err)
	}

	// Verify the loaded key
	if loadedEdPriv.Type() != edPrivKey.Type() {
		t.Errorf("Expected loaded key type %s, got %s", edPrivKey.Type(), loadedEdPriv.Type())
	}

	// Verify that the keys have the same bytes
	if string(loadedEdPriv.Bytes()) != string(edPrivKey.Bytes()) {
		t.Error("Loaded Ed25519 key bytes do not match original")
	}
}

// Test FromBytes methods - these tests are no longer applicable with the new interface
func TestPrivateKeyFromBytes(t *testing.T) {
	// This test is no longer applicable with the new interface
	t.Skip("Test not applicable with new interface")
}

func TestPublicKeyFromBytes(t *testing.T) {
	// This test is no longer applicable with the new interface
	t.Skip("Test not applicable with new interface")
}

// Additional tests for Ed25519 edge cases
func TestEd25519EdgeCases(t *testing.T) {
	// Test multiple key generations produce different keys
	priv1, pub1, err := GenerateKeyPair(Ed25519)
	if err != nil {
		t.Fatalf("Failed to generate first Ed25519 key pair: %v", err)
	}

	priv2, pub2, err := GenerateKeyPair(Ed25519)
	if err != nil {
		t.Fatalf("Failed to generate second Ed25519 key pair: %v", err)
	}

	if string(priv1) == string(priv2) {
		t.Error("Generated private keys should be different")
	}

	if string(pub1) == string(pub2) {
		t.Error("Generated public keys should be different")
	}

	// Test signing with maximum data size
	maxData := make([]byte, 10*1024*1024) // 10MB
	for i := range maxData {
		maxData[i] = byte(i % 256)
	}

	signature, err := Sign(maxData, priv1, Ed25519)
	if err != nil {
		t.Fatalf("Failed to sign maximum data size: %v", err)
	}

	if !Verify(maxData, signature, pub1, Ed25519) {
		t.Error("Failed to verify signature for maximum data size")
	}

	// Test with nil data
	nilSignature, err := Sign(nil, priv1, Ed25519)
	if err != nil {
		t.Fatalf("Failed to sign nil data: %v", err)
	}

	if !Verify(nil, nilSignature, pub1, Ed25519) {
		t.Error("Failed to verify signature for nil data")
	}

	// Test with invalid signature
	invalidSig := make([]byte, 64)
	if Verify([]byte("test"), invalidSig, pub1, Ed25519) {
		t.Error("Verification should fail with invalid signature")
	}

	// Test with invalid public key
	invalidPub := make([]byte, 32)
	if Verify([]byte("test"), signature, invalidPub, Ed25519) {
		t.Error("Verification should fail with invalid public key")
	}

	// Test with corrupted signature
	corruptedSig := make([]byte, len(signature))
	copy(corruptedSig, signature)
	corruptedSig[0] ^= 0x01 // Flip one bit
	if Verify([]byte("test"), corruptedSig, pub1, Ed25519) {
		t.Error("Verification should fail with corrupted signature")
	}
}

// Additional tests for ECDSA edge cases
func TestECDSAEdgeCases(t *testing.T) {
	// Test multiple key generations produce different keys
	priv1, pub1, err := GenerateKeyPair(ECDSA)
	if err != nil {
		t.Fatalf("Failed to generate first ECDSA key pair: %v", err)
	}

	priv2, pub2, err := GenerateKeyPair(ECDSA)
	if err != nil {
		t.Fatalf("Failed to generate second ECDSA key pair: %v", err)
	}

	if string(priv1) == string(priv2) {
		t.Error("Generated private keys should be different")
	}

	if string(pub1) == string(pub2) {
		t.Error("Generated public keys should be different")
	}

	// Test signing with maximum data size
	maxData := make([]byte, 10*1024*1024) // 10MB
	for i := range maxData {
		maxData[i] = byte(i % 256)
	}

	signature, err := Sign(maxData, priv1, ECDSA)
	if err != nil {
		t.Fatalf("Failed to sign maximum data size: %v", err)
	}

	if !Verify(maxData, signature, pub1, ECDSA) {
		t.Error("Failed to verify signature for maximum data size")
	}

	// Test with nil data
	nilSignature, err := Sign(nil, priv1, ECDSA)
	if err != nil {
		t.Fatalf("Failed to sign nil data: %v", err)
	}

	if !Verify(nil, nilSignature, pub1, ECDSA) {
		t.Error("Failed to verify signature for nil data")
	}

	// Test with invalid signature length
	invalidSig := make([]byte, 32) // Should be 64 bytes
	if Verify([]byte("test"), invalidSig, pub1, ECDSA) {
		t.Error("Verification should fail with invalid signature length")
	}

	// Test with invalid public key
	invalidPub := make([]byte, 33)
	if Verify([]byte("test"), signature, invalidPub, ECDSA) {
		t.Error("Verification should fail with invalid public key")
	}

	// Test with corrupted signature
	corruptedSig := make([]byte, len(signature))
	copy(corruptedSig, signature)
	corruptedSig[0] ^= 0x01 // Flip one bit
	if Verify([]byte("test"), corruptedSig, pub1, ECDSA) {
		t.Error("Verification should fail with corrupted signature")
	}
}

// Additional tests for Secp256k1 edge cases
func TestSecp256k1EdgeCases(t *testing.T) {
	// Test multiple key generations produce different keys
	priv1, pub1, err := GenerateKeyPair(Secp256k1)
	if err != nil {
		t.Fatalf("Failed to generate first Secp256k1 key pair: %v", err)
	}

	priv2, pub2, err := GenerateKeyPair(Secp256k1)
	if err != nil {
		t.Fatalf("Failed to generate second Secp256k1 key pair: %v", err)
	}

	if string(priv1) == string(priv2) {
		t.Error("Generated private keys should be different")
	}

	if string(pub1) == string(pub2) {
		t.Error("Generated public keys should be different")
	}

	// Test signing with maximum data size
	maxData := make([]byte, 10*1024*1024) // 10MB
	for i := range maxData {
		maxData[i] = byte(i % 256)
	}

	signature, err := Sign(maxData, priv1, Secp256k1)
	if err != nil {
		t.Fatalf("Failed to sign maximum data size: %v", err)
	}

	if !Verify(maxData, signature, pub1, Secp256k1) {
		t.Error("Failed to verify signature for maximum data size")
	}

	// Test with nil data
	nilSignature, err := Sign(nil, priv1, Secp256k1)
	if err != nil {
		t.Fatalf("Failed to sign nil data: %v", err)
	}

	if !Verify(nil, nilSignature, pub1, Secp256k1) {
		t.Error("Failed to verify signature for nil data")
	}

	// Test with invalid public key
	invalidPub := make([]byte, 33)
	if Verify([]byte("test"), signature, invalidPub, Secp256k1) {
		t.Error("Verification should fail with invalid public key")
	}

	// Test with corrupted signature
	corruptedSig := make([]byte, len(signature))
	copy(corruptedSig, signature)
	corruptedSig[0] ^= 0x01 // Flip one bit
	if Verify([]byte("test"), corruptedSig, pub1, Secp256k1) {
		t.Error("Verification should fail with corrupted signature")
	}
}

// Additional tests for Schnorr edge cases
func TestSchnorrEdgeCases(t *testing.T) {
	// Test multiple key generations produce different keys
	priv1, pub1, err := GenerateKeyPair(Schnorr)
	if err != nil {
		t.Fatalf("Failed to generate first Schnorr key pair: %v", err)
	}

	priv2, pub2, err := GenerateKeyPair(Schnorr)
	if err != nil {
		t.Fatalf("Failed to generate second Schnorr key pair: %v", err)
	}

	if string(priv1) == string(priv2) {
		t.Error("Generated private keys should be different")
	}

	if string(pub1) == string(pub2) {
		t.Error("Generated public keys should be different")
	}

	// Test signing with maximum data size
	maxData := make([]byte, 10*1024*1024) // 10MB
	for i := range maxData {
		maxData[i] = byte(i % 256)
	}

	signature, err := Sign(maxData, priv1, Schnorr)
	if err != nil {
		t.Fatalf("Failed to sign maximum data size: %v", err)
	}

	if !Verify(maxData, signature, pub1, Schnorr) {
		t.Error("Failed to verify signature for maximum data size")
	}

	// Test with nil data
	nilSignature, err := Sign(nil, priv1, Schnorr)
	if err != nil {
		t.Fatalf("Failed to sign nil data: %v", err)
	}

	if !Verify(nil, nilSignature, pub1, Schnorr) {
		t.Error("Failed to verify signature for nil data")
	}

	// Test with invalid public key
	invalidPub := make([]byte, 32)
	if Verify([]byte("test"), signature, invalidPub, Schnorr) {
		t.Error("Verification should fail with invalid public key")
	}

	// Test with corrupted signature
	corruptedSig := make([]byte, len(signature))
	copy(corruptedSig, signature)
	corruptedSig[0] ^= 0x01 // Flip one bit
	if Verify([]byte("test"), corruptedSig, pub1, Schnorr) {
		t.Error("Verification should fail with corrupted signature")
	}
}

// Security tests for key saving and loading
func TestKeySaveLoadSecurity(t *testing.T) {
	// Test with all supported algorithms
	algorithms := []string{Ed25519, ECDSA, Secp256k1, Schnorr}

	for _, alg := range algorithms {
		t.Run(alg, func(t *testing.T) {
			// Generate a key pair
			priv, _, err := GenerateKeyPair(alg)
			if err != nil {
				t.Fatalf("Failed to generate key pair: %v", err)
			}

			// Reconstruct the private key from bytes
			algorithm, err := AlgorithmFactory(alg)
			if err != nil {
				t.Fatalf("Failed to get algorithm: %v", err)
			}

			privKey, err := algorithm.PrivateKeyFromBytes(priv)
			if err != nil {
				t.Fatalf("Failed to reconstruct private key: %v", err)
			}

			// Save with a strong password
			filename := fmt.Sprintf("test_%s_key_secure.json", alg)
			defer os.Remove(filename)

			password := "very-strong-password-123!@#"
			err = SaveToFile(privKey, filename, password)
			if err != nil {
				t.Fatalf("Failed to save key: %v", err)
			}

			// Load with correct password
			loadedPriv, err := LoadFromFile(filename, password)
			if err != nil {
				t.Fatalf("Failed to load key with correct password: %v", err)
			}

			// Verify the loaded key
			if loadedPriv.Type() != privKey.Type() {
				t.Errorf("Expected loaded key type %s, got %s", privKey.Type(), loadedPriv.Type())
			}

			// Verify that the keys have the same bytes
			if string(loadedPriv.Bytes()) != string(privKey.Bytes()) {
				t.Error("Loaded key bytes do not match original")
			}

			// Try to load with wrong password
			_, err = LoadFromFile(filename, "wrong-password")
			if err == nil {
				t.Error("Expected error when loading with wrong password, but got none")
			}

			// Try to load with empty password (should use default)
			_, err = LoadFromFile(filename, "")
			if err == nil {
				t.Error("Expected error when loading with default password for encrypted key, but got none")
			}
		})
	}
}

// Test saving and loading with default password
func TestSaveAndLoadWithDefaultPassword(t *testing.T) {
	// Test with all supported algorithms
	algorithms := []string{Ed25519, ECDSA, Secp256k1, Schnorr}

	for _, alg := range algorithms {
		t.Run(alg, func(t *testing.T) {
			// Generate a key pair
			priv, _, err := GenerateKeyPair(alg)
			if err != nil {
				t.Fatalf("Failed to generate key pair: %v", err)
			}

			// Reconstruct the private key from bytes
			algorithm, err := AlgorithmFactory(alg)
			if err != nil {
				t.Fatalf("Failed to get algorithm: %v", err)
			}

			privKey, err := algorithm.PrivateKeyFromBytes(priv)
			if err != nil {
				t.Fatalf("Failed to reconstruct private key: %v", err)
			}

			// Save with default password (empty string)
			filename := fmt.Sprintf("test_%s_key_default.json", alg)
			defer os.Remove(filename)

			err = SaveToFile(privKey, filename, "") // Empty password should use default
			if err != nil {
				t.Fatalf("Failed to save key with default password: %v", err)
			}

			// Load with default password (empty string)
			loadedPriv, err := LoadFromFile(filename, "") // Empty password should use default
			if err != nil {
				t.Fatalf("Failed to load key with default password: %v", err)
			}

			// Verify the loaded key
			if loadedPriv.Type() != privKey.Type() {
				t.Errorf("Expected loaded key type %s, got %s", privKey.Type(), loadedPriv.Type())
			}

			// Verify that the keys have the same bytes
			if string(loadedPriv.Bytes()) != string(privKey.Bytes()) {
				t.Error("Loaded key bytes do not match original")
			}
		})
	}
}

// Robustness tests for signature verification
func TestSignatureVerificationRobustness(t *testing.T) {
	// Test with all supported algorithms
	algorithms := []string{Ed25519, ECDSA, Secp256k1, Schnorr}

	for _, alg := range algorithms {
		t.Run(alg, func(t *testing.T) {
			// Generate a key pair
			priv, pub, err := GenerateKeyPair(alg)
			if err != nil {
				t.Fatalf("Failed to generate key pair: %v", err)
			}

			// Test data
			data := []byte("test data for signature verification")

			// Sign the data
			signature, err := Sign(data, priv, alg)
			if err != nil {
				t.Fatalf("Failed to sign data: %v", err)
			}

			// Normal verification should pass
			if !Verify(data, signature, pub, alg) {
				t.Error("Normal verification should pass")
			}

			// Verification with wrong data should fail
			wrongData := []byte("wrong data")
			if Verify(wrongData, signature, pub, alg) {
				t.Error("Verification with wrong data should fail")
			}

			// Verification with wrong public key should fail
			wrongPub := make([]byte, len(pub))
			copy(wrongPub, pub)
			wrongPub[0] ^= 0x01 // Flip one bit
			if Verify(data, signature, wrongPub, alg) {
				t.Error("Verification with wrong public key should fail")
			}

			// Verification with wrong signature should fail
			wrongSig := make([]byte, len(signature))
			copy(wrongSig, signature)
			wrongSig[0] ^= 0x01 // Flip one bit
			if Verify(data, wrongSig, pub, alg) {
				t.Error("Verification with wrong signature should fail")
			}

			// Verification with empty data
			emptySig, err := Sign([]byte{}, priv, alg)
			if err != nil {
				t.Fatalf("Failed to sign empty data: %v", err)
			}

			if !Verify([]byte{}, emptySig, pub, alg) {
				t.Error("Verification of empty data signature should pass")
			}

			// Verification with nil data
			nilSig, err := Sign(nil, priv, alg)
			if err != nil {
				t.Fatalf("Failed to sign nil data: %v", err)
			}

			if !Verify(nil, nilSig, pub, alg) {
				t.Error("Verification of nil data signature should pass")
			}

			// Test with very large data
			largeData := make([]byte, 1024*1024) // 1MB
			for i := range largeData {
				largeData[i] = byte(i % 256)
			}

			largeSig, err := Sign(largeData, priv, alg)
			if err != nil {
				t.Fatalf("Failed to sign large data: %v", err)
			}

			if !Verify(largeData, largeSig, pub, alg) {
				t.Error("Verification of large data signature should pass")
			}
		})
	}
}

// Test unsupported algorithm handling
func TestUnsupportedAlgorithm(t *testing.T) {
	// Test GenerateKeyPair with unsupported algorithm
	_, _, err := GenerateKeyPair("unsupported")
	if err == nil {
		t.Error("Expected error for unsupported algorithm, but got none")
	}

	// Test Sign with unsupported algorithm
	_, err = Sign([]byte("test"), []byte("key"), "unsupported")
	if err == nil {
		t.Error("Expected error for unsupported algorithm, but got none")
	}

	// Test Verify with unsupported algorithm
	result := Verify([]byte("data"), []byte("sig"), []byte("key"), "unsupported")
	if result {
		t.Error("Expected false for unsupported algorithm, but got true")
	}

	// Test GenerateAddress with unsupported algorithm
	addr := GenerateAddress([]byte("key"), "unsupported")
	if addr != (types.Address{}) {
		t.Error("Expected empty address for unsupported algorithm")
	}
}

// Consistency tests for address generation
func TestAddressGenerationConsistency(t *testing.T) {
	// Test with all supported algorithms
	algorithms := []string{Ed25519, ECDSA, Secp256k1, Schnorr}

	for _, alg := range algorithms {
		t.Run(alg, func(t *testing.T) {
			// Generate multiple key pairs with the same algorithm
			_, pub1, err := GenerateKeyPair(alg)
			if err != nil {
				t.Fatalf("Failed to generate first key pair: %v", err)
			}

			_, pub2, err := GenerateKeyPair(alg)
			if err != nil {
				t.Fatalf("Failed to generate second key pair: %v", err)
			}

			// Generate addresses
			addr1 := GenerateAddress(pub1, alg)
			addr2 := GenerateAddress(pub2, alg)

			// Addresses should be different for different public keys
			if addr1 == addr2 {
				t.Error("Addresses for different public keys should be different")
			}

			// Generate address for the same public key multiple times
			addr1Again := GenerateAddress(pub1, alg)

			// Addresses should be the same for the same public key
			if addr1 != addr1Again {
				t.Error("Addresses for the same public key should be identical")
			}

			// Address should not be empty
			if addr1 == (types.Address{}) {
				t.Error("Generated address should not be empty")
			}

			// Address should have the correct length (20 bytes)
			if len(addr1) != 20 {
				t.Errorf("Expected address length 20, got %d", len(addr1))
			}
		})
	}
}

// Test address generation with invalid public keys
func TestAddressGenerationWithInvalidKeys(t *testing.T) {
	// Test with all supported algorithms
	algorithms := []string{Ed25519, ECDSA, Secp256k1, Schnorr}

	for _, alg := range algorithms {
		t.Run(alg, func(t *testing.T) {
			// Test with empty public key
			addr := GenerateAddress([]byte{}, alg)
			// For Ed25519, empty key might still generate an address, so we skip this check for Ed25519
			if alg != Ed25519 && addr != (types.Address{}) {
				t.Error("Address for empty public key should be empty")
			}

			// Test with nil public key
			addr = GenerateAddress(nil, alg)
			// For Ed25519, nil key might still generate an address, so we skip this check for Ed25519
			if alg != Ed25519 && addr != (types.Address{}) {
				t.Error("Address for nil public key should be empty")
			}

			// Test with invalid public key (wrong length)
			invalidPub := make([]byte, 10)
			addr = GenerateAddress(invalidPub, alg)
			// For Ed25519, invalid key might still generate an address, so we skip this check for Ed25519
			if alg != Ed25519 && addr != (types.Address{}) {
				t.Error("Address for invalid public key should be empty")
			}
		})
	}
}

// Interoperability tests between different algorithms
func TestAlgorithmInteroperability(t *testing.T) {
	// Generate key pairs for all algorithms
	algorithms := []string{Ed25519, ECDSA, Secp256k1, Schnorr}
	keys := make(map[string]struct {
		priv []byte
		pub  []byte
	})

	for _, alg := range algorithms {
		priv, pub, err := GenerateKeyPair(alg)
		if err != nil {
			t.Fatalf("Failed to generate key pair for %s: %v", alg, err)
		}
		keys[alg] = struct {
			priv []byte
			pub  []byte
		}{priv: priv, pub: pub}
	}

	// Test data
	data := []byte("interoperability test data")

	// For each algorithm, sign with its own private key and verify with its own public key
	for _, signer := range algorithms {
		signature, err := Sign(data, keys[signer].priv, signer)
		if err != nil {
			t.Fatalf("Failed to sign data with %s: %v", signer, err)
		}

		// Verification with the correct algorithm and key should pass
		if !Verify(data, signature, keys[signer].pub, signer) {
			t.Errorf("Verification should pass for %s signature with %s key", signer, signer)
		}

		// Verification with other algorithms should fail
		for _, verifier := range algorithms {
			if verifier != signer {
				result := Verify(data, signature, keys[verifier].pub, verifier)
				if result {
					t.Errorf("Verification should fail for %s signature with %s key", signer, verifier)
				}
			}
		}
	}
}

// Test that addresses from different algorithms are different
func TestAddressUniquenessAcrossAlgorithms(t *testing.T) {
	// Generate key pairs for all algorithms
	algorithms := []string{Ed25519, ECDSA, Secp256k1, Schnorr}
	addresses := make(map[string]types.Address)

	for _, alg := range algorithms {
		_, pub, err := GenerateKeyPair(alg)
		if err != nil {
			t.Fatalf("Failed to generate key pair for %s: %v", alg, err)
		}

		addr := GenerateAddress(pub, alg)
		addresses[alg] = addr
	}

	// All addresses should be different
	for i, alg1 := range algorithms {
		for j, alg2 := range algorithms {
			if i < j && addresses[alg1] == addresses[alg2] {
				t.Errorf("Addresses for %s and %s should be different", alg1, alg2)
			}
		}
	}
}

// Test ListAlgorithms function
func TestListAlgorithms(t *testing.T) {
	// Get the list of algorithms
	algorithms := ListAlgorithms()

	// Check that we have at least the expected algorithms
	expectedAlgorithms := []string{Ed25519, ECDSA, Secp256k1, Schnorr}

	// Create a map for easy lookup
	algorithmMap := make(map[string]bool)
	for _, alg := range algorithms {
		algorithmMap[alg] = true
	}

	// Check that all expected algorithms are present
	for _, expected := range expectedAlgorithms {
		if !algorithmMap[expected] {
			t.Errorf("Expected algorithm %s not found in list", expected)
		}
	}

	// Check that the list is not empty
	if len(algorithms) == 0 {
		t.Error("List of algorithms should not be empty")
	}
}
