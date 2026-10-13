package tui

import "github.com/gabrzb/chess-go/internal/chess"

func (m *model) playSelectedMove() error {
	if m.preparePromotion() {
		return nil
	}

	move := chess.Move{
		From: m.selected,
		To:   m.cursor,
	}

	return m.playMove(move)
}

func (m *model) playMove(move chess.Move) error {
	if _, err := m.game.Play(move); err != nil {
		return err
	}

	m.clearSelection()
	return nil
}
