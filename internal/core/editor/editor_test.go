package editor

import (
	"reflect"
	"testing"
)

func TestNewCmdFromString(t *testing.T) {
	t.Run("No args", testNewCmdFromStringFunc("nano", false, "nano", []string{}))
	t.Run("Whitespace", testNewCmdFromStringFunc(" nano ", false, "nano", []string{}))
	t.Run("With args", testNewCmdFromStringFunc("zed -w -n", false, "zed", []string{"-w", "-n"}))
}

func testNewCmdFromStringFunc(
	value string,
	hasError bool,
	expectName string,
	expectArg []string,
) func(t *testing.T) {
	return func(t *testing.T) {
		cmd, err := newCmdFromString(value)
		if !hasError && err != nil {
			t.Error("Unexpected error occured:", err)
			return
		}
		if expectName != cmd.name {
			t.Errorf("Name does not match. Got: %s,  Expected: %s \n", cmd.name, expectName)
			return
		}
		if !reflect.DeepEqual(expectArg, cmd.arg) {
			t.Errorf("Args do not match. Got: %+v,  Expected: %+v \n", cmd.arg, expectArg)
		}
	}
}
