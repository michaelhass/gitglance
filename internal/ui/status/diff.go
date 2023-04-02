package status

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/michaelhass/gitglance/internal/text"
	"github.com/michaelhass/gitglance/internal/ui/container"
	"github.com/michaelhass/gitglance/internal/ui/styles"
)

var (
	textStyle        = styles.TextSyle.Copy()
	addedTextStyle   = styles.AddedTextStyle.Copy()
	removedTextStyle = styles.RemovedTextStyle.Copy()
)

type Diff struct {
	viewport    viewport.Model
	textBuilder *text.Builder
	keys        diffKeyMap
	err         error
	width       int
	isReady     bool
	isFocused   bool
}

func NewDiff() Diff {
	lineRenderer := func(line string) text.Renderer {
		if strings.HasPrefix(line, "+") {
			return addedTextStyle
		} else if strings.HasPrefix(line, "-") {
			return removedTextStyle
		} else {
			return textStyle
		}
	}

	textBuilder := text.NewBuilder()
	textBuilder.SetLineRenderer(lineRenderer)

	return Diff{textBuilder: textBuilder, keys: newDiffKeyMap()}
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

func (d Diff) UpdateFocus(isFocused bool) (container.Content, tea.Cmd) {
	d.isFocused = isFocused
	return d, nil
}

func (d Diff) View() string {
	if !d.isReady {
		return ""
	}

	return d.viewport.View()
}

func (d Diff) Title() string {
	return "Diff"
}

func (d Diff) SetSize(width, height int) container.Content {
	d.width = width

	if !d.isReady {
		d.isReady = true
		d.viewport = viewport.New(width, height)
	} else {
		d.viewport.Width = width
		d.viewport.Height = height
	}

	// TODO: Fix need for extra padding
	extraPadding := 5
	d.textBuilder.SetLineLength(width - extraPadding)
	d = d.SetContent(d.textBuilder.RawString(), d.err)
	return d
}

func (d Diff) KeyMap() help.KeyMap {
	return d.keys
}

func (d Diff) SetContent(rawDiff string, err error) Diff {
	d.err = err
	d.textBuilder.WriteString(rawDiff)

	if !d.isReady {
		return d
	}

	if d.err != nil {
		d.viewport.SetContent(fmt.Sprint("An error occured:", d.err))
	} else {

		d.viewport.SetContent(d.textBuilder.String())
	}

	return d
}
