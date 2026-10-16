package match

import "github.com/gabrzb/chess-go/internal/chess"

func (m *Model) selectPiece() error {
	moves, err := m.game.LegalMovesFrom(m.cursor)
	if err != nil {
		return err
	}

	m.selected = m.cursor
	m.legalMoves = moves

	return nil
}

func (m *Model) ClearSelection() {
	m.selected = chess.NoSquare
	m.legalMoves = nil
	m.message = ""
	m.promotionMoves = nil
}
