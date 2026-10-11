package attacks

import (
	"github.com/gabrzb/chess-go/internal/chess"
)

func pieceAttacksSquare(board chess.Board, from, to chess.Square) bool {
	piece, ok := board.PieceAt(from)
	if !ok {
		return false
	}

	switch piece.Type {
	case chess.Pawn:
		return PawnAttacksSquare(from, to, piece.Color)
	case chess.Knight:
		return KnightAttacksSquare(from, to)
	case chess.Bishop:
		return BishopAttacksSquare(board, from, to)
	case chess.Rook:
		return RookAttacksSquare(board, from, to)
	case chess.Queen:
		return QueenAttacksSquare(board, from, to)
	case chess.King:
		return KingAttacksSquare(from, to)
	default:
		return false
	}
}

// IsSquareAttacked reports whether a piece of byColor attacks the target.
func IsSquareAttacked(board chess.Board, square chess.Square, byColor chess.Color) bool {
	if !square.Valid() || (byColor != chess.White && byColor != chess.Black) {
		return false
	}

	for from := chess.A1; from < chess.NoSquare; from++ {
		piece, _ := board.PieceAt(from)

		if piece.Type == chess.NoPiece || piece.Color != byColor {
			continue
		}

		if pieceAttacksSquare(board, from, square) {
			return true
		}
	}

	return false
}
