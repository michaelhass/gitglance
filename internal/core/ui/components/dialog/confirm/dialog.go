package confirm

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/michaelhass/gitglance/internal/core/ui/components/dialog"
)

type DialogContent struct {
	Model
}

func NewDialogContent(confirm Model) DialogContent {
	return DialogContent{Model: confirm}
}

func (dc DialogContent) Init() tea.Cmd {
	return dc.Model.Init()
}

func (dc DialogContent) Update(msg tea.Msg) (dialog.Content, tea.Cmd) {
	switch msg.(type) {
	case confirmExecutedMsg:
		return dc, dialog.Close
	default:
		model, cmd := dc.Model.Update(msg)
		dc.Model = model
		return dc, cmd
	}
}

func (dc DialogContent) View() string {
	return dc.Model.View()
}

func (dc DialogContent) SetSize(width, height int) dialog.Content {
	dc.Model = dc.Model.SetSize(width, height)
	return dc
}
