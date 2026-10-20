package rules

import chess "github.com/gabrzb/chess-go/internal/domain"

// InsufficientMaterial recognizes common dead positions. Other blocked positions
// may also be dead; deciding those requires more than counting material.
func InsufficientMaterial(board chess.Board) bool {
	minors, knights, bishopColor := 0, 0, -1
	for square := chess.A1; square <= chess.H8; square++ {
		piece, _ := board.PieceAt(square)
		switch piece.Type {
		case chess.NoPiece, chess.King:
		case chess.Knight:
			minors++
			knights++
		case chess.Bishop:
			minors++
			color := (square.Column() + square.Row()) % 2
			if bishopColor >= 0 && bishopColor != color {
				return false
			}
			bishopColor = color
		default:
			return false
		}
	}
	return minors <= 1 || knights == 0
}

// RepetitionKey excludes clocks and ignores en passant unless a legal capture
// exists. A pinned pawn cannot make two otherwise equal positions different.
func RepetitionKey(position chess.Position) chess.Position {
	position.HalfmoveClock = 0
	position.FullmoveNumber = 0
	if position.EnPassant.Valid() {
		for square := chess.A1; square <= chess.H8; square++ {
			if _, ok := EnPassantMove(position, square); ok {
				return position
			}
		}
		position.EnPassant = chess.NoSquare
	}
	return position
}
