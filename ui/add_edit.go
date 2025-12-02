package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/manifoldco/promptui"

	"rdpctl/model"
)

// AddConnection guides the user through adding a new RDP connection profile to the vault.
func AddConnection(v *model.Vault) error {
	fmt.Println("\n--- Add New Host ---")

	newConn := model.Connection{
		ID:        uuid.New().String(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Prompt for Friendly Name
	namePrompt := promptui.Prompt{
		Label:    "Friendly Name",
		Validate: requireInput,
	}
	name, err := namePrompt.Run()
	if err != nil {
		return fmt.Errorf("prompt failed: %w", err)
	}
	newConn.Name = name

	// Prompt for Host/IP
	hostPrompt := promptui.Prompt{
		Label:    "Host/IP",
		Validate: requireInput,
	}
	host, err := hostPrompt.Run()
	if err != nil {
		return fmt.Errorf("prompt failed: %w", err)
	}
	newConn.Host = host

	// Prompt for Domain (optional)
	domainPrompt := promptui.Prompt{
		Label: "Domain (optional)",
	}
	domain, err := domainPrompt.Run()
	if err != nil {
		return fmt.Errorf("prompt failed: %w", err)
	}
	newConn.Domain = domain

	// Prompt for Username
	usernamePrompt := promptui.Prompt{
		Label:    "Username",
		Validate: requireInput,
	}
	username, err := usernamePrompt.Run()
	if err != nil {
		return fmt.Errorf("prompt failed: %w", err)
	}
	newConn.Username = username

	// Prompt to store password
	storePassPrompt := promptui.Select{ 
		Label: "Store password in vault?",
		Items: []string{"Yes", "No"},
	}
	_, storePassResult, err := storePassPrompt.Run()
	if err != nil {
		return fmt.Errorf("prompt failed: %w", err)
	}
	newConn.StorePassword = (storePassResult == "Yes")

	if newConn.StorePassword {
		passwordPrompt := promptui.Prompt{
			Label:    "Password",
			Mask:     '*',
			Validate: requireInput,
		}
		password, err := passwordPrompt.Run()
		if err != nil {
			return fmt.Errorf("prompt failed: %w", err)
		}
		newConn.Password = password
	}

	// Prompt for Extra Args (optional, comma-separated)
	extraArgsPrompt := promptui.Prompt{
		Label: "Extra xfreerdp arguments (comma-separated, e.g., /cert-ignore)",
	}
	extraArgsInput, err := extraArgsPrompt.Run()
	if err != nil {
		return fmt.Errorf("prompt failed: %w", err)
	}
	if extraArgsInput != "" {
		newConn.ExtraArgs = strings.Split(extraArgsInput, ",")
		for i, arg := range newConn.ExtraArgs {
			newConn.ExtraArgs[i] = strings.TrimSpace(arg)
		}
	}

	v.Connections = append(v.Connections, newConn)
	fmt.Printf("Connection '%s' added successfully!\n", newConn.Name)
	return nil
}

// EditConnection guides the user through editing an existing RDP connection profile in the vault.
func EditConnection(v *model.Vault) error {
	fmt.Println("\n--- Edit Existing Host ---")

	selectedConn, err := SelectConnection(v)
	if err != nil {
		return fmt.Errorf("error selecting connection: %w", err)
	}

	originalID := selectedConn.ID
	originalCreatedAt := selectedConn.CreatedAt

	// Create a copy to edit, so we don't modify the original in case of cancellation
	editedConn := *selectedConn
	editedConn.UpdatedAt = time.Now()

	// Prompt for Friendly Name
	namePrompt := promptui.Prompt{
		Label:     fmt.Sprintf("Friendly Name (current: %s)", editedConn.Name),
		Default:   editedConn.Name,
		Validate:  requireInput,
	}
	name, err := namePrompt.Run()
	if err != nil {
		return fmt.Errorf("prompt failed: %w", err)
	}
	editedConn.Name = name

	// Prompt for Host/IP
	hostPrompt := promptui.Prompt{
		Label:     fmt.Sprintf("Host/IP (current: %s)", editedConn.Host),
		Default:   editedConn.Host,
		Validate:  requireInput,
	}
	host, err := hostPrompt.Run()
	if err != nil {
		return fmt.Errorf("prompt failed: %w", err)
	}
	editedConn.Host = host

	// Prompt for Domain (optional)
	domainPrompt := promptui.Prompt{
		Label:   fmt.Sprintf("Domain (current: %s)", editedConn.Domain),
		Default: editedConn.Domain,
	}
	domain, err := domainPrompt.Run()
	if err != nil {
		return fmt.Errorf("prompt failed: %w", err)
	}
	editedConn.Domain = domain

	// Prompt for Username
	usernamePrompt := promptui.Prompt{
		Label:     fmt.Sprintf("Username (current: %s)", editedConn.Username),
		Default:   editedConn.Username,
		Validate:  requireInput,
	}
	username, err := usernamePrompt.Run()
	if err != nil {
		return fmt.Errorf("prompt failed: %w", err)
	}
	editedConn.Username = username

	// Prompt to store password
	var storePassDefault string
	if editedConn.StorePassword {
		storePassDefault = "Yes"
	} else {
		storePassDefault = "No"
	}
	storePassPrompt := promptui.Select{
		Label:   fmt.Sprintf("Store password in vault? (current: %s)", storePassDefault),
		Items:   []string{"Yes", "No"},
	}
	_, storePassResult, err := storePassPrompt.Run()
	if err != nil {
		return fmt.Errorf("prompt failed: %w", err)
	}
	editedConn.StorePassword = (storePassResult == "Yes")

	// If not storing, clear password. If storing, prompt for it.
	if !editedConn.StorePassword {
		editedConn.Password = "" // Clear stored password if user opts out
	} else {
		passwordPrompt := promptui.Prompt{
			Label:    "Password (leave blank to keep current or if not stored)",
			Mask:     '*',
		}
		password, err := passwordPrompt.Run()
		if err != nil {
			return fmt.Errorf("prompt failed: %w", err)
		}
		if password != "" {
			editedConn.Password = password
		}
	}

	// Prompt for Extra Args (optional, comma-separated)
	extraArgsDefault := strings.Join(editedConn.ExtraArgs, ", ")
	extraArgsPrompt := promptui.Prompt{
		Label:   fmt.Sprintf("Extra xfreerdp arguments (comma-separated, current: %s)", extraArgsDefault),
		Default: extraArgsDefault,
	}
	extraArgsInput, err := extraArgsPrompt.Run()
	if err != nil {
		return fmt.Errorf("prompt failed: %w", err)
	}
	if extraArgsInput != "" {
		editedConn.ExtraArgs = strings.Split(extraArgsInput, ",")
		for i, arg := range editedConn.ExtraArgs {
			editedConn.ExtraArgs[i] = strings.TrimSpace(arg)
		}
	} else {
		editedConn.ExtraArgs = []string{} // Clear extra args if input is empty
	}

	// Find and replace the old connection with the edited one
	for i, conn := range v.Connections {
		if conn.ID == originalID {
			editedConn.ID = originalID // Ensure ID and creation date are preserved
			editedConn.CreatedAt = originalCreatedAt
			v.Connections[i] = editedConn
			fmt.Printf("Connection '%s' updated successfully!\n", editedConn.Name)
			return nil
		}
	}

	return fmt.Errorf("failed to find connection with ID %s to update", originalID)
}

// requireInput is a promptui.ValidateFunc to ensure the input is not empty.
func requireInput(input string) error {
	if strings.TrimSpace(input) == "" {
		return fmt.Errorf("input cannot be empty")
	}
	return nil
}
