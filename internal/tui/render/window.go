package render

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
)

func Window(content string, width, height int) tea.View {
	if lipgloss.Width(content) > width || lipgloss.Height(content) > height {
		content = ansi.Wrap("CHESS\n\nTerminal pequeno demais.\nAmplie a janela para continuar.\n\nEsc: menu · Ctrl+C: sair", max(1, width), " ")
		content = lipgloss.NewStyle().MaxWidth(width).MaxHeight(height).Render(content)
		content = Center(content, width)
	}
	surface := lipgloss.NewStyle().Foreground(lipgloss.Color(Foreground)).Background(lipgloss.Color(Background))
	content = lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, content, lipgloss.WithWhitespaceStyle(surface))
	view := tea.NewView(content)
	view.AltScreen = true
	view.BackgroundColor = lipgloss.Color(Background)
	view.ForegroundColor = lipgloss.Color(Foreground)
	return view
}
