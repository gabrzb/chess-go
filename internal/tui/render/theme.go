package render

import (
	"charm.land/lipgloss/v2"
	"strings"
)

const (
	Background = "#16161D"
	Foreground = "#DCD7BA"
	Muted      = "#72727F"
	Accent     = "#7E9CD8"
	Victory    = "#98BB6C"
	Draw       = "#E6C384"
	Alert      = "#E46876"
)

var Logo = strings.Join([]string{
	" ██████╗██╗  ██╗███████╗███████╗███████╗",
	"██╔════╝██║  ██║██╔════╝██╔════╝██╔════╝",
	"██║     ███████║█████╗  ███████╗███████╗",
	"██║     ██╔══██║██╔══╝  ╚════██║╚════██║",
	"╚██████╗██║  ██║███████╗███████║███████║",
	" ╚═════╝╚═╝  ╚═╝╚══════╝╚══════╝╚══════╝",
}, "\n")

func Paint(text, color string) string {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(color)).Render(text)
}

func CenteredBlock(text string, width int) string {
	left := strings.Repeat(" ", max(0, (width-lipgloss.Width(text))/2))
	lines := strings.Split(text, "\n")
	for index, line := range lines {
		lines[index] = left + line + strings.Repeat(" ", max(0, width-len(left)-lipgloss.Width(line)))
	}
	return strings.Join(lines, "\n")
}
