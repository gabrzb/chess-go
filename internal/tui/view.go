package tui

import tea "charm.land/bubbletea/v2"

func (m model) View() tea.View {
	content := "Chess Go\n\n"

	snapshot, err := m.game.Snapshot()
	if err != nil {
		content += "Failed to load game: " + err.Error()
	} else {
		content += m.renderBoard(snapshot.Board)
	}

	if m.message != "" {
		content += "\n\n" + m.message
	}

	content += "\n\nArrows/WASD: move | Enter: select | Esc: cancel | q/Ctrl+C: quit\n"

	view := tea.NewView(content)
	view.AltScreen = true
	return view
}
