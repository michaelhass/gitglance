package commit

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/michaelhass/gitglance/internal/ui/dialog"
)

// Content is a wrapper to use the commit ui as dialog.Content.
type Content struct {
	Model
}

func NewContent(commit Model) Content {
	return Content{
		Model: commit,
	}
}

func (c Content) Init() tea.Cmd {
	return c.Model.Init()
}

func (c Content) Update(msg tea.Msg) (dialog.Content, tea.Cmd) {
	if _, ok := msg.(ExecutedMsg); ok {
		return c, dialog.Close()
	}

	model, cmd := c.Model.Update(msg)
	c.Model = model

	return c, cmd
}

func (c Content) View() string {
	return c.Model.View()
}

func (c Content) SetSize(width, height int) dialog.Content {
	c.Model = c.Model.SetSize(width, height)
	return c
}
