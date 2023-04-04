package container

import (
	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/michaelhass/gitglance/internal/ui/file"
)

type Content interface {
	Init() tea.Cmd
	Update(msg tea.Msg) (Content, tea.Cmd)
	UpdateFocus(isFocused bool) (Content, tea.Cmd)
	View() string
	Title() string
	SetSize(width, height int) Content
	KeyMap() help.KeyMap
}

type FileListContent struct {
	file.List
}

func NewFileListContent(list file.List) FileListContent {
	return FileListContent{List: list}
}

func (flc FileListContent) Update(msg tea.Msg) (Content, tea.Cmd) {
	list, cmd := flc.List.Update(msg)
	flc.List = list
	return flc, cmd
}

func (flc FileListContent) UpdateFocus(isFocused bool) (Content, tea.Cmd) {
	list, cmd := flc.List.UpdateFocus(isFocused)
	flc.List = list
	return flc, cmd
}

func (flc FileListContent) SetSize(width, height int) Content {
	flc.List = flc.List.SetSize(width, height)
	return flc
}
