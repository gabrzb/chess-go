package rules

import (
	"github.com/gabrzb/chess-go/internal/chess"
	"github.com/gabrzb/chess-go/internal/chess/attacks"
)

// CanCastleKingSide checks king-side castling for the side whose turn it is.
func CanCastleKingSide(position chess.Position) bool {
	color := position.SideToMove
	var base chess.Square
	var right chess.CastlingRights
	var enemyColor chess.Color

	switch color {
	case chess.White:
		base = 0
		right = chess.WhiteKingSide
		enemyColor = chess.Black
	case chess.Black:
		base = 56
		right = chess.BlackKingSide
		enemyColor = chess.White
	default:
		return false
	}

	if !position.CastlingRights.Has(right) {
		return false
	}

	king, _ := position.Board.PieceAt(base + 4)
	if king.Type != chess.King || king.Color != color {
		return false
	}

	rook, _ := position.Board.PieceAt(base + 7)
	if rook.Type != chess.Rook || rook.Color != color {
		return false
	}

	for _, square := range []chess.Square{base + 5, base + 6} {
		piece, _ := position.Board.PieceAt(square)

		if piece.Type != chess.NoPiece {
			return false
		}
	}

	for _, square := range []chess.Square{base + 4, base + 5, base + 6} {
		if attacks.IsSquareAttacked(position.Board, square, enemyColor) {
			return false
		}
	}

	return true
}

func KingSideCastlingMove(position chess.Position) (chess.Move, bool) {
	if !CanCastleKingSide(position) {
		return chess.Move{}, false
	}

	from := chess.Square(4)
	if position.SideToMove == chess.Black {
		from = 60
	}

	return chess.Move{From: from, To: from + 2}, true
}
