// Package commit provides ui to perform a git commit.
package commit

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/michaelhass/gitglance/internal/git"
)

// Execute creates a tea.Cmd to execute a git commit.
func Execute(msg string) tea.Cmd {
	return func() tea.Msg {
		return ExecutedMsg{Err: git.Commit(msg)}
	}
}

// ExecutedMsg is the message to be sent after we performed a git commit.
type ExecutedMsg struct {
	Err error
}

func loadMergeMsg() tea.Msg {
	msg, err := git.MergeMsg()
	if err != nil {
		msg = ""
	}
	return MergeMsgLoaded{msg: msg}
}

type MergeMsgLoaded struct {
	msg string
}
