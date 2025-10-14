package diff

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/michaelhass/gitglance/internal/ui/container"
)

// ContainerContent is a wrapper to use the commit ui as container.ContainerContent.
type ContainerContent struct {
	Model
}

func NewContent(model Model) ContainerContent {
	return ContainerContent{Model: model}
}

func (c ContainerContent) Update(msg tea.Msg) (container.Content, tea.Cmd) {
	model, cmd := c.Model.Update(msg)
	c.Model = model
	return c, cmd
}

func (c ContainerContent) UpdateFocus(isFocused bool) (container.Content, tea.Cmd) {
	model, cmd := c.Model.UpdateFocus(isFocused)
	c.Model = model
	return c, cmd
}

func (c ContainerContent) SetSize(width, height int) container.Content {
	c.Model = c.Model.SetSize(width, height)
	return c
}
