package movegen

import (
	"github.com/gabrzb/chess-go/internal/chess"
)

// RookMoves generates basic pseudolegal rook moves.
func RookMoves(board chess.Board, from chess.Square) []chess.Move {
	rook, ok := board.PieceAt(from)

	if !ok || rook.Type != chess.Rook {
		return nil
	}

	directions := []chess.Direction{
		{ColStep: 1, RowStep: 0},  // Right
		{ColStep: -1, RowStep: 0}, // Left
		{ColStep: 0, RowStep: 1},  // Up
		{ColStep: 0, RowStep: -1}, // Down
	}

	var moves []chess.Move

	for _, dir := range directions {
		directionMoves := SlidingMoves(board, from, dir)
		moves = append(moves, directionMoves...)
	}

	return moves
}

// BishopMoves generates basic pseudolegal bishop moves.
func BishopMoves(board chess.Board, from chess.Square) []chess.Move {
	bishop, ok := board.PieceAt(from)

	if !ok || bishop.Type != chess.Bishop {
		return nil
	}

	directions := []chess.Direction{
		{ColStep: 1, RowStep: 1},   // Up-Right
		{ColStep: -1, RowStep: 1},  // Up-Left
		{ColStep: 1, RowStep: -1},  // Down-Right
		{ColStep: -1, RowStep: -1}, // Down-Left
	}

	var moves []chess.Move

	for _, dir := range directions {
		directionMoves := SlidingMoves(board, from, dir)
		moves = append(moves, directionMoves...)
	}

	return moves
}

// QueenMoves generates basic pseudolegal queen moves.
func QueenMoves(board chess.Board, from chess.Square) []chess.Move {
	queen, ok := board.PieceAt(from)

	if !ok || queen.Type != chess.Queen {
		return nil
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

	var moves []chess.Move

	for _, dir := range directions {
		directionMoves := SlidingMoves(board, from, dir)
		moves = append(moves, directionMoves...)
	}

	return moves
}

// KnightMoves generates basic pseudolegal knight moves.
func KnightMoves(board chess.Board, from chess.Square) []chess.Move {
	knight, ok := board.PieceAt(from)

	if !ok || knight.Type != chess.Knight {
		return nil
	}

	directions := []chess.Direction{
		{ColStep: 2, RowStep: 1},
		{ColStep: 2, RowStep: -1},
		{ColStep: -2, RowStep: 1},
		{ColStep: -2, RowStep: -1},
		{ColStep: 1, RowStep: 2},
		{ColStep: -1, RowStep: 2},
		{ColStep: 1, RowStep: -2},
		{ColStep: -1, RowStep: -2},
	}

	var moves []chess.Move

	for _, dir := range directions {
		directionMoves := JumpMoves(board, from, dir)
		moves = append(moves, directionMoves...)
	}

	return moves
}

// KingMoves generates adjacent king moves before filtering attacked squares.
func KingMoves(board chess.Board, from chess.Square) []chess.Move {
	king, ok := board.PieceAt(from)

	if !ok || king.Type != chess.King {
		return nil
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

	var moves []chess.Move

	for _, dir := range directions {
		directionMoves := JumpMoves(board, from, dir)
		moves = append(moves, directionMoves...)
	}

	return moves
}

// PawnMoves generates basic pawn advances and diagonal captures.
func PawnMoves(board chess.Board, from chess.Square) []chess.Move {
	pawn, ok := board.PieceAt(from)
	if !ok || pawn.Type != chess.Pawn {
		return nil
	}
	var rowStep int
	startRow := chess.A2.Row()
	switch pawn.Color {
	case chess.White:
		rowStep = 1
	case chess.Black:
		rowStep = -1
		startRow = chess.A7.Row()
	default:
		return nil
	}
	destination, ok := from.Offset(0, rowStep)
	if !ok {
		return nil
	}
	occupant, _ := board.PieceAt(destination)
	var moves []chess.Move
	if occupant.Type == chess.NoPiece {
		moves = append(moves, chess.Move{From: from, To: destination})
		if from.Row() == startRow {
			doubleDestination, ok := destination.Offset(0, rowStep)
			if ok {
				doubleOccupant, _ := board.PieceAt(doubleDestination)
				if doubleOccupant.Type == chess.NoPiece {
					moves = append(moves, chess.Move{From: from, To: doubleDestination})
				}
			}
		}
	}
	for _, colStep := range []int{-1, 1} {
		captureDestination, ok := from.Offset(colStep, rowStep)
		if !ok {
			continue
		}
		target, _ := board.PieceAt(captureDestination)
		if target.Type != chess.NoPiece && target.Color != pawn.Color {
			moves = append(moves, chess.Move{From: from, To: captureDestination})
		}
	}
	return moves
}
