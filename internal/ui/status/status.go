package status

import (
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/michaelhass/gitglance/internal/git"
	"github.com/michaelhass/gitglance/internal/ui/container"
	"github.com/michaelhass/gitglance/internal/ui/styles"
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

const (
	filesWidthFactor         float32 = 0.4
	sectionsHorizontalMargin int     = 1
	helpHeight               int     = 1
)

var (
	helpStyle = styles.ShortHelpStyle.Copy()
)

type Model struct {
	workTreeStatus         git.WorkTreeStatus
	sections               [3]container.Model
	help                   help.Model
	keys                   statusKeyMap
	statusErr              error
	focusedSection         section
	lastFocusedFileSection section
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

	help := help.New()
	help.ShowAll = false

	return Model{
		sections: [3]container.Model{
			container.NewModel(NewFileList("Unstaged", unstagedFilesItemHandler, newFilesKeyMap("stage file"))),
			container.NewModel(NewFileList("Staged", stagedFilesItemHandler, newFilesKeyMap("unstage file"))),
			container.NewModel(NewDiff()),
		},
		help: help,
		keys: newStatusKeyMap(),
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
		switch {
		case key.Matches(msg, m.keys.right), key.Matches(msg, m.keys.focusDiff):
			m.focusedSection = diffSection
		case key.Matches(msg, m.keys.left):
			m.focusedSection = m.lastFocusedFileSection
		case key.Matches(msg, m.keys.focusUnstaged):
			m.focusedSection = unstagedSection
			m.lastFocusedFileSection = m.focusedSection
		case key.Matches(msg, m.keys.focusStaged):
			m.focusedSection = stagedSection
			m.lastFocusedFileSection = m.focusedSection
		}
	}

	m.keys = m.updateKeys()

	for i, section := range m.sections {
		updatedSection, cmd := section.UpdateFocus(i == int(m.focusedSection))
		cmds = append(cmds, cmd)

		updatedSection, cmd = updatedSection.Update(msg)
		cmds = append(cmds, cmd)

		m.sections[i] = updatedSection
	}

	return m, tea.Batch(cmds...)
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

	sections := lipgloss.JoinHorizontal(
		lipgloss.Left,
		files,
		" ",
		m.sections[diffSection].View(),
	)

	return lipgloss.JoinVertical(
		lipgloss.Top,
		sections,
		helpStyle.Render(m.help.View(m.keys)),
	)
}

func (m Model) SetSize(width, height int) Model {
	var (
		maxSectionHeight = height - helpHeight

		filesWidth  = int(float32(width) * filesWidthFactor)
		filesHeight = maxSectionHeight / 2

		diffWidth  = width - filesWidth - sectionsHorizontalMargin
		diffHeight = filesHeight * 2 // don't use maxSectionHeight. Avoids layouting issues if uneven.
	)

	m.sections[unstagedSection] = m.sections[unstagedSection].SetSize(filesWidth, filesHeight)
	m.sections[stagedSection] = m.sections[stagedSection].SetSize(filesWidth, filesHeight)
	m.sections[diffSection] = m.sections[diffSection].SetSize(diffWidth, diffHeight)
	return m
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

func (m Model) updateKeys() statusKeyMap {
	keys := m.keys

	switch m.focusedSection {
	case unstagedSection:
		keys.left.SetEnabled(false)
		keys.right.SetEnabled(true)
		keys.focusUnstaged.SetEnabled(false)
		keys.focusStaged.SetEnabled(true)
		keys.focusDiff.SetEnabled(true)
	case stagedSection:
		keys.left.SetEnabled(false)
		keys.right.SetEnabled(true)
		keys.focusUnstaged.SetEnabled(true)
		keys.focusStaged.SetEnabled(false)
		keys.focusDiff.SetEnabled(true)
	case diffSection:
		keys.left.SetEnabled(true)
		keys.right.SetEnabled(false)
		keys.focusUnstaged.SetEnabled(true)
		keys.focusStaged.SetEnabled(true)
		keys.focusDiff.SetEnabled(false)
	}

	return keys
}

func createFileListItems(fileStatusList git.FileStatusList) []FileListItem {
	items := make([]FileListItem, len(fileStatusList))
	for i, fs := range fileStatusList {
		items[i] = NewFileListItem(fs)
	}
	return items
}
