package commit

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/michaelhass/gitglance/internal/git"
)

func Execute(msg string) tea.Cmd {
	return func() tea.Msg {
		return ExecutedMsg{Err: git.Commit(msg)}
	}
}

type ExecutedMsg struct {
	Err error
}
