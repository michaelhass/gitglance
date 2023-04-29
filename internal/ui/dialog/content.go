package dialog

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/michaelhass/gitglance/internal/ui/commit"
)

type Content interface {
	Init() tea.Cmd
	Update(msg tea.Msg) (Content, tea.Cmd)
	View() string
	SetSize(width, height int) Content
}

type CommitContent struct {
	commit.Model
}

func NewCommitContent(commit commit.Model) CommitContent {
	return CommitContent{
		Model: commit,
	}
}

func (cc CommitContent) Init() tea.Cmd {
	return cc.Model.Init()
}

func (cc CommitContent) Update(msg tea.Msg) (Content, tea.Cmd) {
	model, cmd := cc.Model.Update(msg)
	cc.Model = model
	return cc, cmd
}

func (cc CommitContent) View() string {
	return cc.Model.View()
}

func (cc CommitContent) SetSize(width, height int) Content {
	cc.Model = cc.Model.SetSize(width, height)
	return cc
}
