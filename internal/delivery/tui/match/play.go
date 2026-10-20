package match

import chess "github.com/gabrzb/chess-go/internal/domain"

func (m *Model) playSelectedMove() error {
	if m.preparePromotion() {
		return nil
	}

	move := chess.Move{
		From: m.selected,
		To:   m.cursor,
	}

	return m.playMove(move)
}

func (m *Model) playMove(move chess.Move) error {
	if _, err := m.game.Play(move); err != nil {
		return err
	}

	m.ClearSelection()
	return nil
}
