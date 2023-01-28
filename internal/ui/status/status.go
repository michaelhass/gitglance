package status

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/michaelhass/gitglance/internal/git"
	"github.com/michaelhass/gitglance/internal/ui/files"
)

var (
	sectionStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder())
	focusedSectionStyle = lipgloss.NewStyle().
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("69"))
	sectionTitleStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("#874BFD"))
)

type section byte

const (
	unstagedSection section = iota
	stagedSection
	lastSection section = stagedSection
)

type Model struct {
	repo           git.Repository
	status         git.Status
	statusErr      error
	stagedFiles    files.List
	unstagedFiles  files.List
	focusedSection section
	width, height  int
}

func New(repo git.Repository) Model {
	return Model{
		repo:          repo,
		unstagedFiles: files.NewList("Unstaged"),
		stagedFiles:   files.NewList("Staged"),
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
		m.unstagedFiles.SetItems(createListItems(m.status.Unstaged))
		m.stagedFiles.SetItems(createListItems(m.status.Staged))
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

	switch m.focusedSection {
	case unstagedSection:
		m.unstagedFiles.SetIsEnabled(true)
		m.stagedFiles.SetIsEnabled(false)
	case stagedSection:
		m.unstagedFiles.SetIsEnabled(false)
		m.stagedFiles.SetIsEnabled(true)
	}

	var cmds []tea.Cmd

	unstagedFiles, unstagedCmd := m.unstagedFiles.Update(msg)
	m.unstagedFiles = unstagedFiles
	cmds = append(cmds, unstagedCmd)

	stagedFiles, stagedCmd := m.stagedFiles.Update(msg)
	m.stagedFiles = stagedFiles
	cmds = append(cmds, stagedCmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.statusErr != nil {
		return fmt.Sprint(m.statusErr)
	}

	// var (
	// 	unstagedTitle string = "Unstaged"
	// 	stagedTitle   string = "Staged"
	// 	itemIdx       int    = m.itemIdx[m.focusedSection]
	// 	unstaged      string
	// 	staged        string
	// )

	// // fileSectionWidth := (m.width / 2) - 2
	// switch m.focusedSection {
	// case unstagedSection:
	// 	unstaged = focusedSectionStyle.Render(fileStatusListView(m.status.Unstaged, unstagedTitle, itemIdx))
	// 	staged = sectionStyle.Render(fileStatusListView(m.status.Staged, stagedTitle, -1))
	// case stagedSection:
	// 	unstaged = sectionStyle.Render(fileStatusListView(m.status.Unstaged, unstagedTitle, -1))
	// 	staged = focusedSectionStyle.Render(fileStatusListView(m.status.Staged, stagedTitle, itemIdx))
	// }

	// // filesSectionHeight := m.height / 3
	// // filesSection := lipgloss.NewStyle().
	// // 	//	Width(m.width).
	// // 	Height(filesSectionHeight).
	// // 	Render()

	files := lipgloss.JoinHorizontal(lipgloss.Left, m.unstagedFiles.View(), m.stagedFiles.View())
	diff := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Width(m.width - 4).Height(m.height/2 - 4).Render("")
	return lipgloss.JoinVertical(lipgloss.Top, files, diff)
	// // diffSection := sectionStyle.Width(m.width).Height(m.height - filesSectionHeight).Render("Diff")

	// // return lipgloss.JoinVertical(lipgloss.Top, filesSection, verticalSpacer, diffSection)
}

func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.unstagedFiles.SetSize(width/2-1, height/2-2)
	m.stagedFiles.SetSize(width/2-1, height/2-2)
}

func createListItems(fileStatusList git.FileStatusList) []files.ListItem {
	items := make([]files.ListItem, len(fileStatusList))
	for i, fs := range fileStatusList {
		items[i] = files.NewListItem(fs.Path, fmt.Sprintf("[%s]", string(fs.Code)))
	}
	return items
}

// Cmd

func load(repo git.Repository) func() tea.Msg {
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
