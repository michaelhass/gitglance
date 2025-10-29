package confirm

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/michaelhass/gitglance/internal/core/err"
)

type confirmExecutedMsg struct {
	errMsg err.Msg
}

func (c confirmExecutedMsg) isSuccess() bool {
	return c.errMsg == nil || (c.errMsg != nil && c.errMsg.Err() == nil)
}

func executeConfirmCmd(confirmCmd tea.Cmd) tea.Cmd {
	msg := confirmCmd()
	executedMsg := confirmExecutedMsg{}

	if errMsg, ok := msg.(err.Msg); ok && errMsg.Err() != nil {
		executedMsg.errMsg = errMsg
	}

	return tea.Sequence(
		func() tea.Msg { return msg },
		func() tea.Msg { return executedMsg },
	)
}
