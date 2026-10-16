package render

import "charm.land/lipgloss/v2"

var (
	cursorStyle    = lipgloss.NewStyle().Reverse(true)
	selectedStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#E5C07B"))
	legalMoveStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#98C379"))
)
