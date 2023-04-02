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

// func TestBuilderWriteLines(t *testing.T) {
// 	var (
// 		builder = NewBuilder()

// 		rawLines     = []string{"01234567", "12 34 567", "end"}
// 		wrappedLines = make([]Wrapper, len(rawLines))
// 		lineLength   = 5

// 		expect = "01234\n567\n12 34\n567\nend"
// 		got    string
// 	)

// 	for i, line := range rawLines {
// 		wrapper := NewWordWrapper(lineLength)
// 		wrapper.WriteString(line)
// 		wrappedLines[i] = wrapper
// 	}

// 	builder.SetLineLength(lineLength)
// 	builder.WriteLines(wrappedLines)
// 	got = builder.String()

// 	if expect != got {
// 		t.Errorf("[%s] is not equal to [%s]", expect, got)
// 	}
// }
