package rdp

import (
	"fmt"
	"strings"

	"rdpctl/model"
)

// BuildArgs constructs the arguments slice for the xfreerdp command.
func BuildArgs(c *model.Connection, password string) []string {
	args := []string{
		"+clipboard",         // Enable clipboard redirection
		"+dynamic-resolution", // Enable dynamic resolution updates
	}

	// Mandatory arguments
	args = append(args, fmt.Sprintf("/v:%s", c.Host))
	args = append(args, fmt.Sprintf("/u:%s", c.Username))

	if c.Domain != "" {
		args = append(args, fmt.Sprintf("/d:%s", c.Domain))
	}

	// Handle password storage
	if c.StorePassword {
		// If password is stored in the vault, use it directly.
		// NOTE: For MVP, we are passing the password directly via /p:. Day 2 will improve this.
		args = append(args, fmt.Sprintf("/p:%s", password))
	} else if password != "" {
		// If not stored but provided at runtime (e.g., from prompt), use it.
		args = append(args, fmt.Sprintf("/p:%s", password))
	}

	// Add any extra arguments configured by the user
	if len(c.ExtraArgs) > 0 {
		args = append(args, c.ExtraArgs...)
	}

	return args
}

// SanitizeArgsForDisplay removes sensitive information (like passwords) from the argument list
// before displaying or logging them.
func SanitizeArgsForDisplay(args []string) []string {
	sanitizedArgs := make([]string, 0, len(args))
	for _, arg := range args {
		if strings.HasPrefix(arg, "/p:") {
			sanitizedArgs = append(sanitizedArgs, "/p:********")
		} else {
			sanitizedArgs = append(sanitizedArgs, arg)
		}
	}
	return sanitizedArgs
}
