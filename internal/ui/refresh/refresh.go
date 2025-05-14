package refresh

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type Msg time.Time

func Schedule(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(t time.Time) tea.Msg {
		return Msg(t)
	})
}
