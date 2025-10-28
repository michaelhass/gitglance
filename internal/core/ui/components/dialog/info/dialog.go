package info

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/michaelhass/gitglance/internal/core/ui/components/dialog"
)

type DialogContent struct {
	Model

	closeKey key.Binding
}

func NewDialogContent(confirm Model) DialogContent {
	return DialogContent{
		Model: confirm,
		closeKey: key.NewBinding(
			key.WithKeys(tea.KeyEsc.String()),
			key.WithHelp("ESC", "close"),
		),
	}
}

func (dc DialogContent) Init() tea.Cmd {
	return dc.Model.Init()
}

func (dc DialogContent) Update(msg tea.Msg) (dialog.Content, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok && key.Matches(keyMsg, dc.closeKey) {
		return dc, dialog.Close
	}
	m, cmd := dc.Model.Update(msg)
	dc.Model = m
	return dc, cmd
}

func (dc DialogContent) View() string {
	return dc.Model.View()
}

func (dc DialogContent) SetSize(width, height int) dialog.Content {
	dc.Model = dc.Model.SetSize(width, height)
	return dc
}

func (dc DialogContent) Help() []key.Binding {
	return []key.Binding{dc.closeKey}
}
