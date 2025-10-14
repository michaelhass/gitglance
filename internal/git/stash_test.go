package git

import (
	"fmt"
	"testing"
)

func TestGetStashIdx(t *testing.T) {
	tests := []struct {
		input     string
		expectIdx int
		hasError  bool
	}{
		{"stash@{12}", 12, false},
		{"prefix stash@{12}", 12, false},
		{"stash@{12} stash@{3} multiple", 12, false},
		{"stash@{} empty", -1, true},
		{"25", -1, true},
	}

	builder := newDefaultStashEntryBuilder()
	for _, tt := range tests {
		testName := fmt.Sprintf("Reading idx from: `%s`", tt.input)
		t.Run(testName, func(t *testing.T) {
			gotIdx, gotErr := builder.getStashIdxFromLine(tt.input)
			if (tt.hasError && gotErr == nil) ||
				(!tt.hasError && gotErr != nil) {
				t.Error("Unexpected error result:", gotErr)
			}
			if gotIdx != tt.expectIdx {
				t.Errorf("Got idx `%d` but expected `%d.`", gotIdx, tt.expectIdx)
			}
		})
	}
}

func TestGetStashEntryMsg(t *testing.T) {
	tests := []struct {
		input     string
		expectMsg string
		hasError  bool
	}{
		{"stash@{12}", "", true},
		{": ", "", false},
		{"stash@{12}: this is the message", "this is the message", false},
		{"stash@{12}: this is the message: more message", "this is the message: more message", false},
	}

	builder := newDefaultStashEntryBuilder()
	for _, tt := range tests {
		testName := fmt.Sprintf("Reading message from: `%s`", tt.input)
		t.Run(testName, func(t *testing.T) {
			gotMsg, gotErr := builder.getStashMsgFromLine(tt.input)
			if (tt.hasError && gotErr == nil) ||
				(!tt.hasError && gotErr != nil) {
				t.Error("Unexpected error result:", gotErr)
			}
			if gotMsg != tt.expectMsg {
				t.Errorf("Got msg`%s` but expected `%s`.", gotMsg, tt.expectMsg)
			}
		})
	}
}
