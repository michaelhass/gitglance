package text

import (
	"testing"
)

func TestBuilderWriteMultineText(t *testing.T) {
	var (
		text    = "01234567\n12 34 567\nend"
		builder = NewBuilder()

		expect = "01234\n567\n12 34\n567\nend"
		got    string
	)

	builder.SetLineLength(5)
	builder.WriteString(text)
	got = builder.String()

	if expect != got {
		t.Errorf("[%s] is not equal to [%s]", expect, got)
	}
}
