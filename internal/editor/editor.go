package editor

import (
	"errors"
	"os"
	"os/exec"
	"strings"
)

type cmd struct {
	name string
	arg  []string
}

type CreateEditorCmdOption func() (*cmd, error)

func CreateFallbackEditor() *cmd {
	return &cmd{name: "vi", arg: []string{}}
}

func EnvVisual() (*cmd, error) {
	return newCmdFromString(os.Getenv("VISUAL"))
}

func EnvEditor() (*cmd, error) {
	return newCmdFromString(os.Getenv("EDITOR"))
}

func WithCmdString(fn func() (string, error)) CreateEditorCmdOption {
	return func() (*cmd, error) {
		value, err := fn()
		if err != nil {
			value = " "
		}
		return newCmdFromString(value)
	}
}

func newCmdFromString(cmdString string) (*cmd, error) {
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

	return &cmd{name: components[0], arg: components[1:]}, nil
}

func OpenFileCmdDefault(path string, prioritize ...CreateEditorCmdOption) *exec.Cmd {
	defaultCreateEditorOpts := []CreateEditorCmdOption{
		EnvVisual,
		EnvEditor,
	}
	all := append(prioritize, defaultCreateEditorOpts...)
	return OpenFileCmd(path, all...)
}

func OpenFileCmd(path string, opts ...CreateEditorCmdOption) *exec.Cmd {
	editorCmd := CreateFallbackEditor()
	for _, opt := range opts {
		if other, err := opt(); err == nil {
			editorCmd = other
			break
		}
	}
	editorCmd.arg = append(editorCmd.arg, path)
	return exec.Command(editorCmd.name, editorCmd.arg...)
}
