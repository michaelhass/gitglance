package status

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/michaelhass/gitglance/internal/git"
)

var (
	itemStyle           = lipgloss.NewStyle()
	focusedItemStyle    = itemStyle.Foreground(lipgloss.Color("170"))
	sectionStyle        = lipgloss.NewStyle().Padding(1)
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
	itemIdx        [2]int
	focusedSection section
}

func New(repo git.Repository) Model {
	return Model{repo: repo}
}

func (m Model) Init() tea.Cmd {
	return load(m.repo)
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case statusUpdateMsg:
		m.status = msg.status
		m.statusErr = msg.err
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			m.itemIdx[m.focusedSection] = max(
				m.itemIdx[m.focusedSection]-1,
				0,
			)
		case "down":
			m.itemIdx[m.focusedSection] = min(
				m.itemIdx[m.focusedSection]+1,
				m.numItems(m.focusedSection)-1,
			)
		case "tab":
			if m.focusedSection == lastSection {
				m.focusedSection = 0
			} else {
				m.focusedSection += 1
			}
		}
	}
	return m, nil
}

func (m Model) View() string {
	if m.statusErr != nil {
		return fmt.Sprint(m.statusErr)
	}

	var (
		unstagedTitle string = "Unstaged"
		stagedTitle   string = "Staged"
		itemIdx       int    = m.itemIdx[m.focusedSection]
		unstaged      string
		staged        string
	)

	switch m.focusedSection {
	case unstagedSection:
		unstaged = focusedSectionStyle.Render(fileStatusListView(m.status.Unstaged, unstagedTitle, itemIdx))
		staged = sectionStyle.Render(fileStatusListView(m.status.Staged, stagedTitle, -1))
	case stagedSection:
		unstaged = sectionStyle.Render(fileStatusListView(m.status.Unstaged, unstagedTitle, -1))
		staged = focusedSectionStyle.Render(fileStatusListView(m.status.Staged, stagedTitle, itemIdx))
	}

	return lipgloss.JoinHorizontal(lipgloss.Left, unstaged, staged)
}

func fileStatusListView(fsl git.FileStatusList, title string, focusIdx int) string {
	var builder strings.Builder

	builder.WriteString(sectionTitleStyle.Render(title))

	if len(fsl) == 0 {
		builder.WriteString("\n")
		return builder.String()
	}
	builder.WriteString("\n")
	for i, fs := range fsl {
		style := itemStyle
		if i == focusIdx {
			style = focusedItemStyle
		}
		item := style.Render(fmt.Sprintf("[%s] %s", string(fs.Code), fs.Path))
		builder.WriteString("\n")
		builder.WriteString(item)
	}

	return builder.String()
}

func (m Model) numItems(section section) int {
	if section == unstagedSection {
		return len(m.status.Unstaged)
	}
	return len(m.status.Staged)
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

func min(i, j int) int {
	if i < j {
		return i
	}
	return j
}

func max(i, j int) int {
	if i > j {
		return i
	}
	return j
}

// Msg

type statusUpdateMsg struct {
	err    error
	status git.Status
}
