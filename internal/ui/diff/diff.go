package diff

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/michaelhass/gitglance/internal/textwrap"
	"github.com/michaelhass/gitglance/internal/ui/styles"
)

var (
	normalTextStyle  = styles.TextSyle.Copy()
	addedTextStyle   = styles.AddedTextStyle.Copy()
	removedTextStyle = styles.RemovedTextStyle.Copy()
)

type Model struct {
	viewport    viewport.Model
	textBuilder *textwrap.Builder
	keys        KeyMap
	err         error
	width       int
	isReady     bool
	isFocused   bool
}

func New() Model {
	lineRenderer := func(line string) textwrap.Renderer {
		if strings.HasPrefix(line, "+") {
			return addedTextStyle
		} else if strings.HasPrefix(line, "-") {
			return removedTextStyle
		} else {
			return normalTextStyle
		}
	}

	textBuilder := textwrap.NewBuilder()
	textBuilder.SetLineRenderer(lineRenderer)

	return Model{textBuilder: textBuilder, keys: newDiffKeyMap()}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if !m.isFocused {
		return m, nil
	}

	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m Model) UpdateFocus(isFocused bool) (Model, tea.Cmd) {
	m.isFocused = isFocused
	return m, nil
}

func (m Model) View() string {
	if !m.isReady {
		return ""
	}

	return m.viewport.View()
}

func (m Model) Title() string {
	return "Diff"
}

func (m Model) SetSize(width, height int) Model {
	m.width = width

	if !m.isReady {
		m.isReady = true
		m.viewport = viewport.New(width, height)
	} else {
		m.viewport.Width = width
		m.viewport.Height = height
	}

	// TODO: Fix need for extra padding
	extraPadding := 5
	m.textBuilder.SetLineLength(width - extraPadding)
	m = m.SetContent(m.textBuilder.RawString(), m.err)
	return m
}

func (m Model) KeyMap() help.KeyMap {
	return m.keys
}

func (m Model) SetContent(rawDiff string, err error) Model {
	m.err = err
	m.textBuilder.WriteString(rawDiff)

	if !m.isReady {
		return m
	}

	if m.err != nil {
		m.viewport.SetContent(fmt.Sprint("An error occured:", m.err))
	} else {

		m.viewport.SetContent(m.textBuilder.String())
	}

	return m
}
