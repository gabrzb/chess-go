package movegen

import (
	"github.com/gabrzb/chess-go/internal/chess"
)

func RookMoves(board chess.Board, from chess.Square) []chess.Move {
	rook, ok := board.PieceAt(from)

	if !ok || rook.Type != chess.Rook {
		return nil
	}

	directions := []chess.Direction{
		{ColStep: 1, RowStep: 0},
		{ColStep: -1, RowStep: 0},
		{ColStep: 0, RowStep: 1},
		{ColStep: 0, RowStep: -1},
	}

	var moves []chess.Move

	for _, dir := range directions {
		directionMoves := SlidingMoves(board, from, dir)
		moves = append(moves, directionMoves...)
	}

	return moves
}

func BishopMoves(board chess.Board, from chess.Square) []chess.Move {
	bishop, ok := board.PieceAt(from)

	if !ok || bishop.Type != chess.Bishop {
		return nil
	}

	directions := []chess.Direction{
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

// King safety is checked by rules.LegalMoves after simulating the move.
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
