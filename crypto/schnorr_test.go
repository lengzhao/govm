package crypto

import (
	"testing"
)

func TestSchnorrKeyGeneration(t *testing.T) {
	priv, pub, err := GenerateKeyPair(Schnorr)
	if err != nil {
		t.Fatalf("Failed to generate Schnorr key pair: %v", err)
	}

	if priv == nil {
		t.Error("Private key is nil")
	}

	if pub == nil {
		t.Error("Public key is nil")
	}

	// Since we're now returning []byte directly, we can only check the length
	if len(priv) == 0 {
		t.Error("Private key is empty")
	}

	if len(pub) == 0 {
		t.Error("Public key is empty")
	}
}

func TestSchnorrSigningAndVerification(t *testing.T) {
	priv, pub, err := GenerateKeyPair(Schnorr)
	if err != nil {
		t.Fatalf("Failed to generate Schnorr key pair: %v", err)
	}

	data := []byte("test message")
	signature, err := Sign(data, priv, Schnorr)
	if err != nil {
		t.Fatalf("Failed to sign data: %v", err)
	}

	if signature == nil {
		t.Error("Signature is nil")
	}

	valid := Verify(data, signature, pub, Schnorr)
	if !valid {
		t.Error("Failed to verify signature")
	}

	// Test with wrong data
	wrongData := []byte("wrong message")
	valid = Verify(wrongData, signature, pub, Schnorr)
	if valid {
		t.Error("Verification should fail with wrong data")
	}
}

func TestSchnorrAddressGeneration(t *testing.T) {
	_, pub, err := GenerateKeyPair(Schnorr)
	if err != nil {
		t.Fatalf("Failed to generate Schnorr key pair: %v", err)
	}

	address := GenerateAddress(pub, Schnorr)
	if address == [20]byte{} {
		t.Error("Address is empty")
	}
}

// Test FromBytes methods - these tests are no longer applicable with the new interface
func TestSchnorrPrivateKeyFromBytes(t *testing.T) {
	t.Skip("Test not applicable with new interface")
}

func TestSchnorrPublicKeyFromBytesTest(t *testing.T) {
	t.Skip("Test not applicable with new interface")
}

// Additional edge case tests for Schnorr
func TestSchnorrAdditionalEdgeCases(t *testing.T) {
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

	// Test with empty data
	emptyData := []byte{}
	emptySignature, err := Sign(emptyData, priv1, Schnorr)
	if err != nil {
		t.Fatalf("Failed to sign empty data: %v", err)
	}

	if !Verify(emptyData, emptySignature, pub1, Schnorr) {
		t.Error("Failed to verify signature for empty data")
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

	// Test signature verification with wrong public key
	testData := []byte("test data")
	wrongPriv, _, err := GenerateKeyPair(Schnorr)
	if err != nil {
		t.Fatalf("Failed to generate wrong key pair: %v", err)
	}

	wrongSignature, err := Sign(testData, wrongPriv, Schnorr)
	if err != nil {
		t.Fatalf("Failed to sign with wrong key: %v", err)
	}

	if Verify(testData, wrongSignature, pub1, Schnorr) {
		t.Error("Verification should fail with wrong public key")
	}
}

// Test Schnorr address generation edge cases
func TestSchnorrAddressAdditionalEdgeCases(t *testing.T) {
	// Test multiple address generations produce consistent results
	_, pub, err := GenerateKeyPair(Schnorr)
	if err != nil {
		t.Fatalf("Failed to generate key pair: %v", err)
	}

	addr1 := GenerateAddress(pub, Schnorr)
	addr2 := GenerateAddress(pub, Schnorr)

	if addr1 != addr2 {
		t.Error("Address generation should be deterministic")
	}

	// Test address with invalid public key
	invalidPub := make([]byte, 10)
	addr := GenerateAddress(invalidPub, Schnorr)
	if addr != [20]byte{} {
		t.Error("Address for invalid public key should be empty")
	}

	// Test address with empty public key
	emptyAddr := GenerateAddress([]byte{}, Schnorr)
	if emptyAddr != [20]byte{} {
		t.Error("Address for empty public key should be empty")
	}

	// Test address with nil public key
	nilAddr := GenerateAddress(nil, Schnorr)
	if nilAddr != [20]byte{} {
		t.Error("Address for nil public key should be empty")
	}
}
