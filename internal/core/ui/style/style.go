package style

import "github.com/charmbracelet/lipgloss"

var (
	subtleColor = lipgloss.AdaptiveColor{Light: "#999999", Dark: "#5b5b5b"}

	textColor         = lipgloss.AdaptiveColor{}
	Text              = lipgloss.NewStyle().Foreground(textColor)
	inactiveTextColor = textColor
	InactiveText      = lipgloss.NewStyle().Foreground(inactiveTextColor)
	SubtleTextColor   = subtleColor
	SublteText        = lipgloss.NewStyle().Foreground(subtleColor)
	focusTextColor    = lipgloss.Color("170")
	FocusText         = lipgloss.NewStyle().Foreground(focusTextColor)
	addedTextColor    = lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}
	AddedText         = lipgloss.NewStyle().Foreground(addedTextColor)
	removedTextColor  = lipgloss.AdaptiveColor{Light: "ff6166", Dark: "#ff6961"}
	RemovedText       = lipgloss.NewStyle().Foreground(removedTextColor)

	titleBackgroundColor         = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	Title                        = lipgloss.NewStyle().Padding(0, 1).Background(titleBackgroundColor)
	inactiveTitleBackgroundColor = subtleColor
	InactiveTitle                = lipgloss.NewStyle().Padding(0, 1).Background(inactiveTitleBackgroundColor)

	inactiveBorderColor = subtleColor
	InactiveBorder      = lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).BorderForeground(inactiveBorderColor)
	FocusBorderColor    = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	FocusBorder         = lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).BorderForeground(FocusBorderColor)

	ShortHelp = lipgloss.NewStyle().Padding(0, 1, 0, 1)
)
