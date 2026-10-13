package tui

import (
	"slices"

	"github.com/gabrzb/chess-go/internal/chess"
)

func (m *model) selectPiece() error {
	moves, err := m.game.LegalMovesFrom(m.cursor)
	if err != nil {
		return err
	}

	m.selected = m.cursor
	m.legalMoves = moves

	return nil
}

func (m *model) clearSelection() {
	m.selected = chess.NoSquare
	m.legalMoves = nil
	m.message = ""
	m.promotionMoves = nil
}

func (m model) isLegalDestination(square chess.Square) bool {
	return slices.ContainsFunc(m.legalMoves, func(move chess.Move) bool {
		return move.To == square
	})
}
