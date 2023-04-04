package container

import (
	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/michaelhass/gitglance/internal/ui/diff"
	"github.com/michaelhass/gitglance/internal/ui/filelist"
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
	filelist.Model
}

func NewFileListContent(model filelist.Model) FileListContent {
	return FileListContent{Model: model}
}

func (flc FileListContent) Update(msg tea.Msg) (Content, tea.Cmd) {
	model, cmd := flc.Model.Update(msg)
	flc.Model = model
	return flc, cmd
}

func (flc FileListContent) UpdateFocus(isFocused bool) (Content, tea.Cmd) {
	model, cmd := flc.Model.UpdateFocus(isFocused)
	flc.Model = model
	return flc, cmd
}

func (flc FileListContent) SetSize(width, height int) Content {
	flc.Model = flc.Model.SetSize(width, height)
	return flc
}

type DiffContent struct {
	diff.Model
}

func NewDiffContent(model diff.Model) DiffContent {
	return DiffContent{Model: model}
}

func (dc DiffContent) Update(msg tea.Msg) (Content, tea.Cmd) {
	model, cmd := dc.Model.Update(msg)
	dc.Model = model
	return dc, cmd
}

func (dc DiffContent) UpdateFocus(isFocused bool) (Content, tea.Cmd) {
	model, cmd := dc.Model.UpdateFocus(isFocused)
	dc.Model = model
	return dc, cmd
}

func (dc DiffContent) SetSize(width, height int) Content {
	dc.Model = dc.Model.SetSize(width, height)
	return dc
}
