package tui

import tea "charm.land/bubbletea/v2"

func (m model) View() tea.View {
	content := "Chess Go\n\n"

	snapshot, err := m.game.Snapshot()
	if err != nil {
		content += "Failed to load game: " + err.Error()
	} else {
		content += renderBoard(snapshot.Board, m.cursor)
	}

	content += "\n\nq or Ctrl+C to quit\n"

	view := tea.NewView(content)
	view.AltScreen = true
	return view
}
