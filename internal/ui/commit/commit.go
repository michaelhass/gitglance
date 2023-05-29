package commit

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/michaelhass/gitglance/internal/git"
	"github.com/michaelhass/gitglance/internal/ui/container"
	"github.com/michaelhass/gitglance/internal/ui/filelist"
)

type Model struct {
	stagedFileList container.Model
	message        container.Model
	keys           KeyMap
}

func New(stagedFileList git.FileStatusList) Model {
	fileListContent := container.NewFileListContent(
		filelist.New(
			"Staged",
			func(msg tea.Msg) tea.Cmd { return nil },
			filelist.NewKeyMap("Nothing"),
		),
	)

	createListItems := func(fileStatusList git.FileStatusList) []filelist.Item {
		items := make([]filelist.Item, len(fileStatusList))
		for i, fs := range fileStatusList {
			items[i] = filelist.NewItem(fs, string(fs.StagedStatusCode))
		}
		return items
	}

	fileListContent.Model, _ = fileListContent.SetItems(createListItems(stagedFileList))

	messageContainer := container.New(newMessageContent())
	messageContainer, _ = messageContainer.UpdateFocus(true)

	return Model{
		stagedFileList: container.New(fileListContent),
		message:        messageContainer,
		keys:           NewKeyMap(),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var (
		cmds []tea.Cmd
		cmd  tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.toggleFocus):
			m, cmd = m.toggleFocus()
			cmds = append(cmds, cmd)
		case key.Matches(msg, m.keys.commit):
			if mc, ok := m.message.Content().(messageContent); ok {
				return m, Execute(mc.message())
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
