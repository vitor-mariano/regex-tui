package styles

import "github.com/charmbracelet/lipgloss/v2"

var (
	PrimaryColor = lipgloss.Color("12")
	MutedColor   = lipgloss.Color("240")
	LightColor   = lipgloss.Color("15")
	ErrorColor   = lipgloss.Color("9")

	InputContainerStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(MutedColor).
				Padding(0, 1)
	FocusedInputContainerStyle = InputContainerStyle.
					BorderForeground(PrimaryColor)
	ErrorInputContainerStyle = InputContainerStyle.
					BorderForeground(ErrorColor)
)
