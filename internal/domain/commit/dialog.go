package commit

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/michaelhass/gitglance/internal/core/ui/components/dialog"
)

// DialogContent is a wrapper to use the commit ui as dialog.DialogContent.
type DialogContent struct {
	Model
}

func NewContent(commit Model) DialogContent {
	return DialogContent{
		Model: commit,
	}
}

func (dc DialogContent) Init() tea.Cmd {
	return dc.Model.Init()
}

func (dc DialogContent) Update(msg tea.Msg) (dialog.Content, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case ExecutedMsg:
		cmds = append(cmds, dialog.Close)
	default:
		model, cmd := dc.Model.Update(msg)
		dc.Model = model
		cmds = append(cmds, cmd)
	}
	return dc, tea.Batch(cmds...)
}

func (dc DialogContent) View() string {
	return dc.Model.View()
}

func (dc DialogContent) SetSize(width, height int) dialog.Content {
	dc.Model = dc.Model.SetSize(width, height)
	return dc
}
