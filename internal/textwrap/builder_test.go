package textwrap

import (
	"testing"
)

func TestBuilderWriteMultineText(t *testing.T) {
	var (
		text    = "01234567\n12 34 567\nend"
		builder = NewBuilder()

		expect = "01234\n567\n12 34\n 567\nend"
		got    string
	)

	builder.SetLineLength(5)
	builder.WriteString(text)
	got = builder.String()

	if expect != got {
		t.Errorf("[%s] is not equal to [%s]", expect, got)
	}
}

func TestBuilderWriteMultiLineWithTab(t *testing.T) {
	var (
		text    = "func() {\n\tfmt.Println('')"
		builder = NewBuilder()

		expect = "func() {\n    fmt.Println('')"
		got    string
	)

	builder.SetLineLength(20)
	builder.WriteString(text)
	got = builder.String()

	if expect != got {
		t.Errorf("[%s] is not equal to [%s]", expect, got)
	}
}
