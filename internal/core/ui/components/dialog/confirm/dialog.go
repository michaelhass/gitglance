package confirm

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/michaelhass/gitglance/internal/core/ui/components/dialog"
)

type DialogContent struct {
	Model
	errHandler func(tea.Msg) tea.Cmd
}

func NewDialogContent(confirm Model) DialogContent {
	return DialogContent{Model: confirm}
}

func (dc DialogContent) WithErrHandler(errHandler func(tea.Msg) tea.Cmd) DialogContent {
	dc.errHandler = errHandler
	return dc
}

func (dc DialogContent) Init() tea.Cmd {
	return dc.Model.Init()
}

func (dc DialogContent) Update(msg tea.Msg) (dialog.Content, tea.Cmd) {
	switch msg := msg.(type) {
	case confirmExecutedMsg:
		if !msg.isSuccess() && dc.errHandler != nil {
			return dc, dc.errHandler(msg.errMsg)
		}
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
