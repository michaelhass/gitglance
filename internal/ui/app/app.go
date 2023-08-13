// Package app provides the main bubbletea model.
package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/michaelhass/gitglance/internal/ui/dialog"
	"github.com/michaelhass/gitglance/internal/ui/status"
)

// Model is the main bubbletea model of the application.
// It displays multiple submodules and is responsible for
// displaying dialogs.
type Model struct {
	status status.Model // Model to display the git status

	dialog          dialog.Model // A dialog that is shown.
	isDialogShowing bool         // Indicated whether a dialog is showing.

	isReady bool // Indicates if the application is ready / initialized.

	width, height int
}

func New() Model {
	return Model{
		status: status.New(),
	}
}

func (m Model) Init() tea.Cmd {
	return m.status.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

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

func (m Model) View() string {
	if !m.isReady {
		return "loading"
	}

	if m.isDialogShowing {
		return m.dialog.View()
	}
	return m.status.View()
}
