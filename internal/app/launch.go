package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/michaelhass/gitglance/internal/core/logger"
)

type Option func(opts *options)

type options struct {
	logger logger.Logger
}

func newOptions() *options {
	return &options{
		logger: logger.NewEmptyLogger(),
	}
}

func WithLogger(logger logger.Logger) Option {
	return func(opts *options) {
		if logger != nil {
			opts.logger = logger
		}
	}
}

func WithDebugLogger() Option {
	return func(opts *options) {
		debugLogger, err := logger.NewDebugLogger()
		if err == nil {
			opts.logger = debugLogger
		}
	}
}

func Launch(opts ...Option) error {
	appOpts := newOptions()
	for _, opt := range opts {
		opt(appOpts)
	}

	defer appOpts.logger.Close()
	if _, err := tea.NewProgram(newModel(appOpts.logger), tea.WithAltScreen()).Run(); err != nil {
		return err
	}
	return nil
}
