package attacks

import (
	"github.com/gabrzb/chess-go/internal/chess"
)

// PawnAttacksSquare checks the attack geometry of a pawn of the given color.
func PawnAttacksSquare(from, to chess.Square, color chess.Color) bool {
	if !from.Valid() || !to.Valid() {
		return false
	}

	colDiff := to.Column() - from.Column()
	rowDiff := to.Row() - from.Row()

	switch color {
	case chess.White:
		return rowDiff == 1 && (colDiff == 1 || colDiff == -1)
	case chess.Black:
		return rowDiff == -1 && (colDiff == 1 || colDiff == -1)
	default:
		return false
	}
}

// KnightAttacksSquare checks whether the target is a knight jump from the source.
func KnightAttacksSquare(from, to chess.Square) bool {
	if !from.Valid() || !to.Valid() {
		return false
	}

	colDiff := to.Column() - from.Column()
	rowDiff := to.Row() - from.Row()

	if colDiff < 0 {
		colDiff = -colDiff
	}

	if rowDiff < 0 {
		rowDiff = -rowDiff
	}

	return (colDiff == 1 && rowDiff == 2) || (colDiff == 2 && rowDiff == 1)
}

// KingAttacksSquare checks whether the target is adjacent to the source.
func KingAttacksSquare(from, to chess.Square) bool {
	if !from.Valid() || !to.Valid() {
		return false
	}

	colDiff := to.Column() - from.Column()
	rowDiff := to.Row() - from.Row()

	if colDiff < 0 {
		colDiff = -colDiff
	}

	if rowDiff < 0 {
		rowDiff = -rowDiff
	}

	return colDiff <= 1 && rowDiff <= 1 && (colDiff != 0 || rowDiff != 0)
}

// RookAttacksSquare checks rook attacks, including path blockers.
func RookAttacksSquare(board chess.Board, from, to chess.Square) bool {
	rook, ok := board.PieceAt(from)

	if !ok || rook.Type != chess.Rook {
		return false
	}

	directions := []chess.Direction{
		{ColStep: 1, RowStep: 0},  // Right
		{ColStep: -1, RowStep: 0}, // Left
		{ColStep: 0, RowStep: 1},  // Up
		{ColStep: 0, RowStep: -1}, // Down
	}

	for _, dir := range directions {
		if SlidingAttacksSquare(board, from, to, dir) {
			return true
		}
	}

	return false
}

// BishopAttacksSquare checks bishop attacks, including path blockers.
func BishopAttacksSquare(board chess.Board, from, to chess.Square) bool {
	bishop, ok := board.PieceAt(from)

	if !ok || bishop.Type != chess.Bishop {
		return false
	}

	directions := []chess.Direction{
		{ColStep: 1, RowStep: 1},   // Up-Right
		{ColStep: -1, RowStep: 1},  // Up-Left
		{ColStep: 1, RowStep: -1},  // Down-Right
		{ColStep: -1, RowStep: -1}, // Down-Left
	}

	for _, dir := range directions {
		if SlidingAttacksSquare(board, from, to, dir) {
			return true
		}
	}

	return false
}

// QueenAttacksSquare checks queen attacks, including path blockers.
func QueenAttacksSquare(board chess.Board, from, to chess.Square) bool {
	queen, ok := board.PieceAt(from)

	if !ok || queen.Type != chess.Queen {
		return false
	}

	directions := []chess.Direction{
		{ColStep: 1, RowStep: 0},   // Right
		{ColStep: -1, RowStep: 0},  // Left
		{ColStep: 0, RowStep: 1},   // Up
		{ColStep: 0, RowStep: -1},  // Down
		{ColStep: 1, RowStep: 1},   // Up-Right
		{ColStep: -1, RowStep: 1},  // Up-Left
		{ColStep: 1, RowStep: -1},  // Down-Right
		{ColStep: -1, RowStep: -1}, // Down-Left
	}

	for _, dir := range directions {
		if SlidingAttacksSquare(board, from, to, dir) {
			return true
		}
	}

	return false
}

// SlidingAttacksSquare follows one unit direction and includes the first occupied target.
func SlidingAttacksSquare(board chess.Board, from, to chess.Square, dir chess.Direction) bool {
	if !from.Valid() || !to.Valid() || from == to {
		return false
	}
	square, ok := from.Offset(dir.ColStep, dir.RowStep)
	for ok {
		if square == to {
			return true
		}
		occupant, _ := board.PieceAt(square)
		if occupant.Type != chess.NoPiece {
			return false
		}
		square, ok = square.Offset(dir.ColStep, dir.RowStep)
	}
	return false
}
