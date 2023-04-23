package popup

import tea "github.com/charmbracelet/bubbletea"

func ShowPopUp(content Content, displayMode DisplayMode) func() tea.Msg {
	return func() tea.Msg {
		return ShowPopUpMsg{
			PopUp: New(content, displayMode),
		}
	}
}

type ShowPopUpMsg struct {
	PopUp Model
}
