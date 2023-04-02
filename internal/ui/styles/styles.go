package styles

import "github.com/charmbracelet/lipgloss"

var (
	subtleColor = lipgloss.AdaptiveColor{Light: "#999999", Dark: "#5b5b5b"}

	textColor         = lipgloss.AdaptiveColor{}
	TextSyle          = lipgloss.NewStyle().Foreground(textColor)
	inactiveTextColor = textColor
	InactiveTextStyle = lipgloss.NewStyle().Foreground(inactiveTextColor)
	focusTextColor    = lipgloss.Color("170")
	FocusTextStyle    = lipgloss.NewStyle().Foreground(focusTextColor)
	addedTextColor    = lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}
	AddedTextStyle    = lipgloss.NewStyle().Foreground(addedTextColor)
	removedTextColor  = lipgloss.AdaptiveColor{Light: "ff6166", Dark: "#ff6961"}
	RemovedTextStyle  = lipgloss.NewStyle().Foreground(removedTextColor)

	titleBackgroundColor         = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	TitleStyle                   = lipgloss.NewStyle().Padding(0, 1).Background(titleBackgroundColor)
	inactiveTitleBackgroundColor = subtleColor
	InactiveTitleStyle           = lipgloss.NewStyle().Padding(0, 1).Background(inactiveTitleBackgroundColor)

	inactiveBorderColor = subtleColor
	InactiveBorderStyle = lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).BorderForeground(inactiveBorderColor)
	FocusBorderColor    = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	FocusBorderStyle    = lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).BorderForeground(FocusBorderColor)

	ShortHelpStyle = lipgloss.NewStyle().Padding(0, 1, 0, 1)
)
