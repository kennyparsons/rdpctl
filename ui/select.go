package ui

import (
	"fmt"
	"strings"

	"github.com/manifoldco/promptui"

	"rdpctl/model"
)

// SelectConnection prompts the user to select a connection from the vault.
// It returns the selected connection and an error if any.
func SelectConnection(v *model.Vault) (*model.Connection, error) {
	if len(v.Connections) == 0 {
		return nil, fmt.Errorf("no connections available. Please add a new host first.")
	}

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "\U000027A4 {{ .Name | green }} ({{ .Host }})",
		Inactive: "  {{ .Name | faint }} ({{ .Host | faint }})",
		Selected: "{{ .Name | green }} ({{ .Host }})",
	}

	searcher := func(input string, index int) bool {
		connection := v.Connections[index]
		name := strings.ToLower(connection.Name)
		host := strings.ToLower(connection.Host)
		input = strings.ToLower(input)
		return strings.Contains(name, input) || strings.Contains(host, input)
	}

	prompt := promptui.Select{
		Label:     "Select connection",
		Items:     v.Connections,
		Templates: templates,
		Size:      10,
		Searcher:  searcher,
	}

	i, _, err := prompt.Run()

	if err != nil {
		return nil, fmt.Errorf("connection selection failed %w", err)
	}

	return &v.Connections[i], nil
}
