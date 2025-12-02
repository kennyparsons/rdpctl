package vault

import (

	"encoding/json"
	"fmt"
	"io"
	"os"

	"rdpctl/model"
)

const (
	magicBytes = "RDP1"
	version = 1
	magicLen = 4
	versionLen = 1
)

// ReadVaultFile reads an encrypted vault file from disk.
func ReadVaultFile(path string) (salt, nonce, ciphertext []byte, err error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to open vault file: %w", err)
	}
	defer file.Close()

	// Read Magic Bytes
	readMagic := make([]byte, magicLen)
	if _, err := io.ReadFull(file, readMagic); err != nil {
		return nil, nil, nil, fmt.Errorf("failed to read magic bytes: %w", err)
	}
	if string(readMagic) != magicBytes {
		return nil, nil, nil, fmt.Errorf("invalid vault file magic bytes")
	}

	// Read Version
	readVersion := make([]byte, versionLen)
	if _, err := io.ReadFull(file, readVersion); err != nil {
		return nil, nil, nil, fmt.Errorf("failed to read version: %w", err)
	}
	if readVersion[0] != version {
		return nil, nil, nil, fmt.Errorf("unsupported vault version: %d", readVersion[0])
	}

	// Read Salt
	salt = make([]byte, saltLength)
	if _, err := io.ReadFull(file, salt); err != nil {
		return nil, nil, nil, fmt.Errorf("failed to read salt: %w", err)
	}

	// Read Nonce
	nonce = make([]byte, nonceLength)
	if _, err := io.ReadFull(file, nonce); err != nil {
		return nil, nil, nil, fmt.Errorf("failed to read nonce: %w", err)
	}

	// Read Ciphertext
	ciphertext, err = io.ReadAll(file)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to read ciphertext: %w", err)
	}

	return salt, nonce, ciphertext, nil
}

// WriteVaultFile writes an encrypted vault file to disk.
func WriteVaultFile(path string, salt, nonce, ciphertext []byte) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create vault file: %w", err)
	}
	defer file.Close()

	// Write Magic Bytes
	if _, err := file.Write([]byte(magicBytes)); err != nil {
		return fmt.Errorf("failed to write magic bytes: %w", err)
	}

	// Write Version
	if _, err := file.Write([]byte{version}); err != nil {
		return fmt.Errorf("failed to write version: %w", err)
	}

	// Write Salt
	if _, err := file.Write(salt); err != nil {
		return fmt.Errorf("failed to write salt: %w", err)
	}

	// Write Nonce
	if _, err := file.Write(nonce); err != nil {
		return fmt.Errorf("failed to write nonce: %w", err)
	}

	// Write Ciphertext
	if _, err := file.Write(ciphertext); err != nil {
		return fmt.Errorf("failed to write ciphertext: %w", err)
	}

	return nil
}

// DecryptAndUnmarshalVault reads, decrypts, and unmarshals the vault from disk.
func DecryptAndUnmarshalVault(path string, password string) (*model.Vault, error) {
	salt, nonce, ciphertext, err := ReadVaultFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read vault file: %w", err)
	}

	plaintext, err := DecryptVault(salt, nonce, ciphertext, password)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt vault: %w", err)
	}

	v := &model.Vault{}
	if err := json.Unmarshal(plaintext, v); err != nil {
		return nil, fmt.Errorf("failed to unmarshal vault JSON: %w", err)
	}

	return v, nil
}

// MarshalAndEncryptVault marshals the vault, encrypts it, and writes it to disk.
func MarshalAndEncryptVault(path string, v *model.Vault, password string) error {
	jsonData, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("failed to marshal vault: %w", err)
	}

	salt, nonce, ciphertext, err := EncryptVault(jsonData, password)
	if err != nil {
		return fmt.Errorf("failed to encrypt vault: %w", err)
	}

	if err := WriteVaultFile(path, salt, nonce, ciphertext); err != nil {
		return fmt.Errorf("failed to write vault file: %w", err)
	}

	return nil
}

// CreateAndSaveNewVaultFile creates a new, empty vault and saves it to disk.
func CreateAndSaveNewVaultFile(path string, password string) (*model.Vault, error) {
	v := &model.Vault{
		Version:     version,
		Connections: []model.Connection{},
	}

	if err := MarshalAndEncryptVault(path, v, password); err != nil {
		return nil, fmt.Errorf("failed to create and encrypt new vault: %w", err)
	}

	return v, nil
}

