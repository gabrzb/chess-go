package attacks

import (
	"github.com/gabrzb/chess-go/internal/chess"
)

// Diagonal squares count as attacked even when empty.
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

func RookAttacksSquare(board chess.Board, from, to chess.Square) bool {
	rook, ok := board.PieceAt(from)

	if !ok || rook.Type != chess.Rook {
		return false
	}

	directions := []chess.Direction{
		{ColStep: 1, RowStep: 0},
		{ColStep: -1, RowStep: 0},
		{ColStep: 0, RowStep: 1},
		{ColStep: 0, RowStep: -1},
	}

	for _, dir := range directions {
		if SlidingAttacksSquare(board, from, to, dir) {
			return true
		}
	}

	return false
}

func BishopAttacksSquare(board chess.Board, from, to chess.Square) bool {
	bishop, ok := board.PieceAt(from)

	if !ok || bishop.Type != chess.Bishop {
		return false
	}

	directions := []chess.Direction{
		{ColStep: 1, RowStep: 1},
		{ColStep: -1, RowStep: 1},
		{ColStep: 1, RowStep: -1},
		{ColStep: -1, RowStep: -1},
	}

	for _, dir := range directions {
		if SlidingAttacksSquare(board, from, to, dir) {
			return true
		}
	}

	return false
}

func QueenAttacksSquare(board chess.Board, from, to chess.Square) bool {
	queen, ok := board.PieceAt(from)

	if !ok || queen.Type != chess.Queen {
		return false
	}

	directions := []chess.Direction{
		{ColStep: 1, RowStep: 0},
		{ColStep: -1, RowStep: 0},
		{ColStep: 0, RowStep: 1},
		{ColStep: 0, RowStep: -1},
		{ColStep: 1, RowStep: 1},
		{ColStep: -1, RowStep: 1},
		{ColStep: 1, RowStep: -1},
		{ColStep: -1, RowStep: -1},
	}

	for _, dir := range directions {
		if SlidingAttacksSquare(board, from, to, dir) {
			return true
		}
	}

	return false
}

func SlidingAttacksSquare(board chess.Board, from, to chess.Square, dir chess.Direction) bool {
	if !from.Valid() || !to.Valid() || from == to {
		return false
	}
	square, ok := from.Offset(dir.ColStep, dir.RowStep)
	for ok {
		// The first occupied square is attacked even when its piece has the same color.
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
