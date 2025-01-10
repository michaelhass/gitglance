package status

import (
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/michaelhass/gitglance/internal/git"
	"github.com/michaelhass/gitglance/internal/ui/commit"
	"github.com/michaelhass/gitglance/internal/ui/container"
	"github.com/michaelhass/gitglance/internal/ui/dialog"
	"github.com/michaelhass/gitglance/internal/ui/diff"
	"github.com/michaelhass/gitglance/internal/ui/filelist"
	"github.com/michaelhass/gitglance/internal/ui/style"
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
	helpStyle = style.ShortHelp.Copy()
)

type Model struct {
	workTreeStatus git.WorkTreeStatus

	sections [3]container.Model

	help help.Model
	keys KeyMap

	statusErr              error
	focusedSection         section
	lastFocusedFileSection section
}

func New() Model {
	unstagedFilesItemHandler := func(msg tea.Msg) tea.Cmd {
		switch msg := msg.(type) {
		case filelist.SelectItemMsg:
			return stageFile(msg.Item.Path)
		case filelist.FocusItemMsg:
			return diffFile(
				git.DiffOptions{
					FilePath:    msg.Item.Path,
					IsUntracked: msg.Item.IsUntracked(),
				},
			)
		case filelist.SelectAllItemMsg:
			return stageAll()
		case filelist.BottomNoMoreFocusableItems:
			return focusSection(stagedSection)
		default:
			return nil
		}
	}

	stagedFilesItemHandler := func(msg tea.Msg) tea.Cmd {
		switch msg := msg.(type) {
		case filelist.SelectItemMsg:
			return unstageFile(msg.Item.Path)
		case filelist.SelectAllItemMsg:
			return unstageAll()
		case filelist.FocusItemMsg:
			return diffFile(
				git.DiffOptions{
					FilePath:    msg.Item.Path,
					IsStaged:    true,
					IsUntracked: msg.Item.IsUntracked(),
				})
		case filelist.TopNoMoreFocusableItems:
			return focusSection(unstagedSection)
		default:
			return nil
		}
	}

	help := help.New()
	help.ShowAll = false

	unstagedFileList := filelist.NewContent(
		filelist.New("Unstaged", unstagedFilesItemHandler, filelist.NewKeyMap("stage all", "stage file")),
	)
	stagedFileList := filelist.NewContent(
		filelist.New("Staged", stagedFilesItemHandler, filelist.NewKeyMap("unstage all", "unstage file")),
	)
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
	return tea.Batch(cmds...)
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
			content := commit.NewContent(commit.New(
				m.workTreeStatus.CleanedBranchName,
				m.workTreeStatus.StagedFiles()),
			)
			cmds = append(cmds, dialog.Show(content, initializeStatus(), dialog.CenterDisplayMode))
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

	m.help.Width = width - helpStyle.GetHorizontalMargins()

	return m
}

func (m Model) handleStatusUpdateMsg(msg statusUpdateMsg) (Model, tea.Cmd) {
	var cmd tea.Cmd

	m.workTreeStatus = msg.WorkTreeStatus
	m.statusErr = msg.Err

	if section, ok := m.sections[unstagedSection].Content().(filelist.Content); ok {
		section.Model, cmd = section.SetItems(createListItems(m.workTreeStatus.UnstagedFiles(), false))
		m.sections[unstagedSection] = m.sections[unstagedSection].SetContent(section)
	}
	if section, ok := m.sections[stagedSection].Content().(filelist.Content); ok {
		section.Model, cmd = section.SetItems(createListItems(m.workTreeStatus.StagedFiles(), true))
		m.sections[stagedSection] = m.sections[stagedSection].SetContent(section)
	}
	return m, cmd
}

func (m Model) handleLoadedDiffMsg(msg loadedDiffMsg) (Model, tea.Cmd) {
	section, ok := m.sections[diffSection].Content().(diff.Content)
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
		keys.focusUnstaged.SetEnabled(true)
		keys.focusStaged.SetEnabled(true)
		keys.focusDiff.SetEnabled(false)
	}

	return keys
}

func createListItems(fileStatusList git.FileStatusList, isStaged bool) []filelist.Item {
	items := make([]filelist.Item, len(fileStatusList))

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
