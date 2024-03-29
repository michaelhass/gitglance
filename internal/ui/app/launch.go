package app

import (
	tea "github.com/charmbracelet/bubbletea"
)

type LaunchOptions struct{}

func Launch(opt LaunchOptions) error {
	if _, err := tea.NewProgram(newModel(), tea.WithAltScreen()).Run(); err != nil {
		return err
	}
	return nil
}
