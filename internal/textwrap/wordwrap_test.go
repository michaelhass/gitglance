package textwrap

import (
	"testing"
)

func TestWordWrapping(t *testing.T) {
	var (
		text    = "012 3 456 789"
		wrapper = NewWordWrapper(5)

		expect = "012 3\n 456 \n789"
		got    string
	)

	wrapper.WriteString(text)
	got = wrapper.String()

	if expect != got {
		t.Errorf("[%s]\n is not equal to\n[%s]", expect, got)
	}
}

func TestWordWrappingLongWordStart(t *testing.T) {
	var (
		text    = "AAAAAB 012 3 456 789"
		wrapper = NewWordWrapper(5)

		expect = "AAAAA\nB 012\n 3 \n456 \n789"
		got    string
	)

	wrapper.WriteString(text)
	got = wrapper.String()

	if expect != got {
		t.Errorf("[%s] is not equal to [%s]", expect, got)
	}
}

func TestWordWrappingLongWordMiddle(t *testing.T) {
	var (
		text    = "012 aaaaab 3 456 789"
		wrapper = NewWordWrapper(5)

		expect = "012 \naaaaa\nb 3 \n456 \n789"
		got    string
	)

	wrapper.WriteString(text)
	got = wrapper.String()

	if expect != got {
		t.Errorf("[%s] is not equal to [%s]", expect, got)
	}
}

func TestLeadingSpaces(t *testing.T) {
	var (
		text    = "    123    1"
		wrapper = NewWordWrapper(5)

		expect = "    \n123  \n  1"
		got    string
	)

	wrapper.WriteString(text)
	got = wrapper.String()

	if expect != got {
		t.Errorf("[%s]\n is not equal to\n[%s]", expect, got)
	}
}
