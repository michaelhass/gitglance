package stash

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/michaelhass/gitglance/internal/core/ui/components/dialog"
	"github.com/michaelhass/gitglance/internal/core/ui/style"
)

const (
	titleHeight   = 1
	borderPadding = 1
	borderWidth   = 1
)

var (
	titleStyle  = style.Title.Height(titleHeight)
	borderStyle = style.FocusBorder.PaddingLeft(borderPadding).PaddingRight(borderPadding)
)

func NewApplyDialogConent(stashList ListModel) ListDialogContent {
	stashList.listModel, _ = stashList.listModel.UpdateFocus(true)
	return ListDialogContent{ListModel: stashList}
}

func (dc ListDialogContent) Init() tea.Cmd {
	return dc.ListModel.Init()
}

func (dc ListDialogContent) Update(msg tea.Msg) (dialog.Content, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case EntryCmdExecuted:
		cmds = append(cmds, dialog.Close)
	default:
		model, cmd := dc.ListModel.Update(msg)
		dc.ListModel = model
		cmds = append(cmds, cmd)
	}
	return dc, tea.Batch(cmds...)
}

func (dc ListDialogContent) View() string {
	title := titleStyle.Render(dc.ListModel.Title())

	return borderStyle.
		MaxHeight(dc.height).
		MaxWidth(dc.width).
		Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				title,
				"",
				dc.ListModel.View(),
			),
		)
}

func (dc ListDialogContent) SetSize(width, height int) dialog.Content {
	dc.width, dc.height = width, height

	maxContentHeight := height - titleHeight - 1 - borderPadding*2 - borderWidth*2
	maxContentWidth := width - borderPadding*2 - borderWidth*2
	dc.ListModel = dc.ListModel.SetSize(maxContentWidth, maxContentHeight)
	return dc
}

func (dc ListDialogContent) Help() []key.Binding {
	return dc.listModel.KeyMap().ShortHelp()
}
