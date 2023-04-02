package text

type Wrapper interface {
	SetLineLength(lineLength int)
	SetRenderer(renderer Renderer)
	WriteString(s string)
	String() string
	Reset()
}
