package text

type Wrapper interface {
	SetLineLength(lineLength int)
	WriteString(s string)
	String() string
	Reset()
}
