package tui

import (
	"fmt"

	tea "charm.land/bubbletea/v2"

	"github.com/gabrzb/chess-go/internal/chess"
)

func (m model) View() tea.View {
	content := "Chess Go\n\n"

	snapshot, err := m.displayedSnapshot()
	if err != nil {
		content += "Failed to load game: " + err.Error()
	} else {
		content += m.renderBoard(snapshot.Board)
		content += "\n\n" + statusText(snapshot)
		if m.opponent != nil {
			content += "\nAI: " + colorName(m.aiColor)
			if m.thinking {
				content += " (thinking...)"
			}
		}

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
		content += "\n\nPromotion: q=Queen | r=Rook | b=Bishop | n=Knight | Esc: cancel | Ctrl+C: quit\n"
	} else if m.reviewing {
		content += "\n\nLeft/A: previous | Right/D: next | H/Esc: return to game | q/Ctrl+C: quit\n"
	} else if m.isAITurn() && err == nil && snapshot.Status == chess.StatusOngoing {
		if m.thinking {
			content += "\n\nAI is thinking | H: history | q/Ctrl+C: quit\n"
		} else {
			content += "\n\nEnter: retry AI | H: history | q/Ctrl+C: quit\n"
		}
	} else {
		content += "\n\nArrows/WASD: move | Enter: select/play | Esc: cancel | H: history | q/Ctrl+C: quit\n"
	}

	view := tea.NewView(content)
	view.AltScreen = true
	return view
}
