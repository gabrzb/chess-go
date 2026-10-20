package render

import (
	"strings"

	"charm.land/lipgloss/v2"
)

const (
	Background  = "#16161D"
	Foreground  = "#DCD7BA"
	Muted       = "#8B8B98"
	Accent      = "#7E9CD8"
	Victory     = "#98BB6C"
	Draw        = "#E6C384"
	Alert       = "#E46876"
	LightSquare = "#87927C"
	DarkSquare  = "#63715D"
	WhitePiece  = "#FFF9E8"
	BlackPiece  = "#0C100B"
	BoardMarker = "#0C100B"
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
	style := lipgloss.NewStyle().Foreground(lipgloss.Color(color))
	lines := strings.Split(text, "\n")

	for index, line := range lines {
		lines[index] = style.Render(line)
	}

	return strings.Join(lines, "\n")
}

func Center(text string, width int) string {
	return lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render(text)
}
