package tui

import "github.com/gabrzb/chess-go/internal/chess/game"

func (m *model) toggleReview() {
	m.clearSelection()
	m.reviewing = !m.reviewing

	if m.reviewing {
		m.historyIndex = m.game.HistoryLen() - 1
	}
}

func (m model) displayedSnapshot() (game.Snapshot, error) {
	if m.reviewing {
		return m.game.HistorySnapshot(m.historyIndex)
	}

	return m.game.Snapshot()
}

func (m *model) moveHistory(step int) {
	next := m.historyIndex + step

	if next >= 0 && next < m.game.HistoryLen() {
		m.historyIndex = next
	}
}
