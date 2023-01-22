package status

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/michaelhass/gitglance/internal/git"
)

type Model struct {
	repo      git.Repository
	status    git.Status
	statusErr error
}

func New(repo git.Repository) Model {
	return Model{repo: repo}
}

func (m Model) Init() tea.Cmd {
	return Load(m.repo)
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case statusUpdateMsg:
		m.status = msg.status
		m.statusErr = msg.err
	}
	return m, nil
}

func (m Model) View() string {
	if m.statusErr != nil {
		return fmt.Sprint(m.statusErr)
	}
	return fmt.Sprintf("%v", m.status)
}

// Cmd

func Load(repo git.Repository) func() tea.Msg {
	return func() tea.Msg {
		var msg statusUpdateMsg

		wt, err := repo.Worktree()
		if err != nil {
			msg.err = err
			return msg
		}

		status, err := wt.Status()
		if err != nil {
			msg.err = err
			return msg
		}
		msg.status = status
		return msg
	}
}

// Msg

type statusUpdateMsg struct {
	err    error
	status git.Status
}
