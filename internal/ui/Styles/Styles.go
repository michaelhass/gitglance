package styles

import "github.com/charmbracelet/lipgloss"

var (
	TextColor            = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	TitleBackgroundColor = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	TitleStyle           = lipgloss.NewStyle().Padding(0, 1).Background(TitleBackgroundColor)

	BorderColor      = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	FocusBorderColor = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	BorderStyle      = lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).BorderForeground(BorderColor)
	FocusBorderStyle = lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).BorderForeground(FocusBorderColor)
)
