package label

import "github.com/michaelhass/gitglance/internal/core/textwrap"

type MultiLine struct {
	builder *textwrap.Builder
}

func NewDefaultMultiLine() MultiLine {
	return MultiLine{
		builder: textwrap.NewBuilder(),
	}
}

func (ml MultiLine) SetText(text string) MultiLine {
	ml.builder.WriteString(text)
	return ml
}

func (ml MultiLine) SetWidth(width int) MultiLine {
	ml.builder.SetLineLength(width)
	return ml
}

func (ml MultiLine) View() string {
	return ml.builder.String()
}
