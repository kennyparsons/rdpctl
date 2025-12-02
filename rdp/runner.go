package rdp

import (
	"fmt"
	"os"
	"os/exec"

	"rdpctl/model"
)

// Run executes the xfreerdp command with the given connection and password.
func Run(c *model.Connection, password string) error {
	// Build the arguments for xfreerdp
	args := BuildArgs(c, password)

	// Print the sanitized command for user information (excluding sensitive data)
	fmt.Printf("Running xfreerdp %s\n", SanitizeArgsForDisplay(args))

	cmd := exec.Command("xfreerdp", args...)

	// Attach stdin, stdout, and stderr to the current process
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the command
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("xfreerdp command failed: %w", err)
	}

	return nil
}
