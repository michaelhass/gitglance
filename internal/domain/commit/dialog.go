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
	model, cmd := dc.Model.Update(msg)
	dc.Model = model

	return dc, cmd
}

func (dc DialogContent) View() string {
	return dc.Model.View()
}

func (dc DialogContent) SetSize(width, height int) dialog.Content {
	dc.Model = dc.Model.SetSize(width, height)
	return dc
}
