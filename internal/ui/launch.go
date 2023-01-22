package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/michaelhass/gitglance/internal/git"
	"github.com/michaelhass/gitglance/internal/ui/app"
)

type LaunchOptions struct {
	Path    string
	RepoOpt git.RepositoryOpt
}

func LaunchApp(opt LaunchOptions) error {
	repo, err := git.OpenRepository(opt.Path, git.RepositoryOpt{ImplType: git.GoGit})
	if err != nil {
		return err
	}

	model := app.New(repo)
	if _, err := tea.NewProgram(model, tea.WithAltScreen()).Run(); err != nil {
		return err
	}
	return nil
}