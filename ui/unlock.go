package ui

import (
	"fmt"
	"os"

	"github.com/manifoldco/promptui"

	"rdpctl/model"
	"rdpctl/vault"
)

// UnlockFlow handles the process of unlocking the vault, including first-time setup.
// It returns the unlocked vault, the master password used, and an error if any.
func UnlockFlow(vaultPath string) (*model.Vault, string, error) {
	// Check if vault file exists
	if _, err := os.Stat(vaultPath); os.IsNotExist(err) {
		fmt.Println("No vault found. Let's create one.")
		return firstTimeSetup(vaultPath)
	} else if err != nil {
		return nil, "", fmt.Errorf("error checking vault file: %w", err)
	}

	// Vault exists, prompt for password
	return unlockExistingVault(vaultPath)
}

func firstTimeSetup(vaultPath string) (*model.Vault, string, error) {
	fmt.Println("Please set a master password for your new vault.")

	password, err := promptForMasterPassword("Enter master password: ")
	if err != nil {
		return nil, "", err
	}

	confirmPassword, err := promptForMasterPassword("Confirm master password: ")
	if err != nil {
		return nil, "", err
	}

	if password != confirmPassword {
		return nil, "", fmt.Errorf("passwords do not match")
	}

	v, err := vault.CreateNewVault(vaultPath, password)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create new vault: %w", err)
	}

	fmt.Println("New vault created and encrypted successfully!")
	return v, password, nil
}

func unlockExistingVault(vaultPath string) (*model.Vault, string, error) {
	fmt.Println("Vault found. Please enter your master password to unlock.")

	const maxRetries = 3
	for attempts := 0; attempts < maxRetries; attempts++ {
		password, err := promptForMasterPassword("Enter master password: ")
		if err != nil {
			return nil, "", err
		}

		v, err := vault.LoadVault(vaultPath, password)
		if err == nil {
			fmt.Println("Vault unlocked successfully!")
			return v, password, nil
		} else {
			fmt.Printf("Incorrect password. %d attempts remaining.\n", maxRetries-1-attempts)
		}
	}
	return nil, "", fmt.Errorf("too many incorrect password attempts")
}

func promptForMasterPassword(label string) (string, error) {
	prompt := promptui.Prompt{
		Label: label,
		Mask:  '*',
	}
	result, err := prompt.Run()
	
	return result, err
}
