package match

import chess "github.com/gabrzb/chess-go/internal/domain"

func (m *Model) preparePromotion() bool {
	m.promotionMoves = nil

	for _, move := range m.legalMoves {
		if move.To == m.cursor && move.Promotion != chess.NoPiece {
			m.promotionMoves = append(m.promotionMoves, move)
		}
	}

	return len(m.promotionMoves) > 0
}

func (m *Model) choosePromotion(pieceType chess.PieceType) error {
	for _, move := range m.promotionMoves {
		if move.Promotion == pieceType {
			return m.playMove(move)
		}
	}

	return chess.ErrIllegalMove
}

func (m *Model) PromotionKey(key string) {
	if key == "esc" {
		m.ClearSelection()
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
