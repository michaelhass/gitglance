package status

import (
	"fmt"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/michaelhass/gitglance/internal/ui/container"
)

type Diff struct {
	viewport  viewport.Model
	rawDiff   string
	err       error
	isFocused bool
}

func NewDiff() Diff {
	return Diff{}
}

func (d Diff) Init() tea.Cmd {
	return nil
}

func (d Diff) Update(msg tea.Msg) (container.Content, tea.Cmd) {
	if !d.isFocused {
		return d, nil
	}
	var cmd tea.Cmd
	d.viewport, cmd = d.viewport.Update(msg)
	return d, cmd
}

func (d Diff) View() string {
	return d.viewport.View()
}

func (d Diff) Title() string {
	return "Diff"
}

func (d Diff) SetSize(width, height int) container.Content {
	d.viewport = viewport.New(width, height)
	d.viewport.SetContent(d.rawDiff)
	return d
}

func (d Diff) SetIsFocused(isFocused bool) container.Content {
	d.isFocused = isFocused
	return d
}

func (d Diff) SetContent(rawDiff string, err error) Diff {
	d.rawDiff = rawDiff
	d.err = err

	if d.err != nil {
		d.viewport.SetContent(fmt.Sprint("An error occured:", d.err))
	} else {
		d.viewport.SetContent(rawDiff)
	}

	return d
}
