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
)

func (s section) isFileSection() bool {
	return s == stagedSection || s == unstagedSection
}

type Model struct {
	workTreeStatus     git.WorkTreeStatus
	statusErr          error
	sections           [3]container.Model
	focusedSection     section
	lastFocusedSection section
}

func New() Model {
	isUntracked := func(item FileListItem) bool {
		return item.fileStatus.Code == git.Untracked
	}

	unstagedFilesItemHandler := func(msg tea.Msg) tea.Cmd {
		switch msg := msg.(type) {
		case selectItemMsg:
			return stageFile(msg.item.path)
		case focusItemMsg:
			return diff(
				git.DiffOption{
					FilePath:    msg.item.path,
					IsUntracked: isUntracked(msg.item),
				},
			)
		default:
			return nil
		}
	}

	stagedFilesItemHandler := func(msg tea.Msg) tea.Cmd {
		switch msg := msg.(type) {
		case selectItemMsg:
			return unstageFile(msg.item.path)
		case focusItemMsg:
			return diff(
				git.DiffOption{
					FilePath:    msg.item.path,
					IsStaged:    true,
					IsUntracked: isUntracked(msg.item),
				})
		default:
			return nil
		}
	}

	return Model{
		sections: [3]container.Model{
			container.NewModel(NewFileList("Unstaged", unstagedFilesItemHandler)),
			container.NewModel(NewFileList("Staged", stagedFilesItemHandler)),
			container.NewModel(NewDiff()),
		},
	}
}

func (m Model) Init() tea.Cmd {
	cmds := []tea.Cmd{initializeStatus()}
	for _, section := range m.sections {
		cmds = append(cmds, section.Init())
	}
	return tea.Batch(cmds...)
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case initializedMsg:
		var model = m
		model, scmd := model.handleStatusUpdateMsg(msg.statusMsg)
		model, dcmd := model.handleLoadedDiffMsg(msg.diffMsg)
		m = model
		cmds = append(cmds, scmd, dcmd)
	case statusUpdateMsg:
		model, cmd := m.handleStatusUpdateMsg(msg)
		m = model
		cmds = append(cmds, cmd)
	case loadedDiffMsg:
		model, cmd := m.handleLoadedDiffMsg(msg)
		m = model
		cmds = append(cmds, cmd)
	case tea.KeyMsg:
		switch msg.String() {
		case "right":
			m.focusedSection = diffSection
		case "left":
			if m.focusedSection != diffSection || !m.lastFocusedSection.isFileSection() {
				break
			}
			m.focusedSection = m.lastFocusedSection
		case "u":
			if m.focusedSection == unstagedSection {
				break
			}
			m.focusedSection = unstagedSection
			m.lastFocusedSection = m.focusedSection
		case "s":
			if m.focusedSection == stagedSection {
				break
			}
			m.focusedSection = stagedSection
			m.lastFocusedSection = m.focusedSection
		}
	}

	for i, section := range m.sections {
		updatedSection, cmd := section.UpdateFocus(i == int(m.focusedSection))
		cmds = append(cmds, cmd)

		updatedSection, cmd = updatedSection.Update(msg)
		cmds = append(cmds, cmd)

		m.sections[i] = updatedSection
	}

	return m, tea.Batch(cmds...)
}

func (m Model) handleStatusUpdateMsg(msg statusUpdateMsg) (Model, tea.Cmd) {
	m.workTreeStatus = msg.workTreeStatus
	m.statusErr = msg.err
	if section, ok := m.sections[unstagedSection].Content().(FileList); ok {
		section = section.SetFileListItems(createFileListItems(m.workTreeStatus.Unstaged))
		m.sections[unstagedSection] = m.sections[unstagedSection].SetContent(section)
	}
	if section, ok := m.sections[stagedSection].Content().(FileList); ok {
		section = section.SetFileListItems(createFileListItems(m.workTreeStatus.Staged))
		m.sections[stagedSection] = m.sections[stagedSection].SetContent(section)
	}
	return m, nil
}

func (m Model) handleLoadedDiffMsg(msg loadedDiffMsg) (Model, tea.Cmd) {
	section, ok := m.sections[diffSection].Content().(Diff)
	if !ok {
		return m, nil
	}
	section = section.SetContent(msg.diff, msg.err)
	m.sections[diffSection] = m.sections[diffSection].SetContent(section)
	return m, nil
}

func (m Model) View() string {
	if m.statusErr != nil {
		return fmt.Sprint(m.statusErr)
	}

	files := lipgloss.JoinVertical(
		lipgloss.Top,
		m.sections[unstagedSection].View(),
		m.sections[stagedSection].View(),
	)
	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		files,
		" ",
		m.sections[diffSection].View(),
	)
}

func (m Model) SetSize(width, height int) Model {
	filesWidth := int(float32(width) * 0.4)
	filesHeight := height / 2

	diffWidth := width - filesWidth - 1
	diffHeight := filesHeight * 2

	m.sections[unstagedSection] = m.sections[unstagedSection].SetSize(filesWidth, filesHeight)
	m.sections[stagedSection] = m.sections[stagedSection].SetSize(filesWidth, filesHeight)
	m.sections[diffSection] = m.sections[diffSection].SetSize(diffWidth, diffHeight)
	return m
}

func createFileListItems(fileStatusList git.FileStatusList) []FileListItem {
	items := make([]FileListItem, len(fileStatusList))
	for i, fs := range fileStatusList {
		items[i] = NewFileListItem(fs)
	}
	return items
}
