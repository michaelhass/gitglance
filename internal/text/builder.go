package text

import "strings"

type Builder struct {
	wrappedLines []Wrapper
	rawText      string
	lineLength   int
}

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) SetLineLength(lineLength int) {
	b.lineLength = lineLength
}

func (b *Builder) WriteString(s string) {
	b.rawText = s
}

func (b *Builder) WriteLines(wrappedLines []Wrapper) {
	b.wrappedLines = wrappedLines
}

func (b *Builder) String() string {
	var stringBuilder strings.Builder

	if wrappedLines := b.wrappedLines; wrappedLines != nil {
		for i, l := range wrappedLines {
			l.SetLineLength(b.lineLength)
			stringBuilder.WriteString(l.String())
			if i == len(wrappedLines)-1 {
				continue
			}
			stringBuilder.WriteString("\n")
		}
		return stringBuilder.String()
	}

	var (
		wrapper  = NewWordWrapper(b.lineLength)
		rawLines = strings.Split(b.rawText, "\n")
	)

	for i, l := range rawLines {
		wrapper := wrapper
		defer wrapper.Reset()

		wrapper.WriteString(l)
		stringBuilder.WriteString(wrapper.String())
		if i == len(rawLines)-1 {
			continue
		}
		stringBuilder.WriteString("\n")
	}
	return stringBuilder.String()
}
