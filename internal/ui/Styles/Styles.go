package styles

import "github.com/charmbracelet/lipgloss"

var (
	subtleColor = lipgloss.AdaptiveColor{Light: "#999999", Dark: "#5b5b5b"}

	TextColor         = lipgloss.AdaptiveColor{}
	TextSyle          = lipgloss.NewStyle().Foreground(TextColor)
	InactiveTextColor = subtleColor
	InactiveTextStyle = lipgloss.NewStyle().Foreground(InactiveTextColor)
	FocusTextColor    = lipgloss.Color("170")
	FocusTextStyle    = lipgloss.NewStyle().Foreground(FocusTextColor)

	TitleBackgroundColor         = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	TitleStyle                   = lipgloss.NewStyle().Padding(0, 1).Background(TitleBackgroundColor)
	InactiveTitleBackgroundColor = subtleColor
	InactiveTitleStyle           = lipgloss.NewStyle().Padding(0, 1).Background(InactiveTitleBackgroundColor)

	inactiveBorderColor = subtleColor
	InactiveBorderStyle = lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).BorderForeground(inactiveBorderColor)
	FocusBorderColor    = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	FocusBorderStyle    = lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).BorderForeground(FocusBorderColor)
)
