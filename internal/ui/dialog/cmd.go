package dialog

import tea "github.com/charmbracelet/bubbletea"

func Show(content Content, displayMode DisplayMode) func() tea.Msg {
	return func() tea.Msg {
		return ShowMsg{
			PopUp: New(content, displayMode),
		}
	}
}

type ShowMsg struct {
	PopUp Model
}
