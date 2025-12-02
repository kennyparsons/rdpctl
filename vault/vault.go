package vault

import (
	"fmt"
	"os"

	"rdpctl/model"
)

// LoadVault loads an existing vault from the specified path using the provided password.
// It returns the loaded vault and any error encountered.
func LoadVault(path string, password string) (*model.Vault, error) {
	// Check if the vault file exists
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("vault file not found at %s", path)
	}
	if err != nil {
		return nil, fmt.Errorf("error checking vault file: %w", err)
	}

	// Decrypt and unmarshal the vault
	v, err := DecryptAndUnmarshalVault(path, password)
	if err != nil {
		return nil, fmt.Errorf("failed to load and decrypt vault: %w", err)
	}
	return v, nil
}

// SaveVault encrypts and saves the given vault to the specified path.
// It returns an error if the operation fails.
func SaveVault(path string, v *model.Vault, password string) error {
	// Marshal and encrypt the vault, then write to file
	if err := MarshalAndEncryptVault(path, v, password); err != nil {
		return fmt.Errorf("failed to encrypt and save vault: %w", err)
	}
	return nil
}

// CreateNewVault creates a new empty vault and saves it to the specified path.
// It returns the newly created vault and any error encountered.
func CreateNewVault(path string, password string) (*model.Vault, error) {
	v, err := CreateAndSaveNewVaultFile(path, password)
	if err != nil {
		return nil, fmt.Errorf("failed to create and save new vault file: %w", err)
	}
	return v, nil
}