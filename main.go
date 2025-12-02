package main

import (
	"fmt"
	"log"

	"rdpctl/config"
	"rdpctl/ui"
)

func main() {
	// Ensure the configuration directory exists
	configDir, err := config.EnsureConfigDir()
	if err != nil {
		log.Fatalf("Error ensuring config directory: %v", err)
	}

	// Get the full path to the vault file
	vaultPath := config.VaultPath(configDir)

	// Run the unlock flow (handles first-time setup and existing vault unlock)
	v, masterPassword, err := ui.UnlockFlow(vaultPath)
	if err != nil {
		log.Fatalf("Error during vault unlock flow: %v", err)
	}

	// If unlock was successful, enter the main menu loop
	if err := ui.MainMenu(v, masterPassword, vaultPath); err != nil {
		log.Fatalf("Error during main menu: %v", err)
	}
	fmt.Println("Application exited gracefully.")
}
