package exit

import tea "github.com/charmbracelet/bubbletea"

func WithMsg(msg string) tea.Cmd {
	return func() tea.Msg {
		return Msg(msg)
	}
}

type Msg string
