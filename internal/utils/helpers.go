package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"
	"os"
)

func GenerateIdempotencyKey() string {
	// Generate 8 random bytes (which will give us 16 hexadecimal characters)
	bytes := make([]byte, 8)
	_, err := rand.Read(bytes)
	if err != nil {
		log.Fatalf("Failed to generate idempotency key: %v", err)
	}

	// Convert the random bytes to a hexadecimal string
	idempotencyKey := hex.EncodeToString(bytes)

	return idempotencyKey
}

// AES decryption function
func DecryptCookiestring(encryptedCookie string) (string, error) {
	secretKey := os.Getenv("COOKIE_SECRET_KEY") // Store the secret key in environment variables
	// Decode the base64 encoded Cookie
	encryptedData, err := base64.StdEncoding.DecodeString(encryptedCookie)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher([]byte(secretKey))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(encryptedData) < nonceSize {
		return "", fmt.Errorf("encrypted data too short")
	}

	nonce, ciphertext := encryptedData[:nonceSize], encryptedData[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
