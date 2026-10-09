package rules

import (
	"github.com/gabrzb/chess-go/internal/chess"
	"github.com/gabrzb/chess-go/internal/chess/attacks"
)

// FindKing locates a king by color, returning NoSquare when absent or the color is invalid.
func FindKing(board chess.Board, color chess.Color) chess.Square {
	if color != chess.White && color != chess.Black {
		return chess.NoSquare
	}

	for square := chess.Square(0); square < chess.NoSquare; square++ {
		piece, _ := board.PieceAt(square)

		if piece.Type == chess.King && piece.Color == color {
			return square
		}
	}

	return chess.NoSquare
}

// IsInCheck checks attacks on the king of the given color.
func IsInCheck(board chess.Board, color chess.Color) (bool, error) {
	if color != chess.White && color != chess.Black {
		return false, chess.ErrInvalidColor
	}

	kingSquare := FindKing(board, color)
	if kingSquare == chess.NoSquare {
		return false, chess.ErrKingNotFound
	}

	enemyColor := chess.Black
	if color == chess.Black {
		enemyColor = chess.White
	}

	return attacks.IsSquareAttacked(board, kingSquare, enemyColor), nil
}
