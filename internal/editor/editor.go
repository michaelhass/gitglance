package editor

import (
	"errors"
	"os"
	"os/exec"
	"strings"
)

type Cmd struct {
	name string
	arg  []string
}

type CreateEditorCmdOption func() (*Cmd, error)

func newFallbackEditorCmd() *Cmd {
	return &Cmd{name: "vi", arg: []string{}}
}

func EnvVisual() (*Cmd, error) {
	return newCmdFromString(os.Getenv("VISUAL"))
}

func EnvEditor() (*Cmd, error) {
	return newCmdFromString(os.Getenv("EDITOR"))
}

func WithCmdString(fn func() (string, error)) CreateEditorCmdOption {
	return func() (*Cmd, error) {
		value, err := fn()
		if err != nil {
			value = ""
		}
		return newCmdFromString(value)
	}
}

func newCmdFromString(cmdString string) (*Cmd, error) {
	cmdString = strings.TrimSpace(cmdString)
	components := []string{}
	for _, component := range strings.Split(cmdString, " ") {
		if len(component) > 0 {
			components = append(components, component)
		}
	}
	if len(components) == 0 {
		return nil, errors.New("Empty")
	}

	return &Cmd{name: components[0], arg: components[1:]}, nil
}

// OpenFileCmdDefault produces an exec.Cmd to open a path in an editor.
// It tries to prioritize the given list of editors before accessing a list of default editors.
// default editor options with priority:
// 1. env VISUAL
// 2. env EDITOR
func OpenFileCmdDefault(path string, prioritize ...CreateEditorCmdOption) *exec.Cmd {
	defaultCreateEditorOpts := []CreateEditorCmdOption{
		EnvVisual,
		EnvEditor,
	}
	all := append(prioritize, defaultCreateEditorOpts...)
	return OpenFileCmd(path, all...)
}

// OpenFileCmd produces an exec.Cmd to open a path in an editor.
// It will use the first successful CreateEditorCmdOption.
// Otherwise it will fallback to "vi".
func OpenFileCmd(path string, opts ...CreateEditorCmdOption) *exec.Cmd {
	editorCmd := newFallbackEditorCmd()
	for _, opt := range opts {
		if other, err := opt(); err == nil {
			editorCmd = other
			break
		}
	}
	editorCmd.arg = append(editorCmd.arg, path)
	return exec.Command(editorCmd.name, editorCmd.arg...)
}
