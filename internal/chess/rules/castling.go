package rules

import (
	"github.com/gabrzb/chess-go/internal/chess"
	"github.com/gabrzb/chess-go/internal/chess/attacks"
)

func CanCastleKingSide(position chess.Position) bool {
	if position.SideToMove != chess.White {
		return false
	}

	if !position.CastlingRights.Has(chess.WhiteKingSide) {
		return false
	}

	king, _ := position.Board.PieceAt(chess.Square(4))
	if king.Type != chess.King || king.Color != chess.White {
		return false
	}

	rook, _ := position.Board.PieceAt(chess.Square(7))
	if rook.Type != chess.Rook || rook.Color != chess.White {
		return false
	}

	for _, square := range []chess.Square{5, 6} {
		piece, _ := position.Board.PieceAt(square)

		if piece.Type != chess.NoPiece {
			return false
		}
	}

	for _, square := range []chess.Square{4, 5, 6} {
		if attacks.IsSquareAttacked(position.Board, square, chess.Black) {
			return false
		}
	}

	return true
}
