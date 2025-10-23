package refresh

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// Msg produced by Schedule command
type Msg time.Time

// Schedule produces a Msg after the given duration.
func Schedule(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(t time.Time) tea.Msg {
		return Msg(t)
	})
}
