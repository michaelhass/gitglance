package status

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/michaelhass/gitglance/internal/git"
	"github.com/michaelhass/gitglance/internal/ui/container"
)

type section byte

const (
	unstagedSection section = iota
	stagedSection
	diffSection
	lastSection section = diffSection
)

type Model struct {
	repo           *git.Repository
	status         git.Status
	statusErr      error
	sections       [lastSection + 1]container.Model
	focusedSection section
}

func (m Mock) Title() string {
	return m.title
}

func (m Mock) SetIsFocused(isFocused bool) container.Content {
	return m
}

func (m Mock) SetSize(width int, height int) container.Content {
	// TODO: Implement
	return m
}

func (m Mock) Init() tea.Cmd {
	return nil
}

func (m Mock) Update(msg tea.Msg) (container.Content, tea.Cmd) {
	return m, nil // TODO: Implement
}

func (m Mock) View() string {
	return ""
}

type Mock struct {
	title string
}

func New(repo *git.Repository) Model {
	return Model{
		repo: repo,
		sections: [3]container.Model{
			container.NewModel(NewFileList("Unstaged")),
			container.NewModel(NewFileList("Staged")),
			container.NewModel(Mock{title: "Diff"}),
		},
	}
}

func (m Model) Init() tea.Cmd {
	return load(m.repo)
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case statusUpdateMsg:
		m.status = msg.status
		m.statusErr = msg.err
		if list, ok := m.sections[unstagedSection].Content().(FileList); ok {
			list = list.SetFileListItems(createFileListItems(m.status.Unstaged))
			m.sections[unstagedSection] = m.sections[unstagedSection].SetContent(list)
		}
		if list, ok := m.sections[stagedSection].Content().(FileList); ok {
			list = list.SetFileListItems(createFileListItems(m.status.Staged))
			m.sections[stagedSection] = m.sections[stagedSection].SetContent(list)
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			if m.focusedSection == lastSection {
				m.focusedSection = 0
			} else {
				m.focusedSection += 1
			}
		}
	}

	var cmds []tea.Cmd

	for i, section := range m.sections {
		section = section.SetIsFocused(i == int(m.focusedSection))
		updatedSection, cmd := section.Update(msg)
		m.sections[i] = updatedSection
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.statusErr != nil {
		return fmt.Sprint(m.statusErr)
	}

	files := lipgloss.JoinHorizontal(
		lipgloss.Left,
		m.sections[unstagedSection].View(),
		" ",
		m.sections[stagedSection].View(),
	)
	return lipgloss.JoinVertical(lipgloss.Top, files, m.sections[diffSection].View())
}

func (m Model) SetSize(width, height int) Model {
	filesWidth := (width / 2)
	filesHeight := (height / 2)

	diffWidth := width
	diffHeight := filesHeight

	m.sections[unstagedSection] = m.sections[unstagedSection].SetSize(filesWidth, filesHeight)
	m.sections[stagedSection] = m.sections[stagedSection].SetSize(filesWidth, filesHeight)
	m.sections[diffSection] = m.sections[diffSection].SetSize(diffWidth, diffHeight)
	return m
}

func createFileListItems(fileStatusList git.FileStatusList) []FileListItem {
	items := make([]FileListItem, len(fileStatusList))
	for i, fs := range fileStatusList {
		items[i] = NewFileListItem(fs.Path, fmt.Sprintf("[%s]", string(fs.Code)))
	}
	return items
}

// Cmd

func load(repo *git.Repository) func() tea.Msg {
	return func() tea.Msg {
		var msg statusUpdateMsg

		wt, err := repo.Worktree()
		if err != nil {
			msg.err = err
			return msg
		}

		status, err := wt.Status()
		if err != nil {
			msg.err = err
			return msg
		}
		msg.status = status
		return msg
	}
}

// Msg

type statusUpdateMsg struct {
	err    error
	status git.Status
}
