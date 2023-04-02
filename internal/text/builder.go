package text

import "strings"

type LineRenderer func(line string) Renderer

type Builder struct {
	rawText      string
	lines        []string
	lineRenderer LineRenderer
	lineLength   int
}

func NewBuilder() *Builder {
	return &Builder{
		lineRenderer: defaultLineRenderer(),
	}
}

func (b *Builder) SetLineLength(lineLength int) {
	b.lineLength = lineLength
}

func (b *Builder) SetLineRenderer(handler LineRenderer) {
	b.lineRenderer = handler
}

func (b *Builder) WriteString(s string) {
	b.rawText = s
	b.lines = strings.Split(normalizedText(s), "\n")
}

func (b *Builder) RawString() string {
	return b.rawText
}

func (b *Builder) String() string {
	var (
		stringBuilder strings.Builder
		wrapper       = NewWordWrapper(b.lineLength)
	)

	for i, l := range b.lines {
		wrapper := wrapper
		defer wrapper.Reset()

		if renderer := b.lineRenderer(l); renderer != nil {
			wrapper.SetRenderer(b.lineRenderer(l))
		}

		wrapper.WriteString(l)

		stringBuilder.WriteString(wrapper.String())
		if i == len(b.lines)-1 {
			continue
		}
		stringBuilder.WriteString("\n")
	}
	return stringBuilder.String()
}

func defaultLineRenderer() LineRenderer {
	return func(line string) Renderer { return &Passthrough{} }
}

func normalizedText(rawText string) string {
	out := strings.ReplaceAll(rawText, "\r", "\n")
	return strings.ReplaceAll(out, "\t", "    ")
}
