package vault

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"

	"golang.org/x/crypto/argon2"
)

const (
	saltLength = 16
	nonceLength = 12
	keyLength = 32 // AES-256
)

// DeriveKey derives a cryptographic key from the password and salt using Argon2id.
func DeriveKey(password string, salt []byte) ([]byte, error) {
	if len(salt) != saltLength {
		return nil, fmt.Errorf("salt must be %d bytes long", saltLength)
	}

	// Using recommended parameters for Argon2id
	// time: 1, memory: 64MB, threads: 4
	key := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, keyLength)
	return key, nil
}

// EncryptVault encrypts the JSON data using AES-256-GCM with a master password.
func EncryptVault(jsonData []byte, password string) (salt, nonce, ciphertext []byte, err error) {
	salt = make([]byte, saltLength)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, nil, nil, fmt.Errorf("failed to generate salt: %w", err)
	}

	key, err := DeriveKey(password, salt)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to derive key: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to create AES cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	nonce = make([]byte, nonceLength)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, nil, nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	ciphertext = gcm.Seal(nil, nonce, jsonData, nil)
	return salt, nonce, ciphertext, nil
}

// DecryptVault decrypts the ciphertext using AES-256-GCM with the master password.
func DecryptVault(salt, nonce, ciphertext []byte, password string) ([]byte, error) {
	key, err := DeriveKey(password, salt)
	if err != nil {
		return nil, fmt.Errorf("failed to derive key: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt vault: %w", err)
	}

	return plaintext, nil
}
