package ui

import (
	"fmt"

	"github.com/manifoldco/promptui"

	"rdpctl/model"
	"rdpctl/rdp"
	"rdpctl/vault"
)

// MainMenu displays the main menu and handles user selections.
func MainMenu(v *model.Vault, masterPassword string, vaultPath string) error {
	for {
		prompt := promptui.Select{
			Label: "Main Menu",
			Items: []string{
				"Connect to a host",
				"Add new host",
				"Edit existing host",
				"Delete host",
				"Show vault",
				"Quit",
			},
		}

		i, _, err := prompt.Run()
		if err != nil {
			// If the user presses Ctrl+C, it's considered an interrupt and we should exit gracefully.
			if err == promptui.ErrInterrupt {
				fmt.Println("Operation cancelled.")
				return nil
			}
			fmt.Printf("Prompt failed %v\n", err)
			return err
		}

		switch i {
			case 0: // Connect to a host
				selectedConn, err := SelectConnection(v)
				if err != nil {
					// If the user cancelled the selection, continue to main menu
					if err == promptui.ErrInterrupt {
						continue
					}
					fmt.Printf("Error selecting connection: %v\n", err)
					continue
				}

				connPassword := selectedConn.Password
				if !selectedConn.StorePassword || connPassword == "" {
					passwordPrompt := promptui.Prompt{
						Label: "Enter password for " + selectedConn.Username + "@" + selectedConn.Host,
						Mask:  '*',
					}
					result, err := passwordPrompt.Run()
					if err != nil {
						// If the user cancelled the password prompt, continue to main menu
						if err == promptui.ErrInterrupt {
							continue
						}
						fmt.Printf("Password prompt failed %v\n", err)
						continue
					}
					connPassword = result
				}

				fmt.Printf("Connecting to %s...\n", selectedConn.Name)
				if err := rdp.Run(selectedConn, connPassword); err != nil {
					fmt.Printf("RDP connection failed: %v\n", err)
				}
			case 1: // Add new host
				if err := AddConnection(v); err != nil {
					// If the user cancelled the operation, continue to main menu
					if err == promptui.ErrInterrupt {
						continue
					}
					fmt.Printf("Error adding connection: %v\n", err)
				} else {
					if err := vault.SaveVault(vaultPath, v, masterPassword); err != nil {
						fmt.Printf("Error saving vault: %v\n", err)
					}
				}
			case 2: // Edit existing host
				if err := EditConnection(v); err != nil {
					// If the user cancelled the operation, continue to main menu
					if err == promptui.ErrInterrupt {
						continue
					}
					fmt.Printf("Error editing connection: %v\n", err)
				} else {
					if err := vault.SaveVault(vaultPath, v, masterPassword); err != nil {
						fmt.Printf("Error saving vault: %v\n", err)
					}
				}
			case 3: // Delete host
				if err := DeleteConnection(v); err != nil {
					// If the user cancelled the operation, continue to main menu
					if err == promptui.ErrInterrupt {
						continue
					}
					fmt.Printf("Error deleting connection: %v\n", err)
				} else {
					if err := vault.SaveVault(vaultPath, v, masterPassword); err != nil {
						fmt.Printf("Error saving vault: %v\n", err)
					}
				}
			case 4: // Show vault
				if err := ShowVault(v); err != nil {
					fmt.Printf("Error showing vault: %v\n", err)
				}
			case 5: // Quit
				fmt.Println("Goodbye!")
				return nil
		}
	}
}
