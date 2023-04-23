package status

import (
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/michaelhass/gitglance/internal/git"
	uicmd "github.com/michaelhass/gitglance/internal/ui/cmd"
	"github.com/michaelhass/gitglance/internal/ui/commit"
	"github.com/michaelhass/gitglance/internal/ui/container"
	"github.com/michaelhass/gitglance/internal/ui/diff"
	"github.com/michaelhass/gitglance/internal/ui/filelist"
	"github.com/michaelhass/gitglance/internal/ui/styles"
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
	helpStyle = styles.ShortHelpStyle.Copy()
)

type Model struct {
	workTreeStatus git.WorkTreeStatus

	sections [3]container.Model

	commit commit.Model

	help help.Model
	keys KeyMap

	statusErr              error
	focusedSection         section
	lastFocusedFileSection section
}

func New() Model {
	isUntracked := func(item filelist.Item) bool {
		return item.FileStatus.Code == git.Untracked
	}

	unstagedFilesItemHandler := func(msg tea.Msg) tea.Cmd {
		switch msg := msg.(type) {
		case filelist.SelectItemMsg:
			return uicmd.StageFile(msg.Item.Path)
		case filelist.FocusItemMsg:
			return uicmd.Diff(
				git.DiffOption{
					FilePath:    msg.Item.Path,
					IsUntracked: isUntracked(msg.Item),
				},
			)
		default:
			return nil
		}
	}

	stagedFilesItemHandler := func(msg tea.Msg) tea.Cmd {
		switch msg := msg.(type) {
		case filelist.SelectItemMsg:
			return uicmd.UnstageFile(msg.Item.Path)
		case filelist.FocusItemMsg:
			return uicmd.Diff(
				git.DiffOption{
					FilePath:    msg.Item.Path,
					IsStaged:    true,
					IsUntracked: isUntracked(msg.Item),
				})
		default:
			return nil
		}
	}

	help := help.New()
	help.ShowAll = false

	unstagedFileList := container.NewFileListContent(
		filelist.New("Unstaged", unstagedFilesItemHandler, filelist.NewKeyMap("stage file")),
	)
	stagedFileList := container.NewFileListContent(
		filelist.New("Staged", stagedFilesItemHandler, filelist.NewKeyMap("unstage file")),
	)
	diffContent := container.NewDiffContent(diff.New())

	return Model{
		sections: [3]container.Model{
			container.NewModel(unstagedFileList),
			container.NewModel(stagedFileList),
			container.NewModel(diffContent),
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
		// In some cases we do not want to pass the key msg to the sections
		// E.g. switching section focus via arrow keys.
		blockSectionMsgUpdate = false
	)

	switch msg := msg.(type) {
	case InitializedMsg:
		var model = m
		model, scmd := model.handleStatusUpdateMsg(msg.StatusMsg)
		model, dcmd := model.handleLoadedDiffMsg(msg.DiffMsg)
		m = model
		cmds = append(cmds, scmd, dcmd)
	case uicmd.StatusUpdateMsg:
		model, cmd := m.handleStatusUpdateMsg(msg)
		m = model
		cmds = append(cmds, cmd)
	case uicmd.LoadedDiffMsg:
		model, cmd := m.handleLoadedDiffMsg(msg)
		m = model
		cmds = append(cmds, cmd)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.left):
			m.focusedSection = m.lastFocusedFileSection
		case key.Matches(msg, m.keys.right), key.Matches(msg, m.keys.focusDiff):
			m.lastFocusedFileSection = m.focusedSection
			m.focusedSection = diffSection
		case key.Matches(msg, m.keys.up):
			if m.focusedSection != stagedSection {
				break
			}
			section, ok := m.sections[stagedSection].Content().(container.FileListContent)
			if !ok || !section.IsFirstIndexFocused() {
				break
			}
			m.focusedSection = unstagedSection
			blockSectionMsgUpdate = true
		case key.Matches(msg, m.keys.down):
			if m.focusedSection != unstagedSection {
				break
			}
			section, ok := m.sections[unstagedSection].Content().(container.FileListContent)
			if !ok || !section.IsLastIndexFocused() {
				break
			}
			m.focusedSection = stagedSection
			blockSectionMsgUpdate = false
		case key.Matches(msg, m.keys.focusUnstaged):
			m.focusedSection = unstagedSection
		case key.Matches(msg, m.keys.focusStaged):
			m.focusedSection = stagedSection
		}
	}

	m.keys = m.updateKeys()

	for i, section := range m.sections {
		updatedSection, cmd := section.UpdateFocus(i == int(m.focusedSection))
		cmds = append(cmds, cmd)

		if !blockSectionMsgUpdate {
			updatedSection, cmd = updatedSection.Update(msg)
			cmds = append(cmds, cmd)
		}

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

func (m Model) handleStatusUpdateMsg(msg uicmd.StatusUpdateMsg) (Model, tea.Cmd) {
	m.workTreeStatus = msg.WorkTreeStatus
	m.statusErr = msg.Err
	if section, ok := m.sections[unstagedSection].Content().(container.FileListContent); ok {
		section.Model = section.SetItems(createListItems(m.workTreeStatus.Unstaged))
		m.sections[unstagedSection] = m.sections[unstagedSection].SetContent(section)
	}
	if section, ok := m.sections[stagedSection].Content().(container.FileListContent); ok {
		section.Model = section.SetItems(createListItems(m.workTreeStatus.Staged))
		m.sections[stagedSection] = m.sections[stagedSection].SetContent(section)
	}
	return m, nil
}

func (m Model) handleLoadedDiffMsg(msg uicmd.LoadedDiffMsg) (Model, tea.Cmd) {
	section, ok := m.sections[diffSection].Content().(container.DiffContent)
	if !ok {
		return m, nil
	}
	section.Model = section.SetContent(msg.Diff, msg.Err)
	m.sections[diffSection] = m.sections[diffSection].SetContent(section)
	return m, nil
}

func (m Model) updateKeys() KeyMap {
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

func createListItems(fileStatusList git.FileStatusList) []filelist.Item {
	items := make([]filelist.Item, len(fileStatusList))
	for i, fs := range fileStatusList {
		items[i] = filelist.NewItem(fs)
	}
	return items
}
