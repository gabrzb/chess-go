package tui

import tea "charm.land/bubbletea/v2"

func (m model) View() tea.View {
	view := tea.NewView("Chess Go\n\nq or Ctrl+c to quit\n\n")
	view.AltScreen = true
	return view
}
