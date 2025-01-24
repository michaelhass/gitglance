package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/michaelhass/gitglance/internal/ui/logger"
)

type LaunchOptions struct {
	IsDebug bool
}

func Launch(opt LaunchOptions) error {
	logger, err := logger.NewLogger(opt.IsDebug)
	if err != nil {
		return err
	}
	if _, err := tea.NewProgram(newModel(logger), tea.WithAltScreen()).Run(); err != nil {
		return err
	}
	return nil
}
