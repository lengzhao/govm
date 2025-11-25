package main

import (
	"fmt"
	"log"

	"github.com/lengzhao/govm/crypto"
)

func main() {
	// 生成Ed25519密钥对
	fmt.Println("Generating Ed25519 key pair...")
	privKeyBytes, pubKeyBytes, err := crypto.GenerateKeyPair(crypto.Ed25519)
	if err != nil {
		log.Fatal("Failed to generate Ed25519 key pair:", err)
	}

	fmt.Printf("Private key length: %d bytes\n", len(privKeyBytes))
	fmt.Printf("Public key length: %d bytes\n", len(pubKeyBytes))

	// 从字节数据重建密钥对象
	algorithm, err := crypto.AlgorithmFactory(crypto.Ed25519)
	if err != nil {
		log.Fatal("Failed to get algorithm:", err)
	}

	privKey, err := algorithm.PrivateKeyFromBytes(privKeyBytes)
	if err != nil {
		log.Fatal("Failed to reconstruct private key:", err)
	}

	pubKey, err := algorithm.PublicKeyFromBytes(pubKeyBytes)
	if err != nil {
		log.Fatal("Failed to reconstruct public key:", err)
	}

	fmt.Printf("Reconstructed private key type: %s\n", privKey.Type())
	fmt.Printf("Reconstructed public key type: %s\n", pubKey.Type())

	// 使用重建的密钥进行签名和验证
	message := []byte("Hello, Blockchain!")
	signature, err := privKey.Sign(message)
	if err != nil {
		log.Fatal("Failed to sign message:", err)
	}

	fmt.Printf("Signature length: %d bytes\n", len(signature))

	// 验证签名
	valid := pubKey.Verify(message, signature)
	if valid {
		fmt.Println("Signature is valid!")
	} else {
		fmt.Println("Signature is invalid!")
	}

	// 生成地址
	address := pubKey.Address()
	fmt.Printf("Generated address: %x\n", address)

	// Demonstrate the new ListAlgorithms function
	fmt.Println("\n--- Demonstrating ListAlgorithms function ---")
	listAlgorithmsDemo()
}

func listAlgorithmsDemo() {
	// List all available algorithms
	algorithms := crypto.ListAlgorithms()

	fmt.Println("Available algorithms:")
	for _, algorithm := range algorithms {
		fmt.Printf("- %s\n", algorithm)
	}
}
