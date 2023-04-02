package status

import (
	"fmt"
	"strings"

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
	rawDiff     string
	err         error
	width       int
	isReady     bool
	isFocused   bool
}

func NewDiff() Diff {
	return Diff{textBuilder: text.NewBuilder()}
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

	d.textBuilder.SetLineLength(width - 5)
	d = d.SetContent(d.rawDiff, d.err)
	return d
}

func (d Diff) SetContent(rawDiff string, err error) Diff {

	rawDiff = strings.ReplaceAll(rawDiff, "\r", "\n")
	rawDiff = strings.ReplaceAll(rawDiff, "\t", "    ")

	d.err = err
	d.rawDiff = rawDiff

	if !d.isReady {
		return d
	}

	var (
		rawLines = strings.Split(rawDiff, "\n")
		lines    = make([]text.Wrapper, len(rawLines))
	)

	for i, rawLine := range rawLines {
		var wordWrapper = &text.WordWrapper{}
		wordWrapper.WriteString(rawLine)

		if strings.HasPrefix(rawLine, "+") {
			wordWrapper.SetRenderer(addedTextStyle)
		} else if strings.HasPrefix(rawLine, "-") {
			wordWrapper.SetRenderer(removedTextStyle)
		} else {
			wordWrapper.SetRenderer(textStyle)
		}

		lines[i] = wordWrapper
	}

	d.textBuilder.WriteLines(lines)

	if d.err != nil {
		d.viewport.SetContent(fmt.Sprint("An error occured:", d.err))
	} else {

		d.viewport.SetContent(d.textBuilder.String())
	}

	return d
}
