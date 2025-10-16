package textwrap

import (
	"strings"
	"unicode"
)

type WordWrapper struct {
	lineLength int
	runes      []rune
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
	ww.runes = []rune(s)
}

func (ww *WordWrapper) String() string {
	const whiteSpaceLength = 1

	var (
		builder          strings.Builder
		word             []rune
		characterWrapper = NewCharacterWrapper(ww.lineLength)
		lineLength       int
	)

	for _, rune := range ww.runes {
		if !unicode.IsSpace(rune) {
			word = append(word, rune)
			continue
		}

		lineLength = ww.writeNextWord(&builder, characterWrapper, word, lineLength)
		word = nil
		lineLength = ww.writeWhiteSpace(&builder, lineLength)
	}

	_ = ww.writeNextWord(&builder, characterWrapper, word, lineLength)

	if ww.renderer == nil {
		return builder.String()
	}

	return ww.renderer.Render(builder.String())
}

func (ww *WordWrapper) writeWhiteSpace(builder *strings.Builder, lineLength int) int {
	const whiteSpaceLength = 1

	if whiteSpaceLength+lineLength > ww.lineLength {
		builder.WriteString("\n")
		builder.WriteString(" ")
		return whiteSpaceLength
	}

	builder.WriteString(" ")
	return whiteSpaceLength + lineLength
}

func (ww *WordWrapper) writeNextWord(
	builder *strings.Builder,
	characterWrapper *CharacterWrapper,
	word []rune,
	lineLength int,
) int {
	if len(word) <= ww.lineLength {
		return ww.writeWord(builder, word, lineLength)
	}
	return ww.writeLongWord(builder, characterWrapper, word, lineLength)
}

func (ww *WordWrapper) writeWord(builder *strings.Builder, word []rune, lineLength int) int {
	var wordLength = len(word)

	if wordLength+lineLength > ww.lineLength {
		builder.WriteString("\n")
		builder.WriteString(string(word))
		return wordLength
	}

	builder.WriteString(string(word))
	return lineLength + wordLength
}

func (ww *WordWrapper) writeLongWord(
	builder *strings.Builder,
	characterWrapper *CharacterWrapper,
	word []rune,
	lineLength int,
) int {
	characterWrapper.WriteString(string(word))

	for _, line := range characterWrapper.Lines() {
		lineLength = ww.writeWord(builder, []rune(line), lineLength)
	}

	characterWrapper.Reset()
	return lineLength
}

func (ww *WordWrapper) Reset() {
	ww.runes = nil
	ww.renderer = nil
}
