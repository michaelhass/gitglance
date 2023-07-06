package filelist

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/michaelhass/gitglance/internal/ui/container"
)

type Content struct {
	Model
}

func NewContent(model Model) Content {
	return Content{Model: model}
}

func (c Content) Update(msg tea.Msg) (container.Content, tea.Cmd) {
	model, cmd := c.Model.Update(msg)
	c.Model = model
	return c, cmd
}

func (c Content) UpdateFocus(isFocused bool) (container.Content, tea.Cmd) {
	model, cmd := c.Model.UpdateFocus(isFocused)
	c.Model = model
	return c, cmd
}

func (c Content) SetSize(width, height int) container.Content {
	c.Model = c.Model.SetSize(width, height)
	return c
}
