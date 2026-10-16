package application

import "github.com/gabrzb/chess-go/internal/tui/render"

func (m *application) refreshLabels() {
	for len(m.labels)+1 < m.match.Game().HistoryLen() {
		index := len(m.labels) + 1
		before, err := m.match.Game().HistorySnapshot(index - 1)
		if err != nil {
			m.logger.Print(err)
			return
		}
		next, err := m.match.Game().HistorySnapshot(index)
		if err != nil || next.LastMove == nil {
			return
		}
		m.labels = append(m.labels, render.SAN(before, next, *next.LastMove))
	}
}
