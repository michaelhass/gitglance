// Package app provides the main bubbletea model.
package app

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/michaelhass/gitglance/internal/ui/dialog"
	"github.com/michaelhass/gitglance/internal/ui/logger"
	"github.com/michaelhass/gitglance/internal/ui/status"
)

// model is the main bubbletea model of the application.
// It displays multiple sub models and is responsible for
// displaying dialogs.
type model struct {
	status status.Model // Model to display the git status

	dialog          dialog.Model // A dialog that is shown.
	isDialogShowing bool         // Indicated whether a dialog is showing.

	isReady bool // Indicates if the application is ready / initialized.

	width, height int

	logger logger.Logger
}

func newModel(logger logger.Logger) model {
	return model{
		status: status.New(),
		logger: logger,
	}
}

func (m model) Init() tea.Cmd {
	return m.status.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	m.logger.Println(time.Now(), msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.status = m.status.SetSize(msg.Width, msg.Height)
		if m.isDialogShowing {
			m.dialog = m.dialog.SetSize(msg.Width, msg.Height)
		}
		m.isReady = true
	case dialog.ShowMsg:
		dialog := msg.Dialog
		dialog = dialog.SetSize(m.width, m.height)
		m.dialog = dialog
		m.isDialogShowing = true
	case dialog.CloseMsg:
		m.isDialogShowing = false
		cmds = append(cmds, m.dialog.OnCloseCmd())
	}

	if m.isDialogShowing {
		dialog, cmd := m.dialog.Update(msg)
		m.dialog = dialog
		return m, cmd
	}

	status, statusCmd := m.status.Update(msg)
	m.status = status
	cmds = append(cmds, statusCmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if !m.isReady {
		return "loading"
	}

	if m.isDialogShowing {
		return m.dialog.View()
	}
	return m.status.View()
}
