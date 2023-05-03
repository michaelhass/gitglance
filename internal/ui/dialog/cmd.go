package dialog

import tea "github.com/charmbracelet/bubbletea"

func Show(content Content, onCloseCmd tea.Cmd, displayMode DisplayMode) tea.Cmd {
	return func() tea.Msg {
		return ShowMsg{
			Dialog: New(content, onCloseCmd, displayMode),
		}
	}
}

type ShowMsg struct {
	Dialog Model
}

func Close() tea.Cmd {
	return func() tea.Msg {
		return CloseMsg{}
	}
}

type CloseMsg struct{}
