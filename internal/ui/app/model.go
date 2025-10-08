// Package app provides the main bubbletea model.
package app

import (
	"reflect"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/michaelhass/gitglance/internal/ui/dialog"
	"github.com/michaelhass/gitglance/internal/ui/exit"
	"github.com/michaelhass/gitglance/internal/ui/logger"
	"github.com/michaelhass/gitglance/internal/ui/refresh"
	"github.com/michaelhass/gitglance/internal/ui/status"
)

// refreshInterval is the duration when a refresh.Msg is send.
// This message can be used in any view that wants to update its content periodically.
const refreshInterval time.Duration = time.Second * 15

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
	return tea.Sequence(
		m.status.Init(),
		refresh.Schedule(refreshInterval),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	m.logger.Println(reflect.TypeOf(msg))

	switch msg := msg.(type) {
	case exit.Msg:
		return m, tea.Sequence(tea.ExitAltScreen, tea.Println(msg), tea.Quit)
	case tea.KeyMsg:
		m.logger.Println("Key pressed:", msg.String())
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
		cmds = append(cmds, dialog.Init())
	case dialog.CloseMsg:
		m.isDialogShowing = false
		_, cmd := m.dialog.Update(msg)
		cmds = append(cmds, cmd)
	case refresh.Msg:
		cmds = append(cmds, refresh.Schedule(refreshInterval))
	}

	if m.isDialogShowing {
		dialog, cmd := m.dialog.Update(msg)
		m.dialog = dialog
		cmds = append(cmds, cmd)
		return m, tea.Batch(cmds...)
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
