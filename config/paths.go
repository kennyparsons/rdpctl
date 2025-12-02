package config

import (
	"fmt"
	"os"
	"path/filepath"
)

// ConfigDir returns the path to the rdpctl configuration directory (~/.config/rdp/).
func ConfigDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}
	return filepath.Join(homeDir, ".config", "rdp"), nil
}

// EnsureConfigDir ensures the rdpctl configuration directory exists.
func EnsureConfigDir() (string, error) {
	configDir, err := ConfigDir()
	if err != nil {
		return "", err
	}

	if err := os.MkdirAll(configDir, 0700); err != nil {
		return "", fmt.Errorf("failed to create config directory %s: %w", configDir, err)
	}

	return configDir, nil
}

// VaultPath returns the full path to the encrypted vault file.
func VaultPath(dirname string) string {
	return filepath.Join(dirname, "vault.enc")
}
