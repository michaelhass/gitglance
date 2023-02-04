package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/michaelhass/gitglance/internal/git"
	"github.com/michaelhass/gitglance/internal/ui/status"
)

type Model struct {
	repo   *git.Repository
	status status.Model
}

func New(repo *git.Repository) Model {
	return Model{
		repo:   repo,
		status: status.New(repo),
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
		m.status = m.status.SetSize(msg.Width, msg.Height)
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
	return m.status.View()
}
