package status

import (
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/michaelhass/gitglance/internal/core/exit"
	"github.com/michaelhass/gitglance/internal/core/git"
	"github.com/michaelhass/gitglance/internal/core/refresh"
	"github.com/michaelhass/gitglance/internal/core/ui/components/container"
	"github.com/michaelhass/gitglance/internal/core/ui/components/list"
	filelist "github.com/michaelhass/gitglance/internal/core/ui/components/list/file"
	"github.com/michaelhass/gitglance/internal/core/ui/style"
	"github.com/michaelhass/gitglance/internal/domain/diff"
)

type section byte

const (
	unstagedSection section = iota
	stagedSection
	diffSection
)

const (
	filesWidthFactor         float32 = 0.4
	sectionsHorizontalMargin int     = 1
	helpHeight               int     = 1
)

var (
	helpStyle = style.ShortHelp
)

type Model struct {
	workTreeStatus git.WorkTreeStatus

	sections [3]container.Model

	help help.Model
	keys KeyMap

	statusErr              error
	focusedSection         section
	lastFocusedFileSection section

	isInitialized bool
}

func New() Model {
	unstagedFilesItemHandler := func(msg tea.Msg) tea.Cmd {
		switch msg := msg.(type) {
		case list.SelectItemMsg:
			if item, ok := msg.Item.(filelist.Item); ok {
				return stageFile(item.Path)
			}
			return nil
		case list.FocusItemMsg:
			if item, ok := msg.Item.(filelist.Item); ok {
				return diffFile(
					git.DiffOptions{
						FilePath:    item.Path,
						IsUntracked: item.IsUntracked(),
					},
				)
			}
			return nil
		case list.EditItemMsg:
			if item, ok := msg.Item.(filelist.Item); ok {
				return openFile(item.Path)
			}
			return nil
		case list.DeleteItemMsg:
			if item, ok := msg.Item.(filelist.Item); ok {
				return deleteFile(item)
			}
			return nil
		case list.SelectAllItemMsg:
			return stageAll()
		case list.BottomNoMoreFocusableItems:
			return focusSection(stagedSection)
		case list.NoItemsMsg:
			return showEmptyDiff
		default:
			return nil
		}
	}

	stagedFilesItemHandler := func(msg tea.Msg) tea.Cmd {
		switch msg := msg.(type) {
		case list.SelectItemMsg:
			if item, ok := msg.Item.(filelist.Item); ok {
				return unstageFile(item.Path)
			}
			return nil
		case list.SelectAllItemMsg:
			return unstageAll()
		case list.FocusItemMsg:
			if item, ok := msg.Item.(filelist.Item); ok {
				return diffFile(
					git.DiffOptions{
						FilePath:    item.Path,
						IsStaged:    true,
						IsUntracked: item.IsUntracked(),
					})
			}
			return nil
		case list.EditItemMsg:
			if item, ok := msg.Item.(filelist.Item); ok {
				return openFile(item.Path)
			}
			return nil
		case list.TopNoMoreFocusableItems:
			return focusSection(unstagedSection)
		case list.NoItemsMsg:
			return showEmptyDiff
		default:
			return nil
		}
	}

	help := help.New()
	help.ShowAll = false

	unstagedFileList := list.NewContainerContent(
		list.New("Unstaged", unstagedFilesItemHandler, list.NewKeyMap(
			"stage all",
			"stage file",
			"reset file",
		)),
	)
	var stagedFileListKeyMap = list.NewKeyMap(
		"unstage all",
		"unstage file",
		"",
	)
	stagedFileListKeyMap.Delete.SetEnabled(false)

	stagedFileList := list.NewContainerContent(list.New("Staged", stagedFilesItemHandler, stagedFileListKeyMap))
	diffContent := diff.NewContent(diff.New())

	return Model{
		sections: [3]container.Model{
			container.New(unstagedFileList),
			container.New(stagedFileList),
			container.New(diffContent),
		},
		help: help,
		keys: newKeyMap(),
	}
}

func (m Model) Init() tea.Cmd {
	cmds := []tea.Cmd{initializeStatus()}
	for _, section := range m.sections {
		cmds = append(cmds, section.Init())
	}
	return tea.Sequence(cmds...)
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var (
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case initializedMsg:
		var model = m
		model, scmd := model.handleStatusUpdateMsg(msg.statusMsg)
		model, dcmd := model.handleLoadedDiffMsg(msg.diffMsg)
		model.isInitialized = model.statusErr == nil
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
	case focusSectionMsg:
		m = m.focusSection(msg.section)
	case refresh.Msg:
		cmds = append(cmds, refreshStatus())
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.left):
			m = m.focusSection(m.lastFocusedFileSection)
		case key.Matches(msg, m.keys.right), key.Matches(msg, m.keys.focusDiff):
			m = m.focusSection(diffSection)
		case key.Matches(msg, m.keys.focusUnstaged):
			m = m.focusSection(unstagedSection)
		case key.Matches(msg, m.keys.focusStaged):
			m = m.focusSection(stagedSection)
		case key.Matches(msg, m.keys.commit):
			cmds = append(
				cmds,
				showCommitDialog(
					m.workTreeStatus.CleanedBranchName,
					m.workTreeStatus.StagedFiles(),
				),
			)
		case key.Matches(msg, m.keys.refresh):
			cmds = append(cmds, refreshStatus())
		case key.Matches(msg, m.keys.stash):
			cmds = append(cmds, showStashAllConfirmation())
		case key.Matches(msg, key.NewBinding(key.WithKeys("S"))):
			cmds = append(cmds, showApplyStashDialog())
		}
	}

	m.keys = m.updateKeys()

	if !m.isInitialized {
		return m, tea.Batch(cmds...)
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

func (m Model) View() string {
	if !m.isInitialized {
		return "loading..."
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

	m.help.Width = width - helpStyle.GetHorizontalMargins()

	return m
}

func (m Model) handleStatusUpdateMsg(msg statusUpdateMsg) (Model, tea.Cmd) {
	var cmds []tea.Cmd

	m.workTreeStatus = msg.WorkTreeStatus
	m.statusErr = msg.Err
	if !m.isInitialized && msg.Err != nil {
		return m, exit.WithMsg(msg.Err.Error())
	}

	if section, ok := m.sections[unstagedSection].Content().(list.ContainerContent); ok {
		model, cmd := section.SetItems(createListItems(m.workTreeStatus.UnstagedFiles(), false))
		section.Model = model
		cmds = append(cmds, cmd)
		m.sections[unstagedSection] = m.sections[unstagedSection].SetContent(section)
	}
	if section, ok := m.sections[stagedSection].Content().(list.ContainerContent); ok {
		model, cmd := section.SetItems(createListItems(m.workTreeStatus.StagedFiles(), true))
		model = model.SetTitle(fmt.Sprintf("Staged [%s]", m.workTreeStatus.CleanedBranchName))
		section.Model = model
		cmds = append(cmds, cmd)
		m.sections[stagedSection] = m.sections[stagedSection].SetContent(section)
	}

	cmds = append(cmds, tea.SetWindowTitle(m.workTreeStatus.CleanedBranchName))

	return m, tea.Batch(cmds...)
}

func (m Model) handleLoadedDiffMsg(msg loadedDiffMsg) (Model, tea.Cmd) {
	section, ok := m.sections[diffSection].Content().(diff.ContainerContent)
	if !ok {
		return m, nil
	}
	section.Model = section.SetContent(msg.Diff, msg.Err)
	m.sections[diffSection] = m.sections[diffSection].SetContent(section)
	return m, nil
}

func (m Model) focusSection(section section) Model {
	m.lastFocusedFileSection = m.focusedSection
	m.focusedSection = section
	return m
}

func (m Model) updateKeys() KeyMap {
	keys := m.keys
	keys.additionalKeyMap = m.sections[m.focusedSection].Content().KeyMap()

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
		isUnstagedFocusedLast := m.lastFocusedFileSection == unstagedSection
		keys.focusUnstaged.SetEnabled(isUnstagedFocusedLast)
		keys.focusStaged.SetEnabled(!isUnstagedFocusedLast)
		keys.focusDiff.SetEnabled(false)
	}

	return keys
}

func createListItems(fileStatusList git.FileStatusList, isStaged bool) []list.Item {
	items := make([]list.Item, len(fileStatusList))

	accessory := func(fileStatus git.FileStatus, isStaged bool) string {
		if isStaged {
			return string(fileStatus.StagedStatusCode)
		} else {
			return string(fileStatus.UnstagedStatusCode)
		}
	}

	for i, fs := range fileStatusList {
		items[i] = filelist.NewItem(
			fs,
			accessory(fs, isStaged),
		)
	}
	return items
}
