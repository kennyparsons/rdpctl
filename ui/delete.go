package ui

import (
	"fmt"

	"github.com/manifoldco/promptui"

	"rdpctl/model"
)

// DeleteConnection guides the user through deleting an existing RDP connection profile from the vault.
func DeleteConnection(v *model.Vault) error {
	fmt.Println("\n--- Delete Host ---")

	selectedConn, err := SelectConnection(v)
	if err != nil {
		return fmt.Errorf("error selecting connection: %w", err)
	}

	confirmPrompt := promptui.Prompt{
		Label:     fmt.Sprintf("Are you sure you want to delete '%s'? (yes/no)", selectedConn.Name),
		IsConfirm: true,
	}
	_, err = confirmPrompt.Run()

	if err != nil {
		fmt.Println("Deletion cancelled.")
		return nil // User cancelled or input was not 'yes'
	}

	// Find and remove the connection
	for i, conn := range v.Connections {
		if conn.ID == selectedConn.ID {
			v.Connections = append(v.Connections[:i], v.Connections[i+1:]...)
			fmt.Printf("Connection '%s' deleted successfully!\n", selectedConn.Name)
			return nil
		}
	}

	return fmt.Errorf("failed to find connection with ID %s to delete", selectedConn.ID)
}
