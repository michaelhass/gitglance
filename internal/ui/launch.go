package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/michaelhass/gitglance/internal/ui/app"
)

type LaunchOptions struct{}

func LaunchApp(opt LaunchOptions) error {
	model := app.New()
	if _, err := tea.NewProgram(model, tea.WithAltScreen()).Run(); err != nil {
		return err
	}
	return nil
}
