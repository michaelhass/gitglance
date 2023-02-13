package status

import (
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/michaelhass/gitglance/internal/ui/container"
)

type CappedText struct {
	lines []CappedLine
	Limit int
}

func (ct *CappedText) SetString(value string) {
	var lines = strings.Split(value, "\n")
	var cappedLines []CappedLine
	for _, line := range lines {
		cappedLine := newCappedLine(ct.Limit)
		cappedLine.setString(line)
		cappedLines = append(cappedLines, *cappedLine)
	}
	ct.lines = cappedLines
}

func (ct *CappedText) String() string {
	var builder strings.Builder
	for _, line := range ct.lines {
		builder.WriteString(line.String())
		builder.WriteString("\n")
	}
	return builder.String()
}

type CappedLine struct {
	value string
	Limit int
}

func newCappedLine(limit int) *CappedLine {
	return &CappedLine{Limit: limit}
}

func (cl *CappedLine) setString(line string) {
	// var builder strings.Builder
	// for i, r := range line {
	// 	//builder.WriteString(fmt.Sprint(i))
	// 	if i == 0 || i%cl.Limit != 0 {
	// 		builder.WriteRune(r)
	// 		// continue
	// 	}
	// 	builder.WriteString(fmt.Sprintf("%d :: %d \n", i, cl.Limit))
	// }
	runes := []rune(line)
	if cl.Limit < len(runes)-1 {
		runes = runes[:cl.Limit-1]
	}

	cl.value = string(runes)
}

func (cl *CappedLine) String() string {
	return cl.value
}

type Diff struct {
	viewport   viewport.Model
	cappedText CappedText
	rawDiff    string
	err        error
	width      int
	isReady    bool
	isFocused  bool
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
		//	d.writer = wordwrap.NewWriter(width)
		d = d.SetContent(d.rawDiff, d.err)
		//d.viewport.HighPerformanceRendering = true
	} else {
		d.viewport.Width = width
		d.viewport.Height = height
		d = d.SetContent(d.rawDiff, d.err)
	}
	return d
}

func (d Diff) SetIsFocused(isFocused bool) container.Content {
	d.isFocused = isFocused
	return d
}

func (d Diff) SetContent(rawDiff string, err error) Diff {
	d.rawDiff = rawDiff
	d.err = err

	cappedText := CappedText{Limit: d.width}
	cappedText.SetString(rawDiff)

	// if d.err != nil {
	// 	d.viewport.SetContent(fmt.Sprint("An error occured:", d.err))
	// } else {
	d.viewport.SetContent(cappedText.String())
	// }

	return d
}
