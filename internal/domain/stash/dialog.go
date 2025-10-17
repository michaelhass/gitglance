package stash

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/michaelhass/gitglance/internal/core/ui/components/dialog"
	"github.com/michaelhass/gitglance/internal/core/ui/style"
)

const (
	titleHeight       = 1
	borderPadding int = 1
)

var (
	titleStyle  = style.Title.Height(titleHeight)
	borderStyle = style.FocusBorder.PaddingLeft(borderPadding).PaddingRight(borderPadding)
)

func NewApplyDialogConent(stashList StashList) ApplyDialogContent {
	stashList.listModel, _ = stashList.listModel.UpdateFocus(true)
	return ApplyDialogContent{StashList: stashList}
}

func (dc ApplyDialogContent) Init() tea.Cmd {
	return dc.StashList.Init()
}

func (dc ApplyDialogContent) Update(msg tea.Msg) (dialog.Content, tea.Cmd) {
	model, cmd := dc.StashList.Update(msg)
	dc.StashList = model
	return dc, cmd
}

func (dc ApplyDialogContent) View() string {
	title := titleStyle.Render(dc.StashList.Title())
	return borderStyle.
		MaxHeight(dc.height).
		MaxWidth(dc.width).
		Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				title,
				"",
				dc.StashList.View(),
			),
		)
}

func (dc ApplyDialogContent) SetSize(width, height int) dialog.Content {
	dc.width, dc.height = width, height
	maxContentHeight := height - titleHeight - 1 - borderPadding*2
	maxContentWidth := width - borderPadding*2
	dc.StashList = dc.StashList.SetSize(maxContentWidth, maxContentHeight)
	return dc
}

func (dc ApplyDialogContent) Help() []key.Binding {
	return dc.listModel.KeyMap().ShortHelp()
}
