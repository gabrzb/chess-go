package game

func (g *Game) HistoryLen() int {
	return len(g.history) + 1
}

func (g *Game) HistorySnapshot(index int) (Snapshot, error) {
	if index < 0 || index >= g.HistoryLen() {
		return Snapshot{}, ErrInvalidHistoryIndex
	}

	if index == len(g.history) {
		return g.Snapshot()
	}

	snapshot := g.history[index]

	if snapshot.LastMove != nil {
		move := *snapshot.LastMove
		snapshot.LastMove = &move
	}

	return snapshot, nil
}
