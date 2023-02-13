package text

import (
	"reflect"
	"testing"
)

func TestCharacterWrapping(t *testing.T) {
	var (
		text    = "0123456789"
		wrapper = NewCharacterWrapper(4)

		expect = "0123\n4567\n89"
		got    string
	)

	wrapper.WriteString(text)
	got = wrapper.String()

	if expect != got {
		t.Errorf("[%s] is not equal to [%s]", expect, got)
	}
}

func TestCharacterWrappingGetLines(t *testing.T) {
	var (
		text    = "0123456789"
		wrapper = NewCharacterWrapper(4)

		expect = []string{"0123", "4567", "89"}
		got    []string
	)

	wrapper.WriteString(text)
	got = wrapper.Lines()

	if !reflect.DeepEqual(expect, got) {
		t.Errorf("%s is not equal to %s", expect, got)
	}
}
