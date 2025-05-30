package commit

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/michaelhass/gitglance/internal/git"
	"github.com/michaelhass/gitglance/internal/ui/container"
	"github.com/michaelhass/gitglance/internal/ui/list"
	filelist "github.com/michaelhass/gitglance/internal/ui/list/file"
	"github.com/michaelhass/gitglance/internal/ui/textinput"
)

// Model represents the UI to pefrom a commit.
// It shows the staged files to  be included in the commit and
// allows to write a commit message.
type Model struct {
	stagedFileList container.Model
	message        container.Model
	keys           KeyMap
}

func New(branch string, stagedFileList git.FileStatusList) Model {
	fileListContent := list.NewContent(
		list.New(
			"Staged",
			func(msg tea.Msg) tea.Cmd { return nil },
			list.NewKeyMap("Nothing", "Nothing", "Nothing"),
		),
	)

	createListItems := func(fileStatusList git.FileStatusList) []list.Item {
		items := make([]list.Item, len(fileStatusList))
		for i, fs := range fileStatusList {
			items[i] = filelist.NewItem(fs, string(fs.StagedStatusCode))
		}
		return items
	}

	fileListContent.Model, _ = fileListContent.SetItems(createListItems(stagedFileList))

	messageContainer := container.New(
		textinput.NewContent(
			fmt.Sprintf("%s [%s]", "Commit", branch),
			"Enter commit message",
		),
	)
	messageContainer, _ = messageContainer.UpdateFocus(true)

	return Model{
		stagedFileList: container.New(fileListContent),
		message:        messageContainer,
		keys:           NewKeyMap(),
	}
}

func (m Model) Init() tea.Cmd {
	return loadMergeMsg
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var (
		cmds []tea.Cmd
		cmd  tea.Cmd
	)

	switch msg := msg.(type) {
	case MergeMsgLoaded:
		m, cmd = m.setMsg(msg.msg)
		cmds = append(cmds, cmd)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.toggleFocus):
			m, cmd = m.toggleFocus()
			cmds = append(cmds, cmd)
		case key.Matches(msg, m.keys.commit):
			if mc, ok := m.message.Content().(textinput.Content); ok {
				return m, Execute(mc.Text())
			}
		}
	}

	m, cmd = m.updateFocusedContainer(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Top,
		m.stagedFileList.View(),
		m.message.View(),
	)
}

func (m Model) SetSize(width, height int) Model {
	containerHeight := height / 2
	m.stagedFileList = m.stagedFileList.SetSize(width, containerHeight)
	m.message = m.message.SetSize(width, containerHeight)
	return m
}

func (m Model) Help() []key.Binding {
	return []key.Binding{
		m.keys.up,
		m.keys.down,
		m.keys.toggleFocus,
		m.keys.commit,
	}
}

func (m Model) setMsg(msg string) (Model, tea.Cmd) {
	if input, ok := m.message.Content().(textinput.Content); ok {
		input = input.SetValue(msg)
		input = input.SetCursorToStart()
		m.message = m.message.SetContent(input)
	}
	return m, nil
}

func (m Model) toggleFocus() (Model, tea.Cmd) {
	var (
		files    = m.stagedFileList
		filesCmd tea.Cmd

		message    = m.message
		messageCmd tea.Cmd
	)

	files, filesCmd = files.UpdateFocus(!files.IsFocused())
	message, messageCmd = message.UpdateFocus(!message.IsFocused())

	m.stagedFileList = files
	m.message = message

	return m, tea.Batch(filesCmd, messageCmd)
}

func (m Model) updateFocusedContainer(msg tea.Msg) (Model, tea.Cmd) {
	if m.stagedFileList.IsFocused() {
		files, cmd := m.stagedFileList.Update(msg)
		m.stagedFileList = files
		return m, cmd
	}
	message, cmd := m.message.Update(msg)
	m.message = message
	return m, cmd
}
