package tui

import "github.com/gabrzb/chess-go/internal/chess"

func (m *model) preparePromotion() bool {
	m.promotionMoves = nil

	for _, move := range m.legalMoves {
		if move.To == m.cursor && move.Promotion != chess.NoPiece {
			m.promotionMoves = append(m.promotionMoves, move)
		}
	}

	return len(m.promotionMoves) > 0
}

func (m *model) choosePromotion(pieceType chess.PieceType) error {
	for _, move := range m.promotionMoves {
		if move.Promotion == pieceType {
			return m.playMove(move)
		}
	}

	return chess.ErrIllegalMove
}

func (m *model) handlePromotionKey(key string) {
	if key == "esc" {
		m.clearSelection()
		return
	}

	choices := map[string]chess.PieceType{
		"q": chess.Queen,
		"r": chess.Rook,
		"b": chess.Bishop,
		"n": chess.Knight,
	}

	pieceType, ok := choices[key]
	if !ok {
		return
	}

	if err := m.choosePromotion(pieceType); err != nil {
		m.message = err.Error()
	}
}
