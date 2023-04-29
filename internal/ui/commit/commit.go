package commit

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/michaelhass/gitglance/internal/git"
	"github.com/michaelhass/gitglance/internal/ui/container"
	"github.com/michaelhass/gitglance/internal/ui/filelist"
)

type Model struct {
	stagedFileList container.Model
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
			items[i] = filelist.NewItem(fs)
		}
		return items
	}

	fileListContent.Model = fileListContent.SetItems(createListItems(stagedFileList))

	return Model{
		stagedFileList: container.New(fileListContent),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	return m, nil
}

func (m Model) View() string {
	return m.stagedFileList.View()
}

func (m Model) SetSize(width, height int) Model {
	m.stagedFileList = m.stagedFileList.SetSize(width, height/3)
	return m
}
