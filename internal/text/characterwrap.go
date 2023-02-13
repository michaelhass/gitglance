package text

import "strings"

type CharacterWrapper struct {
	lineLength int
	rawValue   string
}

func NewCharacterWrapper(lineLength int) *CharacterWrapper {
	return &CharacterWrapper{lineLength: lineLength}
}

func (cw *CharacterWrapper) SetLineLength(lineLength int) {
	cw.lineLength = lineLength
}

func (cw *CharacterWrapper) WriteString(s string) {
	cw.rawValue = s
}

func (cw *CharacterWrapper) Lines() []string {
	var (
		builder strings.Builder
		lines   []string
	)
	for i, r := range cw.rawValue {
		if i > 0 && i%cw.lineLength == 0 {
			lines = append(lines, builder.String())
			builder.Reset()
		}
		builder.WriteRune(r)
	}

	lines = append(lines, builder.String())
	return lines
}

func (cw *CharacterWrapper) String() string {
	var builder strings.Builder
	for i, r := range cw.rawValue {
		if i > 0 && i%cw.lineLength == 0 {
			builder.WriteString("\n")
		}
		builder.WriteRune(r)
	}
	return builder.String()
}

func (cw *CharacterWrapper) Reset() {
	cw.rawValue = ""
}
