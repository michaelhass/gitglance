package textinput

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/michaelhass/gitglance/internal/ui/container"
)

type ContainerContent struct {
	Model
}

func NewContainerContent(model Model) ContainerContent {
	return ContainerContent{Model: model}
}

func NewContainer(model Model) container.Model {
	containerContent := NewContainerContent(model)
	return container.New(containerContent)
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
