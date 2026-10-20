package match

import "github.com/gabrzb/chess-go/internal/domain/game"

func (m *Model) toggleReview() {
	m.ClearSelection()
	m.reviewing = !m.reviewing

	if m.reviewing {
		m.historyIndex = m.game.HistoryLen() - 1
	}
}

func (m Model) DisplayedSnapshot() (game.Snapshot, error) {
	if m.reviewing {
		return m.game.HistorySnapshot(m.historyIndex)
	}

	return m.game.Snapshot()
}

func (m *Model) moveHistory(step int) {
	next := m.historyIndex + step

	if next >= 0 && next < m.game.HistoryLen() {
		m.historyIndex = next
	}
}
