package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
)

func decodePrivateKey(encodedPrivateKey string) (*rsa.PrivateKey, error) {
	privPEM, err := base64.StdEncoding.DecodeString(encodedPrivateKey)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(privPEM)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block containing private key")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

func generateKeyPair() (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}
	publicKey := &privateKey.PublicKey

	return privateKey, publicKey, nil
}

func encodePrivateKey(privateKey *rsa.PrivateKey) (string, error) {
	privBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: privBytes})

	return base64.StdEncoding.EncodeToString(privPEM), nil
}

func encodePublicKey(publicKey *rsa.PublicKey) (string, error) {
	pubBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return "", err
	}
	pubPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PUBLIC KEY", Bytes: pubBytes})

	return base64.StdEncoding.EncodeToString(pubPEM), nil
}

func encrypt(publicKey *rsa.PublicKey, plaintext []byte) ([]byte, error) {
	return rsa.EncryptPKCS1v15(rand.Reader, publicKey, plaintext)
}

func decrypt(privateKey *rsa.PrivateKey, ciphertext []byte) ([]byte, error) {
	return rsa.DecryptPKCS1v15(rand.Reader, privateKey, ciphertext)
}

func main() {
	// Generate key pair
	privateKey, publicKey, err := generateKeyPair()
	if err != nil {
		fmt.Println("Error generating key pair:", err)
		return
	}

	// Encode and print private key
	encodedPrivateKey, err := encodePrivateKey(privateKey)
	if err != nil {
		fmt.Println("Error encoding private key:", err)
		return
	}
	fmt.Println("Encoded Private Key:")
	fmt.Println(encodedPrivateKey)

	// Encode and print public key
	encodedPublicKey, err := encodePublicKey(publicKey)
	if err != nil {
		fmt.Println("Error encoding public key:", err)
		return
	}
	fmt.Println("\nEncoded Public Key:")
	fmt.Println(encodedPublicKey)

	// Example encryption and decryption
	plaintext := []byte("Won't secure anything, but please do not abuse, api food lookup v1.0.0")

	// Encrypt with public key
	ciphertext, err := encrypt(publicKey, plaintext)
	if err != nil {
		fmt.Println("Error encrypting:", err)
		return
	}

	fmt.Println("\nEncrypted:", base64.StdEncoding.EncodeToString(ciphertext))

	fmt.Println("\nStupid:", base64.StdEncoding.EncodeToString([]byte(`Won't secure anything, but please do not abuse, api food lookup v1.0.0`)))

	// Decrypt with private key
	decrypted, err := decrypt(privateKey, ciphertext)
	if err != nil {
		fmt.Println("Error decrypting:", err)
		return
	}

	fmt.Println("Decrypted:", string(decrypted))
}
