package text

import "strings"

type WordWrapper struct {
	lineLength int
	words      []string
	renderer   Renderer
}

func NewWordWrapper(lineLength int) *WordWrapper {
	return &WordWrapper{lineLength: lineLength}
}

func (ww *WordWrapper) SetRenderer(renderer Renderer) {
	ww.renderer = renderer
}

func (ww *WordWrapper) SetLineLength(lineLength int) {
	ww.lineLength = lineLength
}

func (ww *WordWrapper) WriteString(s string) {
	ww.words = strings.Split(s, " ")
}

func (ww *WordWrapper) String() string {
	var (
		builder          strings.Builder
		characterWrapper = NewCharacterWrapper(ww.lineLength)
		lineLength       int
	)

	for _, word := range ww.words {
		runes := []rune(word)
		if len(runes) < ww.lineLength {
			lineLength = ww.writeWord(&builder, runes, lineLength)
			continue
		}

		characterWrapper.WriteString(word)
		wrappedLines := characterWrapper.Lines()
		characterWrapper.Reset()

		for _, line := range wrappedLines {
			runes = []rune(line)
			lineLength = ww.writeWord(&builder, runes, lineLength)
		}
	}

	if ww.renderer == nil {
		return builder.String()
	}

	return ww.renderer.Render(builder.String())
}

func (ww *WordWrapper) writeWord(builder *strings.Builder, word []rune, lineLength int) int {
	const whiteSpaceLength int = 1
	var wordLength = len(word)

	if lineLength == 0 {
		builder.WriteString(string(word))
		return wordLength
	}

	if wordLength+lineLength > ww.lineLength {
		builder.WriteString("\n")
		builder.WriteString(string(word))
		return wordLength
	}

	builder.WriteString(" ")
	builder.WriteString(string(word))
	return lineLength + wordLength + whiteSpaceLength
}

func (ww *WordWrapper) Reset() {
	ww.words = nil
}
