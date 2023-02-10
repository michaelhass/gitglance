package status

import (
	"errors"
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

type initializedMsg struct {
	statusMsg statusUpdateMsg
	diffMsg   loadedDiffMsg
}

type statusUpdateMsg struct {
	err            error
	workTreeStatus git.WorkTreeStatus
}

type loadedDiffMsg struct {
	err  error
	diff string
}

type Model struct {
	workTreeStatus git.WorkTreeStatus
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
	m.width, m.height = width, height
	return m
}

func (m Mock) Init() tea.Cmd {
	return nil
}

func (m Mock) Update(msg tea.Msg) (container.Content, tea.Cmd) {
	return m, nil // TODO: Implement
}

func (m Mock) View() string {
	if m.diffErr != nil {
		return fmt.Sprint(m.diffErr)
	}
	return lipgloss.NewStyle().MaxWidth(m.width).MaxHeight(m.height).Render(m.diff)
}

type Mock struct {
	title         string
	diffErr       error
	diff          string
	width, height int
}

func New() Model {
	unstagedFilesItemHandler := func(msg tea.Msg) tea.Cmd {
		switch msg := msg.(type) {
		case selectItemMsg:
			return stageFile(msg.item.path)
		case focusItemMsg:
			return diff(git.DiffOption{FilePath: msg.item.path})
		default:
			return nil
		}
	}

	stagedFilesItemHandler := func(msg tea.Msg) tea.Cmd {
		switch msg := msg.(type) {
		case selectItemMsg:
			return unstageFile(msg.item.path)
		case focusItemMsg:
			return diff(git.DiffOption{FilePath: msg.item.path, IsStaged: true})
		default:
			return nil
		}
	}

	return Model{
		sections: [3]container.Model{
			container.NewModel(NewFileList("Unstaged", unstagedFilesItemHandler)),
			container.NewModel(NewFileList("Staged", stagedFilesItemHandler)),
			container.NewModel(Mock{title: "Diff"}),
		},
	}
}

func (m Model) Init() tea.Cmd {
	return initialize()
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
		switch msg.Type {
		case tea.KeyTab:
			if m.focusedSection == lastSection {
				m.focusedSection = 0
			} else {
				m.focusedSection += 1
			}

			fileList, ok := m.sections[m.focusedSection].Content().(FileList)
			if !ok {
				break
			}

			item, err := fileList.FocusedItem()
			if err != nil {
				break
			}
			cmds = append(cmds, diff(git.DiffOption{FilePath: item.path}))
		}
	}

	for i, section := range m.sections {
		section = section.SetIsFocused(i == int(m.focusedSection))
		updatedSection, cmd := section.Update(msg)
		m.sections[i] = updatedSection
		cmds = append(cmds, cmd)
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
	section, ok := m.sections[diffSection].Content().(Mock)
	if !ok {
		return m, nil
	}
	section.diff = msg.diff
	section.diffErr = msg.err
	m.sections[diffSection] = m.sections[diffSection].SetContent(section)
	return m, nil
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
	filesWidth := width / 2
	filesHeight := height / 2
	// Multiply instead of using 'width'
	// Avoids different sizes when width is uneven.
	diffWidth := filesWidth * 2
	diffHeight := filesHeight

	m.sections[unstagedSection] = m.sections[unstagedSection].SetSize(filesWidth, filesHeight)
	m.sections[stagedSection] = m.sections[stagedSection].SetSize(filesWidth, filesHeight)
	m.sections[diffSection] = m.sections[diffSection].SetSize(diffWidth, diffHeight)
	return m
}

func createFileListItems(fileStatusList git.FileStatusList) []FileListItem {
	items := make([]FileListItem, len(fileStatusList))
	for i, fs := range fileStatusList {
		path := fs.Path
		if len(fs.Extra) > 0 {
			path = fmt.Sprintf("%s → %s", path, fs.Extra)
		}
		items[i] = NewFileListItem(path, fmt.Sprintf("[%s]", string(fs.Code)))
	}
	return items
}

// Cmd

func initialize() func() tea.Msg {
	return func() tea.Msg {
		var (
			msg            initializedMsg
			workTreeStatus git.WorkTreeStatus
			unstagedFiles  git.FileStatusList
			err            error
		)

		workTreeStatus, err = git.Status()
		if err != nil {
			msg.statusMsg.err = err
			return msg
		}
		msg.statusMsg.workTreeStatus = workTreeStatus

		unstagedFiles = msg.statusMsg.workTreeStatus.Unstaged
		if len(unstagedFiles) == 0 {
			return msg
		}

		diffMsg, ok := diff(git.DiffOption{FilePath: unstagedFiles[0].Path})().(loadedDiffMsg)
		if !ok {
			diffMsg.err = errors.New("unable to load diff")
			return msg
		}

		msg.diffMsg = diffMsg
		return msg
	}
}

func stageFile(path string) func() tea.Msg {
	return func() tea.Msg {
		var (
			workTreeStatus git.WorkTreeStatus
			msg            statusUpdateMsg
			err            error
		)

		err = git.StageFile(path)
		if err != nil {
			msg.err = err
			return msg
		}

		workTreeStatus, err = git.Status()
		if err != nil {
			msg.err = err
			return msg
		}
		msg.workTreeStatus = workTreeStatus

		return msg
	}
}

func unstageFile(path string) func() tea.Msg {
	return func() tea.Msg {
		var (
			workTreeStatus git.WorkTreeStatus
			msg            statusUpdateMsg
			err            error
		)

		err = git.UnstageFile(path)
		if err != nil {
			msg.err = err
			return msg
		}

		workTreeStatus, err = git.Status()
		if err != nil {
			msg.err = err
			return msg
		}
		msg.workTreeStatus = workTreeStatus

		return msg
	}
}

func diff(opt git.DiffOption) func() tea.Msg {
	return func() tea.Msg {
		var (
			msg  loadedDiffMsg
			err  error
			diff string
		)

		diff, err = git.Diff(opt)
		if err != nil {
			msg.err = err
			return msg
		}
		msg.diff = diff

		return msg
	}
}
