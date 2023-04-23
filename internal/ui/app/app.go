package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/michaelhass/gitglance/internal/ui/popup"
	"github.com/michaelhass/gitglance/internal/ui/status"
)

type Model struct {
	status status.Model

	popUp          popup.Model
	isShowingPopUp bool

	isReady bool

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
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.status = m.status.SetSize(msg.Width, msg.Height)
		if m.isShowingPopUp {
			m.popUp = m.popUp.SetSize(msg.Width, msg.Height)
		}
		m.isReady = true
	case popup.ShowPopUpMsg:
		popUp := msg.PopUp
		popUp = popUp.SetSize(m.width, m.height)
		m.popUp = popUp
		m.isShowingPopUp = true
	}

	if m.isShowingPopUp {
		popUp, cmd := m.popUp.Update(msg)
		m.popUp = popUp
		return m, cmd
	}

	var (
		cmds      []tea.Cmd
		statusCmd tea.Cmd
	)

	m.status, statusCmd = m.status.Update(msg)
	if statusCmd != nil {
		cmds = append(cmds, statusCmd)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if !m.isReady {
		return "loading"
	}

	if m.isShowingPopUp {
		return m.popUp.View()
	}
	return m.status.View()
}
