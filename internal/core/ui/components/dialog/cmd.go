package dialog

import tea "github.com/charmbracelet/bubbletea"

// Show creates a tea.Cmd to show a dialog with the given content.
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

func Close() tea.Msg {
	return CloseMsg{}
}

type CloseMsg struct{}
