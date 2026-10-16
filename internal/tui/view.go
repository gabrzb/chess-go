package tui

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
)

func (m model) View() tea.View {
	content := "Chess Go\n\n"

	snapshot, err := m.displayedSnapshot()
	if err != nil {
		content += "Failed to load game: " + err.Error()
	} else {
		content += m.renderBoard(snapshot.Board)
		content += "\n\n" + statusText(snapshot)

		if m.reviewing {
			content += fmt.Sprintf(
				"\n | Reviewing: position %d/%d",
				m.historyIndex+1,
				m.game.HistoryLen(),
			)
		}
	}

	if m.message != "" {
		content += "\n\n" + m.message
	}

	if len(m.promotionMoves) > 0 {
		content += "\n\nPromotion: q=Queen | r=Rook | b=Bishop | n=Kinght | Esc: cancel | Ctrl+C: quit\n"
	} else if m.reviewing {
		content += "\n\nLeft/A: previous | Right/D: next | H/Esc: cancel | q/Ctrl+C: quit\n"
	} else {
		content += "\n\nArrows/WASD: move | Enter: select/play | Esc: cancel | H: history | q/Ctrl+C: quit\n"
	}

	view := tea.NewView(content)
	view.AltScreen = true
	return view
}
