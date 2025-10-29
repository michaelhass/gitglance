// Package app provides the main bubbletea model.
package app

import (
	"reflect"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/michaelhass/gitglance/internal/core/exit"
	"github.com/michaelhass/gitglance/internal/core/logger"
	"github.com/michaelhass/gitglance/internal/core/refresh"
	"github.com/michaelhass/gitglance/internal/core/ui/components/dialog"
	"github.com/michaelhass/gitglance/internal/page/status"
)

// refreshInterval is the duration when a refresh.Msg is send.
// This message can be used in any view that wants to update its content periodically.
const refreshInterval time.Duration = time.Second * 15

// model is the main bubbletea model of the application.
// It displays multiple sub models and is responsible for
// displaying dialogs.
type model struct {
	// Model to display the git status
	status status.Model
	// dialgs currently on the presentation stack.
	// Only the last dialog will receive messages and will be rendered
	dialogs []dialog.Model

	isReady       bool
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
		if m.isDialogShowing() {
			for i, dialog := range m.dialogs {
				m.dialogs[i] = dialog.SetSize(msg.Width, msg.Height)
			}
		}
		m.isReady = true
	case dialog.ShowMsg:
		dialog := msg.Dialog
		dialog = dialog.SetSize(m.width, m.height)
		m.dialogs = append(m.dialogs, dialog)
		cmds = append(cmds, dialog.Init())
	case dialog.CloseMsg:
		if d, ok := m.topDialog(); ok {
			_, cmd := d.Update(msg)
			m = m.removedTopDialog()
			cmds = append(cmds, cmd)
		}
	case refresh.Msg:
		cmds = append(cmds, refresh.Schedule(refreshInterval))
	}

	if m.isDialogShowing() {
		updatedModel, cmd := m.updatedTopDialog(msg)
		m = updatedModel
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
	if d, ok := m.topDialog(); ok {
		return d.View()
	}
	return m.status.View()
}

func (m model) isDialogShowing() bool {
	return len(m.dialogs) > 0
}

func (m model) topDialog() (dialog.Model, bool) {
	if !m.isDialogShowing() {
		return dialog.Model{}, false
	}
	return m.dialogs[len(m.dialogs)-1], true
}

func (m model) removedTopDialog() model {
	if !m.isDialogShowing() {
		return m
	}
	m.dialogs = m.dialogs[:len(m.dialogs)-1]
	return m
}

func (m model) updatedTopDialog(msg tea.Msg) (model, tea.Cmd) {
	d, ok := m.topDialog()
	if !ok {
		return m, nil
	}
	d, cmd := d.Update(msg)
	m.dialogs[len(m.dialogs)-1] = d
	return m, cmd
}
