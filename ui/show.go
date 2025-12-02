package ui

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"rdpctl/model"
)

// ShowVault displays the contents of the vault in a read-only view.
func ShowVault(v *model.Vault) error {
	fmt.Println("\n--- Vault Contents ---")

	if len(v.Connections) == 0 {
		fmt.Println("Vault is empty.")
		return waitForEnter()
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "Name\tHost\tDomain\tUsername\tPassword\tExtra Args")
	fmt.Fprintln(w, "----\t----\t------\t--------\t--------\t----------")

	for _, conn := range v.Connections {
		passwordDisplay := "(not stored)"
		if conn.StorePassword {
			passwordDisplay = conn.Password
		}
		extraArgs := strings.Join(conn.ExtraArgs, ", ")
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n",
			conn.Name, conn.Host, conn.Domain, conn.Username, passwordDisplay, extraArgs)
	}
	w.Flush()

	return waitForEnter()
}

func waitForEnter() error {
	fmt.Println("\nPress Enter to return to the main menu...")
	_, err := fmt.Scanln()
	return err
}
